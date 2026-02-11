package services

import (
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"math"
	"time"
)

func CalculateTimeBank(userId int, req DTOs.TimeBankRequestDTO) (*DTOs.TimeBankResponseDTO, error) {
	// 1. Convertir la date de début (string -> time.Time)
	layout := "2006-01-02"
	startDateTime, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, err
	}

	// 2. Définir la date de fin du calcul
	calcEndDate := time.Now()

	if req.ExcludeCurrentWeek {
		// On recule jusqu'au lundi de la semaine courante à 00:00
		weekday := int(calcEndDate.Weekday())
		if weekday == 0 {
			weekday = 7
		} // Dimanche = 7
		daysToSubtract := weekday - 1
		calcEndDate = time.Date(calcEndDate.Year(), calcEndDate.Month(), calcEndDate.Day(), 0, 0, 0, 0, calcEndDate.Location()).AddDate(0, 0, -daysToSubtract)
	}

	// 3. Récupérer les heures réellement travaillées via le Repository
	workedSeconds, err := repositories.GetTotalWorkedSeconds(userId, startDateTime, calcEndDate)
	if err != nil {
		return nil, err
	}
	totalWorkedHours := workedSeconds / 3600.0

	// 4. Calculer les heures théoriques attendues
	var totalExpectedHours float64 = 0

	// On ne compte que si la date de fin est après la date de début
	if calcEndDate.After(startDateTime) {
		duration := calcEndDate.Sub(startDateTime)
		weeks := duration.Hours() / (24 * 7)
		totalExpectedHours = weeks * req.HoursPerWeek
	}

	balance := (totalWorkedHours + req.Offset) - totalExpectedHours

	return &DTOs.TimeBankResponseDTO{
		Balance: int(math.Round(balance)),
	}, nil
}
