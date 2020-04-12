package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"zmd_package/meta"
	"zmd_package/util"
)

func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		//返回上传的html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err !=nil{
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST"{
		//接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("failed to get data, err:%s", err.Error())
			return
		}
		defer file.Close()

		fileMeat := meta.FileMeta{
			FileName: head.Filename,
			Location: "E:/code/"+head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeat.Location)
		if err!= nil{
			fmt.Printf("failed to create data, err:%s", err.Error())
			return
		}
		defer newFile.Close()

		newFile.Seek(0, 0)
		fileMeat.FileSha1 = util.FileSha1(newFile)
		fmt.Println(fileMeat.FileSha1)
		meta.UpdateFileMeta(fileMeat)

		fileMeat.FileSize, err = io.Copy(newFile, file)
		if err!= nil{
			fmt.Printf("failed to save data into file, err:%s", err.Error())
			return
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}
}


//上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "Upload finished!")

}


//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	//解析请求参数
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(filehash)
	data, err :=json.Marshal(fmeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}