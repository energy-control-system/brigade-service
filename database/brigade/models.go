package brigade

import "time"

type Brigade struct {
	ID        int       `db:"id"`
	Status    int       `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Member struct {
	BrigadeID   int       `db:"brigade_id"`
	InspectorID int       `db:"inspector_id"`
	AssignedAt  time.Time `db:"assigned_at"`
}
