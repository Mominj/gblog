package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Edited(res http.ResponseWriter, req *http.Request) {

	var blog Blog
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		log.Println("error while req body data  : ", err.Error())
	}
	vars := mux.Vars(req)

	blogid, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("error occurs when convert string to int", err.Error())
	}
	blog.ID = int32(blogid)
	//fmt.Fprintf(res, "User: %+v", user)
	query := `
	UPDATE blogs
			SET  title = :title, message = :message
			WHERE blogs.id = :id
	`

	stmt, err := DB.PrepareNamed(query)
	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if _, err := stmt.Exec(&blog); err != nil {
		log.Println("db error: failed to update data ", err.Error())
		return
	}

}
