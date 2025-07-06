package cluster

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"iflytek.com/weipan4/learn-go/logger/zap"
	nacos1 "iflytek.com/weipan4/learn-go/middleware/nacos"
)

func StartNode(serviceName, groupName string, port int) {
	metaOpts := []nacos1.MetadataOptions{
		nacos1.WithInstanceId(fmt.Sprintf("%s~id", serviceName)),
		nacos1.WithPort(port),
	}
	regInstOpts := []nacos1.RegisterInstanceParamOptions{
		nacos1.WithServiceName(serviceName),
		nacos1.WithInstancePort(uint64(port)),
		nacos1.WithGroupName(groupName),
	}

	err := nacos1.RegisterService(metaOpts, regInstOpts)
	if err != nil {
		zap.GetLogger().Error("register service failed", "service name", serviceName, "group name", groupName)
		return
	}

	// 实现集群订阅
	subSvcOptions := []nacos1.SubScribeParamOptions{
		nacos1.WithSubscribeService(serviceName),
		nacos1.WithSubscribeClusters([]string{"DEFAULT"}),
		nacos1.WithSubscribeGroupName(groupName),
		nacos1.WithSubscribeCallback(nacos1.SubScribeCallback),
	}
	err = nacos1.SubscribeService(subSvcOptions)
	if err != nil {
		zap.GetLogger().Error("subscribe service failed", "service name", serviceName, "group name", groupName)
		return
	}

	zap.GetLogger().Info("register service success", "service name", serviceName, "group name", groupName)
	r := gin.Default()
	r.Run(fmt.Sprintf(":%d", port))
}
