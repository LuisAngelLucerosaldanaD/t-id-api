package user_role

import (
	"database/sql"
	"errors"
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

func newUseRolePsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *UseRole) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.user_role (id ,user_id, role_id, created_at, updated_at) VALUES (:id ,:user_id, :role_id,:created_at, :updated_at) `
	rs, err := s.DB.NamedExec(psqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *UseRole) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.user_role SET user_id = :user_id, role_id = :role_id, updated_at = :updated_at WHERE id = :id `
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
func (s *psql) delete(id string) error {
	const psqlDelete = `DELETE FROM auth.user_role WHERE id = :id `
	m := UseRole{ID: id}
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
func (s *psql) getByID(id string) (*UseRole, error) {
	const psqlGetByID = `SELECT id , user_id, role_id, created_at, updated_at FROM auth.user_role WHERE id = $1 `
	mdl := UseRole{}
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
func (s *psql) getAll() ([]*UseRole, error) {
	var ms []*UseRole
	const psqlGetAll = ` SELECT id , user_id, role_id, created_at, updated_at FROM auth.user_role `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) updateRoleByUserid(m *UseRole) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.user_role SET role_id = :role_id, updated_at = :updated_at WHERE role_id = :user_id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *psql) getByUseID(id string) (*UseRole, error) {
	const psqlGetByUserID = `SELECT id , user_id, role_id, created_at, updated_at FROM auth.user_role WHERE user_id = $1 limit 1`
	mdl := UseRole{}
	err := s.DB.Get(&mdl, psqlGetByUserID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// Update actualiza un registro en la BD
func (s *psql) updateRole(m *UseRole) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.user_role SET role_id = :role_id, updated_at = :updated_at WHERE user_id = :user_id; `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
