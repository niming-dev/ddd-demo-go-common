package gorm

import (
	"errors"

	"gorm.io/gorm"
)

type loggerOption struct {
	Override bool
	Logger   *gormLogger
}

// Apply implement gorm.Option
func (a *loggerOption) Apply(config *gorm.Config) error {
	if a == nil {
		return nil
	} else if config == nil {
		return errors.New("config is nil")
	}

	if config.Logger == nil || a.Override {
		config.Logger = a.Logger
	}
	return nil
}

// AfterInitialize implement gorm.Option
func (a *loggerOption) AfterInitialize(db *gorm.DB) error {
	return nil
}
