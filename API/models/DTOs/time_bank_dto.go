package DTOs

type TimeBankRequestDTO struct {
	StartDate          string  `json:"startDate" binding:"required"`    // Format YYYY-MM-DD
	HoursPerWeek       float64 `json:"hoursPerWeek" binding:"required"` // Ex: 35 ou 40
	Offset             float64 `json:"offset"`                          // Ex: -5 ou 10 (heures à ajouter/retirer)
	ExcludeCurrentWeek bool    `json:"excludeCurrentWeek"`              // true pour ignorer la semaine en cours
}

type TimeBankResponseDTO struct {
	Balance int `json:"balance"` // Le total arrondi (positif ou négatif)
}
