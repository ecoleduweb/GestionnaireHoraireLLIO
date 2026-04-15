package repositories

import (
	"fmt"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/enums"

	"gorm.io/gorm"
)
func GetAvailableManagers(projectId int) ([]*DAOs.User, error) {
	var managers []*DAOs.User
	var project DAOs.Project

	err := database.DB.First(&project, projectId).Error
	if err != nil {
		return nil, DBErrorManager(err)
	}


	err = database.DB.
		Where("role IN ?", []enums.UserRole{enums.ProjectManager, enums.Administrator}).
		Where("id != ?", project.ManagerId).
		Find(&managers).Error


	return managers, DBErrorManager(err)
}
//Transaction pour changer un chargé de projet
func ReassignManager(projectId int, newManagerId int) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
		//Récupère le projet séléctionné
        var project DAOs.Project
        if err := tx.First(&project, projectId).Error; err != nil {
            return DBErrorManager(err)
        }
		//Verifie que le manager séléctionné n'est pas déja attribué au projet
        if project.ManagerId == newManagerId {
            return fmt.Errorf("le manager sélectionné est déjà assigné à ce projet")
        }

        //Supprime le nouveau manager des co-chargés si il est co chargé de projet
        if err := tx.Where("project_id = ? AND user_id = ?", projectId, newManagerId).
            Delete(&DAOs.CoManager{}).Error; err != nil {
            return DBErrorManager(err)
        }

        //Réatribue le manager avec celui qui à été séléctionné
        if err := tx.Model(&project).Update("manager_id", newManagerId).Error; err != nil {
            return DBErrorManager(err)
        }
		//Si la transaction est validée, elle retourne null et s'execute sinon elle retourne une érreur
        return nil
    })
}