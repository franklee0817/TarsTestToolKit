package communicator

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/franklee0817/t3k/backend/constants"
)

var comm *tars.Communicator

func init() {
	if comm == nil {
		comm = tars.NewCommunicator()
		comm.SetProperty(constants.LocatorKeyLocal, constants.ServiceNameServiceRegistry)
	}
}

// GetCommunicator 获取tars.Communicator对象
func GetCommunicator() *tars.Communicator {
	return comm
}

// StringToProxy tars.Communicator.StringToProxy的代理方法
func StringToProxy(servant string, p tars.ProxyPrx) {
	comm.StringToProxy(servant, p)
}

// GetAllEndpoints 获取服务的所有endpoints列表
func GetAllEndpoints(obj string) []string {
	endpoints := tars.GetManager(comm, obj).
		GetAllEndpoint()

	eps := make([]string, 0)
	for _, ep := range endpoints {
		eps = append(eps, fmt.Sprintf("%s -h %s -p %d -t %d", ep.Proto, ep.Host, ep.Port, ep.Timeout))
	}

	return eps
}
