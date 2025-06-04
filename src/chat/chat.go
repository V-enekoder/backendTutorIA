package chat

type PromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type GeminiResponse struct {
	Response string `json:"response"`
}
