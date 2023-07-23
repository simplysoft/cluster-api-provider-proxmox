package instance

import (
	"github.com/sp-yduck/proxmox-go/proxmox"

	"github.com/sp-yduck/cluster-api-provider-proxmox/cloud"
	"github.com/sp-yduck/cluster-api-provider-proxmox/cloud/scope"
)

type Scope interface {
	cloud.Machine
}

type Service struct {
	scope  Scope
	client proxmox.Service
	remote scope.SSHClient
}

func NewService(s Scope) *Service {
	return &Service{
		scope:  s,
		client: *s.CloudClient(),
		remote: *s.RemoteClient(),
	}
}
