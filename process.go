package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/http"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/mqtt"
	"gocv.io/x/gocv"
)

type content map[string]interface{}

func (c content) isDiscard() bool {
	if v, ok := c["imageDiscard"]; ok {
		if b, ok := v.(bool); ok && b {
			return true
		} else if s, ok := v.(string); ok && strings.EqualFold(s, "true") {
			return true
		}
	}
	return false
}

func (c content) location() string {
	if v, ok := c["imageLocation"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (c content) topic() string {
	if v, ok := c["publishTopic"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (c content) qos() uint32 {
	if v, ok := c["publishQOS"]; ok {
		if i, ok := v.(int); ok {
			return uint32(i)
		}
	}
	return 0
}

// Process the image processor
type Process struct {
	cfg ProcessInfo
	bc  *mqtt.Client
	fc  *http.Client

	size image.Point
	mean gocv.Scalar
}

// NewProcess creates a new process
func NewProcess(cfg ProcessInfo, bc *mqtt.Client, fc *http.Client) *Process {
	return &Process{
		cfg:  cfg,
		bc:   bc,
		fc:   fc,
		size: image.Pt(cfg.Before.Width, cfg.Before.Hight),
		mean: gocv.NewScalar(cfg.Before.Mean.V1, cfg.Before.Mean.V2, cfg.Before.Mean.V3, cfg.Before.Mean.V4),
	}
}

// Before processes image before inference
func (p *Process) Before(img gocv.Mat) gocv.Mat {
	return gocv.BlobFromImage(img, p.cfg.Before.Scale, p.size, p.mean, p.cfg.Before.SwapRB, p.cfg.Before.Crop)
}

// After processes image after inference
func (p *Process) After(img gocv.Mat, results gocv.Mat, elapsedTime float64, captureTime time.Time) error {
	log.L().Sugar().Debugf("type: %s total: %d size: %v", results.Type(), results.Total(), results.Size())

	s := time.Now()
	out, err := p.fc.Call(p.cfg.After.Function.Name, results.ToBytes())
	if err != nil {
		return err
	}
	log.L().Sugar().Debugf("[After ][Call     ] elapsed time: %v", time.Since(s))

	if out == nil {
		return nil
	}

	s = time.Now()
	var cnt content
	log.L().Sugar().Debugf(string(out))
	err = json.Unmarshal(out, &cnt)
	if err != nil {
		return err
	}
	log.L().Sugar().Debugf("[After ][Unmarshal] elapsed time: %v", time.Since(s))

	discard := cnt.isDiscard()
	location := cnt.location()
	if !discard && location != "" {
		s = time.Now()
		if !gocv.IMWrite(location, img) {
			os.MkdirAll(path.Dir(location), 0755)
			if !gocv.IMWrite(location, img) {
				return fmt.Errorf("failed to save image: %s", location)
			}
		}
		log.L().Sugar().Debugf("[After ][Write    ] elapsed time: %v", time.Since(s))
	}

	if p.bc == nil || cnt.topic() == "" {
		return nil
	}

	cnt["imageWidth"] = img.Cols()
	cnt["imageHight"] = img.Rows()
	cnt["imageCaptureTime"] = captureTime
	cnt["imageInferenceTime"] = elapsedTime
	cnt["imageProcessTime"] = (time.Since(s)).Seconds() + elapsedTime
	if !discard && location == "" {
		s = time.Now()
		cnt["imageData"] = img.ToBytes()
		log.L().Sugar().Debugf("[After ][ToBytes  ] elapsed time: %v", time.Since(s))
	}

	s = time.Now()
	msgData, err := json.Marshal(cnt)
	if err != nil {
		return errors.Trace(err)
	}
	log.L().Sugar().Debugf("[After ][Marshal  ] elapsed time: %v", time.Since(s))

	s = time.Now()
	err = p.bc.Publish(mqtt.QOS(cnt.qos()), cnt.topic(), msgData, 1, false, false)
	log.L().Sugar().Debugf("[After ][Publish  ] elapsed time: %v", time.Since(s))
	return err
}
