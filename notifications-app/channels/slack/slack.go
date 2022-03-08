package slack

type SlackResponse struct {
	Ok      bool   `json:"ok"`
	Error   string `json:"error"`
	Warning string `json:"warning"`
}
