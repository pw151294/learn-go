package main

import (
	"fmt"
	"iflytek.com/weipan4/learn-go/errorhandling/filelistingserver/filelisting"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const prefix = "/filelist/"

type userError string

func (e userError) Error() string {
	return e.Message()
}

func (e userError) Message() string {
	return string(e)
}

func handlerFileList(w http.ResponseWriter, r *http.Request) error {
	if strings.Index(r.URL.Path, prefix) != 0 {
		//return errors.New("path must start with /list/")
		return userError(fmt.Sprintf("path must start with %s", prefix))
	}
	path := r.URL.Path[len(prefix):]
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	w.Write(all)
	return nil
}

func main() {
	// HandleFunc(pattern string, wshandler func(ResponseWriter, *Request)
	// 将handler参数理解为callback回调接口 在请求发挥并get到响应之后触发
	// ErrorWrapper是在handler的基础上 以handler为参数 又封装了一层新的回调接口
	http.HandleFunc("/", filelisting.ErrorWrapper(handlerFileList))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
