package models

type (
	AppsTable map[string][]AppDetails

	AppDetails struct {
		ID          uint   `json:"id" example:"12"`
		Name        string `json:"name" validate:"required,min=3,max=20" example:"Nextcloud"`
		Description string `json:"description" validate:"required" example:"nextcloud self-hosted cloud"`
		Link        string `json:"link" validate:"required, url" example:"https://nextcloud.lan"`
		Icon        string `json:"icon" validate:"required" example:"nextcloud"`
	}
)
