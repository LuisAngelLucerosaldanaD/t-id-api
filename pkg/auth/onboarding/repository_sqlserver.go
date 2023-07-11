package onboarding

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/models"
)

// sqlServer estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newOnboardingSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *Onboarding) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const sqlInsert = `INSERT INTO auth.onboarding (id ,client_id, request_id, user_id, created_at, updated_at) VALUES (:id ,:client_id, :request_id, :user_id:created_at, :updated_at) `
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) update(m *Onboarding) error {
	date := time.Now()
	m.UpdatedAt = date
	const sqlUpdate = `UPDATE auth.onboarding SET client_id = :client_id, request_id = :request_id, user_id = :user_id, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) delete(id string) error {
	const sqlDelete = `DELETE FROM auth.onboarding WHERE id = :id `
	m := Onboarding{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) getByID(id string) (*Onboarding, error) {
	const sqlGetByID = `SELECT convert(nvarchar(50), id) id , client_id, request_id, user_id, created_at, updated_at FROM auth.onboarding  WITH (NOLOCK)  WHERE id = @id `
	mdl := Onboarding{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) getAll() ([]*Onboarding, error) {
	var ms []*Onboarding
	const sqlGetAll = `SELECT convert(nvarchar(50), id) id , client_id, request_id, user_id, created_at, updated_at FROM auth.onboarding  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
