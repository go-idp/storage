package main

import (
	"github.com/go-idp/storage"
	"github.com/go-idp/storage/config"
	"github.com/go-idp/storage/server"
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "storage",
		Usage:   "Storage Service for IDP",
		Version: storage.Version,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:    "base-dir",
				Usage:   "server base dir",
				Aliases: []string{"b"},
				EnvVars: []string{"BASE_DIR"},
			},
			&cli.StringFlag{
				Name:     "oss-access-key-id",
				Usage:    "OSS Access Key ID",
				EnvVars:  []string{"OSS_ACCESS_KEY_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "oss-access-key-secret",
				Usage:    "OSS Acess Key Secret",
				EnvVars:  []string{"OSS_ACCESS_KEY_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "oss-bucket",
				Usage:    "OSS Bucket",
				EnvVars:  []string{"OSS_BUCKET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "oss-endpoint",
				Usage:    "OSS Endpoint",
				EnvVars:  []string{"OSS_ENDPOINT"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "oss-base-dir",
				Usage:   "OSS Base Dir",
				EnvVars: []string{"OSS_BASE_DIR"},
			},
		},
	})

	app.Command(func(ctx *cli.Context) (err error) {
		return server.Run(&config.Config{
			Port:               ctx.Int("port"),
			BaseDir:            ctx.String("base-dir"),
			OSSAccessKeyID:     ctx.String("oss-access-key-id"),
			OSSAccessKeySecret: ctx.String("oss-access-key-secret"),
			OSSBucket:          ctx.String("oss-bucket"),
			OSSEndpoint:        ctx.String("oss-endpoint"),
			OSSBaseDir:         ctx.String("oss-base-dir"),
		})
	})

	app.Run()
}
