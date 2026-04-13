package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

func (h *Handler) listNotes(w http.ResponseWriter, r *http.Request, dbFn func(*sql.Tx, int) ([]domain.Note, error)) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	notes, err := dbFn(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(notes, w)
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request, dbFn func(*sql.Tx, int, string) (domain.Note, error)) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	body, ok := recieveJSON[struct{ Content string }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	note, err := dbFn(tx, id, body.Content)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(note, w)
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := strconv.Atoi(r.PathValue("nid"))
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	body, ok := recieveJSON[struct{ Content string }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	note, err := db.UpdateNote(tx, noteID, body.Content)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(note, w)
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request, dbFn func(*sql.Tx, int, int) error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	noteID, err := strconv.Atoi(r.PathValue("nid"))
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := dbFn(tx, id, noteID); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
