package handlers

import (
	"net/http"
	"strconv"

	"matstuff/db"
	"matstuff/domain"
)

func (h *Handler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	profiles, err := db.ListProfiles(tx)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(profiles, w)
}

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	body, ok := recieveJSON[struct{ Name string }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	p, err := db.CreateProfile(tx, body.Name)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(p, w)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
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

	p, err := db.GetProfile(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(p, w)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, ok := recieveJSON[domain.Profile](w, r)
	if !ok {
		return
	}
	p.ID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateProfile(tx, p); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(p, w)
}

func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DeleteProfile(tx, id); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListProfileMaterials(w http.ResponseWriter, r *http.Request) {
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

	materials, err := db.ListProfileMaterials(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(materials, w)
}

func (h *Handler) AddMaterialToProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	body, ok := recieveJSON[struct{ MaterialID int }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	mp, err := db.CreateMaterialProfile(tx, body.MaterialID, profileID)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(mp, w)
}

func (h *Handler) ListMaterialProfiles(w http.ResponseWriter, r *http.Request) {
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

	profiles, err := db.ListMaterialProfiles(tx, id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	sendJSON(profiles, w)
}

func (h *Handler) AddProfileToMaterial(w http.ResponseWriter, r *http.Request) {
	materialID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	body, ok := recieveJSON[struct{ ProfileID int }](w, r)
	if !ok {
		return
	}

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	mp, err := db.CreateMaterialProfile(tx, materialID, body.ProfileID)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(mp, w)
}

func (h *Handler) UpdateMaterialProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mp, ok := recieveJSON[domain.MaterialProfile](w, r)
	if !ok {
		return
	}
	mp.ID = id

	tx, ok := beginTx(h.DB, w)
	if !ok {
		return
	}
	defer tx.Rollback()

	if err := db.UpdateMaterialProfile(tx, mp); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	sendJSON(mp, w)
}

func (h *Handler) DeleteMaterialProfile(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DeleteMaterialProfile(tx, id); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}
	if !commitTx(tx, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
