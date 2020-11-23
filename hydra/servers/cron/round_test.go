package cron

import (
	"math"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func TestCounter_Get(t *testing.T) {
	tests := []struct {
		name     string
		executed int
		want     int
	}{
		{name: "Get初始化为0", executed: 0, want: 0},
		{name: "Get初始化<0", executed: -10, want: -10},
		{name: "Get初始化>0 < math.MaxInt32", executed: 10, want: 10},
		{name: "Get初始化 = math.MaxInt32", executed: math.MaxInt32, want: math.MaxInt32},
		{name: "Get初始化 > math.MaxInt32", executed: math.MaxInt32 + 1, want: math.MaxInt32 + 1},
	}
	for _, tt := range tests {
		m := &Counter{executed: tt.executed}
		got := m.Get()
		assert.Equalf(t, tt.want, got, tt.name)
	}
}

func TestCounter_Increase(t *testing.T) {

	tests := []struct {
		name     string
		executed int
		want     int
	}{
		{name: "Increase初始化为0", executed: 0, want: 1},
		{name: "Increase初始化<0", executed: -10, want: -9},
		{name: "Increase初始化>0 < math.MaxInt32", executed: 10, want: 11},
		{name: "Increase初始化 = math.MaxInt32", executed: math.MaxInt32, want: 1},
		{name: "Increase初始化 > math.MaxInt32", executed: math.MaxInt32 + 1, want: 1},
	}
	for _, tt := range tests {
		m := &Counter{executed: tt.executed}
		m.Increase()
		assert.Equalf(t, tt.want, m.Get(), tt.name)
	}
}

func TestRound_Reduce(t *testing.T) {
	tests := []struct {
		name  string
		round int
		want  int
	}{
		{name: "Reduce初始化为0", round: 0, want: -1},
		{name: "Reduce初始化<0", round: -10, want: -11},
		{name: "Reduce初始化>0 ", round: 10, want: 9},
	}
	for _, tt := range tests {
		m := &Round{round: tt.round}
		m.Reduce()
		assert.Equalf(t, tt.want, m.Get(), tt.name)
	}
}

func TestRound_Get(t *testing.T) {
	tests := []struct {
		name  string
		round int
		want  int
	}{
		{name: "RoundGet初始化为0", round: 0, want: 0},
		{name: "RoundGet初始化<0", round: -10, want: -10},
		{name: "RoundGet初始化>0 ", round: 10, want: 10},
	}
	for _, tt := range tests {
		m := &Round{round: tt.round}
		got := m.Get()
		assert.Equalf(t, tt.want, got, tt.name)
	}
}

func TestRound_Update(t *testing.T) {
	tests := []struct {
		name  string
		round int
		v     int
	}{
		{name: "RoundGet初始化为0", round: 0, v: 0},
		{name: "RoundGet初始化<0", round: -10, v: -10},
		{name: "RoundGet初始化>0 ", round: 10, v: 10},
	}
	for _, tt := range tests {
		m := &Round{round: tt.round}
		m.Update(tt.v)
		got := m.Get()
		assert.Equalf(t, tt.v, got, tt.name)
	}
}
