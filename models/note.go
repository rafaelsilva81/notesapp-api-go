package models

// Note represents a note
// @title Note
// @description Represents a note
type Note struct {
	ID               int    `json:"id"`
	Title            string `json:"title" validate:"required"`
	Content          string `json:"content" validate:"required"`
	ShortDescription string `json:"shortDescription" validate:"required,max=20"`
}
