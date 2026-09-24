package db

import (
	"database/sql"
	"matstuff/domain"
)

func ListMaterials(tx *sql.Tx, showAll bool) ([]domain.Material, error) {
	query := `
		SELECT m.id, m.name, m.nickname, m.treatment, m.blurb, m.units,
		       m.length, m.width, m.thickness, m.max_span, m.density,
		       m.max_overhang, m.radius, m.color, u.published, u.archived, m.created, m.modified
		FROM materials.materials m
		JOIN materials.uses u ON u.material_id = m.id
		WHERE true`

	if !showAll {
		query += ` AND u.archived = false AND u.published = true`
	}

	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.Material
	for rows.Next() {
		var m domain.Material
		err := rows.Scan(
			&m.ID, &m.Name, &m.Nickname, &m.Treatment, &m.Blurb, &m.Units,
			&m.Length, &m.Width, &m.Thickness, &m.MaxSpan, &m.Density,
			&m.MaxOverhang, &m.Radius, &m.Color, &m.Published, &m.Archived, &m.Created, &m.Modified,
		)
		if err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}
	return materials, rows.Err()
}

func GetMaterial(tx *sql.Tx, id int) (domain.MaterialWithUses, error) {
	var mwu domain.MaterialWithUses
	err := tx.QueryRow(`
		SELECT m.id, m.name, m.nickname, m.treatment, m.blurb, m.units,
		       m.length, m.width, m.thickness, m.max_span, m.density,
		       m.max_overhang, m.radius, m.color, m.created, m.modified,
		       u.id, u.material_id, u.exterior, u.interior, u.early_access,
		       u.stringers, u.risers, u.treads, u.timber, u.panel,
		       u.published, u.archived, u.handrail, u.created, u.modified
		FROM materials.materials m
		JOIN materials.uses u ON u.material_id = m.id
		WHERE m.id = $1`, id).Scan(
		&mwu.Material.ID, &mwu.Material.Name, &mwu.Material.Nickname,
		&mwu.Material.Treatment, &mwu.Material.Blurb, &mwu.Material.Units,
		&mwu.Material.Length, &mwu.Material.Width, &mwu.Material.Thickness,
		&mwu.Material.MaxSpan, &mwu.Material.Density, &mwu.Material.MaxOverhang,
		&mwu.Material.Radius, &mwu.Material.Color, &mwu.Material.Created, &mwu.Material.Modified,
		&mwu.Uses.ID, &mwu.Uses.MaterialID, &mwu.Uses.Exterior, &mwu.Uses.Interior,
		&mwu.Uses.EarlyAccess, &mwu.Uses.Stringers, &mwu.Uses.Risers, &mwu.Uses.Treads,
		&mwu.Uses.Timber, &mwu.Uses.Panel, &mwu.Uses.Published, &mwu.Uses.Archived,
		&mwu.Uses.Handrail, &mwu.Uses.Created, &mwu.Uses.Modified,
	)
	return mwu, err
}

func UpdateMaterial(tx *sql.Tx, m domain.Material) error {
	_, err := tx.Exec(`
		UPDATE materials.materials
		SET name=$2, nickname=$3, treatment=$4, blurb=$5, units=$6,
		    length=$7, width=$8, thickness=$9, max_span=$10, density=$11,
		    max_overhang=$12, radius=$13, color=$14, modified=CURRENT_TIMESTAMP
		WHERE id=$1`,
		m.ID, m.Name, m.Nickname, m.Treatment, m.Blurb, m.Units,
		m.Length, m.Width, m.Thickness, m.MaxSpan, m.Density,
		m.MaxOverhang, m.Radius, m.Color,
	)
	return err
}

func UpdateUses(tx *sql.Tx, u domain.Uses) error {
	_, err := tx.Exec(`
		UPDATE materials.uses
		SET exterior=$2, interior=$3, early_access=$4, stringers=$5, risers=$6,
		    treads=$7, timber=$8, panel=$9, published=$10, archived=$11,
		    handrail=$12, modified=CURRENT_TIMESTAMP
		WHERE material_id=$1`,
		u.MaterialID, u.Exterior, u.Interior, u.EarlyAccess, u.Stringers,
		u.Risers, u.Treads, u.Timber, u.Panel, u.Published, u.Archived, u.Handrail,
	)
	return err
}

func ListMaterialSuppliers(tx *sql.Tx, materialID int) ([]domain.MaterialSupplier, error) {
	rows, err := tx.Query(`
		SELECT stm.id, stm.priority, stm.supplier_id, stm.material_id,
		       stm.price, (EXTRACT(EPOCH FROM stm.lead_time) / 86400)::int, stm.created, stm.modified, s.name
		FROM materials.suppliers_to_materials stm
		JOIN materials.suppliers s ON s.id = stm.supplier_id
		WHERE stm.material_id = $1
		ORDER BY stm.priority`, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.MaterialSupplier
	for rows.Next() {
		var ms domain.MaterialSupplier
		if err := rows.Scan(
			&ms.ID, &ms.Priority, &ms.SupplierID, &ms.MaterialID,
			&ms.Price, &ms.LeadTime, &ms.Created, &ms.Modified, &ms.SupplierName,
		); err != nil {
			return nil, err
		}
		result = append(result, ms)
	}
	return result, rows.Err()
}

func CreateMaterial(tx *sql.Tx, name string) (domain.MaterialWithUses, error) {
	var mwu domain.MaterialWithUses
	err := tx.QueryRow(`
		INSERT INTO materials.materials (name)
		VALUES ($1)
		RETURNING id, name, nickname, treatment, blurb, units,
		          length, width, thickness, max_span, density, max_overhang, radius,
		          color, created, modified`, name).Scan(
		&mwu.Material.ID, &mwu.Material.Name, &mwu.Material.Nickname,
		&mwu.Material.Treatment, &mwu.Material.Blurb, &mwu.Material.Units,
		&mwu.Material.Length, &mwu.Material.Width, &mwu.Material.Thickness,
		&mwu.Material.MaxSpan, &mwu.Material.Density, &mwu.Material.MaxOverhang,
		&mwu.Material.Radius, &mwu.Material.Color, &mwu.Material.Created, &mwu.Material.Modified,
	)
	if err != nil {
		return mwu, err
	}
	err = tx.QueryRow(`
		INSERT INTO materials.uses (material_id)
		VALUES ($1)
		RETURNING id, material_id, exterior, interior, early_access,
		          stringers, risers, treads, timber, panel, published, archived,
		          handrail, created, modified`, mwu.Material.ID).Scan(
		&mwu.Uses.ID, &mwu.Uses.MaterialID, &mwu.Uses.Exterior, &mwu.Uses.Interior,
		&mwu.Uses.EarlyAccess, &mwu.Uses.Stringers, &mwu.Uses.Risers, &mwu.Uses.Treads,
		&mwu.Uses.Timber, &mwu.Uses.Panel, &mwu.Uses.Published, &mwu.Uses.Archived,
		&mwu.Uses.Handrail, &mwu.Uses.Created, &mwu.Uses.Modified,
	)
	return mwu, err
}

