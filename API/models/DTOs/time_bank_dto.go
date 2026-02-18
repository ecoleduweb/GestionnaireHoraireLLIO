package DTOs

type TimeBankRequestDTO struct {
	StartDate    string  `json:"startDate" binding:"required,datetime=2006-01-02"` // Format YYYY-MM-DD
	HoursPerWeek float64 `json:"hoursPerWeek" binding:"required,gt=0 "`            // Ex: 35 ou 40
	Offset       float64 `json:"offset"`                                           // Ex: -5 ou 10 (heures à ajouter/retirer)
}

type TimeBankResponseDTO struct {
	Balance int `json:"balance"` // Le total arrondi (positif ou négatif)
}

// Pour configurer (GET/POST config)
type TimeBankConfigDTO struct {
	StartDate    string  `json:"startDate" binding:"required,datetime=2006-01-02"`
	HoursPerWeek float64 `json:"hoursPerWeek" binding:"required,min=0"`
	Offset       float64 `json:"offset"`
}

// Pour le résultat du calcul (GET balance)
type TimeBankBalanceDTO struct {
	IsConfigured bool `json:"isConfigured"`
	TimeInBank   *int `json:"timeInBank"`
}
