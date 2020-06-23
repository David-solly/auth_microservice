package models

type AddressConfig struct {
	ConsulAddr          string
	ConsulPort          string
	AdvertiseAddr       string
	AdvertisePort       string
	AdvertiseHealthPort string
}

type ServiceConfig struct {
	ID   string
	Name string
	Tags []string
}
