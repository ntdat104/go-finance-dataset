package response

type Meta struct {
	MessageID string `json:"message_id"`
	Timestamp int64  `json:"timestamp"`
	Datetime  string `json:"datetime"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token,omitempty"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data,omitempty"`
}
