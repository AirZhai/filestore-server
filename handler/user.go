package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

//登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password+pwd_salt))
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked{
		w.Write([]byte("FAILED"))
		return
	}
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes{
		w.Write([]byte("FAILED"))
		return
	}

	//w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location:"http://"+r.Host+"/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix + ts[:8]
}

//查询用户信息
func UserInfoHandler(w http.ResponseWriter, r * http.Request)  {
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//
	//isValidToken := IsTokenValid(username, token)
	//if !isValidToken{
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}

	user, err:=dblayer.GetUserInfo(username)
	if err !=nil{
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func IsTokenValid(username string, token string) bool {
	DBToken, err := dblayer.GetUserToken(username)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	if token == DBToken {
		return true
	}
	return false
}