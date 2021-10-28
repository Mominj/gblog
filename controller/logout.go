package controller

import (
	"net/http"
)

func Logout(res http.ResponseWriter, req *http.Request) {
	//json.NewEncoder(res).Encode(token)

	cookie := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}
	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Doc Get Successful"))

}
