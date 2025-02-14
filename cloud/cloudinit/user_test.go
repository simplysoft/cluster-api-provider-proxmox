package cloudinit_test

import (
	"reflect"
	"testing"

	infrav1 "github.com/sp-yduck/cluster-api-provider-proxmox/api/v1beta1"
	"github.com/sp-yduck/cluster-api-provider-proxmox/cloud/cloudinit"
)

func TestParseUser(t *testing.T) {
	testYaml := `
write_files:
  - path: /run/kubeadm/kubeadm.yaml
    owner: root:root
    permissions: '0640'
    content: |
      asdfasdfasdf
runcmd:
  - 'kubeadm init --config /run/kubeadm/kubeadm.yaml  && echo success > /run/cluster-api/bootstrap-success.complete'
  - "curl -L https://dl.k8s.io/release/v1.27.3/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl"
  - "chmod +x /usr/local/bin/kubectl"
  - "reboot now"
  `
	_, err := cloudinit.ParseUserData(testYaml)
	if err != nil {
		t.Errorf("failed to parse user: %v", err)
	}

	testYaml = `
write_files:
  - path: /run/kubeadm/kubeadm.yaml
owner: root:root
    permissions: '0640'
    content: |
      asdfasdfasdf
  `
	user, err := cloudinit.ParseUserData(testYaml)
	if err == nil {
		t.Errorf("should returns an error. user=%v", *user)
	}
}

func TestGenerateUserYaml(t *testing.T) {
	testYaml := `
write_files:
  - path: /run/kubeadm/kubeadm.yaml
    owner: root:root
    permissions: '0640'
    content: |
      asdfasdfasdf
runcmd:
  - 'kubeadm init --config /run/kubeadm/kubeadm.yaml  && echo success > /run/cluster-api/bootstrap-success.complete'
  - "curl -L https://dl.k8s.io/release/v1.27.3/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl"
  - "chmod +x /usr/local/bin/kubectl"
  - "reboot now"
  `

	uc, err := cloudinit.ParseUserData(testYaml)
	if err != nil {
		t.Fatalf("failed to parse user: %v", err)
	}

	_, err = cloudinit.GenerateUserDataYaml(*uc)
	if err != nil {
		t.Fatalf("generate : %v", err)
	}
}

func TestMergeUsers(t *testing.T) {
	a := infrav1.UserData{
		User:   "override-user",
		RunCmd: []string{"command A", "command B"},
	}
	b := infrav1.UserData{
		User:   "test-user",
		RunCmd: []string{"command C"},
	}
	expected := infrav1.UserData{
		User:   "override-user",
		RunCmd: []string{"command A", "command B", "command C"},
	}
	c, err := cloudinit.MergeUserDatas(&a, &b)
	if err != nil {
		t.Errorf("failed to merge cloud init user data: %v", err)
	}
	if !reflect.DeepEqual(*c, expected) {
		t.Errorf("%v is expected to same as %v", *c, expected)
	}
}
