package db

import (
	"database/sql"
	"fmt"
	"matstuff/domain"
)

func listNotes(tx *sql.Tx, junctionTable, entityCol string, entityID int) ([]domain.Note, error) {
	rows, err := tx.Query(fmt.Sprintf(`
		SELECT n.id, n.content, n.created, n.modified
		FROM common.notes n
		JOIN %s j ON j.note_id = n.id
		WHERE j.%s = $1
		ORDER BY n.created DESC`, junctionTable, entityCol), entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		var n domain.Note
		if err := rows.Scan(&n.ID, &n.Content, &n.Created, &n.Modified); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func createNote(tx *sql.Tx, junctionTable, entityCol string, entityID int, content string) (domain.Note, error) {
	var n domain.Note
	err := tx.QueryRow(`
		INSERT INTO common.notes (content) VALUES ($1)
		RETURNING id, content, created, modified`, content).
		Scan(&n.ID, &n.Content, &n.Created, &n.Modified)
	if err != nil {
		return n, err
	}
	_, err = tx.Exec(fmt.Sprintf(`
		INSERT INTO %s (%s, note_id) VALUES ($1, $2)`, junctionTable, entityCol),
		entityID, n.ID)
	return n, err
}

func deleteNote(tx *sql.Tx, junctionTable, entityCol string, entityID, noteID int) error {
	_, err := tx.Exec(fmt.Sprintf(`
		DELETE FROM %s WHERE %s=$1 AND note_id=$2`, junctionTable, entityCol),
		entityID, noteID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM common.notes WHERE id=$1`, noteID)
	return err
}

func UpdateNote(tx *sql.Tx, id int, content string) (domain.Note, error) {
	var n domain.Note
	err := tx.QueryRow(`
		UPDATE common.notes SET content=$2, modified=CURRENT_TIMESTAMP
		WHERE id=$1
		RETURNING id, content, created, modified`, id, content).
		Scan(&n.ID, &n.Content, &n.Created, &n.Modified)
	return n, err
}

func ListMaterialNotes(tx *sql.Tx, materialID int) ([]domain.Note, error) {
	return listNotes(tx, "materials.materials_notes", "material_id", materialID)
}

func CreateMaterialNote(tx *sql.Tx, materialID int, content string) (domain.Note, error) {
	return createNote(tx, "materials.materials_notes", "material_id", materialID, content)
}

func DeleteMaterialNote(tx *sql.Tx, materialID, noteID int) error {
	return deleteNote(tx, "materials.materials_notes", "material_id", materialID, noteID)
}

func ListBracketNotes(tx *sql.Tx, bracketID int) ([]domain.Note, error) {
	return listNotes(tx, "handrails.brackets_notes", "bracket_id", bracketID)
}

func CreateBracketNote(tx *sql.Tx, bracketID int, content string) (domain.Note, error) {
	return createNote(tx, "handrails.brackets_notes", "bracket_id", bracketID, content)
}

func DeleteBracketNote(tx *sql.Tx, bracketID, noteID int) error {
	return deleteNote(tx, "handrails.brackets_notes", "bracket_id", bracketID, noteID)
}
