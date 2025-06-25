package filelisting

import (
	"fmt"
	"net/http"
	"os"
)

type AppHandler func(writer http.ResponseWriter, request *http.Request) error

func ErrorWrapper(handler AppHandler) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		err := handler(writer, request)
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err.Error())
			if userErr, ok := err.(UserError); ok {
				http.Error(writer, userErr.Message(), http.StatusBadRequest)
				return
			}
			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

type UserError interface {
	error
	Message() string
}
