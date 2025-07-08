package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/config"
	"os"
	"path/filepath"
	"strings"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var engine *xorm.Engine

func GetEngine() *xorm.Engine {
	return engine
}

type CustomTableMapper struct {
	snakeMapper names.SnakeMapper
	tablePrefix string
}

func NewCustomTableMapper() *CustomTableMapper {
	return &CustomTableMapper{
		snakeMapper: names.SnakeMapper{},
		tablePrefix: "BO_",
	}
}

func (c *CustomTableMapper) Obj2Table(objName string) string { // 结构体名称转换成表名
	return c.tablePrefix + strings.ToUpper(c.snakeMapper.Obj2Table(objName))
}

func (c *CustomTableMapper) Table2Obj(tableName string) string { // 表名转换成结构体名称
	tbName := strings.TrimPrefix(tableName, c.tablePrefix)  // BO_NODE_DEPLOY_TASK --> NODE_DEPLOY_TASK
	return c.snakeMapper.Table2Obj(strings.ToLower(tbName)) // NODE_DEPLOY_TASK --> nodeDeployTask
}

func InitEngine() error {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		config.Get().User, config.Get().Pwd, config.Get().Host, config.Get().Port, config.Get().Database)
	newEngine, err := xorm.NewEngine(config.Get().Drive, dns)
	if err != nil {
		return err
	}
	// 设置表名和实体名之间按照驼峰和下划线格式映射
	newEngine.SetTableMapper(NewCustomTableMapper())

	// 设置日志的路径
	if err = os.MkdirAll(filepath.Dir(config.Get().LogPath), os.ModePerm); err != nil {
		return err
	}
	writer, err := os.Create(config.Get().LogPath)
	if err != nil {
		return err
	}
	logger := log.NewSimpleLogger(writer)
	newEngine.SetLogger(logger)
	newEngine.ShowSQL(config.Get().ShowSql)
	newEngine.Logger().SetLevel(log.LOG_INFO)

	engine = newEngine
	return nil
}
