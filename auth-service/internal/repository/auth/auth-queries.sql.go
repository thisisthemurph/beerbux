// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: auth-queries.sql

package auth

import (
	"context"
)

const getUserByID = `-- name: GetUserByID :one
select id, username, hashed_password, created_at, updated_at from users where id = ?
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
select id, username, hashed_password, created_at, updated_at from users where username = ?
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const registerUser = `-- name: RegisterUser :one
insert into users (id, username, hashed_password)
values (?, ?, ?)
returning id, username, hashed_password, created_at, updated_at
`

type RegisterUserParams struct {
	ID             string
	Username       string
	HashedPassword string
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, registerUser, arg.ID, arg.Username, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
update users
set hashed_password = ?,
    updated_at = current_timestamp
where id = ?
returning id, username, hashed_password, created_at, updated_at
`

type UpdateUserPasswordParams struct {
	HashedPassword string
	ID             string
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.HashedPassword, arg.ID)
	return err
}

const updateUserUsername = `-- name: UpdateUserUsername :one
update users
set username = ?,
    updated_at = current_timestamp
where id = ?
returning id, username, hashed_password, created_at, updated_at
`

type UpdateUserUsernameParams struct {
	Username string
	ID       string
}

func (q *Queries) UpdateUserUsername(ctx context.Context, arg UpdateUserUsernameParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserUsername, arg.Username, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userWithUsernameExists = `-- name: UserWithUsernameExists :one
select exists(select 1 from users where username = ?)
`

func (q *Queries) UserWithUsernameExists(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, userWithUsernameExists, username)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}
