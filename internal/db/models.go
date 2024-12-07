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
	ID                  int64
	GameID              int64
	SessionID           sql.NullString
	Role                int64
	CardsInHand         string
	CardsInBattle       string
	CardsInBattleHidden string
	CardsInReserve      string
	Created             string
	Updated             string
}

type Session struct {
	ID      string
	Created string
}
