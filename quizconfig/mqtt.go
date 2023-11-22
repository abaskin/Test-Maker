package quizconfig

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
)

type TransportMQTT_t struct {
	cm         *autopaho.ConnectionManager
	ctx        *context.Context
	mqttCancel context.CancelFunc
}

func NewMQTT(topicHandler func(*paho.Publish)) (*TransportMQTT_t, error) {
	mqtt := &TransportMQTT_t{}
	serverURL, err := url.Parse(BrokerUrl)
	if err != nil {
		return &TransportMQTT_t{}, errors.Wrap(err, "parse URL Failed")
	}

	cliCfg := configMqtt(serverURL, ClientTopic, topicHandler)

	ctx, cancel := context.WithCancel(context.Background())
	mqtt.mqttCancel = cancel
	mqtt.ctx = &ctx

	// Connect to the broker - this will return immediately after initiating the connection process
	cm, err := autopaho.NewConnection(ctx, *cliCfg)
	if err != nil {
		return &TransportMQTT_t{}, errors.Wrap(err, "mqtt new connection failed")
	}

	if err := cm.AwaitConnection(ctx); err != nil { // Should only happen when context is cancelled
		return &TransportMQTT_t{}, errors.Wrap(err, "publisher done (AwaitConnection)")
	}
	mqtt.cm = cm

	log.Println("connected to broker")

	return &TransportMQTT_t{
		mqttCancel: cancel,
		ctx:        &ctx,
		cm:         cm,
	}, nil
}

func configMqtt(serverURL *url.URL, recvTopic string, handler func(*paho.Publish)) *autopaho.ClientConfig {
	return &autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverURL},
		KeepAlive:         30,
		ConnectRetryDelay: 5 * time.Minute,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: recvTopic, QoS: QoS},
				},
			}); err != nil {
				log.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			log.Printf("mqtt subscription: %s\n", recvTopic)
		}, OnConnectError: func(err error) { log.Printf("error whilst attempting connection: %s\n", err) },
		Debug:     Logger{Prefix: "autoPaho"},
		PahoDebug: Logger{Prefix: "paho"},
		ClientConfig: paho.ClientConfig{
			ClientID:      randstr.String(32),
			Router:        paho.NewSingleHandlerRouter(handler),
			OnClientError: func(err error) { log.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
}

func (m *TransportMQTT_t) Send(data *ClientData_t, topic string) error {
	var msg bytes.Buffer
	if err := gob.NewEncoder(&msg).Encode(data); err != nil {
		return errors.Wrap(err, "gob encode error")
	}

	pr, err := m.cm.Publish(*m.ctx, &paho.Publish{
		QoS:     QoS,
		Topic:   topic,
		Payload: msg.Bytes(),
	})

	if err != nil {
		err = errors.Wrap(err, "NewTransport")
	} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
		log.Printf("reason code %d received\n", pr.ReasonCode)
	}

	return err
}

func (m *TransportMQTT_t) SendNoop() error {
	return nil
}

func (m *TransportMQTT_t) ShutDown() {
	// We could cancel the context at this point but will call Disconnect
	// instead (this waits for autopaho to shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	m.cm.Disconnect(ctx)
	cancel()
	m.mqttCancel()
}
