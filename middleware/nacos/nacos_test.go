package nacos

import "testing"

const nacosCfgFile = "/Users/a123/Downloads/go/learn-go/learn-go/middleware/nacos/nacos.json"

func TestSelectAllInstances(t *testing.T) {
	InitNacos(nacosCfgFile)
	instances, err := SelectAllInstances("cmdb", "cmdb")
	if err != nil {
		t.Error(err)
	}
	t.Log(instances)
}
