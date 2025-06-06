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

package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/core/ngt/handler/rest"
)

type router struct {
	handler rest.Handler
	eg      errgroup.Group
	timeout string
}

// New returns REST route&method information from handler interface.
func New(opts ...Option) http.Handler {
	r := new(router)

	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	h := r.handler

	return routing.New(
		routing.WithMiddleware(
			middleware.NewTimeout(
				middleware.WithTimeout(r.timeout),
				middleware.WithErrorGroup(r.eg),
			)),
		routing.WithRoutes([]routing.Route{
			{
				Name: "Index",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/",
				HandlerFunc: h.Index,
			},
			{
				Name: "Search",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/search",
				HandlerFunc: h.Search,
			},
			{
				Name: "Search By ID",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/id/search",
				HandlerFunc: h.SearchByID,
			},
			{
				Name: "LinearSearch",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/linearsearch",
				HandlerFunc: h.LinearSearch,
			},
			{
				Name: "LinearSearch By ID",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/id/linearsearch",
				HandlerFunc: h.LinearSearchByID,
			},
			{
				Name: "Insert",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/insert",
				HandlerFunc: h.Insert,
			},
			{
				Name: "Multiple Insert",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/insert/multi",
				HandlerFunc: h.MultiInsert,
			},
			{
				Name: "Update",
				Methods: []string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				Pattern:     "/update",
				HandlerFunc: h.Update,
			},
			{
				Name: "Multiple Update",
				Methods: []string{
					http.MethodPost,
					http.MethodPatch,
					http.MethodPut,
				},
				Pattern:     "/update/multi",
				HandlerFunc: h.MultiUpdate,
			},
			{
				Name: "Remove",
				Methods: []string{
					http.MethodDelete,
				},
				Pattern:     "/delete",
				HandlerFunc: h.Remove,
			},
			{
				Name: "Multiple Remove",
				Methods: []string{
					http.MethodDelete,
					http.MethodPost,
				},
				Pattern:     "/delete/multi",
				HandlerFunc: h.MultiRemove,
			},
			{
				Name: "Create Index",
				Methods: []string{
					http.MethodPost,
				},
				Pattern:     "/index/create",
				HandlerFunc: h.CreateIndex,
			},
			{
				Name: "Save Index",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/index/save",
				HandlerFunc: h.SaveIndex,
			},
			{
				Name: "GetObject",
				Methods: []string{
					http.MethodGet,
				},
				Pattern:     "/object/{id}",
				HandlerFunc: h.GetObject,
			},
		}...))
}
