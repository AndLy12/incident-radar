package config

import (
	"fmt"
	"time"
)

func validate(cfg *Config) error {
	if cfg.path == "" {
		return fmt.Errorf("path is required")
	}
	var fromTime, toTime time.Time
	var err error

	layout := "2006-01-02"

	if cfg.fromDate != "" {
		fromTime, err = time.Parse(layout, cfg.fromDate)
		if err != nil {
			return err
		}
	}

	if cfg.toDate != "" {
		toTime, err = time.Parse(layout, cfg.toDate)
		if err != nil {
			return err
		}
	}

	if !fromTime.IsZero() && !toTime.IsZero() && !fromTime.Before(toTime) {
		return fmt.Errorf("from date must be before to date")
	}

	return nil
}
