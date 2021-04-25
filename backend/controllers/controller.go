package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/TarsCloud/TarsGo/tars"
	jsoniter "github.com/json-iterator/go"

	"github.com/franklee0817/t3k/backend/constants/errors"
	"github.com/franklee0817/t3k/backend/impl"
)

var app *impl.APIImpl

func init() {
	app = new(impl.APIImpl)
	err := app.Init()
	if err != nil {
		fmt.Printf("apiImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
}

// MarshalErr 构建error结构byte slice
func MarshalErr(errCode int32, msg string) []byte {
	ret := map[string]interface{}{
		"code": errCode,
		"msg":  msg,
	}
	byt, _ := jsoniter.Marshal(ret)

	return byt
}

// ParseErr 将错误转换成[]byte结构用于返回
func ParseErr(err error) []byte {
	if err == nil {
		return MarshalErr(0, "succ")
	}
	code := tars.GetErrorCode(err)
	return MarshalErr(code, err.Error())
}

// UnpackReq 解析http请求body
func UnpackReq(r *http.Request, dest interface{}) ([]byte, error) {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		err := tars.Errorf(errors.ErrCodeInternalErr, "internal error, illegal body receiver")
		tars.GetLogger("").Error(err)
		return ParseErr(err), err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = tars.Errorf(errors.ErrCodeInvalidReq, "read request body failed:%s", err)
		tars.GetLogger("").Error(err)
		return ParseErr(err), err
	}

	err = jsoniter.Unmarshal(body, dest)
	if err != nil {
		err = tars.Errorf(errors.ErrCodeInvalidReq, "invalid request body:%s", err)
		tars.GetLogger("").Error(err)
		return ParseErr(err), err
	}

	return nil, nil
}

func initHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("Content-Type", "application/json")             //返回数据格式是json
}
