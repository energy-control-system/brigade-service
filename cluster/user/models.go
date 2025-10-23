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
	ID          int       `json:"ID"`
	Role        Role      `json:"Role"`
	Surname     string    `json:"Surname"`
	Name        string    `json:"Name"`
	Patronymic  string    `json:"Patronymic"`
	PhoneNumber string    `json:"PhoneNumber"`
	Email       string    `json:"Email"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
