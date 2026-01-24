package core

import (
	"mdnav/internal/conf"
	"mdnav/internal/pkg/zap"
)

type Context struct {
	Log  zap.Logger
	Conf conf.Config
}
