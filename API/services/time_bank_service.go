package services

import (
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/repositories"
	"math"
	"time"
)

func SaveTimeBankConfig(userId int, configDTO DTOs.TimeBankConfigDTO) (*DTOs.TimeBankConfigDTO, error) {
	// Conversion String -> Time (On force le mode LOCAL pour éviter les décalages)
	layout := "2006-01-02"
	startDate, _ := time.ParseInLocation(layout, configDTO.StartDate, time.Local)

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

func GetTimeBankConfig(userId int) (*DTOs.TimeBankConfigDTO, error) {
	user, err := repositories.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if user.TimeBankStartDate == nil || user.TimeBankHoursPerWeek == nil {
		return nil, nil
	}

	return &DTOs.TimeBankConfigDTO{
		StartDate:    user.TimeBankStartDate.Format("2006-01-02"),
		HoursPerWeek: *user.TimeBankHoursPerWeek,
		Offset:       user.TimeBankBalanceOffset,
	}, nil
}

func GetTimeBankBalance(userId int) (*DTOs.TimeBankBalanceDTO, error) {
	user, err := repositories.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if user.TimeBankStartDate == nil || user.TimeBankHoursPerWeek == nil {
		return &DTOs.TimeBankBalanceDTO{
			IsConfigured: false,
			TimeInBank:   nil,
		}, nil
	}

	startDateTime := user.TimeBankStartDate.In(time.Local)
	hoursPerWeek := *user.TimeBankHoursPerWeek
	offset := user.TimeBankBalanceOffset

	now := time.Now().In(time.Local)
	calcEndDate := getLastMondayLocal(now)

	workedSeconds, err := repositories.GetTotalWorkedSeconds(userId, startDateTime, calcEndDate)
	if err != nil {
		return nil, err
	}
	totalWorkedHours := workedSeconds / 3600.0

	var totalExpectedHours float64 = 0

	if calcEndDate.After(startDateTime) {
		weeks := 0
		for t := startDateTime; t.AddDate(0, 0, 7).Before(calcEndDate) || t.AddDate(0, 0, 7).Equal(calcEndDate); t = t.AddDate(0, 0, 7) {
			weeks++
		}
		totalExpectedHours = float64(weeks) * hoursPerWeek
	}

	balance := (totalWorkedHours + offset) - totalExpectedHours

	finalBalance := int(math.Round(balance))
	return &DTOs.TimeBankBalanceDTO{
		IsConfigured: true,
		TimeInBank:   &finalBalance,
	}, nil
}

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
