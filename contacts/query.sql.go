// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: query.sql

package contacts

import (
	"context"
)

const add = `-- name: Add :exec
insert into contacts (fname, lname, phone, email, birthday, address, city, state, zipcode, notes) values
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type AddParams struct {
	Fname    string
	Lname    string
	Phone    string
	Email    string
	Birthday string
	Address  string
	City     string
	State    string
	Zipcode  string
	Notes    string
}

func (q *Queries) Add(ctx context.Context, arg AddParams) error {
	_, err := q.db.ExecContext(ctx, add,
		arg.Fname,
		arg.Lname,
		arg.Phone,
		arg.Email,
		arg.Birthday,
		arg.Address,
		arg.City,
		arg.State,
		arg.Zipcode,
		arg.Notes,
	)
	return err
}

const delete = `-- name: Delete :exec
delete from contacts where id = ?
`

func (q *Queries) Delete(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, delete, id)
	return err
}

const list = `-- name: List :many
SELECT id, fname, lname, phone, email, birthday, address, city, state, zipcode, notes FROM contacts ORDER BY lname
`

func (q *Queries) List(ctx context.Context) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, list)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.Fname,
			&i.Lname,
			&i.Phone,
			&i.Email,
			&i.Birthday,
			&i.Address,
			&i.City,
			&i.State,
			&i.Zipcode,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const update = `-- name: Update :exec
UPDATE contacts set fname=?, lname=?, phone=?, email=?, birthday=?, address=?, city=?, state=?, zipcode=?,notes=?
WHERE id = ?
`

type UpdateParams struct {
	Fname    string
	Lname    string
	Phone    string
	Email    string
	Birthday string
	Address  string
	City     string
	State    string
	Zipcode  string
	Notes    string
	ID       int64
}

func (q *Queries) Update(ctx context.Context, arg UpdateParams) error {
	_, err := q.db.ExecContext(ctx, update,
		arg.Fname,
		arg.Lname,
		arg.Phone,
		arg.Email,
		arg.Birthday,
		arg.Address,
		arg.City,
		arg.State,
		arg.Zipcode,
		arg.Notes,
		arg.ID,
	)
	return err
}