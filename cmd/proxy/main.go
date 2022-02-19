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
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{contextID: contextID}
}

func (*pluginContext) OnTick() {

}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
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
	proxywasm.LogWarnf("http call 1")
	err = ctx.httpCall1()
	if err != nil {
		proxywasm.LogWarnf("request1 error,%s", err)
	}

	proxywasm.LogWarnf("http call 2")
	err = ctx.httpCall2()
	if err != nil {
		proxywasm.LogWarnf("request2 error,%s", err)
	}

	return types.ActionContinue
}

func (ctx *httpHeaders) httpCall1() error {
	headers := [][2]string{
		{":method", "GET"},
		{":path", "/"},
		{":authority", "kzscaler.kzscaler"},
		{":scheme", "http"},
	}

	_, err := proxywasm.DispatchHttpCall("outbound|80||kzscaler.kzscaler.svc.cluster.local",
		headers,
		nil,
		nil,
		1000,
		func(numHeaders, bodySize, numTrailers int) {
			resp, _ := proxywasm.GetHttpCallResponseBody(0, 10000)
			r := string(resp)
			proxywasm.LogDebugf("APISERVER RESPONSE %v", r)
		},
	)
	return err
}

func (ctx *httpHeaders) httpCall2() error {
	headers := [][2]string{
		{":method", "GET"},
		{":path", "/"},
		{":authority", "kzscaler.kzscaler"},
		{":scheme", "http"},
	}

	_, err := proxywasm.DispatchHttpCall("istio-ingressgateway.istio-system",
		headers,
		nil,
		nil,
		1000,
		func(numHeaders, bodySize, numTrailers int) {
			resp, _ := proxywasm.GetHttpCallResponseBody(0, 10000)
			r := string(resp)
			proxywasm.LogDebugf("APISERVER RESPONSE %v", r)
		},
	)
	return err
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	hs, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get response headers: %v", err)
	}
	_ = proxywasm.AddHttpResponseHeader("fffff", "sdsdsddsdssd")
	for _, h := range hs {
		proxywasm.LogWarnf("response header <-- %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
