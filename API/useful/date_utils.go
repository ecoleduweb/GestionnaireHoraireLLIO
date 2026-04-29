package useful

import "time"

// DateToISOString Formatter la date dans le format RFC3339 (YYYY-MM-DDThh:mm:ssZ), notamment utile pour l'API Outlook.
func DateToISOString(date time.Time) string {
	return date.Format(time.RFC3339)
}

func ToStartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func ToEndOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
}
