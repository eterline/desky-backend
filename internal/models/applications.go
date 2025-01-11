package models

type (
	AppsTable map[string][]AppDetails

	AppDetails struct {
		Name        string `json:"name" validate:"required,min=3,max=20"`
		Description string `json:"description" validate:"required"`
		Link        string `json:"link" validate:"required, url"`
		Icon        string `json:"icon" validate:"required"`
	}
)
