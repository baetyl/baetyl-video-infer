package main

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"
	"gocv.io/x/gocv"
)

func main() {
	context.Run(func(ctx context.Context) error {
		var cfg Config
		// load custom config
		err := ctx.LoadCustomConfig(&cfg)
		if err != nil {
			return err
		}
		// create a broker client
		bc, err := ctx.NewBrokerClient()
		if err != nil {
			return err
		}
		bc.Start(nil)
		defer bc.Close()
		//  create inference
		infer, err := NewInfer(cfg.Infer)
		if err != nil {
			return err
		}
		defer infer.Close()
		// create video
		video, err := NewVideo(cfg.Video)
		if err != nil {
			return err
		}
		defer video.Close()
		// create function clients
		fc, err := ctx.NewFunctionHttpClient()
		if err != nil {
			return err
		}
		// create process
		proc := NewProcess(cfg.Process, bc, fc)

		var s time.Time
		img := gocv.NewMat()
		defer img.Close()
		quit := ctx.WaitChan()
		for {
			select {
			case <-quit:
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
