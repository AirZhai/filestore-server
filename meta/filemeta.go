package meta

import (
	mydb "zmd_package/db/mysql"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}

//新增或者更新文件元信息
func UpdateFileMeta(fmeta FileMeta)  {
	fileMetas[fmeta.FileSha1] = fmeta
}

//新增/更新文件
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.Location)
}

//通过sha1获取文件的元信息对象
func GetFileMeta(fileSha1 string)FileMeta{
	return fileMetas[fileSha1]
}

//删除
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}
//从mysql获取文件信息
func GetFileMetaDB(fileSha1 string) (FileMeta,error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err !=nil{
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta,nil
}