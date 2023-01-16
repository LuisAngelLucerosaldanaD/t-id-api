package status_request

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

func newStatusRequestPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *StatusRequest) error {
	const psqlInsert = `INSERT INTO wf.status_request (status, description, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Status,
		m.Description,
		m.UserId,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *StatusRequest) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE wf.status_request SET status = :status, description = :description, user_id = :user_id, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM wf.status_request WHERE id = :id `
	m := StatusRequest{ID: id}
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
func (s *psql) getByID(id int64) (*StatusRequest, error) {
	const psqlGetByID = `SELECT id , status, description, user_id, created_at, updated_at FROM wf.status_request WHERE id = $1 `
	mdl := StatusRequest{}
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
func (s *psql) getAll() ([]*StatusRequest, error) {
	var ms []*StatusRequest
	const psqlGetAll = ` SELECT id , status, description, user_id, created_at, updated_at FROM wf.status_request `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByUserId(userID string) (*StatusRequest, error) {
	const psqlGetByUserID = ` SELECT id , status, description, user_id, created_at, updated_at FROM wf.status_request where user_id = $1 `
	mdl := StatusRequest{}
	err := s.DB.Get(&mdl, psqlGetByUserID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getByStatus(status string) ([]*StatusRequest, error) {
	var ms []*StatusRequest
	const psqlGetByStatus = ` SELECT id , status, description, user_id, created_at, updated_at FROM wf.status_request where status = $1`

	err := s.DB.Select(&ms, psqlGetByStatus, status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// Update actualiza un registro en la BD
func (s *psql) updateStatus(status string, description string, userID string) error {
	date := time.Now()
	m := StatusRequest{
		Status:      status,
		UserId:      userID,
		Description: description,
		UpdatedAt:   date,
	}
	const psqlUpdate = `UPDATE wf.status_request SET status = :status, description = :description, updated_at = :updated_at WHERE user_id = :user_id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
