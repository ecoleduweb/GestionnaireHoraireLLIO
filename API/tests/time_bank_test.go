package tests

import (
	"encoding/json"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeBankScenarios(t *testing.T) {
	testUserID := int(doNotDeleteUser.Id)

	// --- 1. Préparation des dates (Strictement en UTC) ---
	now := time.Now().UTC()

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// Lundi de la semaine en cours
	mondayCurrentWeek := time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, time.UTC)
	// Lundi de la semaine dernière
	mondayLastWeek := mondayCurrentWeek.AddDate(0, 0, -7)
	startDateStr := mondayLastWeek.Format("2006-01-02")

	// --- 2. Nettoyage initial agressif (Environnement stérile) ---
	// Suppression de TOUTES les activités de cet utilisateur pour éviter la pollution des autres tests
	if err := database.DB.Where("user_id = ?", testUserID).Delete(&DAOs.Activity{}).Error; err != nil {
		t.Fatalf("Erreur purge initiale des activités : %v", err)
	}

	if err := database.DB.Model(&doNotDeleteUser).Updates(map[string]interface{}{
		"time_bank_start_date":     nil,
		"time_bank_hours_per_week": nil,
		"time_bank_balance_offset": 0,
	}).Error; err != nil {
		t.Fatalf("Erreur réinitialisation doNotDeleteUser : %v", err)
	}

	// --- 3. Définition du Nettoyage Automatique de fin ---
	t.Cleanup(func() {
		database.DB.Where("name LIKE ?", "Test TimeBank%").Delete(&DAOs.Activity{})
		database.DB.Model(&doNotDeleteUser).Updates(map[string]interface{}{
			"time_bank_start_date":     nil,
			"time_bank_hours_per_week": nil,
			"time_bank_balance_offset": 0,
		})
	})

	// --- SCÉNARIO 1 : Date Invalide ---
	t.Run("InvalidDate", func(t *testing.T) {
		badConfig := map[string]interface{}{
			"startDate":    "99-99-2024",
			"hoursPerWeek": 35,
		}
		w := sendRequest(router, "PUT", "/user/time-bank/config", badConfig, &testUserID)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// --- SCÉNARIO 2 : Pas d'activités (Déficit Total) ---
	t.Run("NoActivities_Deficit", func(t *testing.T) {
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 10, Offset: 0}
		sendRequest(router, "PUT", "/user/time-bank/config", config, &testUserID)

		w := sendRequest(router, "GET", "/user/time-bank", nil, &testUserID)
		require.Equal(t, http.StatusOK, w.Code)

		var resp DTOs.TimeBankBalanceDTO
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.NotNil(t, resp.TimeInBank)

		// 0h faites - 10h attendues = -10
		assert.Equal(t, float64(-10), *resp.TimeInBank)
	})

	// --- INJECTION : 12 heures travaillées la semaine dernière ---
	activities := []DAOs.Activity{
		{
			Name:        "Test TimeBank 1",
			Description: "Auto",
			StartDate:   mondayLastWeek.Add(9 * time.Hour),  // 09:00 UTC
			EndDate:     mondayLastWeek.Add(13 * time.Hour), // 13:00 UTC (4h)
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 2",
			Description: "Auto",
			StartDate:   mondayLastWeek.AddDate(0, 0, 1).Add(9 * time.Hour),
			EndDate:     mondayLastWeek.AddDate(0, 0, 1).Add(13 * time.Hour),
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 3",
			Description: "Auto",
			StartDate:   mondayLastWeek.AddDate(0, 0, 2).Add(9 * time.Hour),
			EndDate:     mondayLastWeek.AddDate(0, 0, 2).Add(13 * time.Hour),
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
	}
	for _, act := range activities {
		if err := database.DB.Create(&act).Error; err != nil {
			t.Fatalf("Erreur insert activité: %v", err)
		}
	}

	// --- SCÉNARIO 3 : Surplus (Fait > Attendu) ---
	t.Run("Surplus", func(t *testing.T) {
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 10, Offset: 0}
		sendRequest(router, "PUT", "/user/time-bank/config", config, &testUserID)

		w := sendRequest(router, "GET", "/user/time-bank", nil, &testUserID)

		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 12h faites - 10h attendues = +2
		assert.Equal(t, float64(2), *resp.TimeInBank)
	})

	// --- SCÉNARIO 4 : Déficit (Fait < Attendu) ---
	t.Run("Deficit", func(t *testing.T) {
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 20, Offset: 0}
		sendRequest(router, "PUT", "/user/time-bank/config", config, &testUserID)

		w := sendRequest(router, "GET", "/user/time-bank", nil, &testUserID)

		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 12h faites - 20h attendues = -8
		assert.Equal(t, float64(-8), *resp.TimeInBank)
	})

	// --- SCÉNARIO 5 : Solde Nul (Fait == Attendu) ---
	t.Run("ZeroBalance", func(t *testing.T) {
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 12, Offset: 0}
		sendRequest(router, "PUT", "/user/time-bank/config", config, &testUserID)

		w := sendRequest(router, "GET", "/user/time-bank", nil, &testUserID)

		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 12h faites - 12h attendues = 0
		assert.Equal(t, float64(0), *resp.TimeInBank)
	})
}
