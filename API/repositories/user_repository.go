package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"
)

func FirstOrCreateUser(user *DAOs.User) (*DAOs.User, error) {
	err := database.DB.Where("email = ?", user.Email).FirstOrCreate(user).Error
	return user, DBErrorManager(err)
}

func GetUserById(id int) (*DAOs.User, error) {
	var user DAOs.User

	err := database.DB.First(&user, id).Error
	return &user, DBErrorManager(err)
}

func GetUserByEmail(email string) (*DAOs.User, error) {
	var user DAOs.User

	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, DBErrorManager(err)
}

func GetAllUsers() ([]*DAOs.User, error) {
	var users []*DAOs.User

	err := database.DB.Find(&users).Error
	return users, DBErrorManager(err)
}

func UpdateUserRole(user *DAOs.User) (*DAOs.User, error) {
	err := database.DB.Save(user).Error
	return user, DBErrorManager(err)
}

func DeleteUserById(id int) error {
	err := database.DB.Delete(&DAOs.User{}, id).Error
	return DBErrorManager(err)
}

func UserHasActivities(userId int) (bool, error) {
	var count int64
	err := database.DB.Model(&DAOs.Activity{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return false, DBErrorManager(err)
	}
	return count > 0, nil
}
func UserHasProjects(userId int) (bool, error) {
	var count int64
	err := database.DB.Model(&DAOs.Project{}).Where("manager_id = ?", userId).Count(&count).Error
	if err != nil {
		return false, DBErrorManager(err)
	}
	return count > 0, nil
}
