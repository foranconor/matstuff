package db

import (
	"database/sql"
	"matstuff/domain"
)

func ListProfiles(tx *sql.Tx) ([]domain.Profile, error) {
	rows, err := tx.Query(`
		SELECT id, name, created, modified
		FROM handrails.profiles
		ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []domain.Profile
	for rows.Next() {
		var p domain.Profile
		if err := rows.Scan(&p.ID, &p.Name, &p.Created, &p.Modified); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

func GetProfile(tx *sql.Tx, id int) (domain.Profile, error) {
	var p domain.Profile
	err := tx.QueryRow(`
		SELECT id, name, created, modified
		FROM handrails.profiles
		WHERE id = $1`, id).Scan(&p.ID, &p.Name, &p.Created, &p.Modified)
	return p, err
}

func CreateProfile(tx *sql.Tx, name string) (domain.Profile, error) {
	var p domain.Profile
	err := tx.QueryRow(`
		INSERT INTO handrails.profiles (name)
		VALUES ($1)
		RETURNING id, name, created, modified`, name).
		Scan(&p.ID, &p.Name, &p.Created, &p.Modified)
	return p, err
}

func UpdateProfile(tx *sql.Tx, p domain.Profile) error {
	_, err := tx.Exec(`
		UPDATE handrails.profiles
		SET name=$2, modified=CURRENT_TIMESTAMP
		WHERE id=$1`,
		p.ID, p.Name)
	return err
}

func DeleteProfile(tx *sql.Tx, id int) error {
	_, err := tx.Exec(`DELETE FROM handrails.profiles WHERE id=$1`, id)
	return err
}

func scanMaterialProfiles(rows *sql.Rows) ([]domain.MaterialProfile, error) {
	defer rows.Close()
	var result []domain.MaterialProfile
	for rows.Next() {
		var mp domain.MaterialProfile
		if err := rows.Scan(
			&mp.ID, &mp.MaterialID, &mp.ProfileID, &mp.Price,
			&mp.Created, &mp.Modified, &mp.MaterialName, &mp.ProfileName,
		); err != nil {
			return nil, err
		}
		result = append(result, mp)
	}
	return result, rows.Err()
}

func ListProfileMaterials(tx *sql.Tx, profileID int) ([]domain.MaterialProfile, error) {
	rows, err := tx.Query(`
		SELECT mp.id, mp.material_id, mp.profile_id, mp.price,
		       mp.created, mp.modified,
		       COALESCE(NULLIF(m.nickname, ''), m.name), p.name
		FROM handrails.material_profiles mp
		JOIN materials.materials m ON m.id = mp.material_id
		JOIN handrails.profiles p ON p.id = mp.profile_id
		WHERE mp.profile_id = $1
		ORDER BY COALESCE(NULLIF(m.nickname, ''), m.name)`, profileID)
	if err != nil {
		return nil, err
	}
	return scanMaterialProfiles(rows)
}

func ListMaterialProfiles(tx *sql.Tx, materialID int) ([]domain.MaterialProfile, error) {
	rows, err := tx.Query(`
		SELECT mp.id, mp.material_id, mp.profile_id, mp.price,
		       mp.created, mp.modified,
		       COALESCE(NULLIF(m.nickname, ''), m.name), p.name
		FROM handrails.material_profiles mp
		JOIN materials.materials m ON m.id = mp.material_id
		JOIN handrails.profiles p ON p.id = mp.profile_id
		WHERE mp.material_id = $1
		ORDER BY p.name`, materialID)
	if err != nil {
		return nil, err
	}
	return scanMaterialProfiles(rows)
}

func CreateMaterialProfile(tx *sql.Tx, materialID, profileID int) (domain.MaterialProfile, error) {
	var mp domain.MaterialProfile
	err := tx.QueryRow(`
		WITH ins AS (
			INSERT INTO handrails.material_profiles (material_id, profile_id, price)
			VALUES ($1, $2, 0)
			RETURNING *
		)
		SELECT ins.id, ins.material_id, ins.profile_id, ins.price,
		       ins.created, ins.modified,
		       COALESCE(NULLIF(m.nickname, ''), m.name), p.name
		FROM ins
		JOIN materials.materials m ON m.id = ins.material_id
		JOIN handrails.profiles p ON p.id = ins.profile_id`,
		materialID, profileID).Scan(
		&mp.ID, &mp.MaterialID, &mp.ProfileID, &mp.Price,
		&mp.Created, &mp.Modified, &mp.MaterialName, &mp.ProfileName,
	)
	return mp, err
}

func UpdateMaterialProfile(tx *sql.Tx, mp domain.MaterialProfile) error {
	_, err := tx.Exec(`
		UPDATE handrails.material_profiles
		SET price=$2, modified=CURRENT_TIMESTAMP
		WHERE id=$1`,
		mp.ID, mp.Price)
	return err
}

func DeleteMaterialProfile(tx *sql.Tx, id int) error {
	_, err := tx.Exec(`DELETE FROM handrails.material_profiles WHERE id=$1`, id)
	return err
}
