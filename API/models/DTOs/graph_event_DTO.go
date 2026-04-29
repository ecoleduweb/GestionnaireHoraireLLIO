package DTOs

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
