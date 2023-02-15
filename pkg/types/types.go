package types

type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
		Server                   string `yaml:"server"`
	}
}

type Context struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster string `yaml:"cluster"`
		User    string `yaml:"user"`
	}
}

type User struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateDate string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
	}
}

type KubeConfig struct {
	ApiVersion     string            `yaml:"apiVersion"`
	Clusters       []Cluster         `yaml:"clusters"`
	Users          []User            `yaml:"users"`
	Contexts       []Context         `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Kind           string            `yaml:"kind"`
	Preferences    map[string]string `yaml:"preferences"`
}
