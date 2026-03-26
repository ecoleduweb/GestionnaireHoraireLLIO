package services

import (
	"llio-api/customs_errors"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"log"

	"github.com/jinzhu/copier"
)

func VerifyActivityJSON(activityDTO *DTOs.ActivityDTO) []DTOs.FieldErrorDTO {
	var errors []DTOs.FieldErrorDTO

	if activityDTO.ProjectId == 0 {
		errors = append(errors, DTOs.FieldErrorDTO{
			Field:   "projectId",
			Message: "Le champ projectId est invalide ou manquant",
		})
	}

	if activityDTO.CategoryId == 0 {
		errors = append(errors, DTOs.FieldErrorDTO{
			Field:   "categoryId",
			Message: "Le champ categoryId est invalide ou manquant",
		})
	}
	// Vérifier que StartDate est avant EndDate
	if activityDTO.StartDate.After(activityDTO.EndDate) {
		errors = append(errors, DTOs.FieldErrorDTO{
			Field:   "startDate",
			Message: "La date de début doit être avant la date de fin",
		})
	}

	return errors
}

func CreateActivity(activityDTO *DTOs.ActivityDTO, currentUserId int) (*DTOs.DetailedActivityDTO, error) {

	activity := &DAOs.Activity{}
	err := copier.Copy(activity, activityDTO)
	if err != nil {
		return nil, err
	}
	activity.UserId = currentUserId

	activityDAOAded, err := repositories.CreateActivity(activity)
	if err != nil {
		return nil, err
	}

	return GetDetailedActivityById(activityDAOAded.Id)
}

func GetUsersActivities(userId int) ([]*DTOs.ActivityDTO, error) {
	activities, err := repositories.GetUsersActivities(userId)
	if err != nil {
		return nil, err
	}

	var activitiesDTOs []*DTOs.ActivityDTO
	for _, activity := range activities {
		activityDTO := &DTOs.ActivityDTO{}
		err = copier.Copy(activityDTO, activity)
		activitiesDTOs = append(activitiesDTOs, activityDTO)
	}

	return activitiesDTOs, err
}

func GetActivityById(id string) (*DTOs.ActivityDTO, error) {
	activity, err := repositories.GetActivityById(id)
	if err != nil {
		return nil, err
	}

	activityDTO := &DTOs.ActivityDTO{}
	err = copier.Copy(activityDTO, activity)

	return activityDTO, err
}

func GetDetailedActivityById(id int) (*DTOs.DetailedActivityDTO, error) {
	activity, err := repositories.GetDetailedActivityById(id)
	if err != nil {
		return nil, err
	}

	detailedActivityDTO := &DTOs.DetailedActivityDTO{
		Id:          activity.Id,
		Name:        activity.Name,
		Description: activity.Description,
		StartDate:   activity.StartDate,
		EndDate:     activity.EndDate,
		UserId:      activity.UserId,
		ProjectId:   activity.ProjectId,
		ProjectName: activity.Project.Name,
		CategoryId:  activity.CategoryId,
	}

	return detailedActivityDTO, err
}

func UpdateActivity(activityDTO *DTOs.ActivityDTO) (*DTOs.DetailedActivityDTO, error) {

	activityDAO := &DAOs.Activity{}
	err := copier.Copy(activityDAO, activityDTO)
	if err != nil {
		return nil, err
	}

	activityDAOUpdated, err := repositories.UpdateActivity(activityDAO)
	if err != nil {
		return nil, err
	}

	return GetDetailedActivityById(activityDAOUpdated.Id)
}

func DeleteActivity(id string, user *DTOs.UserDTO) error {
	canDelete, err := repositories.UserHasPermissionToInteractWithActivities(user, id)

	if err != nil {
		return err
	}

	if canDelete {
		return repositories.DeleteActivity(id)
	}

	return customs_errors.ErrUserForbidden
}

func GetActivitiesFromRange(from string, to string, idUser int) ([]*DTOs.ActivityDTO, error) {
	fromDate := from
	toDate := to

	if from == to {
		toDate = to + " 23:59:59"
		fromDate = from + " 00:00:00"
	}

	activities, err := repositories.GetActivitiesFromRange(fromDate, toDate, idUser)
	if err != nil {
		return nil, err
	}

	var activitiesDTOs []*DTOs.ActivityDTO
	for _, activity := range activities {
		activityDTO := &DTOs.ActivityDTO{}
		err = copier.Copy(activityDTO, activity)
		activitiesDTOs = append(activitiesDTOs, activityDTO)
	}
	if activitiesDTOs == nil {
		log.Printf("Aucune activité trouvée dans la plage de dates spécifiée.")
		return []*DTOs.ActivityDTO{}, err
	}
	return activitiesDTOs, err
}

func GetDetailedActivitiesFromRange(from string, to string, idUser int) ([]*DTOs.DetailedActivityDTO, error) {
	fromDate := from
	toDate := to

	if from == to {
		toDate = to + " 23:59:59"
		fromDate = from + " 00:00:00"
	}

	activities, err := repositories.GetActivitiesFromRange(fromDate, toDate, idUser)
	if err != nil {
		return nil, err
	}

	var detailedActivitiesDTOs []*DTOs.DetailedActivityDTO
	for _, activity := range activities {
		detailedActivityDTO, err := GetDetailedActivityById(activity.Id)
		if err != nil {
			return nil, err
		}
		detailedActivitiesDTOs = append(detailedActivitiesDTOs, detailedActivityDTO)
	}

	if detailedActivitiesDTOs == nil {
		log.Printf("Aucune activité trouvée dans la plage de dates spécifiée.")
		return []*DTOs.DetailedActivityDTO{}, err
	}

	return detailedActivitiesDTOs, err
}
