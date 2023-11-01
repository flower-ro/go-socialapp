package socialserver

import (
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/app"
	"go-socialapp/internal/socialserver/config"
	"go-socialapp/internal/socialserver/options"
)

const commandDesc = `The socialserver used for kind of social app.`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("social Server app", //应用简短描述
		basename, //二进制文件名
		app.WithOptions(opts),
		app.WithDescription(commandDesc), //应用的详细描述。
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {

		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
