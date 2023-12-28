# URL Scenario Prober
This Go project is a URL scenario prober that allows you to define a set of tasks (HTTP requests) to be executed in a specific order and expose their results as prometheus metrics. It's useful for probing complex scenarios where the state is carried over multiple HTTP requests.

# Key Components
- App: The main application structure. It contains the application configuration, tasks, and an orchestrator to manage the tasks.
- Orchestrator: Manages the execution of tasks. It maintains a graph of tasks and their dependencies, and runs tasks in the correct order.
- Task: Represents a single task (HTTP request) to be executed. It contains the details of the request and the status of the task.
- ClientFactory: Creates HTTP clients for executing tasks.
- main.go: The entry point of the application. It parses command-line arguments and starts the application.

# How to Run
To run the application, you need to provide paths to two configuration files: the application configuration file and the tasks configuration file. You can do this using command-line flags when starting the application:

```bash
go run main.go -config path/to/config.yml -task path/to/tasks.yml
```
The config flag specifies the path to the application configuration file, and the task flag specifies the path to the tasks configuration file.

# Configuration
The application configuration file (config.yml) contains settings for the application, such as whether to ignore SSL certificates and the timeout for HTTP requests.

The tasks configuration file (tasks.yml) defines the tasks to be executed. Each task represents an HTTP request and contains details such as the URL, method, headers, and body. Tasks can also specify dependencies on other tasks, in which case they will not be executed until all their dependencies have completed successfully.
For example configuration, please refer to example-tasks.yml

# Prometheus Integration
The application includes integration with Prometheus for monitoring. It exposes a Prometheus endpoint on a configurable port and provides several metrics related to task execution, such as the number of tasks started, succeeded, and failed.