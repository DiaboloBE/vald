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

package rest

import (
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
	"github.com/vdaas/vald/pkg/agent/core/ngt/handler/grpc"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	Exists(w http.ResponseWriter, r *http.Request) (int, error)
	Search(w http.ResponseWriter, r *http.Request) (int, error)
	SearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	Insert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiInsert(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpdate(w http.ResponseWriter, r *http.Request) (int, error)
	Remove(w http.ResponseWriter, r *http.Request) (int, error)
	MultiRemove(w http.ResponseWriter, r *http.Request) (int, error)
	CreateIndex(w http.ResponseWriter, r *http.Request) (int, error)
	SaveIndex(w http.ResponseWriter, r *http.Request) (int, error)
	CreateAndSaveIndex(w http.ResponseWriter, r *http.Request) (int, error)
	GetObject(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	agent grpc.Server
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (*handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]any)
	return json.Handler(w, r, &data, func() (any, error) {
		return dump.Request(nil, data, r)
	})
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.Search(r.Context(), req)
	})
}

func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.SearchByID(r.Context(), req)
	})
}

func (h *handler) LinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.LinearSearch(r.Context(), req)
	})
}

func (h *handler) LinearSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.LinearSearchByID(r.Context(), req)
	})
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.Insert(r.Context(), req)
	})
}

func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.MultiInsert(r.Context(), req)
	})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.Update(r.Context(), req)
	})
}

func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.MultiUpdate(r.Context(), req)
	})
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.Remove(r.Context(), req)
	})
}

func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.MultiRemove(r.Context(), req)
	})
}

func (h *handler) CreateIndex(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Control_CreateIndexRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.CreateIndex(r.Context(), req)
	})
}

func (h *handler) SaveIndex(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Empty
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.SaveIndex(r.Context(), req)
	})
}

func (h *handler) CreateAndSaveIndex(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Control_CreateIndexRequest
	return json.Handler(w, r, &req, func() (any, error) {
		_, err = h.agent.CreateIndex(r.Context(), req)
		if err != nil {
			return nil, err
		}
		return h.agent.SaveIndex(r.Context(), nil)
	})
}

func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_VectorRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.GetObject(r.Context(), req)
	})
}

func (h *handler) Exists(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, &req, func() (any, error) {
		return h.agent.Exists(r.Context(), req)
	})
}
