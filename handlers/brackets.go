package handlers

import (
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

func (h *Handler) CreateBracket(w http.ResponseWriter, r *http.Request) {
	body, ok := recieveJSON[struct{ Name string }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	b, err := db.CreateBracket(tx, body.Name)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(b, w)
}

func (h *Handler) ListBrackets(w http.ResponseWriter, r *http.Request) {
	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	showAll := r.URL.Query().Get("all") == "true"
	brackets, err := db.ListBrackets(tx, showAll)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(brackets, w)
}

func (h *Handler) GetBracket(w http.ResponseWriter, r *http.Request) {
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

	bracket, err := db.GetBracket(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(bracket, w)
}

func (h *Handler) UpdateBracket(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	b, ok := recieveJSON[domain.Bracket](w, r)
	if !ok {
		return
	}
	b.ID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateBracket(tx, b); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(b, w)
}

func (h *Handler) ListBracketNotes(w http.ResponseWriter, r *http.Request) {
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

	notes, err := db.ListBracketNotes(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(notes, w)
}

func (h *Handler) CreateBracketNote(w http.ResponseWriter, r *http.Request) {
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

	note, err := db.CreateBracketNote(tx, id, body.Content)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(note, w)
}

func (h *Handler) UpdateBracketNote(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) DeleteBracketNote(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DeleteBracketNote(tx, id, noteID); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
