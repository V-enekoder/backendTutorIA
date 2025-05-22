package project

type ProjectResponseDTO struct {
    ID          uint    `json:"id"`
    UserID      uint    `json:"user_id"`      
    Name        string  `json:"name"`         
    Address     string  `json:"address"`      
    Summary     string  `json:"summary"`      
    CreatedAt   string  `json:"created_at"`   
    MimeType    string  `json:"mime_type"`    
    FileSize    float64 `json:"file_size"`  
}

type ProjectCreateDTO struct {
    UserID   uint    `form:"user_id" binding:"required"`   
    Name     string  `form:"name" binding:"required"`     
    Address  string  `form:"address"`                     
    Summary  string  `form:"summary" binding:"required"`   
    MimeType string  `form:"mime_type"`                    
    FileSize float64 `form:"file_size"`                   
}

type ProjectUpdateDTO struct {
    Name     string  `json:"name"`              
    Address  string  `json:"address"`           
    Summary  string  `json:"summary"`           
    MimeType string  `json:"mime_type"`        
    FileSize float64 `json:"file_size"`         
}

