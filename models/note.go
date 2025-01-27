package models

// Note represents a note
// @title Title of the note
// @content Content of the note
// @shortDescription Short description of the note
type Note struct {
	ID               int    `json:"id"`
	Title            string `json:"title" validate:"required"`
	Content          string `json:"content" validate:"required"`
	ShortDescription string `json:"shortDescription" validate:"required,max=20"`
}
