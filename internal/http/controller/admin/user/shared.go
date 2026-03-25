package user

import (
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
)

type RoleData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Name      string    `json:"name" mapstructure:"name" validate:"required"`
	IsPublic  bool      `json:"is_public" mapstructure:"is_public" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
}

func (data *RoleData) Fill(m *model.Role) {
	data.Id = m.Id
	data.Name = m.Name
	data.IsPublic = m.IsPublic
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
}

type UserData struct {
	Id        uint       `json:"id" mapstructure:"id" validate:"required"`
	Name      string     `json:"name" mapstructure:"name" validate:"required"`
	Email     string     `json:"email" mapstructure:"email" validate:"required"`
	CreatedAt time.Time  `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt time.Time  `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
	Roles     []RoleData `json:"roles" mapstructure:"roles" validate:"required"`
}

func (data *UserData) Fill(m *model.User) {
	data.Id = m.Id
	data.Name = m.Name
	data.Email = m.Email
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
	data.Roles = make([]RoleData, len(m.Roles))

	for i := 0; i < len(data.Roles); i++ {
		data.Roles[i].Fill(&m.Roles[i])
	}
}
