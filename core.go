package main

import (
	"io"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/headers"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

type Config struct {
	Port            int
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	Endpoint        string
	//
	Directory string
}

func server(cfg *Config) error {
	app := defaults.Default()

	client, err := oss.New(
		cfg.Endpoint,
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
	)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return err
	}

	// @TODO
	// app.Get("/(.*)", func(ctx *zoox.Context) {
	app.Use(func(ctx *zoox.Context) {
		if ctx.Path == "" {
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}

		filepath := ctx.Path[1:]
		if filepath == "" {
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}

		fullpath := fs.JoinPath(cfg.Directory, filepath)

		ctx.Logger.Infof("oss file: %s", fullpath)
		reader, err := bucket.GetObject(fullpath)
		if err != nil {
			ctx.Logger.Errorf("failed to get file path: %s (err: %s)", fullpath, err)
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}
		defer reader.Close()

		ctx.SetCacheControlWithMaxAge(365 * 24 * time.Hour)
		ctx.Set(headers.Vary, "origin")

		if _, err := io.Copy(ctx.Writer, reader); err != nil {
			ctx.Logger.Errorf("failed to send file reader: %s (err: %s)", fullpath, err)
		}
	})

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}
