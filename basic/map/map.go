package main

import (
	"iflytek.com/weipan4/learn-go/logger/zap"
)

func add(k, v int, schema map[int]int) {
	schema[k] = v
}

func main() {
	zap.InitLogger()

	m := make(map[int]int)
	add(1, 1, m)
	zap.GetLogger().Info("create map", "val", m)

	add(2, 1, m)
	zap.GetLogger().Info("add to map", "val", m)

	add(1, 2, m)
	zap.GetLogger().Info("modify map", "val", m)
}
