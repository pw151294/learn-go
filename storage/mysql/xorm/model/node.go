package model

import (
	"database/sql"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/datasource"
	"xorm.io/xorm"
)

type NodeRepository interface {
	Insert(*Node) error
	Update(*Node) error
	Delete(int64) error
	SelectByID(int64) (*Node, error)
	SelectByIds([]int64) ([]*Node, error)
}

type Node struct {
	Id         int64          `xorm:"pk autoincr comment('主键')" json:"activationId"`
	InstanceId sql.NullString `xorm:"unique varchar(36) comment('实例ID')" json:"instanceId"`
}

func NewNodeRepository() NodeRepository {
	return &NodeRepositoryImpl{engine: datasource.GetEngine()}
}

type NodeRepositoryImpl struct {
	engine *xorm.Engine
}

func (r *NodeRepositoryImpl) Insert(node *Node) error {
	_, err := r.engine.Insert(node)
	return err
}

func (r *NodeRepositoryImpl) Update(node *Node) error {
	_, err := r.engine.ID(node.Id).Update(node)
	return err
}

func (r *NodeRepositoryImpl) Delete(id int64) error {
	_, err := r.engine.ID(id).Delete(&Node{})
	return err
}

func (r *NodeRepositoryImpl) SelectByID(id int64) (*Node, error) {
	node := &Node{}
	_, err := r.engine.ID(id).Get(node)
	return node, err
}

func (r *NodeRepositoryImpl) SelectByIds(ids []int64) ([]*Node, error) {
	nodes := make([]*Node, 0)
	err := r.engine.In("id", ids).Find(&nodes)
	return nodes, err
}
