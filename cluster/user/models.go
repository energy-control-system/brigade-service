package user

import "time"

type Role int

const (
	RoleUnknown Role = iota
	RoleInspector
	RoleDispatcher
	RoleSpecialist
)

type User struct {
	ID         int       `json:"id"`
	Role       Role      `json:"role"`
	Surname    string    `json:"surname"`
	Name       string    `json:"name"`
	Patronymic string    `json:"patronymic"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
