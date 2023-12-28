package task

import (
	"fmt"
	"url-scenario-prober/graph"
)

type Graph struct {
	*graph.Graph[Task]
	graphMap map[string]*graph.Graph[Task]
}

const ROOT_GRAPH = "rootGraph"

func (t *Tasks) composeGraph() (*Graph, error) {
	rootGraph := graph.NewGraph[Task](&Task{
		Name:   ROOT_GRAPH,
		status: STATUS_SUCCESS,
	})

	taskGraph := Graph{
		graphMap: map[string]*graph.Graph[Task]{ROOT_GRAPH: rootGraph},
		Graph:    rootGraph,
	}

	dependentTasks := make(map[string][]string)

	for _, task := range t.Tasks {
		if len(task.DependsOn) == 0 {
			err := taskGraph.addTask(task, "")
			if err != nil {
				return nil, err
			}
			continue
		}

		don := task.GetDependsOn()
		for _, dependentTask := range don {
			dependentTasks[dependentTask] = append(dependentTasks[dependentTask], task.GetName())
		}
	}

	for len(dependentTasks) > 0 {
		for taskName, dependentTaskNames := range dependentTasks {
			if _, exists := taskGraph.graphMap[taskName]; exists {

				// add dependent tasks to graph
				for _, dependentTaskName := range dependentTaskNames {
					err := taskGraph.addTask(t.Tasks[dependentTaskName], taskName)
					if err != nil {
						return nil, fmt.Errorf("error adding task %s to graph: %w", dependentTaskName, err)
					}
				}

				delete(dependentTasks, taskName)
			}
		}
	}

	return &taskGraph, nil
}

func (g *Graph) addTask(task *Task, parentTask string) error {
	if parentTask == "" {
		parentTask = ROOT_GRAPH
	}

	parentGraph := g.graphMap[parentTask]
	if parentGraph == nil {
		return fmt.Errorf("parent task %s not found", parentTask)
	}

	newGraph := graph.NewGraph[Task](task)
	parentGraph.AddNode(newGraph)

	g.graphMap[task.Name] = newGraph
	return nil
}

func (g *Graph) GetChildren(taskName string) ([]*Task, error) {
	graph := g.graphMap[taskName]
	if graph == nil {
		return nil, fmt.Errorf("task %s not found", taskName)
	}

	tasks := make([]*Task, 0, 5)
	nodes := graph.GetNodes()
	for i, _ := range nodes {
		task := nodes[i].GetObj()
		if task.GetStatus() == STATUS_NOT_STARTED {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (g *Graph) GetNotStartedTask() *Task {
	t, _ := g.Graph.BFSFirst(func(obj Task) bool {
		if obj.GetStatus() == STATUS_NOT_STARTED {
			return true
		}
		return false
	})

	return t
}

func (g *Graph) SetChildrenAsFailed(taskName string) []Task {
	graph := g.graphMap[taskName]
	if graph == nil {
		return nil
	}

	childrenTask := make([]Task, 0)

	nodes := graph.GetNodes()
	for i, _ := range nodes {
		nodes[i].GetObj().SetStatus(STATUS_FAILED)
		childrenTask = append(childrenTask, *nodes[i].GetObj())
		childrenTask = append(childrenTask, g.SetChildrenAsFailed(nodes[i].GetObj().Name)...)
	}
	return childrenTask
}

func (g *Graph) StillRunning() bool {
	_, match := g.Graph.BFSFirst(func(obj Task) bool {
		if obj.GetStatus() == STATUS_RUNNING {
			return true
		}
		return false
	})

	return match
}
