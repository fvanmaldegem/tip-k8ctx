package kubeconfig

import log "github.com/sirupsen/logrus"

type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
		Server                   string `yaml:"server"`
	}
}

func (c *Cluster) Rename(name string) {
	log.WithFields(log.Fields{
		"oldName": c.Name,
		"newName": name,
	}).Debug("renaming cluster")
	c.Name = name
}
