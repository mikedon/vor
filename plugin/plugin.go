package plugin

import (
	"github.com/trevor-atlas/vor/env"
	"github.com/trevor-atlas/vor/formatters"
	"github.com/trevor-atlas/vor/rest"
)

type Plugin interface {
	Init(
		loader *env.EnvironmentLoader,
		http rest.RequestBuilder,
		formatter *formatters.StringFormatter)
	PreRun(args ...string) (success bool)
	Run(args ...string) (success bool)
	PostRun(args ...string) (success bool)
}
