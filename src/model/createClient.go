package model

type CreateClientModel struct {
	Name         string `json:"name" validate:"required,min=3,max=255"`
	WorkingHours string `json:"workingHours" validate:"required"`
	LeadCapacity int    `json:"leadCapacity" validate:"required,min=1"`
	Priority     int    `json:"priority" validate:"min=0"`
}
