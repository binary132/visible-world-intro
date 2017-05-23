package main

import (
	"fmt"
	models "models"
	"database/sql"
	"gopkg.in/nullbio/null.v6"
	_ "github.com/go-sql-driver/mysql"
)

//go:generate sqlboiler --basedir ./ --tinyint-as-bool --no-tests mysql 

func main() {
	fmt.Println("sqlboiler")
	db, err := sql.Open("mysql", "root:root@/library")
	if err != nil {
		fmt.Println("err", err)
	}
	book := &models.Book{Name: null.StringFrom("Harry Potter and the Prisoner of Azkaban"), Author: null.StringFrom("J. K. Rowling") }
	if err = book.Insert(db); err != nil {
		fmt.Println("failed to insert book", "err", err)
	}
}
