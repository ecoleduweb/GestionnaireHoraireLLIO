package DTOs

// Pour configurer (GET/POST config)
// La date permet de spécifier le format de date, je donne un exemple à suivre qui est 2006-01-02. La date n'est pas utilisée pour le calcul, c'est juste pour la validation du format et pour l'affichage (on stocke en Time dans la BD)
type TimeBankConfigDTO struct {
	StartDate    string  `json:"startDate" binding:"required,datetime=2006-01-02"`
	HoursPerWeek float64 `json:"hoursPerWeek" binding:"required,gt=0"`
	Offset       float64 `json:"offset"`
}

// Pour afficher le solde (GET root)
type TimeBankBalanceDTO struct {
	IsConfigured bool     `json:"isConfigured"`
	TimeInBank   *float64 `json:"timeInBank"`
}
