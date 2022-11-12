package jokit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type AMQP_Protocol string

const (
	AMQP  AMQP_Protocol = "amqp"
	AMQPS AMQP_Protocol = "amqps"
)

type AMQPConfig struct {
	Protocol AMQP_Protocol `json:"amqp_protocol"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Host     string        `json:"host"`
	Port     string        `json:"port"`
	Vhost    string        `json:"vhost"`
}

// MQConnect connects to rabbitmq and returns a pointer to a connection/ error
// this function should be called at the initialization of the app and the connection stored
// Channels should be created/closed on demand from this one connection.
func MQConnect(mqConfig AMQPConfig) (connection *amqp.Connection, err error) {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s%s", mqConfig.Protocol, mqConfig.Username, mqConfig.Password, mqConfig.Host, mqConfig.Port, mqConfig.Vhost)
	Log("%s Connecting to rabbitMq...", LogPrefix)
	connection, err = amqp.Dial(connectionString)
	if err != nil {
		err = fmt.Errorf("%s failed to connect to RabbitMQ, err: %v", LogPrefix, err)
		return nil, err
	}
	Log("%s Connected to rabbitMq on host %s successfully", LogPrefix, mqConfig.Host)
	return
}
