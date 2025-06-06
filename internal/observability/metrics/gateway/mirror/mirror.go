// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mirror

import (
	"context"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	MetricsName        = "gateway_mirror_connecting_target"
	MetricsDescription = "Target to which the mirror gateway is connecting"

	targetAddrKey = "addr"
)

type mirrorMetrics struct {
	m service.Mirror
}

func New(m service.Mirror) metrics.Metric {
	return &mirrorMetrics{
		m: m,
	}
}

func (*mirrorMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        MetricsName,
				Description: MetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (mm *mirrorMetrics) Register(m metrics.Meter) error {
	targetCount, err := m.Int64ObservableGauge(
		MetricsName,
		metrics.WithDescription(MetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			mm.m.RangeMirrorAddr(func(addr string, _ any) bool {
				o.ObserveInt64(targetCount, 1, api.WithAttributes(attribute.String(targetAddrKey, addr)))
				return true
			})
			return nil
		},
		targetCount,
	)
	return err
}
