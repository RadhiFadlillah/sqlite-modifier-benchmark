package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func openDatabase() (db *sqlx.DB, err error) {
	// Connect to database
	db, err = sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(-1)

	// Create transaction
	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}

	// Make sure to rollback if panic ever happened
	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			logrus.Errorln("Database error:", panicErr)
			tx.Rollback()

			db = nil
			err = panicErr
		}
	}()

	// Generate tables
	tx.MustExec(`CREATE TABLE IF NOT EXISTS purchase (
		id         INTEGER NOT NULL,
		qty        INTEGER NOT NULL,
		total      INTEGER NOT NULL,
		input_time TEXT    NOT NULL,
		CONSTRAINT purchase_PK PRIMARY KEY(id))`)

	// Clear table
	tx.MustExec(`DELETE FROM purchase`)

	// Insert random value
	rand.Seed(time.Now().Unix())

	stmtInsert, err := tx.Preparex(
		`INSERT INTO purchase(qty, total, input_time) 
		VALUES (?, ?, ?)`)
	checkError(err)

	maxTime := time.Time(today).UTC()
	current := time.Time(epoch).UTC()
	for {
		if current.After(maxTime) {
			break
		}

		// Create random value
		qty := rand.Intn(11) + 1
		price := rand.Intn(330000) + 10000
		timeIncrease := rand.Intn(600) + 60

		// Round price
		flPrice := float64(price)
		flPrice = math.Round(flPrice/1000) * 1000
		price = int(flPrice)

		// Increase current time
		current = current.Add(time.Duration(timeIncrease) * time.Second)

		stmtInsert.MustExec(
			strconv.Itoa(qty),
			strconv.Itoa(qty*price),
			current.Format("2006-01-02 15:04:05"))
	}

	// Commit transaction
	err = tx.Commit()
	checkError(err)

	return db, err
}
