package main

import (
	Day3_save_query "Day3-save-query"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	engine, err := Day3_save_query.NewEngine("sqlite", "gee.db")
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()
	s := engine.NewSession()
	if _, err := s.Raw("DROP TABLE IF EXISTS User;").Exec(); err != nil {
		log.Fatal(err)
	}
	if _, err := s.Raw("CREATE TABLE User(Name text);").Exec(); err != nil {
		log.Fatal(err)
	}
	var result sql.Result
	var err1 error
	if result, err1 = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec(); err1 != nil {
		log.Fatal(err1)
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exec success, %d affected\n", count)

}
