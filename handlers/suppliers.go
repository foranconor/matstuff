package handlers

import (
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

func (h *Handler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	body, ok := receiveJSON[struct {
		Name      string
		Website   string
		ContactID int
	}](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	s, err := db.CreateSupplier(tx, body.Name, body.Website, body.ContactID)
	if err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(s, w)
}

func (h *Handler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	suppliers, err := db.ListSuppliers(tx)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(suppliers, w)
}

func (h *Handler) GetSupplier(w http.ResponseWriter, r *http.Request) {
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

	supplier, err := db.GetSupplier(tx, id)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(supplier, w)
}

func (h *Handler) ListSupplierMaterials(w http.ResponseWriter, r *http.Request) {
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

	materials, err := db.ListSupplierMaterials(tx, id)
	if err != nil {
		dbErr(w, err)
		return
	}
	sendJSON(materials, w)
}

func (h *Handler) AddSupplierToMaterial(w http.ResponseWriter, r *http.Request) {
	materialID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	body, ok := receiveJSON[struct{ SupplierID int }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	ms, err := db.CreateMaterialSupplier(tx, materialID, body.SupplierID)
	if err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(ms, w)
}

func (h *Handler) AddMaterialToSupplier(w http.ResponseWriter, r *http.Request) {
	supplierID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	body, ok := receiveJSON[struct{ MaterialID int }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	ms, err := db.CreateMaterialSupplier(tx, body.MaterialID, supplierID)
	if err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(ms, w)
}

func (h *Handler) UpdateMaterialSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	ms, ok := receiveJSON[domain.MaterialSupplier](w, r)
	if !ok {
		return
	}
	ms.ID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateMaterialSupplier(tx, ms); err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(ms, w)
}

func (h *Handler) DeleteMaterialSupplier(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DeleteMaterialSupplier(tx, id); err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, "invalid id", http.StatusBadRequest, err)
		return
	}

	s, ok := receiveJSON[domain.Supplier](w, r)
	if !ok {
		return
	}
	s.ID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateSupplier(tx, s); err != nil {
		dbErr(w, err)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(s, w)
}
