// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package polaris Configuration management function of Polaris
package polaris

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcfg"
)

const (
	// DefaultInstanceName is the default instance name for the configuration.
	DefaultInstanceName = "config" // DefaultName is the default instance name for instance usage.
	// DefaultConfigFileName is the default configuration file name.
	DefaultConfigFileName = "config" // DefaultConfigFile is the default configuration file name.
)

var (
	_                      gcfg.Adapter = &AdapterPolaris{}
	localInstances                      = gmap.NewStrAnyMap(true) // Instances map containing configuration instances.
	customConfigContentMap              = gmap.NewStrStrMap(true) // Customized configuration content.
)

// Config is the configuration management object.
type Config struct {
	adapter gcfg.Adapter
}

// SetAdapter sets the adapter of current Config object.
func (c *Config) SetAdapter(adapter gcfg.Adapter) {
	c.adapter = adapter
}

// GetAdapter returns the adapter of current Config object.
func (c *Config) GetAdapter() gcfg.Adapter {
	return c.adapter
}

// Available checks and returns the configuration service is available.
// The optional parameter `pattern` specifies certain configuration resource.
//
// It returns true if configuration file is present in default AdapterFile, or else false.
// Note that this function does not return error as it just does simply check for backend configuration service.
func (c *Config) Available(ctx context.Context, resource ...string) (ok bool) {
	var usedResource string
	if len(resource) > 0 {
		usedResource = resource[0]
	}
	return c.adapter.Available(ctx, usedResource)
}

// Get retrieves and returns value by specified `pattern`.
// It returns all values of current Json object if `pattern` is given empty or string ".".
// It returns nil if no value found by `pattern`.
//
// It returns a default value specified by `def` if value for `pattern` is not found.
func (c *Config) Get(ctx context.Context, pattern string, def ...interface{}) (*gvar.Var, error) {
	var (
		err   error
		value interface{}
	)
	value, err = c.adapter.Get(ctx, pattern)
	if err != nil {
		return nil, err
	}
	if value == nil {
		if len(def) > 0 {
			return gvar.New(def[0]), nil
		}
		return nil, nil
	}
	return gvar.New(value), nil
}

// Data retrieves and returns all configuration data as map type.
func (c *Config) Data(ctx context.Context) (data map[string]interface{}, err error) {
	return c.adapter.Data(ctx)
}

// MustGet acts as function Get, but it panics if error occurs.
func (c *Config) MustGet(ctx context.Context, pattern string, def ...interface{}) *gvar.Var {
	v, err := c.Get(ctx, pattern, def...)
	if err != nil {
		panic(err)
	}
	if v == nil {
		return nil
	}
	return v
}
