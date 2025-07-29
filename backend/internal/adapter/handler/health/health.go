package health

import "github.com/sirupsen/logrus"

type HealthHandler struct {
	Logger *logrus.Logger
}

func NewHealthHandler(logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		Logger: logger,
	}
}
