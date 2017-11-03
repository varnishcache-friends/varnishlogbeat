package beater

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/phenomenes/vago"
	"github.com/phenomenes/varnishlogbeat/config"
)

// Varnishlogbeat implements the Beater interface.
type Varnishlogbeat struct {
	client  publisher.Client
	varnish *vago.Varnish
	config  *vago.Config
}

// New creates a new Varnishlogbeat.
func New(b *beat.Beat, c *common.Config) (beat.Beater, error) {
	cfg := config.DefaultConfig
	if err := c.Unpack(&cfg); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	vb := Varnishlogbeat{
		config: &vago.Config{
			Path:    cfg.Path,
			Timeout: cfg.Timeout,
		},
	}

	return &vb, nil
}

// Run opens a Varnish Shared Memory file and publishes log events.
func (vb *Varnishlogbeat) Run(b *beat.Beat) error {
	var err error

	logp.Info("varnishlogbeat is running! Hit CTRL-C to stop it.")

	vb.varnish, err = vago.Open(vb.config)
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

// harvest reads and parses Varnish log data.
func (vb *Varnishlogbeat) harvest() error {
	tx := make(common.MapStr)
	counter := 1

	vb.varnish.Log("",
		vago.REQ,
		vago.COPT_TAIL|vago.COPT_BATCH,
		func(vxid uint32, tag, _type, data string) int {
			switch _type {
			case "c":
				_type = "client"
			case "b":
				_type = "backend"
			default:
				return 0
			}

			switch tag {
			case "BereqHeader",
				"BerespHeader",
				"ObjHeader",
				"ReqHeader",
				"RespHeader",
				"Timestamp":
				header := strings.SplitN(data, ": ", 2)
				key := header[0]
				var value interface{}
				switch {
				case key == "Content-Length":
					value, _ = strconv.Atoi(header[1])
				case len(header) == 2:
					value = header[1]
				// if the header is too long, header and value might get truncated
				default:
					value = "truncated"
				}
				if _, ok := tx[tag]; ok {
					tx[tag].(common.MapStr)[key] = value
				} else {
					tx[tag] = common.MapStr{key: value}
				}

			case "Length":
				tx[tag], _ = strconv.Atoi(data)

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

// Stop stops processing Varnish events, closes the VSM and publisher client.
func (vb *Varnishlogbeat) Stop() {
	vb.varnish.Stop()
	vb.varnish.Close()
	vb.client.Close()
}
