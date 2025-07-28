package models

type TextItem struct {
	Item
	Content string `json:"content"`
}
