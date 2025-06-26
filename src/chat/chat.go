package chat

type UriParams struct {
	ID int `uri:"id" binding:"required,min=1,max=17"`
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
