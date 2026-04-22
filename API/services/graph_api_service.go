package services

import (
	"llio-api/useful"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func GetCalendarEvents(accessToken string, date time.Time) ([]GraphEvent, error) {
	startOfDay := useful.ToStartOfDay(date)
	endOfDay := useful.ToEndOfDay(date)
	url := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/me/calendarView?startDateTime=%s&endDateTime=%s",
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

type GraphEvent struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Body    struct {
		Content string `json:"content"`
	} `json:"body"`
	Start struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
	IsAllDay bool `json:"isAllDay"`
}
