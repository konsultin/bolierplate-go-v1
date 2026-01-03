package dto

type CreateAnonymousSession_Result struct {
	Session *Session `json:"session,omitempty"`
	Scopes  []string `json:"scopes,omitempty"`
}

type Session struct {
	Token     string `json:"token,omitempty"`
	ExpiredAt int64  `json:"expiredAt,omitempty"`
}
