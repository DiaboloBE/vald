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

package logging

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	grpcKindUnary  = "unary"
	grpcKindStream = "stream"

	rpcCompletedMessage = "rpc completed"
	rpcFailedMessage    = "rpc failed"
)

type AccessLogEntity struct {
	GRPC      *AccessLogGRPCEntity `json:"grpc,omitempty"      yaml:"grpc"`
	StartTime int64                `json:"startTime,omitempty" yaml:"startTime"`
	EndTime   int64                `json:"endTime,omitempty"   yaml:"endTime"`
	Latency   int64                `json:"latency,omitempty"   yaml:"latency"`
	TraceID   string               `json:"traceID,omitempty"   yaml:"traceID"`
	Error     error                `json:"error,omitempty"     yaml:"error"`
}

type AccessLogGRPCEntity struct {
	Kind    string `json:"kind,omitempty"    yaml:"kind"`
	Service string `json:"service,omitempty" yaml:"service"`
	Method  string `json:"method,omitempty"  yaml:"method"`
}

func (e AccessLogEntity) String() (str string) {
	var emsg string
	if e.Error != nil {
		st, ok := status.FromError(e.Error)
		if ok && st != nil {
			emsg = st.String()
		} else {
			emsg = e.Error.Error()
		}
	}
	eb, err := json.Marshal(e)
	if err != nil {
		str = fmt.Sprintf("%#v,\tfailed to json.Marshal(AccessLogEntity) error: %v", e, err)
	} else {
		str = string(eb)
	}

	if emsg != "" {
		return str + ",\terror message: " + emsg
	}
	return str
}

func AccessLogInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		var traceID string

		span := trace.FromContext(ctx)
		if span != nil {
			traceID = span.SpanContext().TraceID().String()
		}

		start := time.Now()

		resp, err = handler(ctx, req)

		end := time.Now()

		service, method := parseMethod(info.FullMethod)

		entity := &AccessLogEntity{
			GRPC: &AccessLogGRPCEntity{
				Kind:    grpcKindUnary,
				Service: service,
				Method:  method,
			},
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Latency:   end.Sub(start).Nanoseconds(),
		}

		if traceID != "" {
			entity.TraceID = traceID
		}

		if err != nil {
			entity.Error = err
			log.Warn(rpcFailedMessage, entity.String())
		} else {
			log.Debug(rpcCompletedMessage, entity.String())
		}

		return resp, err
	}
}

func AccessLogStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		var traceID string

		span := trace.FromContext(ss.Context())
		if span != nil {
			traceID = span.SpanContext().TraceID().String()
		}

		start := time.Now()

		err := handler(srv, ss)

		end := time.Now()

		service, method := parseMethod(info.FullMethod)

		entity := &AccessLogEntity{
			GRPC: &AccessLogGRPCEntity{
				Kind:    grpcKindStream,
				Service: service,
				Method:  method,
			},
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Latency:   end.Sub(start).Nanoseconds(),
		}

		if traceID != "" {
			entity.TraceID = traceID
		}

		if err != nil {
			entity.Error = err
			log.Warn(rpcFailedMessage, entity.String())
		} else {
			log.Debug(rpcCompletedMessage, entity.String())
		}

		return err
	}
}

func parseMethod(fullMethod string) (service, method string) {
	service = path.Dir(fullMethod)[1:]
	method = path.Base(fullMethod)

	return service, method
}
