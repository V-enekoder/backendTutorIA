package document

import "time"

type DocumentResponseDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Resume    string    `json:"resume"`
	Mimetype  string    `json:"mimetype"`
	Size      float64   `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentCreateDTO struct {
	Name     string  `form:"name" binding:"required"`
	Address  string  `form:"address" binding:"required"`
	Resume   string  `form:"resume"`
	Mimetype string  `form:"mimetype" binding:"required"`
	Size     float64 `form:"size" binding:"required"`
	UserID   uint    `form:"user_id" binding:"required"`
}

type DocumentUpdateDTO struct {
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Resume  string  `json:"resume"`
	Size    float64 `json:"size"`
}
