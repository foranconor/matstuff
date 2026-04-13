package db

import (
	"database/sql"
	"matstuff/domain"
)

func ListContacts(tx *sql.Tx) ([]domain.Contact, error) {
	rows, err := tx.Query(`
		SELECT id, name, phone, email, created, modified
		FROM common.contacts
		ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []domain.Contact
	for rows.Next() {
		var c domain.Contact
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Email, &c.Created, &c.Modified,
		); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func GetContact(tx *sql.Tx, id int) (domain.Contact, error) {
	var c domain.Contact
	err := tx.QueryRow(`
		SELECT id, name, phone, email, created, modified
		FROM common.contacts
		WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.Created, &c.Modified)
	return c, err
}

func CreateContact(tx *sql.Tx, name, phone, email string) (domain.Contact, error) {
	var c domain.Contact
	err := tx.QueryRow(`
		INSERT INTO common.contacts (name, phone, email)
		VALUES ($1, $2, $3)
		RETURNING id, name, phone, email, created, modified`,
		name, phone, email).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.Created, &c.Modified,
	)
	return c, err
}

func UpdateContact(tx *sql.Tx, c domain.Contact) error {
	_, err := tx.Exec(`
		UPDATE common.contacts
		SET name=$2, phone=$3, email=$4, modified=CURRENT_TIMESTAMP
		WHERE id=$1`,
		c.ID, c.Name, c.Phone, c.Email,
	)
	return err
}
