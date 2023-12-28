package promhelper

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func RunPrometheus(promPort int) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", promPort), nil)
	if err != nil {
		log.Panicf("failed to start prometheus server: %v", err)
		return
	}
}

var (
	tasksResult = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tasks_result",
		Help: "0 if tasks fail, 1 if tasks succeed",
	}, []string{"name", "url"})

	taskDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "task_duration",
		Help: "Task duration",
	}, []string{"name", "url"})

	startedTasks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "started_tasks",
		Help: "Number of started tasks",
	}, []string{"name", "url"})

	failedTasks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "failed_tasks",
		Help: "Number of failed tasks",
	}, []string{"name", "url"})

	finishedTesks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finished_tasks",
		Help: "Number of finished tasks",
	}, []string{"name", "url", "status"})
)

func TaskSuccess(name, url string) {
	tasksResult.WithLabelValues(name, url).Set(1)
}

func TaskFail(name, url string) {
	tasksResult.WithLabelValues(name, url).Set(0)
}

func TaskDuration(name, url string, duration float64) {
	taskDuration.WithLabelValues(name, url).Observe(duration)
}

func TaskStarted(name, url string) {
	startedTasks.WithLabelValues(name, url).Inc()
}

func CountTaskFailed(name, url string) {
	failedTasks.WithLabelValues(name, url).Inc()
}

func CountTaskFinished(name, url, status string) {
	finishedTesks.WithLabelValues(name, url, status).Inc()
}
