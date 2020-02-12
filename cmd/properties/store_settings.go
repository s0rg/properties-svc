package main

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"
)

// Setting
type Setting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Bundle
type Bundle struct {
	// ID bundle id, use it to set/unset bundles
	ID int `json:"id"`
	// ParentID holds parent's id or 0 if none
	ParentID int `json:"parent_id"`
	// Name bundle name
	Name string `json:"name"`
	// Tag for bundle, empty string if none
	Tag string `json:"tag"`
}

type SettingStore interface {
	Get(ctx context.Context, period time.Time, bundles []int) ([]Setting, error)
	TagsList(ctx context.Context) ([]string, error)
	SettingsList(ctx context.Context) ([]string, error)
	BundlesList(ctx context.Context) ([]Bundle, error)
	BundlesByID(ctx context.Context, bundles []int) ([]Bundle, error)
	BundlesByTag(ctx context.Context, tag string) ([]Bundle, error)
	BundlesByName(ctx context.Context, names []string) ([]Bundle, error)
}

type storeSetting struct {
	db *sql.DB
}

func NewSettingStore(db *sql.DB) *storeSetting {
	return &storeSetting{db: db}
}

// Get returns list of setting values for given bundles at given date.
func (ss *storeSetting) Get(ctx context.Context, when time.Time, bundles []int) (rv []Setting, err error) {
	var query = `
SELECT
   	s.name,
	v.value
FROM 
    bundles b
LEFT JOIN
    bundles_values bv ON 
	    bv.bundle_id = b.id 
		AND 
		bv.created_at < ? 
		AND 
		(bv.expired_at IS NULL OR bv.expired_at > ?)
JOIN
	settings_values v ON v.id = bv.value_id
JOIN
	settings s ON s.id = v.setting_id
WHERE
	b.id IN (` + intArray(bundles) + `)`

	var (
		rows      *sql.Rows
		name, val string
	)

	if rows, err = ss.db.QueryContext(ctx, query, when, when); err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&name, &val); err != nil {
			return
		}

		rv = append(rv, Setting{Name: name, Value: val})
	}

	return rv, rows.Err()
}

// SettingsList returns list of settings names.
func (ss *storeSetting) SettingsList(ctx context.Context) ([]string, error) {
	const query = `
SELECT
	name
FROM
	settings
ORDER BY id`

	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return readStrings(rows)
}

// TagsList returns list of unique non-empty tags.
func (ss *storeSetting) TagsList(ctx context.Context) ([]string, error) {
	const query = `
SELECT DISTINCT
	tag
FROM
	bundles
WHERE
	tag IS NOT NULL`

	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return readStrings(rows)
}

// BundlesList returns list of bundles.
func (ss *storeSetting) BundlesList(ctx context.Context) ([]Bundle, error) {
	var query = `
SELECT
    id,
	parent_id,
	name,
	tag
FROM
	bundles 
ORDER BY id`

	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return readBundles(rows)
}

// BundlesByTag returns list of bundles by tag.
func (ss *storeSetting) BundlesByTag(ctx context.Context, tag string) ([]Bundle, error) {
	var query = `
SELECT
    id,
	parent_id,
	name,
	tag
FROM
	bundles 
WHERE
	tag = ?`

	rows, err := ss.db.QueryContext(ctx, query, tag)
	if err != nil {
		return nil, err
	}

	return readBundles(rows)
}

// BundlesByID returns list of bundles by id.
func (ss *storeSetting) BundlesByID(ctx context.Context, bundles []int) ([]Bundle, error) {
	var query = `
SELECT
    id,
	parent_id,
	name,
	tag
FROM
	bundles 
WHERE
	id IN (` + intArray(bundles) + `)`

	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return readBundles(rows)
}

// BundlesByName returns list of bundles by names.
func (ss *storeSetting) BundlesByName(ctx context.Context, names []string) ([]Bundle, error) {
	var query = `
SELECT
    id,
	parent_id,
	name,
	tag
FROM
	bundles 
WHERE
	name IN (` + strArray(names) + `)`

	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return readBundles(rows)
}

// sql helpers

// strArray takes string slice, returns string usable as SQL `IN` argument
func strArray(a []string) string {
	var s string

	switch len(a) {
	case 0:
	case 1:
		s = a[0]
	default:
		s = strings.Join(a, "','")
	}

	return "'" + s + "'"
}

// intArray takes int slice, returns string usable as SQL `IN` argument
func intArray(a []int) string {
	if len(a) == 0 {
		return "0"
	}

	s := make([]string, len(a))

	for i := 0; i < len(a); i++ {
		s[i] = strconv.Itoa(a[i])
	}

	return strings.Join(s, ",")
}

// readInts reads int slice from rows, and closes them
func readInts(rows *sql.Rows) (rv []int, err error) {
	defer rows.Close()

	var val int

	for rows.Next() {
		if err = rows.Scan(&val); err != nil {
			return nil, err
		}

		rv = append(rv, val)
	}

	return rv, rows.Err()
}

// readStrings reads string slice from rows, and closes them
func readStrings(rows *sql.Rows) (rv []string, err error) {
	defer rows.Close()

	var val string

	for rows.Next() {
		if err = rows.Scan(&val); err != nil {
			return nil, err
		}

		rv = append(rv, val)
	}

	return rv, rows.Err()
}

// readBundles reads Bundle slice from rows, and closes them
func readBundles(rows *sql.Rows) (bundles []Bundle, err error) {
	defer rows.Close()

	var (
		b  Bundle
		ns sql.NullString
	)

	for rows.Next() {
		if err = rows.Scan(&b.ID, &b.ParentID, &b.Name, &ns); err != nil {
			return
		}

		if ns.Valid {
			b.Tag = ns.String
		}

		bundles = append(bundles, b)
	}

	return bundles, rows.Err()
}
