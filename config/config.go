package BaseConfig

import (
	"errors"
	"net/http"
	"strconv"
)

var DBport = "5432"
var DBuser = "dbadmin"
var DBpwd = "pwd123SQL"
var DBhost = "localhost"
var DBname = "default_db"

var ServerPort = ":5000"
var BasePath = "/api"

var BaseForumPath = "/forum"
var BasePostPath = "/post"
var BaseThreadPath = "/thread"
var BaseServicePath = "/service"
var BaseUserPath = "/user"

var Headers = "application/json; charset=utf-8"

var ConflictError = errors.New(strconv.Itoa(http.StatusConflict))
var NotFoundError = errors.New(strconv.Itoa(http.StatusNotFound))
var ErrBadRequestError = errors.New(strconv.Itoa(http.StatusBadRequest))
var ServerError = errors.New(strconv.Itoa(http.StatusInternalServerError))

func GetErrorCode(Err error) int {
	code, err := strconv.Atoi(Err.Error())
	if err != nil {
		code = http.StatusInternalServerError
	}
	return code
}
