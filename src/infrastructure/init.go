package infrastructure

import "go.uber.org/zap"

type Infrastructure struct {
	Logger *zap.Logger
}

func NewInfrastructure(logger *zap.Logger) *Infrastructure {
	return &Infrastructure{
		Logger: logger,
	}
}
