package provider

import (
	"fmt"
	"reflect"
)

type ServiceProvider struct {
	services map[reflect.Type]interface{}
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{
		services: make(map[reflect.Type]interface{}),
	}
}

func (sp *ServiceProvider) Register(iface interface{}, impl interface{}) {
	t := reflect.TypeOf(iface).Elem()
	sp.services[t] = impl
}

func (sp *ServiceProvider) Get(ifaceType reflect.Type) (interface{}, error) {
	if impl, ok := sp.services[ifaceType]; ok {
		return impl, nil
	}
	return nil, fmt.Errorf("service not registered: %s", ifaceType)
}
