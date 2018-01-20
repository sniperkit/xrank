package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/golang/glog"
	"github.com/jacobwpeng/sirius/engine"
	"github.com/jacobwpeng/sirius/server"
)

var config server.AppConfig

func init() {
	flag.StringVar(&config.AcceptClientAddress, "clientaddr", ":9427",
		"Client listening address")
	flag.StringVar(&config.AcceptServerAddress, "serveraddr", ":9428",
		"Server listening address")
	flag.Parse()
}

func ce(err error) {
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	go func() {
		glog.Info(http.ListenAndServe(":6060", nil))
	}()
	loc, err := time.LoadLocation("Asia/Shanghai")
	ce(err)
	clearStart, err := time.ParseInLocation("2006-01-02 15:04:05",
		"2017-03-23 17:18:00", loc)
	ce(err)
	app := server.NewApp(config)
	primaryRankConfig := engine.RankEngineConfig{
		MaxSize: 10,
		ClearPeriod: engine.TimePeriod{
			Start:    clearStart,
			Interval: time.Second * 5,
		},
		NoUpdatePeriod: engine.TimePeriod{
			Start:    clearStart,
			Interval: time.Hour * 24 * 1,
			Duration: time.Hour * 24 * 20,
		},
	}

	snapshotRankConfig := engine.RankEngineConfig{
		MaxSize:       5,
		PrimaryRankID: 1,
		SnapshotPeriod: engine.TimePeriod{
			Start:    clearStart,
			Interval: time.Second * 5,
		},
	}

	app.AddRank(1, primaryRankConfig)
	app.AddRank(2, snapshotRankConfig)
	app.Run()
}
