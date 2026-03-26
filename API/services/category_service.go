package services

import (
	"llio-api/customs_errors"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"log"

	"github.com/jinzhu/copier"
)

func VerifyCreateCategoryJSON(categoryDTO *DTOs.CategoryDTO) []DTOs.FieldErrorDTO {
	var errors []DTOs.FieldErrorDTO

	if categoryDTO.ProjectId == 0 {
		errors = append(errors, DTOs.FieldErrorDTO{
			Field:   "projectId",
			Message: "Le champ projectId est invalide ou manquant",
		})
	}

	return errors
}

func CreateCategory(categoryDTO *DTOs.CategoryDTO) (*DTOs.CategoryDTO, error) {

	category := &DAOs.Category{}
	err := copier.Copy(category, categoryDTO)
	if err != nil {
		return nil, err
	}

	activityDAOAded, err := repositories.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	categoryDTOResponse := &DTOs.CategoryDTO{}
	err = copier.Copy(categoryDTOResponse, activityDAOAded)
	return categoryDTOResponse, err
}

func GetCategories() ([]*DTOs.CategoryDTO, error) {
	categories, err := repositories.GetCategories()
	if err != nil {
		return nil, err
	}

	var categoriesDTOs []*DTOs.CategoryDTO
	for _, category := range categories {
		categoryDTO := &DTOs.CategoryDTO{}
		err = copier.Copy(categoryDTO, category)
		categoriesDTOs = append(categoriesDTOs, categoryDTO)
	}

	return categoriesDTOs, err
}

func GetCategoryById(id string) (*DTOs.CategoryDTO, error) {
	category, err := repositories.GetCategoryById(id)
	if err != nil {
		return nil, err
	}

	categoryDTO := &DTOs.CategoryDTO{}
	err = copier.Copy(categoryDTO, category)

	return categoryDTO, err
}

func UpdateCategory(categoryDTO *DTOs.CategoryDTO) (*DTOs.CategoryDTO, error) {

	categoryDAO := &DAOs.Category{}
	err := copier.Copy(categoryDAO, categoryDTO)
	if err != nil {
		return nil, err
	}

	categoryDAOUpdated, err := repositories.UpdateCategory(categoryDAO)
	if err != nil {
		return nil, err
	}

	categoryDTOResponse := &DTOs.CategoryDTO{}
	err = copier.Copy(categoryDTOResponse, categoryDAOUpdated)
	return categoryDTOResponse, err
}

func GetCategoriesByProjectId(projectId string) ([]*DTOs.CategoryDTO, error) {
	categories, err := repositories.GetCategoriesByProjectId(projectId)
	if err != nil {
		return nil, err
	}

	var categoriesDTOs []*DTOs.CategoryDTO
	for _, category := range categories {
		categoryDTO := &DTOs.CategoryDTO{}
		err = copier.Copy(categoryDTO, category)
		categoriesDTOs = append(categoriesDTOs, categoryDTO)
	}

	return categoriesDTOs, err
}

func DeleteCategory(id string) error {

	activitiesCount, err := repositories.GetActivitiesCountFromCategoryId(id)

	if err != nil {
		log.Printf("Il y a eu une erreur dans la base de données.")
		return err
	}

	if activitiesCount > 0 {
		log.Printf("Il y a déjà une activité liée à cette catégorie.")
		return customs_errors.ErrUserHasActivities
	}

	return repositories.DeleteCategory(id)
}
