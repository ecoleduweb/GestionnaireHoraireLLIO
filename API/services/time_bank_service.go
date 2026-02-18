package services

import (
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"math"
	"time"
)

func CalculateTimeBank(userId int, req DTOs.TimeBankRequestDTO) (*DTOs.TimeBankResponseDTO, error) {
	layout := "2006-01-02"
	startDateTime, _ := time.Parse(layout, req.StartDate)

	// Définir la date de fin du calcul (Lundi de la semaine courante à 00:00 UTC)
	calcEndDate := getLastMonday(time.Now().UTC())

	workedSeconds, err := repositories.GetTotalWorkedSeconds(userId, startDateTime, calcEndDate)
	if err != nil {
		return nil, err
	}
	totalWorkedHours := workedSeconds / 3600.0

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

// getLastMonday retourne le lundi de la semaine courante à 00:00:00 UTC
func getLastMonday(t time.Time) time.Time {
	// je laisse les commentaires pour expliquer la logique, même si c'est un peu verbeux, car la manipulation des dates est souvent source de confusion et d'erreurs.
	// En Go, Sunday = 0, Monday = 1, ... Saturday = 6
	weekday := int(t.Weekday())

	// On veut savoir combien de jours reculer pour tomber sur Lundi (1)
	// Si on est Lundi (1) -> reculer de 0 jour
	// Si on est Mardi (2) -> reculer de 1 jour
	// Si on est Dimanche (0) -> reculer de 6 jours (car c'est la fin de la semaine ISO)

	daysSinceMonday := (weekday + 6) % 7

	lastMonday := t.AddDate(0, 0, -daysSinceMonday)

	// On normalise à 00:00:00 UTC
	return time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, time.UTC)
}
