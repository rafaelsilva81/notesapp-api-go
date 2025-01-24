package services

import (
	"database/sql"
	"notesapp/api/config"
	"notesapp/api/models"

	_ "github.com/mattn/go-sqlite3"
)

type NoteService struct {
	DB *sql.DB
}

func NewNoteService() (*NoteService, error) {
	db := config.GetDatabase()
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
				shortDescription TEXT NOT NULL
    )`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &NoteService{DB: db}, nil
}

func (service *NoteService) GetNotes() (notes []models.Note, err error) {
	rows, err := service.DB.Query("SELECT id, title, content, shortDescription FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note models.Note
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.ShortDescription)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (service *NoteService) CreateNote(note models.Note) error {
	stmt, err := service.DB.Prepare("INSERT INTO notes (id, title, content, shortDescription) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nil, note.Title, note.Content, note.ShortDescription)
	if err != nil {
		return err
	}

	return nil
}

func (service *NoteService) GetNoteById(id int) (note models.Note, err error) {
	row := service.DB.QueryRow("SELECT id, title, content, shortDescription FROM notes WHERE id = ?", id)
	err = row.Scan(&note.ID, &note.Title, &note.Content, &note.ShortDescription)
	if err != nil {
		return note, err
	}

	return note, nil
}

func (service *NoteService) UpdateNote(id int, note models.Note) error {
	stmt, err := service.DB.Prepare("UPDATE notes SET title = ?, content = ?, shortDescription = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(note.Title, note.Content, note.ShortDescription, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *NoteService) DeleteNote(id int) error {
	stmt, err := service.DB.Prepare("DELETE FROM notes WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
