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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/tools/benchmark/operator/handler/rest"
)

type router struct {
	handler rest.Handler
	eg      errgroup.Group
	timeout string
}

// New returns REST route&method information from handler interface.
func New(opts ...Option) http.Handler {
	r := new(router)

	for _, opt := range append(defaultOpts, opts...) {
		opt(r)
	}

	return routing.New(
		routing.WithMiddleware(
			middleware.NewTimeout(
				middleware.WithTimeout(r.timeout),
				middleware.WithErrorGroup(r.eg),
			)),
		routing.WithRoutes([]routing.Route{
			// TODO add REST API interface here
		}...))
}
