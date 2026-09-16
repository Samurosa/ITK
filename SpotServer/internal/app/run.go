package app

import (
	"ITK_Code/m/v2/internal/config"
	"fmt"
)

func Run(cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {

		return fmt.Errorf(
			"error loading config file path: %s, error: %s\n",
			cfgPath,
			err,
		)
	}

	application, err := New(cfg)
	if err != nil {

		return fmt.Errorf(
			"create application failed: %s\n",
			err,
		)
	}

	application.Start()

	application.WaitSignal()

	application.Stop()

	return nil
}
