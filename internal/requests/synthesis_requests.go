package requests

type SynthesisRequest struct {
	Text    string `json:"text" validate:"required"`
	ModelId string `json:"modelId" validate:"required"`
}
