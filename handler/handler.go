package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"zmd_package/meta"
	"zmd_package/util"
	dblayer "zmd_package/db"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传的html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data, err:%s", err.Error())
			return
		}
		defer file.Close()

		fileMeat := meta.FileMeta{
			FileName: head.Filename,
			Location: "E:/code/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeat.Location)
		if err != nil {
			fmt.Printf("failed to create data, err:%s", err.Error())
			return
		}
		defer newFile.Close()

		newFile.Seek(0, 0)
		fileMeat.FileSha1 = util.FileSha1(newFile)
		fmt.Println(fileMeat.FileSha1)
		//meta.UpdateFileMeta(fileMeat)
		meta.UpdateFileMetaDB(fileMeat)

		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeat.FileSha1, fileMeat.FileName, fileMeat.FileSize)
		if suc{
			http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
		}else {
			w.Write([]byte("Upload Failed."))
		}

		//fileMeat.FileSize, err = io.Copy(newFile, file)
		//if err != nil {
		//	fmt.Printf("failed to save data into file, err:%s", err.Error())
		//	return
		//}
		//
		//http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}
}

//上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")

}

//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	//解析请求参数
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//fmeta := meta.GetFileMeta(filehash)
	fmeta, err := meta.GetFileMetaDB(filehash)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func FileQueryHandler(w http.ResponseWriter, r http.Request)  {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := json.Marshal(userFiles)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	//meta.UpdateFileMeta(curFileMeta)
	meta.UpdateFileMetaDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func FileDeleteHandle(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)

}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	if fileMeta == nil{
		resp:=util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}

	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc{
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
			Data: time.Now(),
		}
		w.Write(resp.JSONBytes())
		return
	}else {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
			Data: time.Now(),
		}
		w.Write(resp.JSONBytes())
		return
	}
}