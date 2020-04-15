package handler

import (
	"io/ioutil"
	"net/http"
	dblayer "zmd_package/db"
	"zmd_package/util"
)

const (
	pwd_salt = "#*890"
)

//处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username)<3||len(passwd)<5{
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd+pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc{
		w.Write([]byte("SUCCESS"))
	}else {
		w.Write([]byte("FAILED"))
	}
}
