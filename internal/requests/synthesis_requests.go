package requests

type SynthesisRequest struct {
	Text string `json:"text" validate:"required"`
}
