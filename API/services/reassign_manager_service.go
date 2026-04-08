package services

import (
	"fmt"
	"llio-api/models/DTOs"
	"llio-api/repositories"
)

type ProjectService struct{}

func GetAvailableManagers(projectId int) ([]DTOs.UserDTO, error) {
    usersDAO, err := repositories.GetAvailableManagers(projectId)
    if err != nil {
        return nil, err
    }

    // Conversion en DTO
    var usersDTO []DTOs.UserDTO
    for _, u := range usersDAO {
        usersDTO = append(usersDTO, DTOs.UserDTO{
            Id:        u.Id,
            FirstName: u.FirstName,
            LastName:  u.LastName,
            Email:     u.Email,
            Role:      u.Role,
        })
    }

    return usersDTO, nil
}

func (s *ProjectService) ReassignManager(projectId int, newManagerId int) error {
	availableManagers, err := repositories.GetAvailableManagers(projectId)
	if err != nil {
		return fmt.Errorf("impossible de récupérer les managers disponibles: %w", err)
	}

	found := false
	for _, m := range availableManagers {
		if m.Id == newManagerId {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("le manager sélectionné n'est pas disponible pour ce projet")
	}

	if err := repositories.ReassignManager(projectId, newManagerId); err != nil {
		return fmt.Errorf("échec de la réattribution du manager: %w", err)
	}

	return nil
}