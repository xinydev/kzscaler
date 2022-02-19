// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"strings"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{
		tickMilliseconds: 5 * 1000,
		sched:            NewScheduler(),
	}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	tickMilliseconds uint32
	sched            *Scheduler
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{
		contextID: contextID,
		sched:     ctx.sched,
	}
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	if err := proxywasm.SetTickPeriodMilliSeconds(ctx.tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogWarnf("set tick period milliseconds: %d", ctx.tickMilliseconds)
	return types.OnPluginStartStatusOK
}
func (ctx *pluginContext) OnTick() {
	proxywasm.LogWarnf("sync services")
	if err := ctx.sched.SyncService(); err != nil {
		proxywasm.LogWarnf("sync failed,%s", err)
	}
}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
	sched     *Scheduler
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	err := proxywasm.AddHttpRequestHeader("kzscaler-enabled", "true")
	if err != nil {
		proxywasm.LogCritical("failed to set request header: test")
	}
	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}
	for _, h := range hs {
		proxywasm.LogWarnf("request header --> %s: %s", h[0], h[1])
	}

	return types.ActionContinue
}

type Scheduler struct {
	enabledServices   map[string]bool // service which enabled scale to zero
	zeroStateServices map[string]bool // service which is zero now
	cluster           string
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		enabledServices:   map[string]bool{},
		zeroStateServices: map[string]bool{},
		cluster:           "outbound|80||kzscaler.kzscaler.svc.cluster.local",
	}
}

func (s *Scheduler) SyncService() error {
	extract := func(m map[string]bool) func(string) {
		return func(s string) {
			for _, service := range strings.Split(s, "|") {
				m[service] = true
			}
		}
	}

	// 1. sync enabled service
	err := s.syncRequest("enabled", extract(s.enabledServices))
	if err != nil {
		proxywasm.LogWarnf("get enabled service error,%s", err)
	}
	// 1. sync zero state service
	err = s.syncRequest("zerostate", extract(s.zeroStateServices))
	if err != nil {
		proxywasm.LogWarnf("get zero state service error,%s", err)
	}

	loga, logb := s.printService()
	proxywasm.LogWarnf("sync result,%s;%s", loga, logb)

	return nil
}

func (s *Scheduler) syncRequest(path string, f func(string)) error {
	headers := [][2]string{
		{":method", "GET"},
		{":path", fmt.Sprintf("/%s", path)},
		{":authority", "kzscaler.kzscaler"},
		{":scheme", "http"},
	}

	_, err := proxywasm.DispatchHttpCall(s.cluster,
		headers,
		nil,
		nil,
		1000,
		func(numHeaders, bodySize, numTrailers int) {
			resp, _ := proxywasm.GetHttpCallResponseBody(0, 10000)
			f(string(resp))
			proxywasm.LogWarnf("response:%s", resp)
		},
	)
	return err
}

func (s *Scheduler) printService() (string, string) {
	enabled := make([]string, 0)
	zero := make([]string, 0)
	for k, _ := range s.enabledServices {
		enabled = append(enabled, k)
	}
	for k, _ := range s.zeroStateServices {
		zero = append(zero, k)
	}

	return strings.Join(enabled, ","), strings.Join(zero, ",")

}
