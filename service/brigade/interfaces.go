package brigade

import (
	"brigade-service/cluster/user"
	"context"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type Repository interface {
	CreateBrigade(ctx context.Context, inspectors []Inspector) (Brigade, error)
	GetBrigadeByID(ctx context.Context, id int) (Brigade, error)
	GetAllBrigades(ctx context.Context, page pagination.Pagination) ([]Brigade, error)
	UpdateBrigadeStatus(ctx context.Context, brigadeID int, newStatus Status) error
}

type UserService interface {
	GetUsersByIDs(ctx goctx.Context, userIDs []int) ([]user.User, error)
}
