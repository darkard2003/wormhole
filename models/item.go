package models

type Item struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	FileURL string `json:"file_url"`
}
