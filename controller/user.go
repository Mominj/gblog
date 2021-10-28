package controller

import (
	"encoding/json"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sqlx.DB
var DBErr error

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=momin1234 dbname=gblog sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr.Error())
	}
}

type User struct {
	ID       int32  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u User) validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.Required.Error("This Field Is Required"),
			is.Alpha.Error("The field may only contain letters and spaces."),
			validation.Length(3, 20),
		),
		validation.Field(&u.Email,
			validation.Required.Error("This Field Is Required"),
			validation.Length(3, 100),
			is.Email,
		),
		validation.Field(&u.Password,
			validation.Required.Error("This Field Is Required"),
			validation.Length(6, 100),
		),
	)
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	var user User

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Println("error while req body data  : ", err.Error())
	}
	vErr := user.validate()
	if vErr != nil {
		errorjson, err := json.Marshal(vErr)
		if err != nil {
			log.Println("error occur when convert in json ", err.Error())
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write((errorjson))
		return
		//fmt.Fprintf(res, "User: %+v", vErr)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error while encrypted password: ", err.Error())
	}
	user.Password = string(hash)

	query := `
		INSERT INTO users(
			name,
			email,
			password
		)
		VALUES(
			:name,
			:email,
			:password
		)
		RETURNING id
	`
	var id int32
	stmt, err := DB.PrepareNamed(query)
	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if err := stmt.Get(&id, user); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}
}
