package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Frontend struct {
	Model
	EnvironmentId *int
	Token         string
	ZId           string
	PublicName    *string
	UrlTemplate   *string
	Reserved      bool
}

func (str *Store) CreateFrontend(envId int, f *Frontend, tx *sqlx.Tx) (int, error) {
	stmt, err := tx.Prepare("insert into frontends (environment_id, token, z_id, public_name, url_template, reserved) values ($1, $2, $3, $4, $5, $6) returning id")
	if err != nil {
		return 0, errors.Wrap(err, "error preparing frontends insert statement")
	}
	var id int
	if err := stmt.QueryRow(envId, f.Token, f.ZId, f.PublicName, f.UrlTemplate, f.Reserved).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "error executing frontends insert statement")
	}
	return id, nil
}

func (str *Store) CreateGlobalFrontend(f *Frontend, tx *sqlx.Tx) (int, error) {
	stmt, err := tx.Prepare("insert into frontends (token, z_id, public_name, url_template, reserved) values ($1, $2, $3, $4, $5) returning id")
	if err != nil {
		return 0, errors.Wrap(err, "error preparing global frontends insert statement")
	}
	var id int
	if err := stmt.QueryRow(f.Token, f.ZId, f.PublicName, f.UrlTemplate, f.Reserved).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "error executing global frontends insert statement")
	}
	return id, nil
}

func (str *Store) GetFrontend(id int, tx *sqlx.Tx) (*Frontend, error) {
	i := &Frontend{}
	if err := tx.QueryRowx("select * from frontends where id = $1", id).StructScan(i); err != nil {
		return nil, errors.Wrap(err, "error selecting frontend by id")
	}
	return i, nil
}

func (str *Store) FindFrontendWithToken(token string, tx *sqlx.Tx) (*Frontend, error) {
	i := &Frontend{}
	if err := tx.QueryRowx("select frontends.* from frontends where token = $1", token).StructScan(i); err != nil {
		return nil, errors.Wrap(err, "error selecting frontend by name")
	}
	return i, nil
}

func (str *Store) FindFrontendWithZId(zId string, tx *sqlx.Tx) (*Frontend, error) {
	i := &Frontend{}
	if err := tx.QueryRowx("select frontends.* from frontends where z_id = $1", zId).StructScan(i); err != nil {
		return nil, errors.Wrap(err, "error selecting frontend by ziti id")
	}
	return i, nil
}

func (str *Store) FindFrontendPubliclyNamed(publicName string, tx *sqlx.Tx) (*Frontend, error) {
	i := &Frontend{}
	if err := tx.QueryRowx("select frontends.* from frontends where public_name = $1", publicName).StructScan(i); err != nil {
		return nil, errors.Wrap(err, "error selecting frontend by public_name")
	}
	return i, nil
}

func (str *Store) FindFrontendsForEnvironment(envId int, tx *sqlx.Tx) ([]*Frontend, error) {
	rows, err := tx.Queryx("select frontends.* from frontends where environment_id = $1", envId)
	if err != nil {
		return nil, errors.Wrap(err, "error selecting frontends by environment_id")
	}
	var is []*Frontend
	for rows.Next() {
		i := &Frontend{}
		if err := rows.StructScan(i); err != nil {
			return nil, errors.Wrap(err, "error scanning frontend")
		}
		is = append(is, i)
	}
	return is, nil
}

func (str *Store) FindPublicFrontends(tx *sqlx.Tx) ([]*Frontend, error) {
	rows, err := tx.Queryx("select frontends.* from frontends where environment_id is null and reserved = true")
	if err != nil {
		return nil, errors.Wrap(err, "error selecting public frontends")
	}
	var frontends []*Frontend
	for rows.Next() {
		frontend := &Frontend{}
		if err := rows.StructScan(frontend); err != nil {
			return nil, errors.Wrap(err, "error scanning frontend")
		}
		frontends = append(frontends, frontend)
	}
	return frontends, nil
}

func (str *Store) UpdateFrontend(fe *Frontend, tx *sqlx.Tx) error {
	sql := "update frontends set environment_id = $1, token = $2, z_id = $3, public_name = $4, url_template = $5, reserved = $6, updated_at = current_timestamp where id = $7"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "error preparing frontends update statement")
	}
	_, err = stmt.Exec(fe.EnvironmentId, fe.Token, fe.ZId, fe.PublicName, fe.UrlTemplate, fe.Reserved, fe.Id)
	if err != nil {
		return errors.Wrap(err, "error executing frontends update statement")
	}
	return nil
}

func (str *Store) DeleteFrontend(id int, tx *sqlx.Tx) error {

	stmt, err := tx.Prepare("delete from frontends where id = $1")
	if err != nil {
		return errors.Wrap(err, "error preparing frontends delete statement")
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return errors.Wrap(err, "error executing frontends delete statement")
	}
	return nil
}
