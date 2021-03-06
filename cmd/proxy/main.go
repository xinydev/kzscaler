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
	"strconv"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
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
		tickMilliseconds: 10 * 1000,
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
	m := proxywasm.DefineCounterMetric("kzscaler")
	return &httpHeaders{
		contextID: contextID,
		sched:     ctx.sched,
		reqCnt:    &m,
	}
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	// read config
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}
	proxywasm.LogWarnf("plugin config: %s", string(data))

	configs := strings.Split(string(data), "&")

	// we are in INBOUND proxy,only export metric
	if len(configs) == 0 {
		return types.OnPluginStartStatusOK
	}

	// we are in OUTBOUND proxy
	ctx.sched.SetCluster(configs[0])
	if err := proxywasm.SetTickPeriodMilliSeconds(ctx.tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogWarnf("set tick period milliseconds: %d", ctx.tickMilliseconds)

	return types.OnPluginStartStatusOK
}
func (ctx *pluginContext) OnTick() {
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
	reqCnt    *proxywasm.MetricCounter
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	err := proxywasm.AddHttpRequestHeader("kzscaler-enabled", "true")
	if err != nil {
		proxywasm.LogCritical("failed to set request header: test")
	}
	authority, _ := proxywasm.GetHttpRequestHeader(":authority")
	proxywasm.LogWarnf("request auth:%s", authority)

	// export metrics
	ctx.reqCnt.Increment(1)

	act, err := ctx.sched.RequestService(authority, ctx.contextID)

	return act
}

type Scheduler struct {
	cluster  string
	services map[string]int
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		services: map[string]int{},
		//cluster:           "outbound|80||kzscaler.kzscaler.svc.cluster.local",
		cluster: "mock_service",
	}
}

// SetCluster set envoy cluster for requesting
func (s *Scheduler) SetCluster(c string) {
	s.cluster = c
}

func (s *Scheduler) SyncService() error {
	return makeRequest(
		s.cluster,
		"/service",
		"kzscaler.kzscaler",
		func(bytes []byte) {
			// envoy wasm does not support json
			// services:  service1%10&service2%10
			for _, svc := range strings.Split(string(bytes), "&") {
				svcParts := strings.Split(svc, "%")
				cnt, _ := strconv.Atoi(svcParts[1])
				s.services[svcParts[0]] = cnt
			}
		})
}

func (s *Scheduler) RequestService(name string, cid uint32) (types.Action, error) {
	// check service is zero-scale enabled or not

	name = strings.Split(name, ":")[0]

	if v, ok := s.services[name]; ok {
		if v == 0 {
			// need to call scale up first
			proxywasm.LogWarnf("scale up first,%s", name)

			err := makeRequest(
				s.cluster,
				fmt.Sprintf("/scale_up/%s", name),
				"kzscaler.kzscaler",
				func(bytes []byte) {
					if err := proxywasm.SetEffectiveContext(cid); err != nil {
						proxywasm.LogCriticalf("kzscaler callback set error:%s", err)

					}
					if err := proxywasm.ResumeHttpRequest(); err != nil {
						proxywasm.LogCriticalf("kzscaler callback resume error:%s", err)
					}
				})
			return types.ActionPause, err
		}
	}

	return types.ActionContinue, nil
}

func makeRequest(cluster, path, authority string, f func([]byte)) error {
	headers := [][2]string{
		{":method", "GET"},
		{":path", path},
		{":authority", authority},
		{":scheme", "http"},
	}

	_, err := proxywasm.DispatchHttpCall(cluster,
		headers,
		nil,
		nil,
		1000*15,
		func(numHeaders, bodySize, numTrailers int) {
			resp, _ := proxywasm.GetHttpCallResponseBody(0, 10000)
			f(resp)
		},
	)
	if err != nil {
		proxywasm.LogErrorf("request error:%s/n", cluster, path, authority)
	}
	return err
}
