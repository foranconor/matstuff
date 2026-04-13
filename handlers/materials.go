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
	body, ok := recieveJSON[struct{ Name string }](w, r)
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
		http.Error(w, dbError, http.StatusInternalServerError)
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
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(materials, w)
}

func (h *Handler) GetMaterial(w http.ResponseWriter, r *http.Request) {
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

	material, err := db.GetMaterial(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(material, w)
}

func (h *Handler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	m, ok := recieveJSON[domain.Material](w, r)
	if !ok {
		return
	}
	m.ID = id

	validTreatments := map[string]bool{"UT": true, "FDA": true, "MCA": true, "H1.2": true, "H3.2": true, "H4": true, "H5": true}
	validUnits := map[string]bool{"mm": true, "sheet": true}
	if !validTreatments[m.Treatment] {
		http.Error(w, "invalid treatment", http.StatusBadRequest)
		return
	}
	if !validUnits[m.Units] {
		http.Error(w, "invalid units", http.StatusBadRequest)
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateMaterial(tx, m); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
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
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	u, ok := recieveJSON[domain.Uses](w, r)
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
		http.Error(w, dbError, http.StatusInternalServerError)
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
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	suppliers, err := db.ListMaterialSuppliers(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(suppliers, w)
}

func (h *Handler) ListMaterialNotes(w http.ResponseWriter, r *http.Request) {
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

	notes, err := db.ListMaterialNotes(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(notes, w)
}

func (h *Handler) CreateMaterialNote(w http.ResponseWriter, r *http.Request) {
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

	note, err := db.CreateMaterialNote(tx, id, body.Content)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(note, w)
}

func (h *Handler) UpdateMaterialNote(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) DeleteMaterialNote(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DeleteMaterialNote(tx, id, noteID); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
