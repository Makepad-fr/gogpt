package gogpt

type ModelInfo struct {
	Slug                  string                `json:"slug"`
	MaxTokens             int                   `json:"max_tokens"`
	Title                 string                `json:"title"`
	Description           string                `json:"description"`
	Tags                  []string              `json:"tags"`
	QualitativeProperties QualitativeProperties `json:"qualitative_properties"`
}

type QualitativeProperties struct {
	Reasoning   []int `json:"reasoning"`
	Speed       []int `json:"speed"`
	Conciseness []int `json:"conciseness"`
}

type ModelsResponse struct {
	Models []ModelInfo `json:"models"`
}
