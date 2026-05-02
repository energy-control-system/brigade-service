package brigade

import (
	"brigade-service/cluster/user"
	"context"
	"errors"
	"testing"

	"github.com/sunshineOfficial/golib/goctx"
)

type mockUserService struct {
	users []user.User
	err   error
}

func (m *mockUserService) GetUsersByIDs(_ goctx.Context, _ []int) ([]user.User, error) {
	return m.users, m.err
}

type mockRepository struct {
	brigade Brigade
	err     error
}

func (m *mockRepository) CreateBrigade(_ context.Context, inspectors []Inspector) (Brigade, error) {
	if m.err != nil {
		return Brigade{}, m.err
	}
	m.brigade.Inspectors = inspectors
	return m.brigade, nil
}

func (m *mockRepository) GetBrigadeByID(_ context.Context, _ int) (Brigade, error) {
	return m.brigade, m.err
}

func (m *mockRepository) GetAllBrigades(_ context.Context) ([]Brigade, error) {
	return nil, m.err
}

func (m *mockRepository) UpdateBrigadeStatus(_ context.Context, _ int, _ Status) error {
	return m.err
}

func TestCreateBrigade_Success(t *testing.T) {
	svc := NewService(
		&mockRepository{brigade: Brigade{ID: 1, Status: StatusIdle}},
		&mockUserService{users: []user.User{
			{ID: 10, Role: user.RoleInspector, Surname: "Иванов", Name: "Иван"},
			{ID: 20, Role: user.RoleInspector, Surname: "Петров", Name: "Пётр"},
		}},
	)

	b, err := svc.CreateBrigade(goctx.Wrap(context.Background()), CreateBrigadeRequest{
		InspectorIDs: []int{10, 20},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID != 1 {
		t.Fatalf("expected brigade ID = 1, got %d", b.ID)
	}
	if len(b.Inspectors) != 2 {
		t.Fatalf("expected 2 inspectors, got %d", len(b.Inspectors))
	}
}

func TestCreateBrigade_UserServiceError(t *testing.T) {
	svc := NewService(
		&mockRepository{},
		&mockUserService{err: errors.New("user service unavailable")},
	)

	_, err := svc.CreateBrigade(goctx.Wrap(context.Background()), CreateBrigadeRequest{
		InspectorIDs: []int{10},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateBrigade_UserCountMismatch(t *testing.T) {
	svc := NewService(
		&mockRepository{},
		&mockUserService{users: []user.User{
			{ID: 10, Role: user.RoleInspector},
		}},
	)

	_, err := svc.CreateBrigade(goctx.Wrap(context.Background()), CreateBrigadeRequest{
		InspectorIDs: []int{10, 20},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateBrigade_UserNotInspector(t *testing.T) {
	svc := NewService(
		&mockRepository{},
		&mockUserService{users: []user.User{
			{ID: 10, Role: user.RoleInspector},
			{ID: 20, Role: user.RoleDispatcher},
		}},
	)

	_, err := svc.CreateBrigade(goctx.Wrap(context.Background()), CreateBrigadeRequest{
		InspectorIDs: []int{10, 20},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateBrigade_RepositoryError(t *testing.T) {
	svc := NewService(
		&mockRepository{err: errors.New("database error")},
		&mockUserService{users: []user.User{
			{ID: 10, Role: user.RoleInspector},
		}},
	)

	_, err := svc.CreateBrigade(goctx.Wrap(context.Background()), CreateBrigadeRequest{
		InspectorIDs: []int{10},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
