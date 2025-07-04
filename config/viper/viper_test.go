package main

import (
	"encoding/json"
	"github.com/spf13/viper"
	"iflytek.com/weipan4/learn-go/config/viper/config"
	"testing"
)

func TestViperReadConfig(t *testing.T) {
	// 读取配置
	v := viper.New()
	v.SetConfigName("st-agent")
	v.SetConfigType("json")
	v.SetConfigFile("/Users/a123/Downloads/go/learn-go/learn-go/config/viper/config/st-agent.json")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}

	// 读取配置
	bytes, err := json.Marshal(v.AllSettings())
	if err != nil {
		t.Fatal(err)
	}
	cfg := &config.AgentConfig{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// 打印配置
	marshal, _ := json.Marshal(config.AgentCfg)
	t.Log("agent config: ", string(marshal))
}
