package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	envDBUsers    = "APP_DB_USERS"
	envDBSettings = "APP_DB_SETTINGS"
	envAddr       = "APP_ADDR"
	retryAttempts = 3
)

func mustGetEnv(key string) (val string) {
	if val = os.Getenv(key); val == "" {
		log.Fatal("no env value for:", key)
	}
	return
}

func connectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}

func retry(times int, fn func() error) (err error) {
	for i := 0; i < times; i++ {
		if err = fn(); err != nil {
			return
		}
		if i < times {
			time.Sleep(time.Second)
		}
	}
	return
}

func serve(addr, userDSN, settingDSN string) (err error) {
	var (
		uDB *sql.DB
		sDB *sql.DB
	)

	connUsers := func() (err error) {
		uDB, err = connectDB(userDSN)
		return
	}

	closeUsers := func() {
		if err := uDB.Close(); err != nil {
			log.Println("user-db close fail:", err)
		}
	}

	connSettings := func() (err error) {
		sDB, err = connectDB(settingDSN)
		return
	}

	closeSettings := func() {
		if err := sDB.Close(); err != nil {
			log.Println("setting-db close fail:", err)
		}
	}

	if err = retry(retryAttempts, connUsers); err != nil {
		return fmt.Errorf("user-db connect fail: %w", err)
	}

	defer closeUsers()

	if err = retry(retryAttempts, connSettings); err != nil {
		return fmt.Errorf("setting-db connect fail: %w", err)
	}

	defer closeSettings()

	srv := newService(addr, uDB, sDB)

	log.Println("serving at:", addr)

	return srv.Serve()
}

func main() {
	userDSN := mustGetEnv(envDBUsers)
	settingDSN := mustGetEnv(envDBSettings)

	var addr string

	if addr = os.Getenv(envAddr); addr == "" {
		addr = "0.0.0.0:8080"
	}

	if err := serve(addr, userDSN, settingDSN); err != nil {
		log.Fatal(err)
	}
}
