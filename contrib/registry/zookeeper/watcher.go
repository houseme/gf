package zookeeper

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
)

var _ gsvc.Watcher = &watcher{}

type watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
	event  chan struct{}
	set    *serviceSet
}

// Proceed is used to watch the key.
func (w *watcher) Proceed() (services []*gsvc.Service, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
	case <-w.event:
	}
	ss, ok := w.set.services.Load().([]*gsvc.Service)
	if ok {
		services = append(services, ss...)
	}
	return
}

// Close the watcher.
func (w *watcher) Close() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
