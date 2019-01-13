package egoxml

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	open()
}

func open() {
	if db != nil {
		return
	}
	database, err := sql.Open("sqlite3", "./savegame.db")
	if err != nil {
		fmt.Print(err)
	}
	db = database
	db.Ping()
}

func create() {
	if db == nil {
		open()
	}
	_, err := db.Exec(`CREATE TABLE log(Id INTEGER NOT NULL PRIMARY KEY,Time NUMERIC(18,15) NOT NULL, Title VARCHAR(23) NOT NULL,AttrText VARCHAR(181) NOT NULL,Faction VARCHAR(11),Money NUMERIC(10,2) NOT NULL,Ship VARCHAR(7));`)
	if err != nil && err.Error() != "table log already exists" {
		fmt.Println(err.Error())
	}
}

func bulkInsert(lastEntry float64, entries []Entry) error {
	inserts := 0
	query := "INSERT INTO log(Time,Title,AttrText,Faction,Money,Ship) VALUES "
	values := make([]interface{}, 0)

	batchQuery := query
	for index, entry := range entries {
		if entry.Time > lastEntry {
			if inserts != 0 {
				batchQuery += ", "
			}
			batchQuery += "(?,?,?,?,?,?)"
			values = append(values, entry.Time, entry.Title, entry.AttrText, entry.Faction, entry.Money, entry.Ship)
			inserts++
			if inserts > 150 || index == len(entries)-1 {
				batchQuery += ";"
				err := transactionalInsert(batchQuery, values)
				if err != nil {
					return err
				}
				inserts = 0
				values = make([]interface{}, 0)
				batchQuery = query
			}
		}
	}
	return nil
}

func transactionalInsert(query string, values []interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Could not roll back: %v\n", rollbackErr)
		}
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Could not roll back: %v\n", rollbackErr)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Could not roll back: %v\n", rollbackErr)
		}
		return err
	}
	return nil
}

func fetchLastEntry() float64 {
	rows, err := db.Query("SELECT Time FROM log ORDER BY Time DESC LIMIT 1;")
	rows.Next()
	var lastEntry float64
	err = rows.Scan(&lastEntry)
	if err != nil {
		lastEntry = 0
	}
	rows.Close()
	return lastEntry
}

func fetchLogs() []Entry {
	rows, _ := db.Query("SELECT Time, Title, AttrText, Faction, Money, Ship FROM log WHERE Ship != \"\" ORDER BY Time;")
	defer rows.Close()
	var entries []Entry

	for rows.Next() {
		entry := Entry{}
		if err := rows.Scan(&entry.Time, &entry.Title, &entry.AttrText, &entry.Faction, &entry.Money, &entry.Ship); err != nil {
			log.Fatal(err)
		} else {
			entries = append(entries, entry)
		}
	}

	return entries
}
