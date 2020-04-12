package main

import (
	"fmt"
	"net/http"
	"zmd_package/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta/", handler.GetFileMetaHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil{
		fmt.Printf("failed to start server, err:%s", err.Error())
		return
	}
}
