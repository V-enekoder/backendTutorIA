package document

import (
	"errors"

	"github.com/V-enekoder/backendTutorIA/src/schema"
)

func CreateDocumentService(docDTO DocumentCreateDTO) (uint, error) {
	// Verificar si el usuario existe
	if exists, err := DocumentExistsByFieldRepository("user_id", docDTO.UserID, 0); err != nil {
		return 0, err
	} else if !exists {
		return 0, errors.New("Usuario no existe")
	}

	doc := schema.Document{
		Name:     docDTO.Name,
		Path:     docDTO.Address,
		Resume:   docDTO.Resume,
		Mimetype: docDTO.Mimetype,
		Size:     docDTO.Size,
		UserID:   docDTO.UserID,
	}

	return CreateDocumentRepository(doc)
}

func GetDocumentByIdService(id uint) (DocumentResponseDTO, error) {
	doc, err := GetDocumentByIdRepository(id)
	if err != nil {
		return DocumentResponseDTO{}, err
	}

	return DocumentResponseDTO{
		ID:        doc.ID,
		UserID:    doc.UserID,
		Name:      doc.Name,
		Address:   doc.Path,
		Resume:    doc.Resume,
		Mimetype:  doc.Mimetype,
		Size:      doc.Size,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}

func GetDocumentsByUserService(userID uint) ([]DocumentResponseDTO, error) {
	docs, err := GetDocumentsByUserRepository(userID)
	if err != nil {
		return nil, err
	}

	var response []DocumentResponseDTO
	for _, doc := range docs {
		response = append(response, DocumentResponseDTO{
			ID:        doc.ID,
			UserID:    doc.UserID,
			Name:      doc.Name,
			Address:   doc.Path,
			Resume:    doc.Resume,
			Mimetype:  doc.Mimetype,
			Size:      doc.Size,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}
	return response, nil
}

func UpdateDocumentService(id uint, docDTO DocumentUpdateDTO) error {
	return UpdateDocumentRepository(id, docDTO)
}

func DeleteDocumentByIdService(id uint) error {
	return DeleteDocumentByIdRepository(id)
}
