package services

import (
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/repositories"

	"github.com/jinzhu/copier"

	"strconv"
)

func FirstOrCreateUser(userDTO *DTOs.UserDTO) (*DTOs.UserDTO, error) {

	user := &DAOs.User{}
	err := copier.Copy(user, userDTO)
	if err != nil {
		return nil, err
	}
	userDAOAdded, err := repositories.FirstOrCreateUser(user)
	if err != nil {
		return nil, err
	}

	userDTOResponse := &DTOs.UserDTO{}
	err = copier.Copy(userDTOResponse, userDAOAdded)
	return userDTOResponse, err
}

func GetUserById(id string) (*DTOs.UserDTO, error) {
	user, err := repositories.GetUserById(id)
	if err != nil {
		return nil, err
	}

	userDTO := &DTOs.UserDTO{}
	err = copier.Copy(userDTO, user)

	return userDTO, err
}

func GetUserByEmail(email string) (*DTOs.UserDTO, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	userDTO := &DTOs.UserDTO{}
	err = copier.Copy(userDTO, user)

	return userDTO, err
}

func GetAllUsers(userRoles []enums.UserRole) ([]*DTOs.UserDTO, error) {

	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var usersDTOs []*DTOs.UserDTO

	for _, user := range users {
		if len(userRoles) > 0 {
			roleMatch := false
			for _, role := range userRoles {
				if user.Role == role {
					roleMatch = true
					break
				}
			}
			if !roleMatch {
				continue
			}
		}

		userDTO := &DTOs.UserDTO{}
		if copyErr := copier.Copy(userDTO, user); copyErr != nil {
			return nil, copyErr
		}
		usersDTOs = append(usersDTOs, userDTO)
	}

	return usersDTOs, nil
}

func UpdateUserRole(userDTO *DTOs.UserDTO) (*DTOs.UserDTO, error) {
	// Get the existing user
	existingUser, err := repositories.GetUserById(strconv.Itoa(userDTO.Id))
	if err != nil {
		return nil, err
	}

	// Only update the role
	existingUser.Role = userDTO.Role

	// Save the updated user
	updatedUser, err := repositories.UpdateUserRole(existingUser)
	if err != nil {
		return nil, err
	}

	// Convert back to DTO
	userDTOResponse := &DTOs.UserDTO{}
	err = copier.Copy(userDTOResponse, updatedUser)

	return userDTOResponse, err
}

func DeleteUserById(id string) (*DTOs.UserDTO, error) {
	// Check if the user exists
	userDAO, err := repositories.GetUserById(id)
	if err != nil {
		return nil, err
	}
	if userDAO == nil {
		return nil, nil // User not found
	}

	// Delete the user
	err = repositories.DeleteUserById(id)
	if err != nil {
		return nil, err
	}

	// Return the deleted user
	userDTO := &DTOs.UserDTO{}
	err = copier.Copy(userDTO, userDAO)
	if err != nil {
		return nil, err
	}
	return userDTO, nil
}

func UserHasActivities(userId int) (bool, error) {
	// Check if the user has any activities
	activities, err := repositories.GetUsersActivities(userId)
	if err != nil {
		return false, err
	}
	if len(activities) > 0 {
		return true, nil
	}

	return false, nil
}

func UserHasProjects(userId int) (bool, error) {
	projects, err := repositories.GetProjectsByManagerId(userId)
	if err != nil {
		return false, err
	}
	if len(projects) > 0 {
		return true, nil
	}

	return false, nil
}
