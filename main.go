package main

import (
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "alioss-cdn",
		Usage:   "Ali OSS CDN Server",
		Version: Version,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:     "access-key-id",
				Usage:    "OSS Acess Key ID",
				EnvVars:  []string{"ACCESS_KEY_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "access-key-secret",
				Usage:    "OSS Acess Key Secret",
				EnvVars:  []string{"ACCESS_KEY_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "OSS Bucket",
				EnvVars:  []string{"BUCKET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "endpoint",
				Usage:    "OSS Endpoint",
				EnvVars:  []string{"ENDPOINT"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "directory",
				Usage:   "OSS Directory",
				EnvVars: []string{"DIRECTORY"},
			},
		},
	})

	app.Command(func(ctx *cli.Context) (err error) {
		return server(&Config{
			Port:            ctx.Int("port"),
			AccessKeyID:     ctx.String("access-key-id"),
			AccessKeySecret: ctx.String("access-key-secret"),
			Bucket:          ctx.String("bucket"),
			Endpoint:        ctx.String("endpoint"),
			Directory:       ctx.String("directory"),
		})
	})

	app.Run()
}
