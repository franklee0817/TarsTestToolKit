package data

import "encoding/json"

type AppTreeNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	PID  string `json:"pid"`
}

type WebResp struct {
	Code int             `json:"ret_code"`
	Msg  string          `json:"err_msg"`
	Data json.RawMessage `json:"data"`
}

type WebServerInfo struct {
	ID     int    `json:"id"`
	App    string `json:"application"`
	Server string `json:"server_name"`
	Node   string `json:"node_name"`
	Type   string `json:"server_type"`
}

type WebAddTaskReq struct {
	Serial bool      `json:"serial"`
	Items  []WebTask `json:"items"`
}

type WebTask struct {
	ServerID int    `json:"server_id"`
	Command  string `json:"command"`
}

type WebQueryTaskResp struct {
	Status int `json:"status"`
}
