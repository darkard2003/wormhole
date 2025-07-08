package models

type Channel struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Protected   bool   `json:"protected"`
	Password    string `json:"password"`
}
