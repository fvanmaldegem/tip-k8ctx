package kubeconfig

import log "github.com/sirupsen/logrus"

type User struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateDate string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
	}
}

func (u *User) Rename(name string) {
	log.WithFields(log.Fields{
		"oldName": u.Name,
		"newName": name,
	}).Debug("renaming user")
	u.Name = name
}
