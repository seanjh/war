// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Game struct {
	ID      int64
	Code    string
	Created string
}

type GameSession struct {
	GameID    int64
	SessionID int64
	Role      sql.NullString
	Created   interface{}
}
