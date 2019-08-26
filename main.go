package main

import (
	"fmt"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	epoch = time.Date(2010, 01, 01, 01, 00, 00, 00, time.Local)
	today = time.Now()
)

func main() {
	// Open database
	db, err := openDatabase()
	checkError(err)
	defer db.Close()

	// Calculate max number of days
	maxNDays := int(math.Round(today.Sub(epoch).Hours() / 24))
	for nDays := 10; nDays < maxNDays; nDays *= 2 {
		// Select using "localtime" modifier
		r1, d1 := localtimeSelect(db, nDays)

		// Select using "hours" modifier
		r2, d2 := hoursSelect(db, nDays)

		if r1 != r2 {
			fmt.Printf("got different number of rows, localtime %d hours %d\n", r1, r2)
			continue
		}

		fmt.Printf("N Days    : %d\n", nDays)
		fmt.Printf("Rows      : %d\n", r1)
		fmt.Printf("Localtime : %.3f s\n", d1)
		fmt.Printf("Hours     : %.3f s\n", d2)

		if d2 > d1 {
			speedDiff := d2 / d1
			fmt.Printf("Localtime is %.2fx faster than Hours\n", speedDiff)
		} else {
			speedDiff := d1 / d2
			fmt.Printf("Hours is %.2fx faster than Localtime\n", speedDiff)
		}

		fmt.Println()
	}
}

func localtimeSelect(db *sqlx.DB, nDays int) (rows int64, duration float64) {
	start := time.Now()
	lastDay := epoch.AddDate(0, 0, nDays)

	err := db.Get(&rows, `SELECT COUNT(*) FROM purchase
		WHERE DATE(input_time, "localtime") >= ?
		AND DATE(input_time, "localtime") <= ?`,
		epoch.Format("2006-01-02"),
		lastDay.Format("2006-01-02"))
	checkError(err)

	finish := time.Now()
	return rows, finish.Sub(start).Seconds()
}

func hoursSelect(db *sqlx.DB, nDays int) (rows int64, duration float64) {
	start := time.Now()
	lastDay := epoch.AddDate(0, 0, nDays)
	sqliteHours := getSqliteHours()

	err := db.Get(&rows, `SELECT COUNT(*) FROM purchase
		WHERE DATE(input_time, ?) >= ?
		AND DATE(input_time, ?) <= ?`,
		sqliteHours, epoch.Format("2006-01-02"),
		sqliteHours, lastDay.Format("2006-01-02"))
	checkError(err)

	finish := time.Now()
	return rows, finish.Sub(start).Seconds()
}

func getSqliteHours() string {
	_, offset := time.Now().Zone()
	offset = offset / 3600

	strFormat := "%d hours"
	if offset > 0 {
		strFormat = "+%d hours"
	}

	return fmt.Sprintf(strFormat, offset)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
