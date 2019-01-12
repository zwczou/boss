package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"zwczou/boss/bossd"

	log "github.com/sirupsen/logrus"
)

var cfgName string

func main() {
	flag.StringVar(&cfgName, "cfg", "contrib/bossd.yaml", "配置文件路径")
	flag.Parse()

	opts := bossd.NewOption()
	err := opts.Load(cfgName)
	if err != nil {
		log.Fatal(err)
	}
	opts.Print()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	srv := bossd.NewBossServer(opts)
	srv.Main()
	<-signalChan
	srv.Exit()
}
