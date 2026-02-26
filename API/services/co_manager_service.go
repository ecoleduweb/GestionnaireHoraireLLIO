package services

import (
	"llio-api/customs_errors"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/repositories"
	"strconv"

	"github.com/jinzhu/copier"
)

func AddCoManager(coManagerDTO *DTOs.CoManagerDTO, author *DTOs.UserDTO) (*DTOs.CoManagerDTO, error) {
	project, err := GetProjectById(strconv.Itoa(coManagerDTO.ProjectId))
	if err != nil {
		return nil, err
	}

	targetUser, err := GetUserById(coManagerDTO.UserId)
	if err != nil {
		return nil, err
	}

	if targetUser.Role < enums.ProjectManager {
		return nil, customs_errors.ErrUserRoleBelowManager
	}

	if coManagerDTO.UserId == project.ManagerId {
		return nil, customs_errors.ErrUserIsManager
	}

	// Ici, on vérifie les permissions pour que seulement les personnes
	// suivantes puissent ajouter un co-chargé de projet :
	//
	// - Tous les administrateurs
	// - Tout chargé de projet étant chargé ou co-chargé de ce projet en spécifique.
	//
	// Toutes les autres personnes (utilisateur régulier ou chargé de projet non-attribué)
	//	ne peuvent pas modifier les co-chargés du projet
	if author.Role < enums.Administrator {
		if project.ManagerId != author.Id {
			isCoManager, err := repositories.IsUserCoManager(project.Id, author.Id)
			if err != nil {
				return nil, err
			}

			if !isCoManager {
				return nil, customs_errors.ErrUserForbidden
			}
		}
	}

	coManagerDAO := &DAOs.CoManager{}
	err = copier.Copy(coManagerDAO, coManagerDTO)
	if err != nil {
		return nil, err
	}

	alreadyCoManager, err := repositories.IsUserCoManager(coManagerDAO.ProjectId, coManagerDAO.UserId)
	if err != nil {
		return nil, err
	}
	if alreadyCoManager {
		return nil, customs_errors.ErrUserAlreadyCoManager
	}

	coManagerAdded, err := repositories.AddCoManagerToProject(coManagerDAO)
	if err != nil {
		return nil, err
	}

	coManagerDTOResponse := &DTOs.CoManagerDTO{}
	err = copier.Copy(coManagerDTOResponse, coManagerAdded)
	return coManagerDTOResponse, err
}
