package fxgorm

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"gorm.io/gorm"
)

type DbProbe struct {
	db *gorm.DB
}

func NewDbProbe(db *gorm.DB) *DbProbe {
	return &DbProbe{
		db: db,
	}
}

func (p *DbProbe) Name() string {
	return "database"
}

func (p *DbProbe) Check(ctx context.Context) *fxhealthchecker.ProbeResult {

	success := true
	message := "database ping success"

	db, err := p.db.DB()
	if err != nil {
		success = false
		message = fmt.Sprintf("database error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		success = false
		message = fmt.Sprintf("database ping error: %v", err)
	}

	return fxhealthchecker.NewProbeResult(success, message)
}
