package users_rol

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

func newUsersRolPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *UsersRol) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.users_rol (id ,user_id, role_id, created_at, updated_at) VALUES (:id ,:user_id, :role_id,:created_at, :updated_at) `
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
func (s *psql) update(m *UsersRol) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.users_rol SET user_id = :user_id, role_id = :role_id, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.users_rol WHERE id = :id `
	m := UsersRol{ID: id}
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
func (s *psql) getByID(id string) (*UsersRol, error) {
	const psqlGetByID = `SELECT id , user_id, role_id, created_at, updated_at FROM auth.users_rol WHERE id = $1 `
	mdl := UsersRol{}
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
func (s *psql) getAll() ([]*UsersRol, error) {
	var ms []*UsersRol
	const psqlGetAll = ` SELECT id , user_id, role_id, created_at, updated_at FROM auth.users_rol `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByUserID(userID string) (*UsersRol, error) {
	const psqlGetByUserID = `SELECT id , user_id, role_id, created_at, updated_at FROM auth.users_rol WHERE user_id = $1 `
	mdl := UsersRol{}
	err := s.DB.Get(&mdl, psqlGetByUserID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
