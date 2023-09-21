package transport

import (
	"github.com/google/uuid"
	"github.com/undefeel/cloud-storage-backend/internal/services"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ConvertTransportUserToServicesUser(us User) *services.User {
	return &services.User{
		Id:       uuid.New(),
		Login:    us.Login,
		Password: us.Password,
	}
}
