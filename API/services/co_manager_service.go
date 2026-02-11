package services

import (
	"llio-api/customs_errors"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"strconv"

	"github.com/jinzhu/copier"
)

func VerifyAddCoManagerJSON(coManagerDTO *DTOs.CoManagerDTO) []DTOs.FieldErrorDTO {
	var errors []DTOs.FieldErrorDTO

	if coManagerDTO.UserId == 0 {
		errors = append(errors, DTOs.FieldErrorDTO{
			Field:   "userId",
			Message: "Le champ userId est invalide ou manquant",
		})
	}

	if coManagerDTO.ProjectId == 0 {
		errors = append(errors, DTOs.FieldErrorDTO{
			Field:   "projectId",
			Message: "Le champ projectId est invalide ou manquant",
		})
	}

	return errors
}

func AddCoManager(coManagerDTO *DTOs.CoManagerDTO, authorId int) (*DTOs.CoManagerDTO, error) {
	project, err := GetProjectById(strconv.Itoa(coManagerDTO.ProjectId))
	if err != nil {
		return nil, err
	}

	coManagerFromDB, err := GetUserById(coManagerDTO.UserId)
	if err != nil {
		return nil, err
	}

	if coManagerFromDB.Role < 1 {
		return nil, customs_errors.ErrUserRoleBelowManager
	}

	if coManagerDTO.UserId == project.ManagerId {
		return nil, customs_errors.ErrUserIsManager
	}

	if project.ManagerId != authorId {
		isCoManager, err := CheckIfUserIsCoManager(authorId, project.Id)
		if err != nil {
			return nil, err
		}

		if !isCoManager {
			return nil, customs_errors.ErrUserForbidden
		}
	}

	coManagerDAO := &DAOs.CoManager{}
	err = copier.Copy(coManagerDAO, coManagerDTO)
	if err != nil {
		return nil, err
	}

	coManagerAdded, err := repositories.AddCoManagerToProject(coManagerDAO)
	if err != nil {
		return nil, err
	}

	coManagerDTOResponse := &DTOs.CoManagerDTO{}
	err = copier.Copy(coManagerDTOResponse, coManagerAdded)
	return coManagerDTOResponse, err
}

func CheckIfUserIsCoManager(userId int, projectId int) (bool, error) {
	_, err := GetProjectById(strconv.Itoa(projectId))
	if err != nil {
		return false, err
	}
	return repositories.IsUserCoManager(projectId, userId)
}
