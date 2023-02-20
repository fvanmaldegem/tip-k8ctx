package kubeconfig

import log "github.com/sirupsen/logrus"

type Context struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster string `yaml:"cluster"`
		User    string `yaml:"user"`
	}
}

func (c *Context) Rename(name string) {
	log.WithFields(log.Fields{
		"oldName": c.Name,
		"newName": name,
	}).Debug("renaming context")
	c.Name = name
	c.Name = name
}
