package filelisting

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type testingUserError string

func (e testingUserError) Error() string {
	return e.Message()
}

func (e testingUserError) Message() string {
	return string(e)
}

func errPanic(writer http.ResponseWriter, request *http.Request) error {
	panic("PANIC")
	return nil
}

func errUserError(writer http.ResponseWriter, request *http.Request) error {
	return testingUserError("user error")
}

func errNotFound(writer http.ResponseWriter, request *http.Request) error {
	return os.ErrNotExist
}

func ErrPermission(writer http.ResponseWriter, request *http.Request) error {
	return os.ErrPermission
}

func errUnknown(writer http.ResponseWriter, request *http.Request) error {
	return errors.New("unknown error")
}

func noErr(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

var tests = []struct {
	handler AppHandler
	code    int
	message string
}{
	//{errPanic, 500, ""},
	{errUserError, 400, "user error"},
	{errNotFound, 404, "Not Found"},
	{ErrPermission, 403, "Forbidden"},
	{errUnknown, 500, "Internal Server Error"},
	{noErr, 200, ""},
}

func TestErrorWrapper(t *testing.T) {
	for _, tt := range tests {
		f := ErrorWrapper(tt.handler)
		resp := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "https://wwww.imooc.com", nil)
		f(resp, request)
		message := strings.Trim(resp.Body.String(), "\n")
		if resp.Code != tt.code || message != tt.message {
			t.Errorf("expect %d %s; got %d %s", tt.code, tt.message, resp.Code, message)
		}
	}
}

func TestErrWrapperInServer(t *testing.T) {
	for _, tt := range tests {
		f := ErrorWrapper(tt.handler)
		server := httptest.NewServer(http.HandlerFunc(f))
		resp, _ := http.Get(server.URL)

		body, _ := ioutil.ReadAll(resp.Body)
		message := strings.Trim(string(body), "\n")
		if resp.StatusCode != tt.code || message != tt.message {
			t.Errorf("expect %d %s; got %d %s", tt.code, tt.message, resp.StatusCode, message)
		}
	}
}
