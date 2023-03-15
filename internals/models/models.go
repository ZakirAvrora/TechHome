package models

type Link struct {
	ID          string `json:"id,omitempty" db:"redirect_id"`
	ActiveLink  string `json:"active_link" db:"active_link"`
	HistoryLink string `json:"history_link" db:"history_link"`
}
