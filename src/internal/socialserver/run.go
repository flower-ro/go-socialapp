package socialserver

import (
	"go-socialapp/internal/socialserver/config"
)

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := createSocialServer(cfg)
	if err != nil {
		return err
	}
	return server.PrepareRun().Run()
}
