package model

import (
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/cluster"
)

type Any interface { }

type OvoResponse struct {
	Status string
	Code string
	Data Any
}

type OvoKVRequest struct {
	Key string
	Data []byte
	Collection string
	TTL int
	Hash int
}

type OvoKVUpdateRequest struct {
	Key string
	NewKey string
	Data []byte
	NewData []byte
	Hash int
	NewHash int
}

type OvoKVResponse struct {
	Key string
	Data []byte
}

type OvoKVKeys struct {
	Keys []string
}

type OvoTopologyNode struct {
	Name string
	HashRange []int
	Host string
	Port int
}

type OvoTopology struct {
	Nodes []*OvoTopologyNode
}

func NewOvoResponse(status string, code string, data Any) *OvoResponse {
	return &OvoResponse{Status:status, Code:code, Data: data}
}

func NewOvoKVResponse(obj *storage.MetaDataObj) *OvoKVResponse {
	var rsp = &OvoKVResponse{Key:obj.Key, Data:obj.Data}
	return rsp
}

func NewMetaDataObj(req *OvoKVRequest) *storage.MetaDataObj {
	var obj = new(storage.MetaDataObj)
	obj.Key = req.Key
	obj.Data = req.Data
	obj.Collection = req.Collection
	obj.TTL = req.TTL
	obj.Hash = req.Hash
	return obj
}

func NewMetaDataUpdObj(req *OvoKVUpdateRequest) *storage.MetaDataUpdObj {
	var obj = new(storage.MetaDataUpdObj)
	obj.Key = req.Key
	obj.NewKey = req.NewKey
	obj.Data = req.Data
	obj.NewData = req.NewData
	obj.Hash = req.Hash
	obj.NewHash = req.NewHash
	return obj
}

func NewOvoTopologyNode(node *cluster.ClusterTopologyNode) *OvoTopologyNode{
	return &OvoTopologyNode{Name:node.Node.Name,HashRange:node.Node.HashRange,Host:node.Node.ExtHost,Port:node.Node.Port}
}

func NewOvoTopology(topology *cluster.ClusterTopology) *OvoTopology{
	ret := &OvoTopology{Nodes:make([]*OvoTopologyNode,0)}
	for _,node := range topology.Nodes {
		ret.Nodes = append(ret.Nodes, NewOvoTopologyNode(node))
	}
	return ret
}