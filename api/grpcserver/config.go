// Package config provides configuration for GRPC and HTTP api servers
package grpcserver

import (
	"time"
)

type Config struct {
	PublicServices  []Service
	PublicListener  string `mapstructure:"grpc-public-listener"`
	PrivateServices []Service
	PrivateListener string `mapstructure:"grpc-private-listener"`
	PostServices    []Service
	PostListener    string    `mapstructure:"grpc-post-listener"`
	TLSServices     []Service `mapstructure:"grpc-tls-services"`
	TLSListener     string    `mapstructure:"grpc-tls-listener"`
	TLSCACert       string    `mapstructure:"grpc-tls-ca-cert"`
	TLSCert         string    `mapstructure:"grpc-tls-cert"`
	TLSKey          string    `mapstructure:"grpc-tls-key"`
	GrpcSendMsgSize int       `mapstructure:"grpc-send-msg-size"`
	GrpcRecvMsgSize int       `mapstructure:"grpc-recv-msg-size"`
	JSONListener    string    `mapstructure:"grpc-json-listener"`

	SmesherStreamInterval time.Duration `mapstructure:"smesherstreaminterval"`
}

type Service = string

const (
	Admin                    Service = "admin"
	Debug                    Service = "debug"
	GlobalState              Service = "global"
	Mesh                     Service = "mesh"
	Transaction              Service = "transaction"
	Activation               Service = "activation"
	Smesher                  Service = "smesher"
	Post                     Service = "post"
	Node                     Service = "node"
	ActivationV2Alpha1       Service = "activation_v2alpha1"
	ActivationStreamV2Alpha1 Service = "activation_stream_v2alpha1"
)

// DefaultConfig defines the default configuration options for api.
func DefaultConfig() Config {
	return Config{
		PublicServices:        []Service{GlobalState, Mesh, Transaction, Node, Activation, ActivationV2Alpha1},
		PublicListener:        "0.0.0.0:9092",
		PrivateServices:       []Service{Admin, Smesher, Debug, ActivationStreamV2Alpha1},
		PrivateListener:       "127.0.0.1:9093",
		PostServices:          []Service{Post},
		PostListener:          "127.0.0.1:9094",
		TLSServices:           []Service{Post},
		TLSListener:           "",
		JSONListener:          "",
		GrpcSendMsgSize:       1024 * 1024 * 10,
		GrpcRecvMsgSize:       1024 * 1024 * 10,
		SmesherStreamInterval: time.Second,
	}
}

// DefaultTestConfig returns the default config for tests.
func DefaultTestConfig() Config {
	conf := DefaultConfig()
	conf.PublicListener = "127.0.0.1:0"
	conf.PrivateListener = "127.0.0.1:0"
	conf.PostListener = "127.0.0.1:0"
	conf.JSONListener = ""
	conf.TLSListener = ""
	return conf
}
