package user

type UserResponseDTO struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	NumberOfDocuments uint   `json:"number_of_documents,omitempty"`
	NumberOfProjects  uint   `json:"number_of_projects,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type UserCreateDTO struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserUpdateDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserUpdatePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
