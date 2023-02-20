package kubeconfig

import (
	"errors"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	ErrClusterNotFound = errors.New("cluster can not be found in kubeconfig")
	ErrContextNotFound = errors.New("context can not be found in kubeconfig")
	ErrUserNotFound    = errors.New("user can not be found in kubeconfig")
)

type KubeConfig struct {
	ApiVersion     string            `yaml:"apiVersion"`
	Clusters       []Cluster         `yaml:"clusters"`
	Users          []User            `yaml:"users"`
	Contexts       []Context         `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Kind           string            `yaml:"kind"`
	Preferences    map[string]string `yaml:"preferences"`
}

func (k *KubeConfig) AppendCluster(c Cluster) {
	k.Clusters = append(k.Clusters, c)
}

func (k *KubeConfig) AppendContext(c Context) {
	k.Contexts = append(k.Contexts, c)
}

func (k *KubeConfig) AppendUser(u User) {
	k.Users = append(k.Users, u)
}

func (k *KubeConfig) AppendKubeConfig(o KubeConfig) {
	for _, v := range o.Clusters {
		k.AppendCluster(v)
	}

	for _, v := range o.Contexts {
		k.AppendContext(v)
	}

	for _, v := range o.Users {
		k.AppendUser(v)
	}
}

func (k *KubeConfig) GetClusterByName(name string) (*Cluster, error) {
	for _, c := range k.Clusters {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, ErrClusterNotFound
}

func (k *KubeConfig) GetContextByName(name string) (*Context, error) {
	for _, c := range k.Contexts {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, ErrContextNotFound
}

func (k *KubeConfig) GetUserByName(name string) (*User, error) {
	for _, u := range k.Users {
		if u.Name == name {
			return &u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (k *KubeConfig) OverrideClusterByName(name string, cluster Cluster) error {
	c, err := k.GetClusterByName(name)
	if err != nil {
		return err
	}

	c.Cluster = cluster.Cluster
	return nil
}

func (k *KubeConfig) OverrideContextByName(name string, context Context) error {
	c, err := k.GetContextByName(name)
	if err != nil {
		return err
	}

	c.Context = context.Context
	return nil
}

func (k *KubeConfig) OverrideUserByName(name string, user User) error {
	c, err := k.GetUserByName(name)
	if err != nil {
		return err
	}

	c.User = user.User
	return nil
}

func (k *KubeConfig) WriteToFile(path string) error {
	v, err := yaml.Marshal(k)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, v, 0666)
	if err != nil {
		return err
	}

	return nil
}

func NewFromPath(p string) (KubeConfig, error) {
	k := KubeConfig{}
	f, err := os.ReadFile(p)
	if err != nil {
		return k, fmt.Errorf("could not open file '%s': %w", p, err)
	}

	return New(f)
}

func New(i []byte) (KubeConfig, error) {
	k := KubeConfig{}
	err := yaml.Unmarshal(i, &k)
	if err != nil {
		return k, fmt.Errorf("could not unmarshal kubeconfig: %w", err)
	}

	return k, nil
}

func NewFromDefault() (KubeConfig, string, error) {
	k := KubeConfig{}
	var err error
	var p string

	p, found := os.LookupEnv("KUBECONFIG")
	if found {
		k, err = NewFromPath(p)
		return k, p, err
	}

	// Try to load from ~/.kube/config
	home, err := os.UserHomeDir()
	if err != nil {
		return k, "", err
	}

	p = path.Join(home, "/.kube/config")
	k, err = NewFromPath(p)
	return k, p, err
}
