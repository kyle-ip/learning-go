package naming

import (
	"fmt"

	"github.com/klintcheng/kim"
)

// "ID": "qa-dfirst-zfirst-tgateway-172.16.235.145-0-8000",
// "Service": "tgateway",
// "Tags": [
// "ZONE:qa-dfirst-zfirst",
// "TMC_REGION:SH",
// "TMC_DOMAIN:g002-qa.tutormeetplus.com"
// ],
// "Address": "172.16.235.145",
// "Port": 8000,

//Service define a Service

// DefaultService Service Impl
type DefaultService struct {
	Id        string
	Name      string
	Address   string
	Port      int
	Protocol  string
	Namespace string
	Tags      []string
	Meta      map[string]string
}

// NewEntry NewEntry
func NewEntry(id, name, protocol string, address string, port int) kim.ServiceRegistration {
	return &DefaultService{
		Id:       id,
		Name:     name,
		Address:  address,
		Port:     port,
		Protocol: protocol,
	}
}

// ID returns the ServiceImpl ID
func (e *DefaultService) ServiceID() string {
	return e.Id
}

// Name Name
func (e *DefaultService) ServiceName() string { return e.Name }

// Namespace Namespace
func (e *DefaultService) GetNamespace() string { return e.Namespace }

// Address Address
func (e *DefaultService) PublicAddress() string {
	return e.Address
}

func (e *DefaultService) PublicPort() int { return e.Port }

// Protocol Protocol
func (e *DefaultService) GetProtocol() string { return e.Protocol }

func (e *DefaultService) DialURL() string {
	if e.Protocol == "tcp" {
		return fmt.Sprintf("%s:%d", e.Address, e.Port)
	}
	return fmt.Sprintf("%s://%s:%d", e.Protocol, e.Address, e.Port)
}

// Tags Tags
func (e *DefaultService) GetTags() []string { return e.Tags }

// Meta Meta
func (e *DefaultService) GetMeta() map[string]string { return e.Meta }

func (e *DefaultService) String() string {
	return fmt.Sprintf("Id:%s,Name:%s,Address:%s,Port:%d,Ns:%s,Tags:%v,Meta:%v", e.Id, e.Name, e.Address, e.Port, e.Namespace, e.Tags, e.Meta)
}
