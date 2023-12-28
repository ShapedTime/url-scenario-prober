package app

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-scenario-prober/clientFactory"
	"url-scenario-prober/promhelper"
	"url-scenario-prober/task"
)

type App struct {
	appConfig    AppConfig
	tasks        task.Tasks
	orchestrator *task.Orchestrator
}

func NewApp(appConfigFile, tasksConfigFile string) (*App, error) {
	config, err := LoadAppConfig(appConfigFile)
	if err != nil {
		return nil, err
	}
	log.Println("Loaded config")

	tasks, err := task.LoadTasks(tasksConfigFile)
	if err != nil {
		return nil, err
	}
	log.Println("Loaded tasks")

	orchestrator := task.NewOrchestrator(tasks)

	a := &App{
		appConfig:    *config,
		tasks:        tasks,
		orchestrator: orchestrator,
	}

	a.Init()

	go promhelper.RunPrometheus(a.appConfig.PrometheusPort)

	return a, nil
}

func (a *App) Init() {
	if a.appConfig.IgnoreCertificates {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	a.orchestrator.SetClientFactory(clientFactory.NewClientFactory(a.appConfig.TimeoutSeconds))
}

func (a *App) Run() {
	timer := time.NewTicker(time.Duration(a.appConfig.ScrapeSeconds) * time.Second)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	a.orchestrator.Run()
	for {
		select {
		case <-timer.C:
			a.orchestrator.Run()
		case <-sig:
			return
		}
	}
}
