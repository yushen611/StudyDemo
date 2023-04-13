package consul

import (
	"fmt"
	"math/rand"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	client *consulapi.Client
}

func NewConsulClient() (*ConsulClient, error) {
	config := consulapi.DefaultConfig()
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
	// services, _, err := c.client.Catalog().Service(serviceName, "", nil)
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}
	if len(services) == 0 {
		return "", fmt.Errorf("service not found")
	}
	// 随机选择一个服务实例
	if len(services) > 0 {
		index := rand.Intn(len(services))
		service := services[index].Service
		address := fmt.Sprintf("%v:%v", service.Address, service.Port)
		return address, nil
	}

	return "", fmt.Errorf("no healthy instances found for service %s", serviceName)
}
