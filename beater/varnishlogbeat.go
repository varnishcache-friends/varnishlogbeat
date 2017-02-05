package beater

import (
	"fmt"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/phenomenes/vago"
	"github.com/phenomenes/varnishlogbeat/config"
)

type Varnishlogbeat struct {
	done    chan struct{}
	config  config.Config
	client  publisher.Client
	varnish *vago.Varnish
}

// New creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	vb := Varnishlogbeat{
		//done:   make(chan struct{}),
		config: config,
	}

	return &vb, nil
}

func (vb *Varnishlogbeat) Run(b *beat.Beat) error {
	var err error

	logp.Info("varnishlogbeat is running! Hit CTRL-C to stop it.")

	vb.varnish, err = vago.Open(vb.config.Path)
	if err != nil {

		return err
	}

	vb.client = b.Publisher.Connect()
	err = vb.harvest()
	if err != nil {
		logp.Err("%s", err)
	}

	return err
}

func (vb *Varnishlogbeat) harvest() error {
	tx := make(common.MapStr)
	counter := 1

	vb.varnish.Log("", vago.REQ, func(vxid uint32, tag, _type, data string) int {
		switch _type {
		case "c":
			_type = "client"
		case "b":
			_type = "backend"
		default:
			return 0
		}

		switch tag {
		case "BereqHeader", "BerespHeader", "ObjHeader", "ReqHeader", "RespHeader":
			header := strings.SplitN(data, ": ", 2)
			k := header[0]
			v := header[1]
			if _, ok := tx[tag]; ok {
				tx[tag].(common.MapStr)[k] = v
			} else {
				tx[tag] = common.MapStr{k: v}
			}
		case "End":
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"count":      counter,
				"type":       _type,
				"vxid":       vxid,
				"tx":         tx,
			}
			vb.client.PublishEvent(event)
			counter++
			logp.Info("Event sent")

			// destroy and re-create the map
			tx = nil
			tx = make(common.MapStr)
		default:
			tx[tag] = data
		}

		return 0
	})

	return nil
}

func (vb *Varnishlogbeat) Stop() {
	vb.varnish.Stop()
	vb.varnish.Close()
	vb.client.Close()
}
