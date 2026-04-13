package db

import (
	"database/sql"
	"matstuff/domain"
)

func ListSuppliers(tx *sql.Tx) ([]domain.Supplier, error) {
	rows, err := tx.Query(`
		SELECT s.id, s.name, s.website, s.contact_id, c.name, s.created, s.modified
		FROM materials.suppliers s
		JOIN common.contacts c ON c.id = s.contact_id
		ORDER BY s.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []domain.Supplier
	for rows.Next() {
		var s domain.Supplier
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Website, &s.ContactID, &s.ContactName,
			&s.Created, &s.Modified,
		); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, rows.Err()
}

func GetSupplier(tx *sql.Tx, id int) (domain.SupplierWithContact, error) {
	var swc domain.SupplierWithContact
	err := tx.QueryRow(`
		SELECT s.id, s.name, s.website, s.contact_id, s.created, s.modified,
		       c.id, c.name, c.phone, c.email, c.created, c.modified
		FROM materials.suppliers s
		JOIN common.contacts c ON c.id = s.contact_id
		WHERE s.id = $1`, id).Scan(
		&swc.Supplier.ID, &swc.Supplier.Name, &swc.Supplier.Website,
		&swc.Supplier.ContactID, &swc.Supplier.Created, &swc.Supplier.Modified,
		&swc.Contact.ID, &swc.Contact.Name, &swc.Contact.Phone,
		&swc.Contact.Email, &swc.Contact.Created, &swc.Contact.Modified,
	)
	return swc, err
}

func CreateSupplier(tx *sql.Tx, name, website string, contactID int) (domain.Supplier, error) {
	var s domain.Supplier
	err := tx.QueryRow(`
		INSERT INTO materials.suppliers (name, website, contact_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, website, contact_id, created, modified`,
		name, website, contactID).Scan(
		&s.ID, &s.Name, &s.Website, &s.ContactID, &s.Created, &s.Modified,
	)
	return s, err
}

func ListSupplierMaterials(tx *sql.Tx, supplierID int) ([]domain.MaterialSupplier, error) {
	rows, err := tx.Query(`
		SELECT stm.id, stm.priority, stm.supplier_id, stm.material_id,
		       stm.price, (EXTRACT(EPOCH FROM stm.lead_time) / 86400)::int, stm.created, stm.modified, COALESCE(NULLIF(m.nickname, ''), m.name)
		FROM materials.suppliers_to_materials stm
		JOIN materials.materials m ON m.id = stm.material_id
		WHERE stm.supplier_id = $1
		ORDER BY stm.priority, COALESCE(NULLIF(m.nickname, ''), m.name)`, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.MaterialSupplier
	for rows.Next() {
		var ms domain.MaterialSupplier
		if err := rows.Scan(
			&ms.ID, &ms.Priority, &ms.SupplierID, &ms.MaterialID,
			&ms.Price, &ms.LeadTime, &ms.Created, &ms.Modified, &ms.MaterialName,
		); err != nil {
			return nil, err
		}
		result = append(result, ms)
	}
	return result, rows.Err()
}

func CreateMaterialSupplier(tx *sql.Tx, materialID, supplierID int) (domain.MaterialSupplier, error) {
	var ms domain.MaterialSupplier
	err := tx.QueryRow(`
		WITH ins AS (
			INSERT INTO materials.suppliers_to_materials (material_id, supplier_id, priority, price)
			VALUES ($1, $2, 1, 0)
			RETURNING *
		)
		SELECT ins.id, ins.priority, ins.supplier_id, ins.material_id, ins.price,
		       (EXTRACT(EPOCH FROM ins.lead_time) / 86400)::int,
		       ins.created, ins.modified, s.name, m.name
		FROM ins
		JOIN materials.suppliers s ON s.id = ins.supplier_id
		JOIN materials.materials m ON m.id = ins.material_id`,
		materialID, supplierID).Scan(
		&ms.ID, &ms.Priority, &ms.SupplierID, &ms.MaterialID, &ms.Price,
		&ms.LeadTime, &ms.Created, &ms.Modified, &ms.SupplierName, &ms.MaterialName,
	)
	return ms, err
}

func UpdateMaterialSupplier(tx *sql.Tx, ms domain.MaterialSupplier) error {
	_, err := tx.Exec(`
		UPDATE materials.suppliers_to_materials
		SET priority=$2, price=$3, lead_time=$4 * INTERVAL '1 day', modified=CURRENT_TIMESTAMP
		WHERE id=$1`,
		ms.ID, ms.Priority, ms.Price, ms.LeadTime,
	)
	return err
}

func DeleteMaterialSupplier(tx *sql.Tx, id int) error {
	_, err := tx.Exec(`DELETE FROM materials.suppliers_to_materials WHERE id=$1`, id)
	return err
}

func UpdateSupplier(tx *sql.Tx, s domain.Supplier) error {
	_, err := tx.Exec(`
		UPDATE materials.suppliers
		SET name=$2, website=$3, contact_id=$4, modified=CURRENT_TIMESTAMP
		WHERE id=$1`,
		s.ID, s.Name, s.Website, s.ContactID,
	)
	return err
}
