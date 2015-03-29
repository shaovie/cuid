package main

import (
	"math/rand"
	"time"

	"ilog"
)

func initSvc() (err error) {
	if err = ilog.InitLog("log"); err != nil {
		return err
	}

	// Init randdom generater
	rand.Seed(time.Now().UnixNano())
	return err
}
