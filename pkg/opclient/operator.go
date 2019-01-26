package opclient

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
	em "zwczou/gobase/middleware"
	"zwczou/operator/pkg/def"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	log "github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	extra.RegisterFuzzyDecoders()
	em.SetNamingStrategy(em.LowerCaseWithUnderscores)
	em.RegisterTimeAsFormatCodec("2006-01-02 15:04:05")
}

type Packet = def.Packet

type OperatorClient struct {
	sync.RWMutex
	*websocket.Conn
	connected  int32
	host       string
	wsUrl      string
	packetChan chan *Packet
}

func New(host string) *OperatorClient {
	r := strings.NewReplacer("http", "ws")
	c := &OperatorClient{
		host:       host,
		wsUrl:      r.Replace(host + "/ws"),
		packetChan: make(chan *Packet, 10240),
	}
	return c
}

func (client *OperatorClient) Connect() error {
	client.Lock()
	defer client.Unlock()

	if atomic.LoadInt32(&client.connected) == 0 {
		c, _, err := websocket.DefaultDialer.Dial(client.wsUrl, nil)
		if err != nil {
			return err
		}
		client.Conn = c
		atomic.StoreInt32(&client.connected, 1)
	}
	return nil
}

func (client *OperatorClient) IsClosed() bool {
	return atomic.LoadInt32(&client.connected) == 0
}

func (client *OperatorClient) Close() error {
	if atomic.LoadInt32(&client.connected) == 1 {
		err := client.Conn.Close()
		client.Conn = nil
		atomic.StoreInt32(&client.connected, 0)
		return err
	}
	return nil
}

func (client *OperatorClient) PacketChan() <-chan *Packet {
	return client.packetChan
}

func (client *OperatorClient) Loop() (err error) {
	fields := log.Fields{
		"wsurl":  client.wsUrl,
		"client": "operator",
	}

	err = client.Connect()
	if err != nil {
		fields["error"] = err.Error()
		log.WithFields(fields).Warn("connect error")
		return
	}
	defer client.Close()
	log.WithFields(fields).Info("connected")

	for {
		client.SetReadDeadline(time.Now().Add(time.Minute * 9))
		_, msg, err := client.ReadMessage()
		if err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Warn("read message error")
			return err
		}
		var packet Packet
		err = json.Unmarshal(msg, &packet)
		if err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Warn("json decode error")
		}

		// 如果队列满了丢掉数据
		select {
		case client.packetChan <- &packet:
		default:
			log.WithFields(fields).Debug("packet losing")
		}
	}
	return nil
}
