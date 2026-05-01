package services

import (
	"fmt"
	"llio-api/models/DTOs"
	"llio-api/repositories"

	"github.com/jinzhu/copier"
)

type ProjectService struct{}

func GetAvailableManagers(projectId int) ([]DTOs.UserDTO, error) {
    usersDAO, err := repositories.GetAvailableManagers(projectId)
    if err != nil {
        return nil, err
    }

    var usersDTO []DTOs.UserDTO

    err = copier.Copy(&usersDTO, &usersDAO)
    if err != nil {
        return nil, err
    }

    return usersDTO, nil
}
func ReassignManager(projectId int, newManagerId int) error {
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