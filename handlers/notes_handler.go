package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notesapp/api/models"
	"notesapp/api/services"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var service *services.NoteService

func InitializeNotesHandler(router *mux.Router) {
	s, err := services.NewNoteService()
	if err != nil {
		panic(err)
	}

	service = s

	router.HandleFunc("/notes", handleGetAllnotes).Methods("GET")
	router.HandleFunc("/notes", handleCreateNote).Methods("POST")
	router.HandleFunc("/notes/{id}", handleGetNoteById).Methods("GET")
	router.HandleFunc("/notes/{id}", handleUpdateNote).Methods("PUT", "PATCH")
	router.HandleFunc("/notes/{id}", handleDeleteNote).Methods("DELETE")

}

// GetAllNotes returns all notes
// @Summary Get all notes
// @Description This method returns all notes from the database
// @Tags Notes
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Note
// @Failure 400
// @Router /notes [get]
func handleGetAllnotes(writer http.ResponseWriter, req *http.Request) {
	notes, err := service.GetNotes()

	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(notes)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
}

// GetNoteById returns a note by id
// @Summary Get a note by id
// @Description This method returns a note from the database by id
// @Tags Notes
// @Accept  json
// @Produce  json
// @Param id path int true "Note id"
// @Success 200 {object} models.Note
// @Failure 400
// @Router /notes/{id} [get]
func handleGetNoteById(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req) // map[string]string -- id:value

	id, err := strconv.Atoi(vars["id"])

	// Validar se o id é um número inteiro
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	note, err := service.GetNoteById(id)

	// Valida se o id existe
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(note)

	// Valida se o objeto Note foi serializado corretamente
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
}

// CreateNote creates a new note
// @Summary Create a new note
// @Description This method creates a new note in the database
// @Tags Notes
// @Accept  json
// @Produce  json
// @Param note body models.Note true "Note"
// @Success 201 {object} models.Note
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /notes [post]
func handleCreateNote(writer http.ResponseWriter, req *http.Request) {
	var note models.Note

	// Validar se o corpo da requisição é válido (json)
	err := json.NewDecoder(req.Body).Decode(&note)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(note)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(writer, fmt.Sprintf("%s", errors), http.StatusBadRequest)
		return
	}

	service.CreateNote(note)
	writer.WriteHeader(http.StatusCreated)
}

// UpdateNote updates a note
// @Summary Update a note
// @Description This method updates a note in the database
// @Tags Notes
// @Accept  json
// @Produce  json
// @Param id path int true "Note id"
// @Param note body models.Note true "Note"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /notes/{id} [put]
func handleUpdateNote(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req) // map[string]string -- id:value

	id, err := strconv.Atoi(vars["id"])

	// Validar se o id é um número inteiro
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	var newNote models.Note

	// Obter o objeto Note pelo id
	note, err := service.GetNoteById(id)

	// Valida se o id existe
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	// Alterar o objeto note com os dados recebidos
	err = json.NewDecoder(req.Body).Decode(&newNote)

	// Validar se o corpo da requisição é válido (json)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	/*
	* Alterar os campos do objeto note pelos dados recebidos
	 */
	if newNote.Title != "" {
		note.Title = newNote.Title
	}

	if newNote.Content != "" {
		note.Content = newNote.Content
	}

	if newNote.ShortDescription != "" {
		note.ShortDescription = newNote.ShortDescription
	}

	// Salvar o objeto Note atualizado
	service.UpdateNote(id, note)
	writer.WriteHeader(http.StatusOK)
}

// DeleteNote deletes a note
// @Summary Delete a note
// @Description This method deletes a note from the database
// @Tags Notes
// @Accept  json
// @Produce  json
// @Param id path int true "Note id"
// @Success 200
// @Failure 400
// @Router /notes/{id} [delete]
func handleDeleteNote(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req) // map[string]string -- id:value

	id, err := strconv.Atoi(vars["id"])

	// Validar se o id é um número inteiro
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	err = service.DeleteNote(id)

	// Valida se o id existe
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
