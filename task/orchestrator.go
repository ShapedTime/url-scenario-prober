package task

import (
	"log"
	"sync"
	"time"
	"url-scenario-prober/clientFactory"
	"url-scenario-prober/promhelper"
	"url-scenario-prober/vars"
)

type Orchestrator struct {
	tasks         Tasks
	doneChannel   chan string
	vars          *vars.Vars
	status        Status
	statusMux     *sync.Mutex
	clientFactory *clientFactory.ClientFactory
	tasksGraph    *Graph
}

func NewOrchestrator(tasks Tasks) *Orchestrator {
	g, err := tasks.composeGraph()
	if err != nil {
		log.Fatalf("error composing graph: %s", err)
	}

	return &Orchestrator{
		tasks:         tasks,
		doneChannel:   make(chan string),
		vars:          vars.NewVars(),
		status:        STATUS_NOT_STARTED,
		clientFactory: clientFactory.NewClientFactory(10),
		statusMux:     &sync.Mutex{},
		tasksGraph:    g,
	}
}

func (o *Orchestrator) SetClientFactory(clientFactory *clientFactory.ClientFactory) {
	o.clientFactory = clientFactory
}

func (o *Orchestrator) Run() {
	var err error
	o.statusMux.Lock()
	status := o.status
	o.statusMux.Unlock()

	if err != nil {
		log.Fatalf("error composing graph: %s", err)
	}

	if status != STATUS_RUNNING {
		log.Println("Starting orchestrator")
		o.statusMux.Lock()
		o.status = STATUS_RUNNING
		o.statusMux.Unlock()
		o.resetTaskStatus()

		go o.listenDoneTasks()

		childrenOfRoot, err := o.tasksGraph.GetChildren(ROOT_GRAPH)
		if err != nil {
			log.Fatalf("error getting children of root: %s", err)
		}

		for _, g := range childrenOfRoot {
			o.runTask(g)
		}
	}
}

func (o *Orchestrator) listenDoneTasks() {
	for {
		select {
		case tName := <-o.doneChannel:
			task := o.tasks.GetTask(tName)
			taskStatus := task.GetStatus()

			if taskStatus != STATUS_SUCCESS {
				childrenTask := o.tasksGraph.SetChildrenAsFailed(tName)

				for i, _ := range childrenTask {
					promhelper.TaskFail(childrenTask[i].GetName(), childrenTask[i].GetUrl())
					promhelper.CountTaskFailed(childrenTask[i].GetName(), childrenTask[i].GetUrl())
				}

				if !o.tasksGraph.StillRunning() {
					o.statusMux.Lock()
					o.status = STATUS_SUCCESS
					o.statusMux.Unlock()
					return
				}
				break
			}

			nextTasks, err := o.tasksGraph.GetChildren(tName)
			if err != nil {
				log.Fatalf("error getting children of %s: %s", tName, err)
			}

			if nextTasks == nil || len(nextTasks) == 0 {
				nextTask := o.tasksGraph.GetNotStartedTask()
				if nextTask == nil {
					o.statusMux.Lock()
					o.status = STATUS_SUCCESS
					o.statusMux.Unlock()
					return
				}

				nextTasks = []*Task{nextTask}
			}

			for i, _ := range nextTasks {
				o.runTask(nextTasks[i])
			}
		}
	}
}

func (o *Orchestrator) runTask(task *Task) {
	//log.Println("calling runTask for", task.GetName())
	task.SetStatus(STATUS_RUNNING)
	client := o.clientFactory.NewHttpClient()
	go func() {
		defer func() { o.doneChannel <- task.GetName() }()

		promhelper.TaskStarted(task.GetName(), task.GetUrl())

		//log.Println("running task", task.GetName(), "with status", task.GetStatus().String())
		start := time.Now()
		task.Run(o.vars, client)
		dur := time.Since(start)
		promhelper.TaskDuration(task.GetName(), task.GetUrl(), dur.Seconds())
		//log.Println("task", task.GetName(), "finished", "with status", task.GetStatus().String())

		if task.GetStatus() == STATUS_SUCCESS {
			//log.Println("Task", task.GetName(), "succeeded")
			promhelper.TaskSuccess(task.GetName(), task.GetUrl())
		} else if task.GetStatus() == STATUS_FAILED || task.GetStatus() == STATUS_FAILED_UNEXPECTED {
			log.Printf("Task %s failed (%s) with message: %s", task.GetName(), task.GetStatus().String(), task.GetStatusMessage())
			promhelper.TaskFail(task.GetName(), task.GetUrl())
			promhelper.CountTaskFailed(task.GetName(), task.GetUrl())
		} else {
			log.Printf("This should not happen: task %s has status %s", task.GetName(), task.GetStatus().String())
		}

		promhelper.CountTaskFinished(task.GetName(), task.GetUrl(), task.GetStatus().String())
		//log.Println("runTask for", task.GetName(), "finished", "with status", task.GetStatus().String())
	}()
}

func (o *Orchestrator) resetTaskStatus() {
	for i, _ := range o.tasks.Tasks {
		o.tasks.Tasks[i].SetStatus(STATUS_NOT_STARTED)
	}
}
