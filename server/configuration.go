package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/maxzerbini/ovo/cluster"
)

const (
	DefaultConfPath string = "./conf/severconf.json"
)

var (
	starNames []string = []string{
		"mizar", "rigel", "acamar", "akrab", "adhara", "aladfar", "alcor", "aldebaran",
		"algol", "alphard", "alphecca", "alshain", "altais", "antares", "asterion", "arcturus",
		"bellatrix", "betelgeuse", "brachium", "canopus", "castor", "cebalrai", "cheleb", "corcaroli",
		"denebola", "dheneb", "diadem", "edasich", "electra", "elnath", "fomalhaut", "furud",
		"gacrux", "gemma", "gienah", "gomeisa", "hadar", "haldus", "hydrobius", "jabbah",
		"kajam", "kitalpha", "kornephoros", "kuma", "lesath", "mahasim", "maia", "markab",
		"megrez", "menkab", "merak", "merope", "mintaka", "miram", "mirzam", "muliphein",
		"muscida", "naos", "nekkar", "nembus", "nusakan", "okul", "peacock", "pherkad",
		"pleione", "polaris", "porrima", "procyon", "propus", "rasalas", "rastaban", "regor",
		"regulus", "rotanev", "rukbat", "sabik", "sadalbari", "sadatoni", "saiph", "sarin",
		"seginus", "sheliak", "situla", "sirius", "sterope", "sulafat", "syrma", "thabit",
		"tania", "tazaret", "tureis", "unuk", "vega", "wasat", "wezen", "yildun",
		"zaniah", "zavijava", "zosma", "zubenelgenubi", "torcularis", "alioth", "alkalurops", "botein"}
)

type ServerConf struct {
	ServerNode  *cluster.ClusterTopologyNode
	Topology    cluster.ClusterTopology
	Debug       bool
	tmpPath     string
	HttpBindAll bool
	TcpBindAll  bool
}

func (cnf *ServerConf) Init(tmpPath string) {
	cnf.ServerNode.StartDate = time.Now()
	cnf.ServerNode.Node.State = cluster.Active
	if cnf.ServerNode.Twins == nil {
		cnf.ServerNode.Twins = make([]string, 0)
	}
	if cnf.ServerNode.Stepbrothers == nil {
		cnf.ServerNode.Stepbrothers = make([]string, 0)
	}
	if cnf.ServerNode.Node.Name == "" {
		cnf.ServerNode.Node.Name = GetRandomName()
	}
	ip := GetServerIP()
	if cnf.ServerNode.Node.Host == "" {
		cnf.HttpBindAll = true
		cnf.ServerNode.Node.Host = ip
		log.Printf("Setting Host %s", ip)
	}
	if cnf.ServerNode.Node.APIHost == "" {
		cnf.TcpBindAll = true
		cnf.ServerNode.Node.APIHost = ip
		log.Printf("Setting APIHost %s", ip)
	}
	cnf.ServerNode.UpdateDate = time.Now()
	cluster.SetCurrentNode(cnf.ServerNode, &cnf.Topology)
	cnf.tmpPath = tmpPath
}

func (cnf *ServerConf) WriteTmp() {
	WriteConfiguration(cnf.tmpPath, cnf)
}

func LoadConfiguration(path string) ServerConf {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("Configuration file not found at %s", path)
		os.Exit(1)
	}
	var jsontype ServerConf
	json.Unmarshal(file, &jsontype)
	return jsontype
}

func WriteConfiguration(path string, conf *ServerConf) {
	data, _ := json.Marshal(conf)
	e := ioutil.WriteFile(path, data, 0x666)
	if e != nil {
		log.Printf("Configuration file write error at %s\r\n", path)
	}
}

func GetRandomName() string {
	index := rand.Int31n(int32(len(starNames)))
	num := rand.Int31n(1000)
	name := []string{starNames[index], strconv.Itoa(int(num))}
	return strings.Join(name, "-")
}

func GetServerIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Panic("Can't acess the network interfaces.")
	}
	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Panic("Can't read the IP.")
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				//log.Printf("Evaluating IPNet: %v\r\n", v.IP)
				if (v.IP.To4() != nil) && !v.IP.IsLinkLocalUnicast() && !v.IP.IsLoopback() && !v.IP.IsMulticast() {
					ip = v.IP
				}
			case *net.IPAddr:
				//log.Printf("Evaluating IPAddr: %v\r\n", v.IP)
				if (v.IP.To4() != nil) && !v.IP.IsLinkLocalUnicast() && !v.IP.IsLoopback() && !v.IP.IsMulticast() {
					ip = v.IP
				}
			}
		}
	}
	// process IP address
	return ip.String()
}
