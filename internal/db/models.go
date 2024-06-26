// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Note struct {
	ID                pgtype.Int8
	Content           string
	CreatedAt         pgtype.Timestamptz
	ExpiresAt         pgtype.Timestamp
	ExpiresAtTimezone string
	KeyHash           []byte
}
