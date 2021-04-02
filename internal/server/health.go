package server

import (
	pbhealth "github.com/kic/health/pkg/proto/health"
)

type HealthTrackingService struct {
	pbhealth.UnimplementedHealthTrackingServer
}

func NewHealthService() *HealthTrackingService {
	return &HealthTrackingService{}
}

