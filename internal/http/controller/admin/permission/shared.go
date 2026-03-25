package permission

import (
	"time"

	gormadapter "github.com/casbin/gorm-adapter/v3"
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

type HttpApiData struct {
	Path   string `json:"path" mapstructure:"path" validate:"required"`
	Method string `json:"method" mapstructure:"method" validate:"required"`
}

func (data *HttpApiData) Fill(casbinRule gormadapter.CasbinRule) {
	data.Path = casbinRule.V1
	data.Method = casbinRule.V2
}

type RoleData struct {
	Id          uint             `json:"id" mapstructure:"id" validate:"required"`
	Name        string           `json:"name" mapstructure:"name" validate:"required"`
	IsPublic    bool             `json:"is_public" mapstructure:"is_public" validate:"required"`
	CreatedAt   time.Time        `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt   time.Time        `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
	Permissions []PermissionData `json:"permissions" mapstructure:"permissions" validate:"required"`
}

func (data *RoleData) Fill(m *model.Role) {
	data.Id = m.Id
	data.Name = m.Name
	data.IsPublic = m.IsPublic
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
	data.Permissions = make([]PermissionData, len(m.Permissions))

	for i := 0; i < len(m.Permissions); i++ {
		data.Permissions[i].Fill(&m.Permissions[i])
	}
}
