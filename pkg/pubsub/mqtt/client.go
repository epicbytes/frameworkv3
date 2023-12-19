package mqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

type MQTTClient struct {
	ctx             context.Context
	client          mqtt.Client
	URI             string
	ClientID        string
	SubscribeTopics map[string]byte
	onConnect       func(ctx context.Context, client mqtt.Client) error
	onCallback      func(ctx context.Context, client mqtt.Client, message mqtt.Message)
}

func (t *MQTTClient) OnConnect(fn func(ctx context.Context, client mqtt.Client) error) {
	t.onConnect = fn
}

func (t *MQTTClient) OnCallback(fn func(ctx context.Context, client mqtt.Client, message mqtt.Message)) {
	t.onCallback = fn
}

func (t *MQTTClient) GetClient() mqtt.Client {
	return t.client
}

func (t *MQTTClient) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL MQTT")
	mqtOpt := mqtt.NewClientOptions()
	mqtOpt.AddBroker(t.URI)
	mqtOpt.SetClientID(t.ClientID)
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

func (t *MQTTClient) Ping(context.Context) error {
	return nil
}

func (t *MQTTClient) Close() error {
	log.Debug().Msg("CLOSE MQTT connection")
	t.client.Disconnect(0)
	return nil
}
