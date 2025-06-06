//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package grpc

import (
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type Option func(*server)

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithReplicationCount(3),
	WithStreamConcurrency(runtime.GOMAXPROCS(-1) * 10),
	WithMultiConcurrency(runtime.GOMAXPROCS(-1) * 10),
	WithTimeout("5s"),
	WithName(func() string {
		name, err := os.Hostname()
		if err != nil {
			log.Warn(err)
		}
		return name
	}()),
	WithIP(net.LoadLocalIP()),
}

// WithIP returns the option to set the IP for server.
func WithIP(ip string) Option {
	return func(s *server) {
		if len(ip) != 0 {
			s.ip = ip
		}
	}
}

// WithName returns the option to set the name for server.
func WithName(name string) Option {
	return func(s *server) {
		if len(name) != 0 {
			s.name = name
		}
	}
}

func WithGateway(g service.Gateway) Option {
	return func(s *server) {
		if g != nil {
			s.gateway = g
		}
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(s *server) {
		if eg != nil {
			s.eg = eg
		}
	}
}

func WithTimeout(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 10
		}
		s.timeout = d
	}
}

func WithReplicationCount(rep int) Option {
	return func(s *server) {
		if rep > 0 {
			s.replica = rep
		}
	}
}

func WithStreamConcurrency(c int) Option {
	return func(s *server) {
		if c > 1 {
			s.streamConcurrency = c
		}
	}
}

func WithMultiConcurrency(c int) Option {
	return func(s *server) {
		if c > 1 {
			s.multiConcurrency = c
		}
	}
}
