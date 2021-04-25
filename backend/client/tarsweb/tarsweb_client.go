package tarsweb

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	jsoniter "github.com/json-iterator/go"

	"github.com/franklee0817/t3k/backend/config"
	"github.com/franklee0817/t3k/backend/constants"
	"github.com/franklee0817/t3k/backend/constants/errors"
	"github.com/franklee0817/t3k/backend/data"
)

type servCache struct {
	data map[string]interface{}
	mux  *sync.RWMutex
}

var sc *servCache

func init() {
	if sc == nil {
		sc = &servCache{
			data: make(map[string]interface{}),
			mux:  new(sync.RWMutex),
		}
	}
}

func getWebHostAndTicket() (webHost, ticket string) {
	webHost = config.GetKV("web_host")
	ticket = config.GetKV("web_token")

	return
}

func httpGet(ctx context.Context, url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return byt, nil
}

func httpPost(ctx context.Context, url string, body []byte) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second,
	}
	r := bytes.NewBuffer(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return byt, nil
}

func IsTaskFinished(ctx context.Context, taskNo string) (bool, error) {
	url, ticket := getWebHostAndTicket()
	url = fmt.Sprintf("%s/api/task?task_no=%s&ticket=%s", url, taskNo, ticket)
	byt, err := httpGet(ctx, url)
	tars.GetLogger("").Debugf("request tars web /api/task?task_no=%s to query task, resp:%v", taskNo, string(byt))
	if err != nil {
		return false, err
	}

	var obj data.WebResp
	err = jsoniter.Unmarshal(byt, &obj)
	if err != nil {
		return false, err
	}
	if obj.Code != 200 {
		return false, tars.Errorf(int32(obj.Code), obj.Msg)
	}
	var resp data.WebQueryTaskResp
	err = jsoniter.Unmarshal(obj.Data, &resp)
	if err != nil {
		return false, err
	}
	if resp.Status == constants.TaskStatusFailed {
		return false, tars.Errorf(errors.ErrCodeInternalErr, "restart service failed")
	}

	return resp.Status == constants.TaskStatusSucc, nil
}

// Restart 重启服务
func Restart(ctx context.Context, servID int) (string, error) {
	url, ticket := getWebHostAndTicket()
	url = fmt.Sprintf("%s/api/add_task?ticket=%s", url, ticket)
	task := data.WebAddTaskReq{
		Serial: true,
		Items: []data.WebTask{
			{
				ServerID: servID,
				Command:  "restart",
			},
		},
	}

	body, _ := jsoniter.Marshal(task)
	byt, err := httpPost(ctx, url, body)
	tars.GetLogger("").Debugf("request tars web /api/add_task to restart serv, resp:%v", string(byt))

	if err != nil {
		return "", err
	}

	var obj data.WebResp
	err = jsoniter.Unmarshal(byt, &obj)
	if err != nil {
		return "", err
	}
	if obj.Code != 200 {
		return "", tars.Errorf(int32(obj.Code), obj.Msg)
	}
	taskNo := string(obj.Data)

	return strings.Trim(taskNo, "\""), nil
}

func FindServerInfo(ctx context.Context, treeNodeID string, app, server string) (*data.WebServerInfo, error) {
	cacheKey := fmt.Sprintf("info:%s:%s:%s", treeNodeID, app, server)
	sc.mux.RLock()
	if v, ok := sc.data[cacheKey]; ok {
		if c, ok := v.(*data.WebServerInfo); ok {
			sc.mux.RUnlock()
			return c, nil
		}
	}
	sc.mux.RUnlock()
	url, ticket := getWebHostAndTicket()
	url = fmt.Sprintf("%s/api/server_list?tree_node_id=%s&ticket=%s", url, treeNodeID, ticket)
	byt, err := httpGet(ctx, url)
	tars.GetLogger("").Debugf("request tars web /api/server_list?tree_node_id=%s, resp:%v", treeNodeID, string(byt))
	if err != nil {
		return nil, err
	}

	var obj data.WebResp
	err = jsoniter.Unmarshal(byt, &obj)
	if err != nil {
		return nil, err
	}
	if obj.Code != 200 {
		return nil, tars.Errorf(int32(obj.Code), obj.Msg)
	}

	servs := make([]data.WebServerInfo, 0)
	err = jsoniter.Unmarshal(obj.Data, &servs)
	if err != nil {
		return nil, tars.Errorf(errors.ErrCodeInvalidResp, "receive invalid server info format from tars web:%v", string(obj.Data))
	}

	for _, serv := range servs {
		if serv.App == app && serv.Server == server {
			sc.mux.Lock()
			sc.data[cacheKey] = &serv
			sc.mux.Unlock()
			return &serv, nil
		}
	}

	return nil, tars.Errorf(errors.ErrCodeCommon, "not found")
}

// FindAppNode 请求tarsweb接口获取指定app的信息
func FindAppNode(ctx context.Context, app string) (*data.AppTreeNode, error) {
	cacheKey := fmt.Sprintf("node:%s", app)
	sc.mux.RLock()
	if v, ok := sc.data[cacheKey]; ok {
		if c, ok := v.(*data.AppTreeNode); ok {
			sc.mux.RUnlock()
			return c, nil
		}
	}
	sc.mux.RUnlock()
	url, ticket := getWebHostAndTicket()
	url = fmt.Sprintf("%s/api/tree?type=1&ticket=%s", url, ticket)

	byt, err := httpGet(ctx, url)
	tars.GetLogger("").Debugf("request tars web /api/tree?type=1, resp:%v", string(byt))
	if err != nil {
		return nil, err
	}
	var obj data.WebResp
	err = jsoniter.Unmarshal(byt, &obj)
	if err != nil {
		return nil, err
	}
	if obj.Code != 200 {
		return nil, tars.Errorf(int32(obj.Code), obj.Msg)
	}
	nodes := make([]data.AppTreeNode, 0)
	err = jsoniter.Unmarshal(obj.Data, &nodes)
	if err != nil {
		return nil, tars.Errorf(errors.ErrCodeInvalidResp, "receive invalid data format from tars web:%v", string(obj.Data))
	}
	for _, node := range nodes {
		if node.Name == app {
			sc.mux.Lock()
			sc.data[cacheKey] = &node
			sc.mux.Unlock()
			return &node, nil
		}
	}

	return nil, tars.Errorf(errors.ErrCodeCommon, "not found")
}
