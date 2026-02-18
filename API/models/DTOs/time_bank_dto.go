package DTOs

type TimeBankRequestDTO struct {
	StartDate    string  `json:"startDate" binding:"required,datetime=2006-01-02"` // Format YYYY-MM-DD
	HoursPerWeek float64 `json:"hoursPerWeek" binding:"required,gt=0 "`            // Ex: 35 ou 40
	Offset       float64 `json:"offset"`                                           // Ex: -5 ou 10 (heures à ajouter/retirer)
}

type TimeBankResponseDTO struct {
	Balance int `json:"balance"` // Le total arrondi (positif ou négatif)
}
