package project

import (
	"errors"
	"time"

	"github.com/V-enekoder/backendTutorIA/src/schema"
    "github.com/V-enekoder/backendTutorIA/src/user"
)


func CreateProjectService(projectDTO ProjectCreateDTO) (uint, error) {

     if _, err := user.GetUserByIdService(projectDTO.UserID); err != nil {
        return 0, errors.New("usuario no existe")
     }


    project := schema.Project{
        UserID:    projectDTO.UserID,
        Name:      projectDTO.Name,
        Address:   projectDTO.Address,
        Summary:   projectDTO.Summary,
        MimeType:  projectDTO.MimeType,
        FileSize:  projectDTO.FileSize,
    }

    id, err := CreateProjectRepository(project)
    if err != nil {
        return 0, err
    }

    return id, nil
}

func GetProjectByIdService(id uint) (ProjectResponseDTO, error) {
    project, err := GetProjectByIdRepository(id)
    if err != nil {
        return ProjectResponseDTO{}, err
    }

    projectResponse := ProjectResponseDTO{
        ID:        project.ID,
        UserID:    project.UserID,
        Name:      project.Name,
        Address:   project.Address,
        Summary:   project.Summary,
        CreatedAt: project.CreatedAt.Format(time.RFC3339),
        MimeType:  project.MimeType,
        FileSize:  project.FileSize,
    }

    return projectResponse, nil
}

func UpdateProjectService(id uint, projectDTO ProjectUpdateDTO) error {
    return UpdateProjectRepository(id, projectDTO)
}

func DeleteProjectByIdService(id uint) error {
    return DeleteProjectByIdRepository(id)
}