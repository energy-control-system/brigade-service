package brigade

import (
	"brigade-service/service/brigade"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/pagination"
)

var (
	//go:embed sql/add_member.sql
	addMemberSQL string

	//go:embed sql/create_brigade.sql
	createBrigadeSQL string

	//go:embed sql/get_all_brigades.sql
	getAllBrigadesSQL string

	//go:embed sql/get_all_members.sql
	getAllMembersSQL string

	//go:embed sql/get_brigade_by_id.sql
	getBrigadeByIDSQL string

	//go:embed sql/get_members_by_brigade.sql
	getMembersByBrigadeSQL string

	//go:embed sql/update_brigade_status.sql
	updateBrigadeStatusSQL string
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateBrigade(ctx context.Context, inspectors []brigade.Inspector) (brigade.Brigade, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return brigade.Brigade{}, fmt.Errorf("r.db.BeginTxx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	var dbBrigade Brigade
	err = tx.GetContext(ctx, &dbBrigade, createBrigadeSQL)
	if err != nil {
		err = fmt.Errorf("tx.GetContext: %w", err)
		return brigade.Brigade{}, err
	}

	members := MapInspectorsToMembers(inspectors, dbBrigade.ID)

	_, err = tx.NamedExecContext(ctx, addMemberSQL, members)
	if err != nil {
		err = fmt.Errorf("tx.NamedExecContext: %w", err)
		return brigade.Brigade{}, err
	}

	err = tx.SelectContext(ctx, &members, getMembersByBrigadeSQL, dbBrigade.ID)
	if err != nil {
		err = fmt.Errorf("tx.SelectContext: %w", err)
		return brigade.Brigade{}, err
	}

	memberMap := make(map[int]Member, len(members))
	for _, member := range members {
		memberMap[member.InspectorID] = member
	}

	b := MapBrigadeFromDB(dbBrigade)
	b.Inspectors = make([]brigade.Inspector, 0, len(inspectors))

	for _, inspector := range inspectors {
		member, ok := memberMap[inspector.ID]
		if !ok {
			err = fmt.Errorf("member not found for inspector ID %d", inspector.ID)
			return brigade.Brigade{}, err
		}

		inspector.AssignedAt = member.AssignedAt

		b.Inspectors = append(b.Inspectors, inspector)
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("tx.Commit: %w", err)
		return brigade.Brigade{}, err
	}

	return b, err
}

func (r *Repository) GetBrigadeByID(ctx context.Context, id int) (brigade.Brigade, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return brigade.Brigade{}, fmt.Errorf("r.db.BeginTxx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	var dbBrigade Brigade
	err = tx.GetContext(ctx, &dbBrigade, getBrigadeByIDSQL, id)
	if err != nil {
		err = fmt.Errorf("tx.GetContext: %w", err)
		return brigade.Brigade{}, err
	}

	var members []Member
	err = tx.SelectContext(ctx, &members, getMembersByBrigadeSQL, dbBrigade.ID)
	if err != nil {
		err = fmt.Errorf("tx.SelectContext: %w", err)
		return brigade.Brigade{}, err
	}

	b := MapBrigadeFromDB(dbBrigade)
	b.Inspectors = MapMembersToInspectors(members)

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("tx.Commit: %w", err)
		return brigade.Brigade{}, err
	}

	return b, err
}

func (r *Repository) GetAllBrigades(ctx context.Context, page pagination.Pagination) ([]brigade.Brigade, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("r.db.BeginTxx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	var dbBrigades []Brigade
	err = tx.SelectContext(ctx, &dbBrigades, getAllBrigadesSQL, page.LimitArg(), page.Offset)
	if err != nil {
		err = fmt.Errorf("tx.SelectContext: %w", err)
		return nil, err
	}
	if len(dbBrigades) == 0 {
		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("tx.Commit: %w", err)
			return nil, err
		}
		return []brigade.Brigade{}, nil
	}

	brigadeIDs := make([]int, 0, len(dbBrigades))
	for _, dbBrigade := range dbBrigades {
		brigadeIDs = append(brigadeIDs, dbBrigade.ID)
	}

	var members []Member
	err = tx.SelectContext(ctx, &members, getAllMembersSQL, brigadeIDs)
	if err != nil {
		err = fmt.Errorf("tx.SelectContext: %w", err)
		return nil, err
	}

	memberMap := make(map[int][]Member, len(dbBrigades))
	for _, member := range members {
		memberMap[member.BrigadeID] = append(memberMap[member.BrigadeID], member)
	}

	brigades := MapBrigadesFromDB(dbBrigades)
	for i, b := range brigades {
		brigadeMembers, ok := memberMap[b.ID]
		if !ok || len(brigadeMembers) == 0 {
			err = fmt.Errorf("members not found for brigade ID %d", b.ID)
			return nil, err
		}

		brigades[i].Inspectors = MapMembersToInspectors(brigadeMembers)
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("tx.Commit: %w", err)
		return nil, err
	}

	return brigades, err
}

func (r *Repository) UpdateBrigadeStatus(ctx context.Context, brigadeID int, newStatus brigade.Status) error {
	_, err := r.db.ExecContext(ctx, updateBrigadeStatusSQL, brigadeID, newStatus)
	if err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}

	return nil
}
