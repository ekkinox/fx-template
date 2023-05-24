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
	message := "db probe ping success"

	db, err := p.db.DB()
	if err != nil {
		success = false
		message = fmt.Sprintf("db probe error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		success = false
		message = fmt.Sprintf("db probe ping error: %v", err)
	}

	return fxhealthchecker.NewProbeResult(success, message)
}
