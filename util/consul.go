package util

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
	"log"
	"strconv"
	configEnv "gokitdemo/config"

)

var ConsulClient *consulapi.Client
var (
	ServiceID   string
	ServiceName string
	ServicePort int
)

func init() {
	config := consulapi.DefaultConfig()
	config.Address = configEnv.CONSUL_ADDRESS
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
	ServiceID = "userservice" + uuid.New()
}

func SetServiceNameAndPort(name string, port int) {
	ServiceName = name
	ServicePort = port
}
func RegService() {

	req := consulapi.AgentServiceRegistration{
		ID:      ServiceID,
		Name:    ServiceName,
		Address: configEnv.SERVICE_IP,
		Port:    ServicePort,
		Tags:    []string{"primary"},
	}

	check := consulapi.AgentServiceCheck{
		Interval: "5s",
		HTTP:     "http://"+configEnv.SERVICE_IP+":" + strconv.Itoa(ServicePort) + "/health",
	}
	req.Check = &check

	err := ConsulClient.Agent().ServiceRegister(&req)
	if err != nil {
		log.Fatal(err)
	}
}

func Unregservice() {
	ConsulClient.Agent().ServiceDeregister(ServiceID)
}
