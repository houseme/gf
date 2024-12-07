// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package consul implements service Registry and Discovery using consul.
package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
)

const (
	// DefaultTTL is the default TTL for service registration
	DefaultTTL = 20 * time.Second

	// DefaultHealthCheckInterval is the default interval for health check
	DefaultHealthCheckInterval = 10 * time.Second
)

var (
	_ gsvc.Registry = (*Registry)(nil)
)

// Registry implements gsvc.Registry interface using consul.
type Registry struct {
	client  *api.Client       // Consul client
	address string            // Consul address
	options map[string]string // Additional options
}

// Option is the configuration option type for registry.
type Option func(r *Registry)

// WithAddress sets the address for consul client.
func WithAddress(address string) Option {
	return func(r *Registry) {
		r.address = address
	}
}

// WithToken sets the ACL token for consul client.
func WithToken(token string) Option {
	return func(r *Registry) {
		r.options["token"] = token
	}
}

// New creates and returns a new Registry.
func New(opts ...Option) (gsvc.Registry, error) {
	r := &Registry{
		address: "127.0.0.1:8500",
		options: make(map[string]string),
	}

	// Apply options
	for _, opt := range opts {
		opt(r)
	}

	// Create consul config
	config := api.DefaultConfig()
	config.Address = r.address
	if token, ok := r.options["token"]; ok {
		config.Token = token
	}

	// Create consul client
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	r.client = client

	return r, nil
}

// Register registers a service to consul.
func (r *Registry) Register(ctx context.Context, service gsvc.Service) (gsvc.Service, error) {
	metadata, err := json.Marshal(service.GetMetadata())
	if err != nil {
		return nil, gerror.Wrap(err, "failed to marshal service metadata")
	}

	endpoints := service.GetEndpoints()
	if len(endpoints) == 0 {
		return nil, gerror.New("no endpoints found in service")
	}

	// Create service ID
	serviceID := fmt.Sprintf("%s-%s-%s:%d", service.GetName(), service.GetVersion(), endpoints[0].Host(), endpoints[0].Port())

	// Create registration
	reg := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    service.GetName(),
		Tags:    []string{service.GetVersion()},
		Meta:    map[string]string{"metadata": string(metadata)},
		Address: endpoints[0].Host(),
		Port:    endpoints[0].Port(),
	}

	// Add health check
	checkID := fmt.Sprintf("service:%s", serviceID)
	reg.Check = &api.AgentServiceCheck{
		CheckID:                        checkID,
		TTL:                            DefaultTTL.String(),
		DeregisterCriticalServiceAfter: "1m",
	}

	// Register service
	if err := r.client.Agent().ServiceRegister(reg); err != nil {
		return nil, gerror.Wrap(err, "failed to register service")
	}

	// Start TTL health check
	if err := r.client.Agent().PassTTL(checkID, ""); err != nil {
		// Try to deregister service if health check fails
		_ = r.client.Agent().ServiceDeregister(serviceID)
		return nil, gerror.Wrap(err, "failed to pass TTL health check")
	}

	// Start TTL health check goroutine
	go r.ttlHealthCheck(serviceID)

	return service, nil
}

// Deregister deregisters a service from consul.
func (r *Registry) Deregister(ctx context.Context, service gsvc.Service) error {
	endpoints := service.GetEndpoints()
	if len(endpoints) == 0 {
		return gerror.New("no endpoints found in service")
	}

	// Create service ID
	serviceID := fmt.Sprintf("%s-%s-%s:%d", service.GetName(), service.GetVersion(), endpoints[0].Host(), endpoints[0].Port())

	return r.client.Agent().ServiceDeregister(serviceID)
}

// ttlHealthCheck maintains the TTL health check for a service
func (r *Registry) ttlHealthCheck(serviceID string) {
	ticker := time.NewTicker(DefaultHealthCheckInterval)
	defer ticker.Stop()

	checkID := fmt.Sprintf("service:%s", serviceID)
	for range ticker.C {
		if err := r.client.Agent().PassTTL(checkID, ""); err != nil {
			return
		}
	}
}

// GetAddress returns the consul address
func (r *Registry) GetAddress() string {
	return r.address
}

// Watch creates a watcher according to the key prefix.
func (r *Registry) Watch(ctx context.Context, key string) (gsvc.Watcher, error) {
	return newWatcher(r.client, r, key), nil
}