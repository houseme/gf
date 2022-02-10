package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/gconv"
)

var _ gsvc.Watcher = (*watcher)(nil)

type watcher struct {
	serviceName string
	clusters    []string
	groupName   string
	ctx         context.Context
	cancel      context.CancelFunc
	watchChan   chan struct{}
	cli         naming_client.INamingClient
	kind        string
}

// newWatcher creates a watcher.
func newWatcher(ctx context.Context, cli naming_client.INamingClient, serviceName, groupName, kind string, clusters []string) (*watcher, error) {
	w := &watcher{
		serviceName: serviceName,
		clusters:    clusters,
		groupName:   groupName,
		cli:         cli,
		kind:        kind,
		watchChan:   make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)

	e := w.cli.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		Clusters:    clusters,
		GroupName:   groupName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			w.watchChan <- struct{}{}
		},
	})
	return w, e
}

// Proceed returns the next result.
func (w *watcher) Proceed() (services []*gsvc.Service, err error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
	}
	res, err := w.cli.GetService(vo.GetServiceParam{
		ServiceName: w.serviceName,
		GroupName:   w.groupName,
		Clusters:    w.clusters,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*gsvc.Service, 0, len(res.Hosts))
	for _, in := range res.Hosts {
		kind := w.kind
		if k, ok := in.Metadata["kind"]; ok {
			kind = k
		}
		items = append(items, &gsvc.Service{
			ID:        in.InstanceId,
			Name:      res.Name,
			Version:   in.Metadata["version"],
			Metadata:  gconv.Map(in.Metadata),
			Endpoints: []string{fmt.Sprintf("%s://%s:%d", kind, in.Ip, in.Port)},
		})
	}
	return items, nil
}

// Close closes the watcher.
func (w *watcher) Close() error {
	w.cancel()
	return w.cli.Unsubscribe(&vo.SubscribeParam{
		ServiceName: w.serviceName,
		GroupName:   w.groupName,
		Clusters:    w.clusters,
	})
}
