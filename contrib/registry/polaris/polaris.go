// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package polaris implements service Registry and Discovery using polaris.
package polaris

import (
	"fmt"
	"strings"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	_ gsvc.Registrar = &Registry{}
)

// _instanceIDSeparator Instance id Separator.
const _instanceIDSeparator = "-"

type options struct {
	// required, namespace in polaris
	Namespace string

	// required, service access token
	ServiceToken string

	// optional, protocol in polaris. Default value is nil, it means use protocol config in service
	Protocol *string

	// service weight in polaris. Default value is 100, 0 <= weight <= 10000
	Weight int

	// service priority. Default value is 0. The smaller the value, the lower the priority
	Priority int

	// To show service is healthy or not. Default value is True.
	Healthy bool

	// Heartbeat enable .Not in polaris . Default value is True.
	Heartbeat bool

	// To show service is isolate or not. Default value is False.
	Isolate bool

	// TTL timeout. if node needs to use heartbeat to report,required. If not set,server will throw ErrorCode-400141
	TTL int

	// optional, Timeout for single query. Default value is global config
	// Total is (1+RetryCount) * Timeout
	Timeout time.Duration

	// optional, retry count. Default value is global config
	RetryCount int
}

// Option The option is a polaris option.
type Option func(o *options)

// Registry is polaris registry.
type Registry struct {
	opt      options
	provider api.ProviderAPI
	consumer api.ConsumerAPI
}

// WithNamespace with the Namespace option.
func WithNamespace(namespace string) Option {
	return func(o *options) { o.Namespace = namespace }
}

// WithServiceToken with ServiceToken option.
func WithServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithProtocol with the Protocol option.
func WithProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithWeight with the Weight option.
func WithWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithHealthy with the Healthy option.
func WithHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithIsolate with the Isolate option.
func WithIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithTTL with the TTL option.
func WithTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithTimeout the Timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithRetryCount with RetryCount option.
func WithRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

// WithHeartbeat with the Heartbeat option.
func WithHeartbeat(heartbeat bool) Option {
	return func(o *options) { o.Heartbeat = heartbeat }
}

// NewRegistry create a new registry.
func NewRegistry(provider api.ProviderAPI, consumer api.ConsumerAPI, opts ...Option) (r *Registry) {
	op := options{
		Namespace:    "default",
		ServiceToken: "",
		Protocol:     nil,
		Weight:       0,
		Priority:     0,
		Healthy:      true,
		Heartbeat:    true,
		Isolate:      false,
		TTL:          0,
		Timeout:      0,
		RetryCount:   0,
	}
	for _, option := range opts {
		option(&op)
	}
	return &Registry{
		opt:      op,
		provider: provider,
		consumer: consumer,
	}
}

// NewRegistryWithConfig new a registry with config.
func NewRegistryWithConfig(conf config.Configuration, opts ...Option) (r *Registry) {
	provider, err := api.NewProviderAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	consumer, err := api.NewConsumerAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	return NewRegistry(provider, consumer, opts...)
}

func instancesToServiceInstances(instances []model.Instance) []*gsvc.Service {
	serviceInstances := make([]*gsvc.Service, 0, len(instances))
	for _, instance := range instances {
		if instance.IsHealthy() {
			serviceInstances = append(serviceInstances, instanceToServiceInstance(instance))
		}
	}
	return serviceInstances
}

func instanceToServiceInstance(instance model.Instance) *gsvc.Service {
	metadata := instance.GetMetadata()
	// Usually, it won't fail in goframe if register correctly
	kind := ""
	if k, ok := metadata["kind"]; ok {
		kind = k
	}

	name := ""
	names := strings.Split(instance.GetService(), _instanceIDSeparator)
	if names != nil && len(names) > 4 {
		return &gsvc.Service{
			Prefix:     names[0],
			Deployment: names[1],
			Namespace:  names[2],
			Name:       names[3],
			Version:    metadata["version"],
			Metadata:   gconv.Map(metadata),
			Endpoints:  []string{fmt.Sprintf("%s:%d", instance.GetHost(), instance.GetPort())},
			Separator:  _instanceIDSeparator,
		}
	}

	return &gsvc.Service{
		Name:      name,
		Version:   metadata["version"],
		Metadata:  gconv.Map(metadata),
		Endpoints: []string{fmt.Sprintf("%s://%s:%d", kind, instance.GetHost(), instance.GetPort())},
		Separator: _instanceIDSeparator,
	}
}
