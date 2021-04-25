package controllers

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/franklee0817/t3k/backend/tars-protocol/apitars"
)

// HandleDoPerfTest app.DoPerfTest http handler
func HandleDoPerfTest(w http.ResponseWriter, r *http.Request) {
	initHeader(w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	req := apitars.PerfTestReq{}
	byt, err := UnpackReq(r, &req)
	if err != nil {
		w.Write(byt)
		return
	}

	rsp, err := app.DoPerfTest(context.Background(), &req)
	if err != nil {
		w.Write(ParseErr(err))
		return
	}

	rspByt, _ := jsoniter.Marshal(rsp)
	w.Write(rspByt)
}
