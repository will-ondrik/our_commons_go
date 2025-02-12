package dtos

type Task struct {
	Type               string
	Url                string
	ExtractFromElement string
}

type Tasks struct {
	Contract    Task
	Hospitality Task
	Travel      Task
}

type TaskErrors struct {
	Contract    bool
	Hospitality bool
	Travel      bool
}

type RedoTask struct {
	MpHtml     MpHtml
	Tasks      Tasks
	TaskErrors TaskErrors
}
