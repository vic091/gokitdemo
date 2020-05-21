package Services

import (
	"context"
	"encoding/json"
	"errors"
	mymux "github.com/gorilla/mux"
	"gokitdemo/util"
	"net/http"
	"strconv"
)

func DecodeUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	//if r.URL.Query().Get("uid") != "" {
	//	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
	//	return UserRequest{Uid: uid}, nil
	//}

	vars := mymux.Vars(r)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			Uid:    uid,
			Method: r.Method,
		}, nil
	}

	return nil, errors.New("参数错误")

}

func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// 自定义错误处理
func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {

	contentType, body := "text/plain;charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type", contentType)
	if myerr,ok := err.(*util.MyError);ok{
		w.WriteHeader(myerr.Code)
	} else {
		w.WriteHeader(404)
	}

	w.Write(body)

}
