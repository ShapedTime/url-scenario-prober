package main

import (
	"flag"
	"url-scenario-prober/app"
)

func main() {
	appConfig := flag.String("config", "config.yml", "Path to config file")
	tasksConfig := flag.String("task", "tasks.yml", "Path to tasks config file")

	flag.Parse()

	a, err := app.NewApp(*appConfig, *tasksConfig)
	if err != nil {
		panic(err)
	}

	a.Run()
}
