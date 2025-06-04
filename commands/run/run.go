package run 

// commands should implement the `Run()` method
// that calls the sendIPC
type RunCommand interface {
	Register()
	Run() error
	//ParseArgs() error
}

