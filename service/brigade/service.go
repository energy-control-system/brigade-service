package brigade

import (
	"brigade-service/cluster/task"
	"brigade-service/cluster/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/sunshineOfficial/golib/pagination"
)

const kafkaSubscribeTimeout = 2 * time.Minute

type Service struct {
	repository  Repository
	userService UserService
}

func NewService(repository Repository, userService UserService) *Service {
	return &Service{
		repository:  repository,
		userService: userService,
	}
}

func (s *Service) CreateBrigade(ctx goctx.Context, request CreateBrigadeRequest) (Brigade, error) {
	users, err := s.userService.GetUsersByIDs(ctx, request.InspectorIDs)
	if err != nil {
		return Brigade{}, fmt.Errorf("get users: %w", err)
	}

	if len(users) != len(request.InspectorIDs) {
		return Brigade{}, fmt.Errorf("want %d users, got %d", len(request.InspectorIDs), len(users))
	}

	inspectors := make([]Inspector, 0, len(users))
	for _, u := range users {
		if u.Role != user.RoleInspector {
			return Brigade{}, fmt.Errorf("user %d is not inspector: role = %d", u.ID, u.Role)
		}

		inspectors = append(inspectors, MapUserToInspector(u))
	}

	b, err := s.repository.CreateBrigade(ctx, inspectors)
	if err != nil {
		return Brigade{}, fmt.Errorf("create brigade in repository: %w", err)
	}

	return b, nil
}

func (s *Service) GetBrigadeByID(ctx goctx.Context, id int) (Brigade, error) {
	b, err := s.repository.GetBrigadeByID(ctx, id)
	if err != nil {
		return Brigade{}, fmt.Errorf("get brigade by id from repository: %w", err)
	}

	userIDs := make([]int, 0, len(b.Inspectors))
	for _, u := range b.Inspectors {
		userIDs = append(userIDs, u.ID)
	}

	users, err := s.userService.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return Brigade{}, fmt.Errorf("get users: %w", err)
	}

	if len(users) != len(b.Inspectors) {
		return Brigade{}, fmt.Errorf("want %d users, got %d", len(b.Inspectors), len(users))
	}

	userMap := make(map[int]user.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	for i, inspector := range b.Inspectors {
		u, ok := userMap[inspector.ID]
		if !ok {
			return Brigade{}, fmt.Errorf("inspector %d not found", inspector.ID)
		}

		fullInspector := MapUserToInspector(u)
		fullInspector.AssignedAt = inspector.AssignedAt

		b.Inspectors[i] = fullInspector
	}

	return b, nil
}

func (s *Service) GetAllBrigades(ctx goctx.Context, page pagination.Pagination) ([]Brigade, error) {
	if err := page.Validate(); err != nil {
		return nil, fmt.Errorf("validate pagination: %w", err)
	}

	brigades, err := s.repository.GetAllBrigades(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("get all brigades from repository: %w", err)
	}

	userIDs := make([]int, 0, len(brigades))
	for _, b := range brigades {
		for _, u := range b.Inspectors {
			userIDs = append(userIDs, u.ID)
		}
	}

	users, err := s.userService.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	userMap := make(map[int]user.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	for i, b := range brigades {
		for j, inspector := range b.Inspectors {
			u, ok := userMap[inspector.ID]
			if !ok {
				return nil, fmt.Errorf("inspector %d not found", inspector.ID)
			}

			fullInspector := MapUserToInspector(u)
			fullInspector.AssignedAt = inspector.AssignedAt

			brigades[i].Inspectors[j] = fullInspector
		}
	}

	return brigades, nil
}

func (s *Service) ArchiveBrigade(ctx goctx.Context, id int) error {
	if err := s.repository.UpdateBrigadeStatus(ctx, id, StatusArchived); err != nil {
		return fmt.Errorf("archive brigade in repository: %w", err)
	}

	return nil
}

func (s *Service) SubscriberOnTaskEvent(mainCtx context.Context, log golog.Logger) gokafka.Subscriber {
	return func(message gokafka.Message, err error) {
		ctx, cancel := context.WithTimeout(mainCtx, kafkaSubscribeTimeout)
		defer cancel()

		if err != nil {
			log.Errorf("got error on task event: %v", err)
			return
		}

		var event task.Event
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Errorf("failed to unmarshal task event: %v", err)
			return
		}

		switch event.Type {
		case task.EventTypeAdd:
			err = s.handleAddedTask(ctx, event.Task)
		case task.EventTypeStart:
			err = s.handleStartedTask(ctx, event.Task)
		case task.EventTypeFinish:
			err = s.handleFinishedTask(ctx, event.Task)
		default:
			err = fmt.Errorf("unknown event type: %v", event.Type)
		}

		if err != nil {
			log.Errorf("failed to handle task event (type = %d): %v", event.Type, err)
			return
		}
	}
}

func (s *Service) handleAddedTask(ctx context.Context, t task.Task) error {
	return nil
}

func (s *Service) handleStartedTask(ctx context.Context, t task.Task) error {
	if t.Status != task.StatusInWork {
		return fmt.Errorf("invalid task status: %v", t.Status)
	}

	if t.BrigadeID == nil {
		return errors.New("missing brigadeID")
	}

	err := s.repository.UpdateBrigadeStatus(ctx, *t.BrigadeID, StatusOnTask)
	if err != nil {
		return fmt.Errorf("update brigade status: %v", err)
	}

	return nil
}

func (s *Service) handleFinishedTask(ctx context.Context, t task.Task) error {
	if t.Status != task.StatusDone {
		return fmt.Errorf("invalid task status: %v", t.Status)
	}

	if t.BrigadeID == nil {
		return errors.New("missing brigadeID")
	}

	err := s.repository.UpdateBrigadeStatus(ctx, *t.BrigadeID, StatusIdle)
	if err != nil {
		return fmt.Errorf("update brigade status: %v", err)
	}

	return nil
}
