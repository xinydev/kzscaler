package main

import (
	"knative.dev/pkg/injection/sharedmain"

	"github.com/kzscaler/kzscaler/pkg/reconciler/zeroscaler"
)

func main() {
	sharedmain.Main("controller",
		zeroscaler.NewController,
	)
}
