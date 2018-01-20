package server

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/golang/glog"
	"github.com/jacobwpeng/sirius/engine"
)

type App struct {
	wg                sync.WaitGroup
	doneChan          chan struct{}
	config            AppConfig
	dispatcher        *Dispatcher
	tcpClients        []*TCPClient
	nextDynamicRankID uint32
	ranks             map[uint32]engine.RankEngine
	tcpClientListener *net.TCPListener
}

func NewApp(config AppConfig) *App {
	return &App{
		doneChan:          make(chan struct{}),
		config:            config,
		nextDynamicRankID: 1,
		ranks:             make(map[uint32]engine.RankEngine),
	}
}

func (app *App) AddRank(rankID uint32,
	rankConfig engine.RankEngineConfig) error {
	app.ranks[rankID] = engine.NewRankEngine(rankConfig)
	return nil
}

func (app *App) Run() {
	var err error
	if app.dispatcher, err = NewDispatcher(app.ranks); err != nil {
		glog.Fatal(err)
	}
	app.dispatcher.Start()
	app.wg.Add(1)
	go app.AcceptClientConnections()
	app.wg.Add(1)
	go app.AcceptServerConnections()
	app.WaitForExit()
	close(app.doneChan)
	if app.tcpClientListener != nil {
		app.tcpClientListener.Close()
	}
	app.dispatcher.Stop()
	for _, tcpClient := range app.tcpClients {
		tcpClient.StopAndWait()
	}
	app.wg.Wait()
}

func (app *App) AcceptClientConnections() {
	defer app.wg.Done()
	l, err := net.Listen("tcp", app.config.AcceptClientAddress)
	if err != nil {
		glog.Fatal(err)
	}
	listener, _ := l.(*net.TCPListener)
	app.tcpClientListener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-app.doneChan:
				return
			default:
			}
			glog.Fatal(err)
		}
		glog.V(2).Infof("New connection %s", conn.RemoteAddr())
		tcpConn, _ := conn.(*net.TCPConn)
		client := NewTCPClient(app.dispatcher, tcpConn)
		app.tcpClients = append(app.tcpClients, client)
		go client.Run()
	}
}

func (app *App) AcceptServerConnections() {
	defer app.wg.Done()
}

func (app *App) WaitForExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	glog.Infof("Signal %s", <-ch)
}
