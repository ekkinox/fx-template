package fxgorm

import (
	"context"
	"fmt"

	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"gorm.io/gorm"
)

type GormProbe struct {
	db *gorm.DB
}

func NewGormProbe(db *gorm.DB) *GormProbe {
	return &GormProbe{
		db: db,
	}
}

func (p *GormProbe) Name() string {
	return "gorm"
}

func (p *GormProbe) Check(ctx context.Context) *fxhealthchecker.HealthCheckerProbeResult {
	db, err := p.db.DB()
	if err != nil {
		return fxhealthchecker.NewHealthCheckerProbeResult(false, fmt.Sprintf("database error: %v", err))
	}

	err = db.Ping()
	if err != nil {
		return fxhealthchecker.NewHealthCheckerProbeResult(false, fmt.Sprintf("database ping error: %v", err))
	}

	return fxhealthchecker.NewHealthCheckerProbeResult(true, "database ping success")
}
