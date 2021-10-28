package controller

import (
	"encoding/json"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type Blog struct {
	ID      int32   `db:"id"`
	UserID  int32   `db:"userid"`
	Title   string  `db:"title"`
	Message string  `db:"message"`
	Cat     []int32 `db:"cat"`
}

type Postcat struct {
	ID     int32 `db:"id"`
	BlogID int32 `db:"blogid"`
	CatID  int32 `db:"catid"`
}

func Createblog(res http.ResponseWriter, req *http.Request) {
	var blog Blog
	c, err := req.Cookie("token")
	if err != nil {
		res.Write([]byte("unauthorized : " + err.Error() + "\n"))
	} else {
		value := c.Value
		//res.Write([]byte("cookie has : " + value + "\n"))
		err := json.NewDecoder(req.Body).Decode(&blog)
		if err != nil {
			log.Println("error while req body data  : ", err.Error())
		}
		//fmt.Fprintf(res, "User: %+v", blog)

		tokenString := value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("secretkey"), nil
		})
		if err != nil {
			log.Println("error while req body data  : ", err.Error())
		}
		if !token.Valid {
			res.WriteHeader(http.StatusUnauthorized)
		}
		blog.UserID = int32(claims["id"].(float64))
		//fmt.Fprintf(res, "id-------------------  :  %+v", claims["id"])

		query := `
		INSERT INTO blogs(
			userid,
			title,
			message
		)
		VALUES(
			:userid,
			:title,
			:message
		)
		RETURNING id
	`
		var id int32
		stmt, err := DB.PrepareNamed(query)
		if err != nil {
			log.Println("db error: failed prepare ", err.Error())
			return
		}

		if err := stmt.Get(&id, blog); err != nil {
			log.Println("db error: failed to insert data ", err.Error())
			return
		}

		for i := 0; i < len(blog.Cat); i++ {
			var postcat Postcat
			postcat.BlogID = blog.UserID
			postcat.CatID = blog.Cat[i]
			quer := `
			INSERT INTO postcategory(
				blogid,
				catid
	
			)
			VALUES(
				:blogid,
				:catid
			)
			RETURNING id
		`
			stm, err := DB.PrepareNamed(quer)
			if err != nil {
				log.Println("db error: failed prepare ", err.Error())
				return
			}

			if err := stm.Get(&id, postcat); err != nil {
				log.Println("db error: failed to insert data ", err.Error())
				return
			}
		}

	}
}
