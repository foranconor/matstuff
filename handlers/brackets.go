package handlers

import (
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

func (h *Handler) CreateBracket(w http.ResponseWriter, r *http.Request) {
	body, ok := receiveJSON[struct{ Name string }](w, r)
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
		dbErr(w, err)
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
		dbErr(w, err)
		return
	}
	sendJSON(brackets, w)
}

func (h *Handler) GetBracket(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	bracket, err := db.GetBracket(tx, id)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(bracket, w)
}

func (h *Handler) UpdateBracket(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	b, ok := receiveJSON[domain.Bracket](w, r)
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
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(b, w)
}

func (h *Handler) ListBracketNotes(w http.ResponseWriter, r *http.Request) {
	h.listNotes(w, r, db.ListBracketNotes)
}

func (h *Handler) CreateBracketNote(w http.ResponseWriter, r *http.Request) {
	h.createNote(w, r, db.CreateBracketNote)
}

func (h *Handler) UpdateBracketNote(w http.ResponseWriter, r *http.Request) {
	h.updateNote(w, r)
}

func (h *Handler) DeleteBracketNote(w http.ResponseWriter, r *http.Request) {
	h.deleteNote(w, r, db.DeleteBracketNote)
}
