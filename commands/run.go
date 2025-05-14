package run 

// commands should implement the `Run()` method
// that calls the sendIPC
type Run interface {
	Run() (string, error)
}

type Ipc interface {
	ProcessResponse() string
	ParseArgs() error
}
