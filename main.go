package main

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/http"
	"github.com/baetyl/baetyl-go/v2/log"
	"gocv.io/x/gocv"
)

func main() {
	context.Run(func(ctx context.Context) error {
		var cfg Config
		err := ctx.LoadCustomConfig(&cfg)
		if err != nil {
			return errors.Trace(err)
		}

		bcfg, err := ctx.NewSystemBrokerClientConfig()
		if err != nil {
			return errors.Trace(err)
		}
		bcli, err := ctx.NewBrokerClient(bcfg)
		if err != nil {
			return errors.Trace(err)
		}
		bcli.Start(nil)
		defer bcli.Close()

		infer, err := NewInfer(cfg.Process.Infer)
		if err != nil {
			return errors.Trace(err)
		}
		defer infer.Close()

		video, err := NewVideo(cfg.Video)
		if err != nil {
			return errors.Trace(err)
		}
		defer video.Close()

		fcli, err := ctx.NewFunctionHttpClient()
		if err != nil {
			return errors.Trace(err)
		}

		var s time.Time
		img := gocv.NewMat()
		defer img.Close()
		proc := NewProcess(cfg.Process, bcli, fcli)
		for {
			select {
			case <-ctx.WaitChan():
				return nil
			default:
			}
			s = time.Now()
			err := video.Read(&img)
			if err != nil {
				ctx.Log().Error("failed to read video", log.Error(err))
				time.Sleep(time.Second) // TODO: configured
				continue
			}
			if img.Empty() {
				continue
			}
			ctx.Log().Sugar().Debugf("[Read  ] elapsed time: %v", time.Since(s))

			t := time.Now()
			blob := proc.Before(img)
			ctx.Log().Sugar().Debugf("[Before] elapsed time: %v", time.Since(t))

			s = time.Now()
			prob := infer.Run(blob)
			ctx.Log().Sugar().Debugf("[Infer ] elapsed time: %v", time.Since(s))

			s = time.Now()
			err = proc.After(img, prob, infer.GetElapsedTime(), t)
			if err != nil {
				ctx.Log().Error("failed to process infer result", log.Error(err))
			}
			ctx.Log().Sugar().Debugf("[After ] elapsed time: %v", time.Since(s))

			prob.Close()
			blob.Close()
		}
	})
}

// TODO: remove
func newFunctionClient(cfg http.ClientConfig) (*http.Client, error) {
	ops, err := cfg.ToClientOptions()
	if err != nil {
		return nil, err
	}
	return http.NewClient(ops), nil
}
