package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
)

type mqttClient struct {
	ctx             context.Context
	client          mqtt.Client
	cfg             *config.Config
	SubscribeTopics map[string]byte
	onConnect       func(ctx context.Context, client mqtt.Client) error
	onCallback      func(ctx context.Context, client mqtt.Client, message mqtt.Message)
}

type MQTTClient interface {
	OnConnect(fn func(ctx context.Context, client mqtt.Client) error)
	OnCallback(fn func(ctx context.Context, client mqtt.Client, message mqtt.Message))
	GetClient() mqtt.Client
	GetSubscriptions() map[string]byte
}

func New(cfg *config.Config, subscribeTopics map[string]byte) runtime.Task {
	return &mqttClient{
		cfg:             cfg,
		SubscribeTopics: subscribeTopics,
	}
}

func (t *mqttClient) OnConnect(fn func(ctx context.Context, client mqtt.Client) error) {
	t.onConnect = fn
}

func (t *mqttClient) OnCallback(fn func(ctx context.Context, client mqtt.Client, message mqtt.Message)) {
	t.onCallback = fn
}

func (t *mqttClient) GetClient() mqtt.Client {
	return t.client
}
func (t *mqttClient) GetSubscriptions() map[string]byte {
	return t.SubscribeTopics
}

func (t *mqttClient) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL MQTT")
	mqtOpt := mqtt.NewClientOptions()
	mqtOpt.AddBroker(t.cfg.MQTTClient.URI)
	mqtOpt.SetClientID(t.cfg.MQTTClient.ClientId)
	if len(t.cfg.MQTTClient.Username) > 0 {
		mqtOpt.SetUsername(t.cfg.MQTTClient.Username)
	}
	if len(t.cfg.MQTTClient.Password) > 0 {
		mqtOpt.SetPassword(t.cfg.MQTTClient.Password)
	}
	t.client = mqtt.NewClient(mqtOpt)

	if token := t.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Errs("MQTT Client", []error{token.Error()}).Send()
	}

	if t.onConnect != nil {
		err = t.onConnect(t.ctx, t.client)
		if err != nil {
			return err
		}
	}

	if t.SubscribeTopics != nil && t.onCallback != nil {
		t.client.SubscribeMultiple(t.SubscribeTopics, func(client mqtt.Client, message mqtt.Message) {
			t.onCallback(t.ctx, client, message)
		})
	}

	return nil
}

func (t *mqttClient) Ping(context.Context) error {
	return nil
}

func (t *mqttClient) Close() error {
	log.Debug().Msg("CLOSE MQTT connection")
	t.client.Disconnect(0)
	return nil
}
