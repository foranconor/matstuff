package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func receiveJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var out T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "Can't read the request body"
		log.Println(msg, "error", err)
		http.Error(w, msg, http.StatusBadRequest)
		return out, false
	}
	err = json.Unmarshal(body, &out)
	if err != nil {
		msg := "Can't unmarshal JSON"
		log.Println(msg, "error", err)
		http.Error(w, msg, http.StatusBadRequest)
		return out, false
	}
	return out, true
}

func sendJSON[T any](data T, w http.ResponseWriter) bool {
	js, err := json.Marshal(data)
	if err != nil {
		msg := "Can't marshal JSON"
		log.Println(msg, "error", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return false
	}
	_, err = w.Write(js)
	if err != nil {
		msg := "Can't write JSON"
		log.Println(msg, "error", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return false
	}
	return true
}

const dbError = "Database error"

func sendError(w http.ResponseWriter, msg string, status int, errs ...error) {
	if len(errs) > 0 && errs[0] != nil {
		log.Println(msg, "error", errs[0])
	} else {
		log.Println(msg)
	}
	http.Error(w, msg, status)
}

func dbErr(w http.ResponseWriter, err error) {
	sendError(w, dbError, http.StatusInternalServerError, err)
}

func beginTx(db *sql.DB, w http.ResponseWriter) (*sql.Tx, bool) {
	tx, err := db.Begin()
	if err != nil {
		msg := "Failed to begin database transaction"
		log.Println(msg, "error", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return nil, false
	}
	return tx, true
}

func commitTx(tx *sql.Tx, w http.ResponseWriter) bool {
	err := tx.Commit()
	if err != nil {
		msg := "Failed to commit database transaction"
		log.Println(msg, "error", err)

		http.Error(w, dbError, http.StatusInternalServerError)
		return false
	}
	return true
}
