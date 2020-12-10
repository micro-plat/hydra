// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package go2sky

import (
	"math/rand"
	"time"
)

type Sampler interface {
	IsSampled(operation string) (sampled bool)
}

type ConstSampler struct {
	decision bool
}

// NewConstSampler creates a ConstSampler.
func NewConstSampler(sample bool) *ConstSampler {
	s := &ConstSampler{
		decision: sample,
	}
	return s
}

// IsSampled implements IsSampled() of Sampler.
func (s *ConstSampler) IsSampled(operation string) bool {
	return s.decision
}

type RandomSampler struct {
	samplingRate float64
	rand         *rand.Rand
	threshold    int
}

// IsSampled implements IsSampled() of Sampler.
func (s *RandomSampler) IsSampled(operation string) bool {
	return s.threshold >= s.rand.Intn(100)
}

func (s *RandomSampler) init() {
	s.rand = rand.New(rand.NewSource(time.Now().Unix()))
	s.threshold = int(s.samplingRate * 100)
}

func NewRandomSampler(samplingRate float64) *RandomSampler {
	s := &RandomSampler{
		samplingRate: samplingRate,
	}
	s.init()
	return s
}
