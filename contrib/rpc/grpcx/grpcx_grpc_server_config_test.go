// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Grpcx_Grpc_Server(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := Server.New()
		s.Start()
		defer s.Stop()
		time.Sleep(time.Millisecond * 100)
		t.Assert(len(s.services) != 0, true)
	})
}

func Test_Grpcx_Grpc_Server_Address(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := Server.NewConfig()
		c.Address = "127.0.0.1:0"
		s := Server.New(c)
		s.Start()
		defer s.Stop()
		time.Sleep(time.Millisecond * 100)
		t.Assert(len(s.services) != 0, true)
		t.Assert(gstr.Contains(s.services[0].GetEndpoints().String(), "127.0.0.1:"), true)
	})
}

func Test_Grpcx_Grpc_Server_Config(t *testing.T) {
	cfg := Server.NewConfig()
	addr := "10.0.0.29:80"
	cfg.Endpoints = []string{
		addr,
	}
	// cfg set one endpoint
	gtest.C(t, func(t *gtest.T) {
		s := Server.New(cfg)
		s.doServiceRegister()
		for _, svc := range s.services {
			t.Assert(svc.GetEndpoints().String(), addr)
		}
	})
	// cfg set more endpoints
	addr = "10.0.0.29:80,10.0.0.29:81"
	cfg.Endpoints = []string{
		"10.0.0.29:80",
		"10.0.0.29:81",
	}
	gtest.C(t, func(t *gtest.T) {
		s := Server.New(cfg)
		s.doServiceRegister()
		for _, svc := range s.services {
			t.Assert(svc.GetEndpoints().String(), addr)
		}
	})
}
