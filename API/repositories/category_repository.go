package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"
	"strconv"
)

func CreateCategory(category *DAOs.Category) (*DAOs.Category, error) {
	err := database.DB.Create(category).Error
	return category, DBErrorManager(err)
}

func GetCategories() ([]*DAOs.Category, error) {
	var categories []*DAOs.Category
	err := database.DB.Find(&categories).Error
	return categories, DBErrorManager(err)
}

func GetCategoryById(id string) (*DAOs.Category, error) {
	var category DAOs.Category

	err := database.DB.First(&category, id).Error
	return &category, DBErrorManager(err)
}

func UpdateCategory(categoryDAO *DAOs.Category) (*DAOs.Category, error) {
	err := database.DB.Updates(categoryDAO).Error
	return categoryDAO, DBErrorManager(err)
}

func GetCategoriesByProjectId(projectId string) ([]*DAOs.Category, error) {
	var categories []*DAOs.Category

	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		return nil, err
	}

	err = database.DB.Where("project_id = ?", projectIdInt).Find(&categories).Error
	return categories, DBErrorManager(err)
}

func DeleteCategory(id string) error {
	return DBErrorManager(database.DB.Delete(&DAOs.Category{}, id).Error)
}
