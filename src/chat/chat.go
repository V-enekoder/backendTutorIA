package chat

import "time"

type UriParams struct {
	ID     uint `uri:"id" binding:"required,min=1,max=17"`
	UserID uint `uri:"user_id" binding:"required,min=1"`
}

type PromptBody struct {
	Prompt string `json:"prompt" binding:"required"`
}
type GeminiResponse struct {
	Response string `json:"response"`
}

type FilePromptRequest struct {
	ID_context int    `uri:"id" binding:"required,min=1,max=20"`
	Prompt     string `form:"prompt" binding:"required"`
}

type GeminiFileData struct {
	MIMEType string
	Data     []byte
}

type ChatResponseDTO struct {
	ID        uint      `json:"id"`
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}
