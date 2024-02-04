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

	return &MQTT{
		Client: client,
		Logger: logger,
	}
}

func (m *MQTT) StartMqtt() error {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		m.Logger.Fatal("MQTT Client", zap.Error(token.Error()))
		return token.Error()
	}
	return nil
}

func (m *MQTT) StopMqtt() error {
	m.Client.Disconnect(0)
	return nil
}
