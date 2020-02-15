package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/fxamacker/cbor"
)

// UserSettings holds two arrays: bundle and settings ids, and time,
// when settings must be treated as expired (can be nil for long-time sets).
type UserSettings struct {
	Bundles []int      `json:"bundles"`
	Expire  *time.Time `json:"expire,omitempty"`
}

type UserStore interface {
	Get(ctx context.Context, userID int, when time.Time) (s UserSettings, err error)
	Set(ctx context.Context, userID int, s UserSettings) error
}

type storeUser struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return &storeUser{db: db}
}

// Get returns UserSettings for given user and time.
func (su *storeUser) Get(ctx context.Context, userID int, when time.Time) (s UserSettings, err error) {
	const query = `
	SELECT
		settings,
        expires_at
	FROM
		user_settings
	WHERE
		user_id = ?
		AND
		created_at <= ?
		AND
		(expires_at IS NULL OR expires_at > ?)
	ORDER BY
		created_at DESC
	LIMIT 1
	`

	var (
		buf []byte
		snt sql.NullTime
	)

	err = su.db.QueryRowContext(ctx, query, userID, when, when).Scan(&buf, &snt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil // its OK to return empty, if none found.
		}

		return
	}

	if snt.Valid {
		s.Expire = &snt.Time
	}

	err = cbor.Unmarshal(buf, &s.Bundles)

	return s, err
}

// Set sets new settings for user.
func (su *storeUser) Set(ctx context.Context, userID int, s UserSettings) (err error) {
	const query = `
INSERT INTO user_settings
	(user_id, settings, expires_at)
VALUES
	(?, ?, ?)`

	var buf []byte

	buf, err = cbor.Marshal(s.Bundles, cbor.EncOptions{Canonical: true})
	if err != nil {
		return
	}

	log.Println("cbor size:", len(buf))

	_, err = su.db.ExecContext(ctx, query, userID, buf, s.Expire)

	return
}
