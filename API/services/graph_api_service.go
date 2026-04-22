package services

import (
	"llio-api/models/DTOs"
	"llio-api/useful"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const GraphApiBaseUrl = "https://graph.microsoft.com/v1.0"

func GetCalendarEvents(accessToken string, date time.Time) ([]DTOs.GraphEvent, error) {
	startOfDay := useful.ToStartOfDay(date)
	endOfDay := useful.ToEndOfDay(date)
	url := fmt.Sprintf(
		"%s/me/calendarView?startDateTime=%s&endDateTime=%s",
		GraphApiBaseUrl,
		useful.DateToISOString(startOfDay),
		useful.DateToISOString(endOfDay),
	)

	body, err := GraphApiGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}

	var result struct {
		Value []GraphEvent `json:"value"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Value, nil
}

func IsJWTExpired(tokenString string) (bool, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid JWT format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return false, err
	}

	return time.Now().Unix() > claims.Exp, nil
}
