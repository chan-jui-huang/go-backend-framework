package httpapi

import (
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
)

type HttpApiData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Method    string    `json:"method" mapstructure:"method" validate:"required"`
	Path      string    `json:"path" mapstructure:"path" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
}

func (data *HttpApiData) Fill(m *model.HttpApi) {
	data.Id = m.Id
	data.Method = m.Method
	data.Path = m.Path
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
}

type Rule struct {
	Object string `json:"object" binding:"required,startswith=/"`
	Action string `json:"action" binding:"required,oneof=GET POST PUT PATCH DELETE"`
}
