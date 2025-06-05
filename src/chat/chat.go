package chat

type PromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type GeminiResponse struct {
	Response string `json:"response"`
}

type FilePromptRequest struct {
	Prompt string `form:"prompt" binding:"required"`
}

type GeminiFileData struct {
	MIMEType string
	Data     []byte
}
