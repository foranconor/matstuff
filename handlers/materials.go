package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	body, ok := receiveJSON[struct{ Name string }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	mwu, err := db.CreateMaterial(tx, body.Name)
	if err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(mwu, w)
}

func (h *Handler) ListMaterials(w http.ResponseWriter, r *http.Request) {
	showAll := r.URL.Query().Get("all") == "true"

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	materials, err := db.ListMaterials(tx, showAll)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(materials, w)
}

func (h *Handler) GetMaterial(w http.ResponseWriter, r *http.Request) {
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

	material, err := db.GetMaterial(tx, id)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(material, w)
}

func (h *Handler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	m, ok := receiveJSON[domain.Material](w, r)
	if !ok {
		return
	}
	m.ID = id

	validTreatments := map[string]bool{"UT": true, "FDA": true, "MCA": true, "H1.2": true, "H3.2": true, "H4": true, "H5": true}
	validUnits := map[string]bool{"mm": true, "sheet": true}
	if !validTreatments[m.Treatment] {
		sendError(w, "invalid treatment", http.StatusBadRequest)
		return
	}
	if !validUnits[m.Units] {
		sendError(w, "invalid units", http.StatusBadRequest)
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateMaterial(tx, m); err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(m, w)
}

func (h *Handler) UpdateUses(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	u, ok := receiveJSON[domain.Uses](w, r)
	if !ok {
		return
	}
	u.MaterialID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateUses(tx, u); err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(u, w)
}

func (h *Handler) ListMaterialSuppliers(w http.ResponseWriter, r *http.Request) {
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

	suppliers, err := db.ListMaterialSuppliers(tx, id)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(suppliers, w)
}

func (h *Handler) ListMaterialNotes(w http.ResponseWriter, r *http.Request) {
	h.listNotes(w, r, db.ListMaterialNotes)
}

func (h *Handler) CreateMaterialNote(w http.ResponseWriter, r *http.Request) {
	h.createNote(w, r, db.CreateMaterialNote)
}

func (h *Handler) UpdateMaterialNote(w http.ResponseWriter, r *http.Request) {
	h.updateNote(w, r)
}

func (h *Handler) DeleteMaterialNote(w http.ResponseWriter, r *http.Request) {
	h.deleteNote(w, r, db.DeleteMaterialNote)
}
