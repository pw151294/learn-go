package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/iris/config"
	"iflytek.com/weipan4/learn-go/net/iris/pkg/resp"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	"net/http"
	"strings"
	"sync"
	"time"
)

type InstallHost struct {
	Instance string `json:"instance"`
}

type Request struct {
	InstanceList []string `json:"instanceList"`
}

type UpStreamNode string

const (
	UpStreamGseHeaderKey = "UpStream-Gse"
)

func FanOutFanIn(ctx iris.Context) {
	//.判断上游节点是否是其他的gse节点 如果是的话说明该节点是被重定向的目标节点 无需再重定向到其他节点
	usGse := ctx.GetHeader(UpStreamGseHeaderKey)
	if usGse != "" {
		zap.GetLogger().Info("receive agent operation task distributed by gse node", "upstream", usGse)
		ctx.Next()
		return
	}

	// 获取请求参数
	req := Request{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp.E(fmt.Sprintf("invalid request params: %s", err.Error()), ""))
		return
	}
	insIdSet := make(map[string]struct{})
	for _, ins := range req.InstanceList {
		if ins != "" {
			insIdSet[ins] = struct{}{}
		}
	}

	// 查询出所有Key
	redisCli := go_redis.GetClient()
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	keys, _, err := redisCli.ScanType(timeoutCtx, 0, gseConnectionPrefix+"*", 10000, "hash").Result()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(resp.E(fmt.Sprintf("get hash keys of gse connection failed: %s", err.Error()), ""))
		return
	}

	curKey := gseConnectionPrefix + fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port)
	//curKey := fmt.Sprintf(gseConnectionPrefix + fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port))
	// fan out
	var (
		wg   sync.WaitGroup
		errs []error
	)
	for _, key := range keys {
		if key == curKey {
			continue
		}
		// 筛选出请求参数中属于key的实例id
		hKeys, err := redisCli.HKeys(context.Background(), key).Result()
		if err != nil {
			zap.GetLogger().Warn("get instanceIds connected with gse node failed",
				"gse node", strings.TrimPrefix(key, gseConnectionPrefix))
			continue
		}
		insIds := make([]string, 0)
		for _, hk := range hKeys {
			if _, ok := insIdSet[hk]; ok {
				insIds = append(insIds, hk)
			}
		}
		if len(insIds) == 0 {
			continue
		}

		wg.Add(1)
		go func(key string, instanceIds []string) {
			defer wg.Done()
			// 广播(fan out)插件操作任务至别的gse节点
			gseSvc := strings.TrimPrefix(key, gseConnectionPrefix)
			url := fmt.Sprintf("http://%s", gseSvc+ctx.Path())
			zap.GetLogger().Info("distribute agent operation task",
				"src", fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port), "target", gseSvc)
			response, err := resty.New().R().
				SetBody(Request{InstanceList: instanceIds}).
				SetHeader("ScriptContent-Type", "application/json").
				SetHeader(UpStreamGseHeaderKey, gseSvc).
				Post(url)
			if err != nil {
				zap.GetLogger().Error("send agent operation task failed", "message", err)
				errs = append(errs, errors.New(fmt.Sprintf("send agent operation task to  gse service %s failed: %s", gseSvc, err.Error())))
				return
			}
			result := resp.ResultV2{}
			if err := json.Unmarshal(response.Body(), &result); err != nil {
				zap.GetLogger().Error("json unmarshal result failed", "message", err)
				errs = append(errs, errors.New(fmt.Sprintf("parse response from gse node %s failed; %s", gseSvc, err.Error())))
				return
			}
			if result.Code == resp.Fail {
				errs = append(errs, errors.New(fmt.Sprintf("agent operation by gse node %s failed: %s", gseSvc, result.Msg)))
			}
		}(key, insIds)
	}

	// 等待任务完成
	wg.Wait()

	// 执行当前节点的请求
	ctx.Record()
	ctx.Next()
	result := resp.ResultV2{}
	if err := json.Unmarshal(ctx.Recorder().Body(), &result); err != nil {
		zap.GetLogger().Error("json unmarshal result failed", "message", err)
		errs = append(errs, errors.New(fmt.Sprintf("parse result from gse node %s failed %s",
			strings.TrimPrefix(curKey, gseConnectionPrefix), err.Error())))
	}

	// fan in
	err = errors.Join(errs...)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(resp.E(err.Error(), ""))
	} else {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(resp.S("operation for agent success!", ""))
	}
}
