// Copyright 2018 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package balancer implements client balancer.
package balancer

import (
	"github.com/zhiyunliu/etcd/clientv3/balancer/picker"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	_ "google.golang.org/grpc/resolver/dns"         // register DNS resolver
	_ "google.golang.org/grpc/resolver/passthrough" // register passthrough resolver
)

// Config defines balancer configurations.
type Config struct {
	// Policy configures balancer policy.
	Policy picker.Policy

	// Picker implements gRPC picker.
	// Leave empty if "Policy" field is not custom.
	// TODO: currently custom policy is not supported.
	// Picker picker.Picker

	// Name defines an additional name for balancer.
	// Useful for balancer testing to avoid register conflicts.
	// If empty, defaults to policy name.
	Name string

	// Logger configures balancer logging.
	// If nil, logs are discarded.
	Logger *zap.Logger
}

// RegisterBuilder creates and registers a builder. Since this function calls balancer.Register, it
// must be invoked at initialization time.
func RegisterBuilder(cfg Config) {
	bb := newBuilder(cfg)
	balancer.Register(bb)

	cfg.Logger.Debug(
		"registered balancer",
		zap.String("policy", cfg.Policy.String()),
		zap.String("name", cfg.Name),
	)
}

func newBuilder(cfg Config) balancer.Builder {
	builder := base.NewBalancerBuilder(cfg.Name, picker.NewPickBuilder(cfg.Policy), base.Config{HealthCheck: true})
	return builder
}
