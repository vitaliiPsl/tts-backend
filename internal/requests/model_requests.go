package requests

type ModelRequest struct {
	Url      string `json:"url" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Language string `json:"language" validate:"required"`
}
