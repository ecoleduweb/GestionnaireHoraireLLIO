package DTOs

// Pour configurer (GET/POST config)
type TimeBankConfigDTO struct {
	StartDate    string  `json:"startDate" binding:"required,datetime=2006-01-02"`
	HoursPerWeek float64 `json:"hoursPerWeek" binding:"required,gt=0"`
	Offset       float64 `json:"offset"`
}

// Pour afficher le solde (GET root)
type TimeBankBalanceDTO struct {
	IsConfigured bool `json:"isConfigured"`
	TimeInBank   *int `json:"timeInBank"`
}
