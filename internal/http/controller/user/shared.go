package user

import (
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
)

type PermissionData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Name      string    `json:"name" mapstructure:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
}

func (data *PermissionData) Fill(m *model.Permission) {
	data.Id = m.Id
	data.Name = m.Name
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
}

type RoleData struct {
	Id          uint             `json:"id" mapstructure:"id" validate:"required"`
	Name        string           `json:"name" mapstructure:"name" validate:"required"`
	CreatedAt   time.Time        `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt   time.Time        `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
	Permissions []PermissionData `json:"permissions" mapstructure:"permissions" validate:"required"`
}

func (data *RoleData) Fill(m *model.Role) {
	data.Id = m.Id
	data.Name = m.Name
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
	data.Permissions = make([]PermissionData, len(m.Permissions))

	for i := 0; i < len(m.Permissions); i++ {
		data.Permissions[i].Fill(&m.Permissions[i])
	}
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

type TokenData struct {
	AccessToken string `json:"access_token" mapstructure:"access_token" validate:"required"`
}

func (data *TokenData) Fill(accessToken string) {
	data.AccessToken = accessToken
}
