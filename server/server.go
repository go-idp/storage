package server

import (
	"io"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-idp/storage"
	"github.com/go-idp/storage/config"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/datetime"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/headers"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

func Run(cfg *config.Config) error {
	logger.Infof("server config: %+v", cfg)

	runningAt := datetime.Now().Format("YYYY-MM-DD HH:mm:ss")

	// BaseDir is the base directory of the server
	//		/ => ""
	//		/abc/ => /abc
	// remove suffix slash
	if matched := regexp.Match("/$", cfg.BaseDir); matched {
		cfg.BaseDir = cfg.BaseDir[:len(cfg.BaseDir)-1]
	}

	// OSSBaseDir is the base directory of the OSS
	//		/ => ""
	//		/abc => abc
	//	 /abc/ => abc/
	// remove prefix slash
	if matched := regexp.Match("^/", cfg.OSSBaseDir); matched {
		cfg.OSSBaseDir = cfg.OSSBaseDir[1:]
	}

	app := defaults.Default()

	client, err := oss.New(
		cfg.OSSEndpoint,
		cfg.OSSAccessKeyID,
		cfg.OSSAccessKeySecret,
	)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(cfg.OSSBucket)
	if err != nil {
		return err
	}

	counts := safe.NewInt64()
	app.Get(fmt.Sprintf("%s/*", cfg.BaseDir), func(ctx *zoox.Context) {
		if ctx.Path == "" {
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}

		filepath := ctx.Path[len(cfg.BaseDir)+1:]
		if filepath == "" {
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}

		counts.Inc(1)

		osspath := fs.JoinPath(cfg.OSSBaseDir, filepath)
		ctx.Logger.Infof("match: %s => oss:%s", ctx.Path, osspath)

		reader, err := bucket.GetObject(osspath)
		if err != nil {
			ctx.Logger.Errorf("failed to get file path: %s (osspath: %s)", err, osspath)
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}
		defer reader.Close()

		ctx.SetCacheControlWithMaxAge(365 * 24 * time.Hour)
		ctx.Set(headers.Vary, "origin")

		if _, err := io.Copy(ctx.Writer, reader); err != nil {
			ctx.Logger.Errorf("failed to read: %s (osspath: %s)", err, osspath)
		}
	})

	app.Get("/", func(ctx *zoox.Context) {
		ctx.JSON(200, zoox.H{
			"name":    "storage service for idp",
			"version": storage.Version,
			"status": map[string]any{
				"counts":     counts.Get(),
				"running_at": runningAt,
			},
		})
	})

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}
