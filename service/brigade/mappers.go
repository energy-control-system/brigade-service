package brigade

import (
	"brigade-service/cluster/user"
)

func MapUserToInspector(u user.User) Inspector {
	return Inspector{
		ID:          u.ID,
		Surname:     u.Surname,
		Name:        u.Name,
		Patronymic:  u.Patronymic,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
