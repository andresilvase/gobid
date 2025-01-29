// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package pgstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users ("user_name", "email", "password_hash", "bio")
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateUserParams struct {
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Bio          string `json:"bio"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UserName,
		arg.Email,
		arg.PasswordHash,
		arg.Bio,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUserById = `-- name: GetUserById :one
SELECT
  id,
  user_name,
  password_hash,
  email,
  bio,
  created_at,
  updated_at
FROM
  users
WHERE
  id = $1
`

type GetUserByIdRow struct {
	ID           uuid.UUID          `json:"id"`
	UserName     string             `json:"user_name"`
	PasswordHash []byte             `json:"password_hash"`
	Email        string             `json:"email"`
	Bio          string             `json:"bio"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.UserName,
		&i.PasswordHash,
		&i.Email,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
