package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	"hello/api"

	"models"

	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/julienschmidt/httprouter"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sqlx.Open("mysql", "root:root@tcp(localhost:3306)/library")
	if err != nil { panic(err.Error()) }

	r := httprouter.New()
	Shelf{db}.Bind(r)

	http.ListenAndServe("localhost:8083", r)

	select{}
}

type Shelf struct { *sqlx.DB }

var _ = api.API(Shelf{})

func (s Shelf) Bind(r *httprouter.Router) error {
	r.GET("/shelves", s.GetAll)
	r.GET("/shelf/:id", s.Get)

	return nil
}

func (s Shelf) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        o, e1 := models.Shelves(s.DB).All()
	if e1 != nil {
		fmt.Println("error", e1)
	}	

	if err := json.NewEncoder(w).Encode(o); err != nil {
		panic(err.Error())
	}
}

func (s Shelf) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	qms := []qm.QueryMod{qm.Where("id = ?", id)}
        o, e1 := models.Shelves(s.DB, qms...).One()
	if e1 != nil {
		fmt.Println("error", e1)
	}	

	if err := json.NewEncoder(w).Encode(o); err != nil {
		panic(err.Error())
	}
}
