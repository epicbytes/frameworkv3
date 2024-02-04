package v1

import (
	"github.com/epicbytes/frameworkv3/v1/config"
	"github.com/epicbytes/frameworkv3/v1/context"
	"github.com/epicbytes/frameworkv3/v1/logger"
	"go.uber.org/fx"
)

var StandardModules = fx.Options(
	context.NewModule(),
	logger.NewModule(),
	config.NewModule(),
)
