package admin

import (
	"sync"
	"zwczou/gobase/container"

	"github.com/codegangsta/inject"
)

type adminServer struct {
	sync.RWMutex
	inject.Injector
	name     string
	exitChan chan struct{}
}

func newAdminServer() *adminServer {
	s := &adminServer{
		Injector: inject.New(),
		name:     "admin.system",
		exitChan: make(chan struct{}),
	}
	return s
}

func (as *adminServer) Name() string {
	return as.name
}

func (as *adminServer) Load(app *container.Container) error {
	return nil
}

func (as *adminServer) Exit() {
	close(as.exitChan)
}

func init() {
	container.Pre(newAdminServer())
}
