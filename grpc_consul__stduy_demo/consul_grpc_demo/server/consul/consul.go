package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	client *consulapi.Client
}

func NewConsulClient() (*ConsulClient, error) {
	config := consulapi.DefaultConfig()
	// config.Address = "192.168.3.2:8500" // 指定 Consul 地址
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulClient{client: client}, nil
}

func (c *ConsulClient) RegisterService(serviceID, serviceName, serviceHost string, servicePort int) error {
	service := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceHost,
		Port:    servicePort,
	}
	return c.client.Agent().ServiceRegister(service)
}

func (c *ConsulClient) DiscoverService(serviceName string) (string, error) {
	services, _, err := c.client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return "", err
	}
	if len(services) == 0 {
		return "", fmt.Errorf("service not found")
	}
	return fmt.Sprintf("%s:%d", services[0].ServiceAddress, services[0].ServicePort), nil
}
