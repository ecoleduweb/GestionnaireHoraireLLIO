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

	err = canUserManageProject(project, author)
	if err != nil {
		return nil, err
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

func DeleteCoManager(projectId int, coManagerToDeleteUserId int, author *DTOs.UserDTO) error {
	project, err := GetProjectById(strconv.Itoa(projectId))
	if err != nil {
		return err
	}

	err = canUserManageProject(project, author)
	if err != nil {
		return err
	}

	isCoManager, err := repositories.IsUserCoManager(projectId, coManagerToDeleteUserId)
	if err != nil {
		return err
	}
	if !isCoManager {
		return customs_errors.ErrSelectedUserIsNotCoManager
	}

	return repositories.DeleteCoManager(projectId, coManagerToDeleteUserId)
}

func canUserManageProject(project *DTOs.ProjectDTO, author *DTOs.UserDTO) error {
	if author.Role >= enums.Administrator || project.ManagerId == author.Id {
		return nil
	}

	isCoManager, err := repositories.IsUserCoManager(project.Id, author.Id)
	if err != nil {
		return err
	}
	if !isCoManager {
		return customs_errors.ErrUserForbidden
	}

	return nil
}
