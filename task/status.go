package task

type Status int

const (
	STATUS_NOT_STARTED Status = iota
	STATUS_RUNNING
	STATUS_SUCCESS
	STATUS_FAILED
	STATUS_PENDING
	STATUS_FAILED_UNEXPECTED
)

func (t *Task) GetStatus() Status {
	t.Lock()
	status := t.status
	t.Unlock()
	return status
}

func (t *Task) SetStatus(status Status) {
	//log.Println("Setting status", status, "for task", t.GetName())
	t.Lock()
	t.status = status
	t.Unlock()
	//log.Println("Status set", t.GetStatus(), "for task", t.GetName())
}

func (s Status) String() string {
	switch s {
	case STATUS_NOT_STARTED:
		return "NOT_STARTED"
	case STATUS_RUNNING:
		return "RUNNING"
	case STATUS_SUCCESS:
		return "SUCCESS"
	case STATUS_FAILED:
		return "FAILED"
	case STATUS_PENDING:
		return "PENDING"
	case STATUS_FAILED_UNEXPECTED:
		return "FAILED_UNEXPECTED"
	}

	return "SHOULDN'T_HAPPEN"
}
