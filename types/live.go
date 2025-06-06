package types

type LiveUpdate struct {
	Username  string `json:"github_username"`
	Message   string `json:"message"`
	EventType string `json:"event_type"`
	Timestamp int64  `json:"time"` // time is in unix.milliseconds for less size
}
