package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/codingeasygo/pdservice/discover"
	"github.com/codingeasygo/util/xprop"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Printf("pdservice %v version\n", Version)
		return
	}
	confPath := "conf/pdservice.properties"
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}
	wd, _ := os.Getwd()
	fmt.Printf("starting pdservice with working on %v\n", wd)
	cfg := xprop.NewConfig()
	err := cfg.Load(confPath)
	if err != nil {
		panic(err)
	}
	cfg.Print()
	listenAddr := cfg.StrDef(":9231", "listen")
	refreshTime := cfg.Int64Def(10000, "refresh_time")
	triggerAdded := cfg.StrDef("", "trigger_added")
	triggerRemoved := cfg.StrDef("", "trigger_removed")
	triggerUpdated := cfg.StrDef("", "trigger_updated")
	priview := cfg.StrDef("", "preview")
	server := discover.NewDiscover()
	server.TriggerBash = cfg.StrDef("bash", "trigger_bash")
	server.DockerFinder = cfg.StrDef("", "trigger_finder")
	server.DockerCert = cfg.StrDef("certs", "docker_cert")
	server.DockerAddr = cfg.StrDef("tcp://127.0.0.1:2376", "docker_addr")
	server.DockerHost = cfg.StrDef("127.0.0.1", "docker_host")
	server.DockerClearDelay = time.Duration(cfg.Int64Def(0, "docker_clear_delay")) * time.Minute
	server.DockerClearExc = cfg.ArrayStrDef(nil, "docker_clear_exc")
	server.DockerPruneDelay = time.Duration(cfg.Int64Def(0, "docker_prune_delay")) * time.Minute
	server.DockerPruneExc = cfg.ArrayStrDef(nil, "docker_prune_exc")
	server.HostSuff = cfg.StrDef("", "host_suffix")
	server.HostProto = cfg.StrDef("https", "host_proto")
	server.HostSelf = cfg.StrDef("https", "host_self")
	server.SrvPrefix = cfg.StrDef("/_s", "srv_prefix")
	if len(priview) > 0 {
		server.Preview, err = template.ParseFiles(priview)
		if err != nil {
			panic(err)
		}
	}
	discover.SetLogLevel(cfg.IntDef(30, "log"))
	server.StartRefresh(time.Duration(refreshTime)*time.Millisecond, triggerAdded, triggerRemoved, triggerUpdated)
	err = http.ListenAndServe(listenAddr, server)
	if err != nil {
		panic(err)
	}
}
