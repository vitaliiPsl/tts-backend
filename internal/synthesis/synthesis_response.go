package synthesis

type SynthesisResponse struct {
	Samples      []float32 `json:"samples"`
	SamplingRate int       `json:"sampling_rate"`
}
