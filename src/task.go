package src

type Task interface {
	Start() error
	Stop() error
	Input(data interface{})
	Output(output func(data interface{}))
}
