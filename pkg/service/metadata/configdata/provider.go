package configdata

type CNode struct {
	instanceID string
}

type ConfigDataService interface {
	RegisterCNodes() error
	DeregisterCNodes() error
	UpdateCNodes() error
}
