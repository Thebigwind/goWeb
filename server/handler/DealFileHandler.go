package handler

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"

	. "github.com/bigwind/goWeb/common"
)

/*
 处理upload 逻辑
表单中增加enctype="multipart/form-data"
服务端调用r.ParseMultipartForm,把上传的文件存储在内存和临时文件中
使用r.FormFile获取文件句柄，然后对文件进行存储等处理。
*/

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func DownLoadFileHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	locker := GetGlobeLocker()
	locker.RLock()
	defer locker.RUnlock()
	err := req.ParseForm()
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	/*
		fpath := req.Form.Get("fpath")
		if fpath == "" {
		}
		fpath = path.Join(UPLOAD_DIR, fpath)
	*/
	fpath := "/root/whisper_data_result.txt"
	err = CheckReadFile(fpath)
	if err != nil {
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	file, err := os.Open(fpath)
	if err != nil {
		w.Header().Set("Content-Type", " text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	fileName := path.Base(fpath)
	//	fileName = url.QueryEscape(fileName) //for chinese characters
	//	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	fstat, _ := file.Stat()
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fstat.Size()))
	io.Copy(w, file)
}
