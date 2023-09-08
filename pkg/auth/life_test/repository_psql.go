package life_test

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

func newLifeTestPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *LifeTest) error {
	const psqlInsert = `INSERT INTO cfg.life_test (client_id, max_num_validation, request_id, expired_at, user_id, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.ClientId,
		m.MaxNumTest,
		m.RequestId,
		m.ExpiredAt,
		m.UserID,
		m.Status,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *LifeTest) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.life_test SET client_id = :client_id, max_num_validation = :max_num_validation, request_id = :request_id, expired_at = :expired_at, user_id = :user_id, status = :status, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM cfg.life_test WHERE id = :id `
	m := LifeTest{ID: id}
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
func (s *psql) getByID(id int64) (*LifeTest, error) {
	const psqlGetByID = `SELECT id , client_id, max_num_validation, request_id, expired_at, user_id, status, created_at, updated_at FROM cfg.life_test WHERE id = $1 `
	mdl := LifeTest{}
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
func (s *psql) getAll() ([]*LifeTest, error) {
	var ms []*LifeTest
	const psqlGetAll = ` SELECT id , client_id, max_num_validation, request_id, expired_at, user_id, status, created_at, updated_at FROM cfg.life_test `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByClientIDAndRequestID consulta un registro por su ID
func (s *psql) getByClientIDAndRequestID(clientIid int64, requestID string) (*LifeTest, error) {
	const psqlGetByClientID = `SELECT id , client_id, max_num_validation, request_id, expired_at, user_id, status, created_at, updated_at FROM cfg.life_test WHERE client_id = %d and request_id = '%s' limit 1`
	mdl := LifeTest{}
	query := fmt.Sprintf(psqlGetByClientID, clientIid, requestID)
	err := s.DB.Get(&mdl, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) updateStatus(m *LifeTest) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.life_test SET status = :status, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAllByUserId(userID string) ([]*LifeTest, error) {
	var ms []*LifeTest
	const psqlGetAll = ` SELECT id , client_id, max_num_validation, request_id, expired_at, user_id, status, created_at, updated_at FROM cfg.life_test where user_id = $1;`

	err := s.DB.Select(&ms, psqlGetAll, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByUserID(userId string) (*LifeTest, error) {
	const psqlGetByID = `SELECT id , client_id, max_num_validation, request_id, expired_at, user_id, status, created_at, updated_at FROM cfg.life_test WHERE id = $1 order by id desc limit 1`
	mdl := LifeTest{}
	err := s.DB.Get(&mdl, psqlGetByID, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}