package server

import (
	"time"
	"github.com/maxzerbini/ovo/cluster"
	"encoding/json"
    "io/ioutil"
	"os"
	"log"
)

const CONF_PATH string = "./conf/severconf.json"

type ServerConf struct {
	ServerNode cluster.ClusterTopologyNode
	Topology cluster.ClusterTopology
}

func ( cnf *ServerConf) Init() { 
	cnf.ServerNode.StartDate = time.Now()
	cnf.Topology.AddNode(&cnf.ServerNode)
}

func LoadConfiguration(path string) ServerConf {
	file, e := ioutil.ReadFile(path)
    if e != nil {
		log.Fatalf("Configuration file not found at %s", path)
        os.Exit(1)
    }
    var jsontype ServerConf
    json.Unmarshal(file, &jsontype)
	return jsontype;
}

func WriteConfiguration(path string, conf ServerConf) {
	data, _ := json.Marshal(conf)
	e := ioutil.WriteFile(path, data, 0x666)
    if e != nil {
		log.Fatalf("Configuration file write error at %s", path)
    }
}
