package work_validation

import (
	"database/sql"
	"fmt"
	"time"

	"check-id-api/internal/models"
	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newWorkValidationPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *WorkValidation) error {
	const psqlInsert = `INSERT INTO wf.work_validation (status, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Status,
		m.UserId,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *WorkValidation) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE wf.work_validation SET status = :status, user_id = :user_id, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) delete(id int64) error {
	const psqlDelete = `DELETE FROM wf.work_validation WHERE id = :id `
	m := WorkValidation{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByID(id int64) (*WorkValidation, error) {
	const psqlGetByID = `SELECT id , status, user_id, created_at, updated_at FROM wf.work_validation WHERE id = $1 `
	mdl := WorkValidation{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*WorkValidation, error) {
	var ms []*WorkValidation
	const psqlGetAll = ` SELECT id , status, user_id, created_at, updated_at FROM wf.work_validation `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByUserId(userID string) (*WorkValidation, error) {
	const psqlGetByID = `SELECT id , status, user_id, created_at, updated_at FROM wf.work_validation WHERE user_id = $1 `
	mdl := WorkValidation{}
	err := s.DB.Get(&mdl, psqlGetByID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getByStatus(status string) ([]*WorkValidation, error) {
	var ms []*WorkValidation
	const psqlGetAll = ` SELECT id , status, user_id, created_at, updated_at FROM wf.work_validation where status = $1`

	err := s.DB.Select(&ms, psqlGetAll, status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
