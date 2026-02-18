package services

import (
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"math"
	"time"
)

// 1. Sauvegarder la configuration
func SaveTimeBankConfig(userId int, configDTO DTOs.TimeBankConfigDTO) (*DTOs.TimeBankConfigDTO, error) {
	// Conversion String -> Time (On force le mode LOCAL pour éviter les décalages)
	layout := "2006-01-02"
	startDate, _ := time.ParseInLocation(layout, configDTO.StartDate, time.Local)

	// Préparation de l'objet DAO partiel
	user := &DAOs.User{
		Id:                    userId,
		TimeBankStartDate:     &startDate,
		TimeBankHoursPerWeek:  &configDTO.HoursPerWeek,
		TimeBankBalanceOffset: configDTO.Offset,
	}

	_, err := repositories.UpdateUserTimeBankConfig(user)
	if err != nil {
		return nil, err
	}

	return &configDTO, nil
}

// 2. Récupérer la configuration
func GetTimeBankConfig(userId int) (*DTOs.TimeBankConfigDTO, error) {
	user, err := repositories.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// Si l'utilisateur n'a pas encore configuré sa banque
	if user.TimeBankStartDate == nil || user.TimeBankHoursPerWeek == nil {
		return nil, nil
	}

	return &DTOs.TimeBankConfigDTO{
		StartDate:    user.TimeBankStartDate.Format("2006-01-02"),
		HoursPerWeek: *user.TimeBankHoursPerWeek,
		Offset:       user.TimeBankBalanceOffset,
	}, nil
}

// 3. Calculer le solde (Lecture depuis DB + Calcul Local)
func GetTimeBankBalance(userId int) (*DTOs.TimeBankBalanceDTO, error) {
	user, err := repositories.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// Vérification de la configuration
	if user.TimeBankStartDate == nil || user.TimeBankHoursPerWeek == nil {
		return &DTOs.TimeBankBalanceDTO{
			IsConfigured: false,
			TimeInBank:   nil,
		}, nil
	}

	// --- LOGIQUE DE CALCUL (SIMPLE & LOCAL) ---

	// 1. Récupération des paramètres (En Local)
	startDateTime := user.TimeBankStartDate.In(time.Local)
	hoursPerWeek := *user.TimeBankHoursPerWeek
	offset := user.TimeBankBalanceOffset

	// 2. Date de fin : Maintenant (Local)
	// On recule jusqu'au dernier Lundi minuit pour avoir des semaines complètes
	now := time.Now().In(time.Local)
	calcEndDate := getLastMondayLocal(now)

	// 3. Récupération des heures travaillées
	workedSeconds, err := repositories.GetTotalWorkedSeconds(userId, startDateTime, calcEndDate)
	if err != nil {
		return nil, err
	}
	totalWorkedHours := workedSeconds / 3600.0

	// 4. Calcul des heures attendues
	var totalExpectedHours float64 = 0

	// Comme tout est en Local, la différence de temps est exacte
	if calcEndDate.After(startDateTime) {
		duration := calcEndDate.Sub(startDateTime)
		weeks := duration.Hours() / (24 * 7)
		totalExpectedHours = weeks * hoursPerWeek
	}

	// 5. Résultat
	balance := (totalWorkedHours + offset) - totalExpectedHours

	finalBalance := int(math.Round(balance))
	return &DTOs.TimeBankBalanceDTO{
		IsConfigured: true,
		TimeInBank:   &finalBalance,
	}, nil
}

// Utilitaire pour trouver le lundi minuit local
func getLastMondayLocal(t time.Time) time.Time {
	weekday := int(t.Weekday()) // Dimanche=0, Lundi=1...
	if weekday == 0 {
		weekday = 7
	}
	daysToSubtract := weekday - 1

	// On retourne une date "pure" (00:00:00) en Local
	lastMonday := t.AddDate(0, 0, -daysToSubtract)
	return time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, time.Local)
}
