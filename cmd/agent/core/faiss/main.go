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

package main

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/core/faiss/config"
	"github.com/vdaas/vald/pkg/agent/core/faiss/usecase"
)

const (
	maxVersion = "v0.0.10"
	minVersion = "v0.0.0"
	name       = "agent faiss"
)

func main() {
	if err := safety.RecoverFunc(func() error {
		return runner.Do(
			context.Background(),
			runner.WithName(name),
			runner.WithVersion(info.Version, maxVersion, minVersion),
			runner.WithConfigLoader(func(path string) (any, *config.GlobalConfig, error) {
				cfg, err := config.NewConfig(path)
				if err != nil {
					return nil, nil, errors.Wrap(err, "failed to load "+name+"'s configuration")
				}
				return cfg, &cfg.GlobalConfig, nil
			}),
			runner.WithDaemonInitializer(func(cfg any) (runner.Runner, error) {
				return usecase.New(cfg.(*config.Data))
			}),
		)
	})(); err != nil {
		log.Fatal(err, info.Get())
		return
	}
}
