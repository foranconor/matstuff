package handlers

import (
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	body, ok := recieveJSON[struct {
		Name  string
		Phone string
		Email string
	}](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	c, err := db.CreateContact(tx, body.Name, body.Phone, body.Email)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(c, w)
}

func (h *Handler) ListContacts(w http.ResponseWriter, r *http.Request) {
	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	contacts, err := db.ListContacts(tx)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(contacts, w)
}

func (h *Handler) GetContact(w http.ResponseWriter, r *http.Request) {
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

	contact, err := db.GetContact(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(contact, w)
}

func (h *Handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	c, ok := recieveJSON[domain.Contact](w, r)
	if !ok {
		return
	}
	c.ID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateContact(tx, c); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(c, w)
}
