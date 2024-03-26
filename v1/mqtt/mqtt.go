package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type MQTT struct {
	Client          mqtt.Client
	Logger          *zap.Logger
	SubscribeTopics map[string]byte
	Done            chan struct{}
}

func NewMqtt(logger *zap.Logger, config *Config) *MQTT {
	mqtOpt := mqtt.NewClientOptions()
	mqtOpt.AddBroker(config.Host)
	mqtOpt.SetClientID(config.ClientId)
	if len(config.User) > 0 {
		mqtOpt.SetUsername(config.User)
	}
	if len(config.Password) > 0 {
		mqtOpt.SetPassword(config.Password)
	}
	client := mqtt.NewClient(mqtOpt)

	st := make(map[string]byte)
	if len(config.Subscriptions) > 0 {
		for _, topic := range config.Subscriptions {
			st[topic] = 0
		}
	}

	return &MQTT{
		Client:          client,
		Logger:          logger,
		SubscribeTopics: st,
	}
}

func (m *MQTT) StartMqtt() error {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		m.Logger.Fatal("MQTT Client", zap.Error(token.Error()))
		return token.Error()
	}

	if len(m.SubscribeTopics) > 0 {
		if token := m.Client.SubscribeMultiple(m.SubscribeTopics, nil); token.Wait() && token.Error() != nil {
			m.Logger.Fatal("MQTT Client Subscribe", zap.Error(token.Error()))
			return token.Error()
		}
	}

	return nil
}

func (m *MQTT) StopMqtt() error {
	m.Client.Disconnect(0)
	return nil
}
