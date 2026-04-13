package db

import (
	"database/sql"
	"matstuff/domain"
)

func ListBrackets(tx *sql.Tx, showAll bool) ([]domain.Bracket, error) {
	query := `
		SELECT b.id, b.name, b.nickname, b.supplier_id, COALESCE(b.price, 0),
		       b.published, b.archived, b.created, COALESCE(s.name, '')
		FROM handrails.brackets b
		LEFT JOIN materials.suppliers s ON s.id = b.supplier_id
		WHERE true`

	if !showAll {
		query += ` AND b.archived = false AND b.published = true`
	}

	query += ` ORDER BY b.name`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brackets []domain.Bracket
	for rows.Next() {
		var b domain.Bracket
		var supplierID sql.NullInt64
		if err := rows.Scan(
			&b.ID, &b.Name, &b.Nickname, &supplierID, &b.Price,
			&b.Published, &b.Archived, &b.Created, &b.SupplierName,
		); err != nil {
			return nil, err
		}
		if supplierID.Valid {
			id := int(supplierID.Int64)
			b.SupplierID = &id
		}
		brackets = append(brackets, b)
	}
	return brackets, rows.Err()
}

func GetBracket(tx *sql.Tx, id int) (domain.Bracket, error) {
	var b domain.Bracket
	var supplierID sql.NullInt64
	err := tx.QueryRow(`
		SELECT b.id, b.name, b.nickname, b.supplier_id, COALESCE(b.price, 0),
		       b.published, b.archived, b.created, COALESCE(s.name, '')
		FROM handrails.brackets b
		LEFT JOIN materials.suppliers s ON s.id = b.supplier_id
		WHERE b.id = $1`, id).
		Scan(&b.ID, &b.Name, &b.Nickname, &supplierID, &b.Price,
			&b.Published, &b.Archived, &b.Created, &b.SupplierName)
	if supplierID.Valid {
		id := int(supplierID.Int64)
		b.SupplierID = &id
	}
	return b, err
}

func UpdateBracket(tx *sql.Tx, b domain.Bracket) error {
	_, err := tx.Exec(`
		UPDATE handrails.brackets
		SET name=$2, nickname=$3, supplier_id=$4, price=$5, published=$6, archived=$7
		WHERE id=$1`,
		b.ID, b.Name, b.Nickname, b.SupplierID, b.Price, b.Published, b.Archived,
	)
	return err
}

func CreateBracket(tx *sql.Tx, name string) (domain.Bracket, error) {
	var b domain.Bracket
	err := tx.QueryRow(`
		INSERT INTO handrails.brackets (name)
		VALUES ($1)
		RETURNING id, name, COALESCE(nickname, ''), supplier_id, COALESCE(price, 0),
		          published, archived, created`,
		name).Scan(&b.ID, &b.Name, &b.Nickname, &b.SupplierID, &b.Price,
		&b.Published, &b.Archived, &b.Created)
	return b, err
}

