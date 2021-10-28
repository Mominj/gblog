package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Category struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=momin1234 dbname=gblog sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error occur when database conneting", DBErr)
	}
}

func CategoryCreate(res http.ResponseWriter, req *http.Request) {
	var category Category

	err := json.NewDecoder(req.Body).Decode(&category)
	if err != nil {
		log.Println("error while req body data  : ", err.Error())
	}

	query := `
		INSERT INTO categories(
			name
		)
		VALUES(
			:name
		)
		RETURNING id
	`
	var id int32
	stmt, err := DB.PrepareNamed(query)
	if err != nil {
		log.Println("db error: failed prepare wh", err.Error())
		return
	}

	if err := stmt.Get(&id, category); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}

}
