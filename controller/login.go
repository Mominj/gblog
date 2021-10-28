package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
type Cookie struct {
	Name   string
	Value  string
	MaxAge int
}

var mySigninKey = []byte("secretkey")

func GenerateJWT(userID int32) (string, error) {
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 100).Unix()

	tokenString, err := token.SignedString(mySigninKey)
	if err != nil {
		fmt.Printf("something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil

}
func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=momin1234 dbname=new sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr)
	}
}

func (l LoginForm) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email,
			validation.Required.Error("Email is required"),
		),
		validation.Field(&l.Password,
			validation.Required.Error("Password is required"),
			validation.Length(6, 16).Error("Password must be 6 to 16 characters length"),
		),
	)
}
func Login(res http.ResponseWriter, req *http.Request) {

	var loginForm LoginForm

	err := json.NewDecoder(req.Body).Decode(&loginForm)
	if err != nil {
		log.Println("error occur while login req body data  : ", err.Error())
	}
	if vErr := loginForm.Validate(); vErr != nil {
		errorjson, err := json.Marshal(vErr)
		if err != nil {
			log.Println("error occur when convert in json ", err.Error())
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write((errorjson))
		return
	}

	var user User

	query := `SELECT  id, email, password FROM users where email = $1`
	if err := DB.Get(&user, query, loginForm.Email); err != nil {
		if err != nil {
			log.Println("error while getting user from db: ", err.Error())
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		log.Println("password do not match: ", err.Error())
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		fmt.Printf("error token")
	}
	//json.NewEncoder(res).Encode(token)

	cookie := http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 300,
	}
	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Doc Get Successful"))
}

func ReadCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("momin")
	if err != nil {
		w.Write([]byte("error in reading cookie : " + err.Error() + "\n"))
	} else {
		value := c.Value
		w.Write([]byte("cookie has : " + value + "\n"))
	}
}
