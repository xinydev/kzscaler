package main

import (
	"github.com/kzscaler/kzscaler/pkg/reconciler/zeroscaler"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("controller",
		zeroscaler.NewController,
	)
}
