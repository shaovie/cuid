package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"prof"
	"syscall"

	"ilog"
	"lru"
	"math/rand"
)

// Launch args
var (
	sConfigPath = ""
	showVersion = false
	toProfile   = false
	toDaemon    = false
)

func usage() {
	var usageStr = `
    Server options:
    -c  FILE                Configuration file
    -p                      To profile
    -d                      To daemon mode

    Common options:
    -h                      Show this message
    -v                      Show version
    `
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
func parseFlag() {
	flag.StringVar(&sConfigPath, "c", "svc.json", "Configuration file.")
	flag.BoolVar(&showVersion, "v", showVersion, "Show version.")
	flag.BoolVar(&toProfile, "p", toProfile, "To profile.")
	flag.BoolVar(&toDaemon, "d", toDaemon, "To daemon mode.")

	flag.Usage = usage
	flag.Parse()
}
func signalHandle() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP)
	for {
		sig := <-ch
		if sig == syscall.SIGHUP {
			ilog.Log("signal hup")
		}
	}
}
func main() {
	parseFlag()

	if showVersion {
		fmt.Printf("version %s\n", "1.0.0")
		os.Exit(0)
	}

	if toDaemon {
		daemon()
	}

	//
	outputPid("pid")

	go signalHandle()

	if toProfile {
		prof.StartProf()
	}

	if err := initSvc(); err != nil {
		fmt.Printf("Error - init [%s]\n", err.Error())
		os.Exit(1)
	}
	ilog.Log("-------------launch ok ------------------")
	ilog.Log("-------------launch ok ------------------")

	lc := lru.New(1000)
	var rn int32
	for i := 0; i < 10000; i++ {
		rn = rand.Int31n(1000000)
		lc.Add(i, rn)
	}
	for i := 0; i < 1000; i++ {
		v, ok := lc.Get(i)
		if ok {
			ilog.Log("%d:%d", i, v.(int32))
		}
	}

	select {}

	os.Exit(1)
}
