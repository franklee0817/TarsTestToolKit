package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/TarsCloud/TarsGo/tars"
	jsoniter "github.com/json-iterator/go"

	"github.com/franklee0817/t3k/backend/constants/errors"
	"github.com/franklee0817/t3k/backend/tars-protocol/apitars"
)

func readIntFromQuery(req *http.Request, name string, required bool) (int, error) {
	val := req.URL.Query().Get(name)
	val = strings.Trim(val, "\"")
	if val == "" && !required {
		return 0, nil
	}
	if val == "" && required {
		return 0, tars.Errorf(errors.ErrCodeParam, "invalid param %s=%v", name, val)
	}
	intVal, err := strconv.Atoi(val)
	if err != nil && required {
		return 0, tars.Errorf(errors.ErrCodeParam, "param %s is required as a integer", name)
	}

	return intVal, nil
}

// HandleGetTestDetail 查询测试明细
func HandleGetTestDetail(w http.ResponseWriter, r *http.Request) {
	initHeader(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	testID, err := readIntFromQuery(r, "test_id", true)
	if err != nil {
		tars.GetLogger("").Error(err)
		w.Write(ParseErr(err))
		return
	}

	ts, err := readIntFromQuery(r, "timestamp", false)
	if err != nil {
		tars.GetLogger("").Error(err)
		w.Write(ParseErr(err))
		return
	}

	rsp, err := app.GetTestDetail(context.Background(), uint32(testID), uint32(ts))
	if err != nil {
		w.Write(ParseErr(err))
		return
	}

	rspByt, _ := jsoniter.Marshal(rsp)
	w.Write(rspByt)
}

// HandleGetTestHistories 查询测试历史
func HandleGetTestHistories(w http.ResponseWriter, r *http.Request) {
	initHeader(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	pageStr := r.URL.Query().Get("page")
	pageStr = strings.Trim(pageStr, "\"")
	var page, pageSize int
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			w.Write(ParseErr(err))
			return
		}
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	pageSizeStr = strings.Trim(pageSizeStr, "\"")
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			w.Write(ParseErr(err))
			return
		}
	}

	req := new(apitars.QueryTestHistoryReq)
	req.Page = uint32(page)
	req.PageSize = uint32(pageSize)
	rsp, err := app.GetTestHistories(context.Background(), req)
	if err != nil {
		w.Write(ParseErr(err))
		return
	}

	rspByt, _ := jsoniter.Marshal(rsp)
	w.Write(rspByt)
}
