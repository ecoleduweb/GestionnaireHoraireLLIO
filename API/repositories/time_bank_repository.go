package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"
	"time"
)

// Calcule la somme des heures travaillées (en secondes) pour une période donnée
func GetTotalWorkedSeconds(userId int, startDate time.Time, endDate time.Time) (float64, error) {
	var totalSeconds float64

	// On additionne la différence entre fin et début de chaque activité
	// COALESCE(..., 0) évite de retourner NULL si aucune activité n'est trouvée
	err := database.DB.Model(&DAOs.Activity{}).
		Select("COALESCE(SUM(TIMESTAMPDIFF(SECOND, start_date, end_date)), 0)").
		Where("user_id = ? AND start_date >= ? AND start_date < ?", userId, startDate, endDate).
		Scan(&totalSeconds).Error

	return totalSeconds, DBErrorManager(err)
}
