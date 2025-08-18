package main

import (
	"context"
	"encoding/json"
	"fmt"
	mysql "iflytek.com/weipan4/learn-go/storage/mysql/xorm/config"
	ds "iflytek.com/weipan4/learn-go/storage/mysql/xorm/datasource"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/model"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"testing"
	"time"
	"xorm.io/builder"
)

const (
	configFile = "/Users/a123/Downloads/go/learn-go/learn-go/storage/mysql/xorm/config/mysq.json"
	redisFile  = "/Users/a123/Downloads/go/learn-go/learn-go/storage/redis/go-redis/config/redis.json"
	port       = 8881
)

func TestBuildRequest(t *testing.T) {
	// 初始化MySQL和Redis数据源
	redis.InitConfig(redisFile)
	mysql.InitConfig(configFile)
	go_redis.InitRedis()
	ds.InitEngine()
	ds.GetEngine().Sync2(new(model.Node))

	// 查询出所有的节点实例
	nodeIds := []int64{99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83}
	repository := model.NewNodeRepository()
	nodes, err := repository.SelectByIds(nodeIds)
	if err != nil {
		t.Fatal(err)
	}
	insIds := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if node.InstanceId.Valid {
			insIds = append(insIds, node.InstanceId.String)
		}
	}
	bytes, _ := json.Marshal(insIds)
	t.Logf("instance ids: %s", string(bytes))

	now := time.Now()
	redisCli := go_redis.GetClient()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	for i := range insIds {
		insId := insIds[i]
		hash := fmt.Sprintf("GSE:CONNECTIONS:172.30.34.73:%d", port+(i%4))
		connMeta := struct {
			InstanceId string `json:"instanceId"`
			LastPing   string `json:"lastPing"`
			ExpireTime string `json:"expireTime"`
		}{
			InstanceId: insId,
			LastPing:   now.Format("2006-01-02 15:04:05"),
			ExpireTime: now.Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		}
		bytes, err := json.Marshal(connMeta)
		if err != nil {
			t.Errorf("marshal connMeta failed: %v", err)
			continue
		}
		val := string(bytes)
		err = redisCli.HSet(ctx, hash, insId, val).Err()
		if err != nil {
			t.Errorf("hset connMeta failed: %v", err)
		}
	}
}

func TestSelectNodes(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()
	ds.GetEngine().Sync2(new(model.Node))

	nodeIds := []int64{99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83}
	repository := model.NewNodeRepository()
	nodes, err := repository.SelectByIds(nodeIds)
	if err != nil {
		t.Fatal(err)
	}

	// 构建出映射关系
	type InstallHost struct {
		NodeId int64 `json:"node_id"`
	}
	hosts := make([]*InstallHost, 0)
	insIds := make([]string, 0)
	for _, node := range nodes {
		hosts = append(hosts, &InstallHost{NodeId: node.Id})
		if node.InstanceId.Valid {
			insIds = append(insIds, node.InstanceId.String)
		}
	}
	bytes, err := json.Marshal(hosts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("request: %s", string(bytes))
}

func TestSelectByIds(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()
	ds.GetEngine().Sync2(new(model.Node))

	nodeIds := []int64{99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83}
	nodeList := make([]*model.Node, 0)
	err := ds.GetEngine().Table(new(model.Node)).
		Select("id, instance_id").In("id", nodeIds).Find(&nodeList)
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := json.Marshal(nodeList)
	t.Logf("request: %s", string(bytes))
}

func TestExecSql(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()
	ds.GetEngine().Sync2(new(model.Node))

	// 执行SQL查询
	sql := "SELECT id, instance_id FROM BO_NODE WHERE id IN (99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83)"
	rows, err := ds.GetEngine().QueryString(sql)
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := json.Marshal(rows)
	nodes := make([]*model.Node, 0)
	if err := json.Unmarshal(bytes, &nodes); err != nil {
		t.Fatal("unmarshal failed:", err)
	}
	for _, node := range nodes {
		t.Log(node.Id, node.InstanceId)
	}
}

func TestBuildSql(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()
	ds.GetEngine().Sync2(new(model.Node))

	sql, args, err := builder.Select("count(distinct ID)").
		From("BO_JOB_RESULT").
		Where(builder.Gt{"CREATE_TIME": time.Now().Add(-30 * 24 * time.Hour).Format("2006-01-02 15:04:05")}).
		ToSQL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	// 实际执行SQL
	engine := ds.GetEngine()
	var count int64
	has, err := engine.SQL(sql, args...).Get(&count)
	if err != nil {
		t.Fatalf("execute sql failed: %v", err)
	}
	if !has {
		t.Log("no result found")
	} else {
		t.Logf("Query Result: %d", count)
	}
}

// 查询出3个月内执行次数超过3次的任务数
func TestExecSql1(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()

	sql, args, err := builder.Select("count(*)").
		From(builder.As(builder.Select("TASK_ID, count(TASK_ID) as count").
			From("BO_JOB_RESULT").
			Where(builder.Gt{"CREATE_TIME": time.Now().AddDate(0, -3, 1).Format("2006-01-02 15:04:05")}).
			And(builder.Eq{"WORKSPACE_ID": 1}).
			GroupBy("TASK_ID"), "tc")).
		Where(builder.Gte{"count": 3}).
		ToSQL()
	if err != nil {
		t.Fatalf("build sql failed: %v", err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	var count int64
	has, err := ds.GetEngine().SQL(sql, args...).Get(&count)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Logf("no result found")
	} else {
		t.Logf("Query Result: %d", count)
	}
}

// 查询出近七天内的任务执行情况
func TestExecSql2(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()
	subQuery := builder.Select("(UNIX_TIMESTAMP(CREATE_TIME) + 8 * 3600) DIV 86400 AS idx, COUNT(DISTINCT JOB_RESULT_ID) AS cnt").
		From("BO_JOB_RESULT_ITEM").
		Where(builder.Eq{"WORKSPACE_ID": 1}).
		Where(builder.Gte{"CREATE_TIME": time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")}).
		GroupBy("idx")
	sql, args, err := builder.Select("UNIXTIME(idx * 86400, '%Y-%m-%d') AS create_time").
		Select("cnt").
		From(builder.As(subQuery, "t")).
		OrderBy("idx Asc, cnt DESC").
		ToSQL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	var results []struct {
		CreateTime string `xorm:"create_time"`
		Count      int64  `xorm:"cnt"`
	}
	err = ds.GetEngine().SQL(sql, args...).Find(&results)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range results {
		t.Logf("CreateTime: %s, Count: %d", result.CreateTime, result.Count)
	}
}

// 展示近7天的任务类型分布
// select CATE, count(ID)
// from BO_JOB_TASK
// where WORKSPACE_ID = 1
// group by CATE;
func TestExecSql3(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()

	sql, args, err := builder.Select("CATE, count(ID) as count").
		From("BO_JOB_TASK").
		Where(builder.Eq{"WORKSPACE_ID": 1}).
		GroupBy("CATE").
		ToSQL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	var results []struct {
		Cate  string `xorm:"cate"`
		Count int64  `xorm:"count"`
	}
	err = ds.GetEngine().SQL(sql, args...).Find(&results)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range results {
		t.Logf("Category: %s, Count: %d", result.Cate, result.Count)
	}
}

// 展示任务状态占比
func TestExecSql4(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()

	sql, args, err := builder.Select("if(max(STATUS) = 0, 0, 1) as status, count(distinct JOB_RESULT_ID) as count").
		From("BO_JOB_RESULT r").
		Join("", "BO_JOB_RESULT_ITEM ri", "r.ID = ri.JOB_RESULT_ID").
		Where(builder.Eq{"r.WORKSPACE_ID": 1}).
		ToSQL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	var results []struct {
		Status int64 `xorm:"status"`
		Count  int64 `xorm:"count"`
	}
	err = ds.GetEngine().SQL(sql, args...).Find(&results)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range results {
		t.Logf("Status: %d, Count: %d", result.Status, result.Count)
	}
}

// 查询最近执行的10条任务信息
func TestExecSql5(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()

	sql, args, err := builder.Dialect(builder.MYSQL).
		Select("r.ID, r.TASK_NAME, if(max(STATUS) = 0, 0, 1) as status, max(ri.CREATE_TIME) as exec_time, max(ri.EXECUTE_TIME) as exec_duration").
		From("BO_JOB_RESULT r").
		Join("", "BO_JOB_RESULT_ITEM ri", "r.ID = ri.JOB_RESULT_ID").
		Where(builder.Eq{"r.WORKSPACE_ID": 1}).
		GroupBy("r.ID, r.TASK_NAME").
		OrderBy("exec_time DESC").
		Limit(10, 0).
		ToSQL()

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	var results []struct {
		Id           string `xorm:"id"`
		TaskName     string `xorm:"task_name"`
		Status       int64  `xorm:"status"`
		ExecTime     string `xorm:"exec_time"`
		ExecDuration string `xorm:"exec_duration"`
	}
	err = ds.GetEngine().SQL(sql, args...).Find(&results)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range results {
		t.Logf("ID: %s, TaskName: %s, Status: %d, ExecTime: %s, ExecDuration: %s", result.Id, result.TaskName, result.Status, result.ExecTime, result.ExecDuration)
	}
}

func TestExecSql6(t *testing.T) {
	mysql.InitConfig(configFile)
	ds.InitEngine()

	sql, args, err := builder.Select("count(distinct INSTANCE_ID)").
		From("BO_JOB_RESULT_ITEM").
		Where(builder.Eq{"WORKSPACE_ID": 1}).
		ToSQL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated SQL: %s", sql)
	t.Logf("SQL Args: %v", args)

	var count int64
	has, err := ds.GetEngine().SQL(sql, args...).Get(&count)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Logf("no result found")
	} else {
		t.Logf("Query Result: %d", count)
	}
}
