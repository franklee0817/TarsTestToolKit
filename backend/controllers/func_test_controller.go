package controllers

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// HandleDoFuncTest 启动功能测试
func HandleDoFuncTest(w http.ResponseWriter, r *http.Request) {
	initHeader(w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rsp, err := app.DoFuncTest(context.Background())
	if err != nil {
		w.Write(ParseErr(err))
		return
	}

	rspByt, _ := jsoniter.Marshal(rsp)
	w.Write(rspByt)
}
