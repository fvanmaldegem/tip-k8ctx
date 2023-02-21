package app

import (
	"fmt"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/fvanmaldegem/tip-k8ctx/pkg/editor"
	"github.com/fvanmaldegem/tip-k8ctx/pkg/kubeconfig"
	log "github.com/sirupsen/logrus"
)

type App struct {
	existingConfigPath string
	existingConfig     kubeconfig.KubeConfig
	newConfigPath      string
	newConfig          kubeconfig.KubeConfig
	forceOverride      bool
	newContextName     string
}

func New(kubeConfigPath, newConfigPath, newContextName string, force bool) App {
	var err error
	var k kubeconfig.KubeConfig

	if newContextName == "" {
		log.Fatal("please specify a name for the context")
	}

	if kubeConfigPath == "" {
		k, kubeConfigPath, err = kubeconfig.NewFromDefault()
		if err != nil {
			log.WithError(err).Fatal("could not load kubeconfig")
		}
	} else {
		k, err = kubeconfig.NewFromPath(kubeConfigPath)
		if err != nil {
			log.WithError(err).WithField("path", kubeConfigPath).Fatal("could not load kubeconfig")
		}
	}

	var kn kubeconfig.KubeConfig
	if newConfigPath == "" {
		kn, err = kubeconfig.NewFromYaml(editor.OpenAndRead())
		if err != nil {
			log.WithError(err).Fatal("could not load kubeconfig")
		}
	} else {
		kn, err = kubeconfig.NewFromPath(newConfigPath)
		if err != nil {
			log.WithError(err).WithField("path", kubeConfigPath).Fatal("could not load kubeconfig")
		}
	}

	return App{
		existingConfig:     k,
		existingConfigPath: kubeConfigPath,
		newConfig:          kn,
		newConfigPath:      newConfigPath,
		forceOverride:      force,
		newContextName:     newContextName,
	}
}

func (a *App) Run() {
	if a.doesNewNameAlreadyExist() {
		var err error
		override := a.forceOverride

		if !a.forceOverride {
			msg := fmt.Sprintf("Do you want to overwrite the context %s", a.newContextName)
			c := confirmation.New(msg, confirmation.No)
			override, err = c.RunPrompt()
			if err != nil {
				log.WithError(err).Fatal("could not prompt user")
			}
		}

		if override {
			a.existingConfig.OverrideClusterByName(a.newContextName, a.newConfig.Clusters[0])

			a.newConfig.Contexts[0].Context.User = a.newContextName
			a.newConfig.Contexts[0].Context.Cluster = a.newContextName
			a.existingConfig.OverrideContextByName(a.newContextName, a.newConfig.Contexts[0])

			a.existingConfig.OverrideUserByName(a.newContextName, a.newConfig.Users[0])

			a.existingConfig.WriteToFile(a.existingConfigPath)
		}

		return
	}

	log.WithField("newName", a.newContextName).Debug("preparing/renaming new config")
	a.newConfig.Clusters[0].Rename(a.newContextName)
	a.newConfig.Users[0].Rename(a.newContextName)
	a.newConfig.Contexts[0].Rename(a.newContextName)
	a.newConfig.Contexts[0].Context.Cluster = a.newContextName
	a.newConfig.Contexts[0].Context.User = a.newContextName

	log.WithField("newName", a.newContextName).Debug("adding new config")
	a.existingConfig.AppendCluster(a.newConfig.Clusters[0])
	a.existingConfig.AppendUser(a.newConfig.Users[0])
	a.existingConfig.AppendContext(a.newConfig.Contexts[0])

	log.WithField("path", a.existingConfigPath).Debug("saving config")
	a.existingConfig.WriteToFile(a.existingConfigPath)
}

func (a *App) doesNewNameAlreadyExist() bool {
	if _, err := a.existingConfig.GetClusterByName(a.newContextName); err == nil {
		log.WithError(err).WithField("name", a.newContextName).Debug("name already exists in clusters")
		return true
	}

	if _, err := a.existingConfig.GetUserByName(a.newContextName); err == nil {
		log.WithField("name", a.newContextName).Debug("name already exists in users")
		return true
	}

	if _, err := a.existingConfig.GetContextByName(a.newContextName); err == nil {
		log.WithField("name", a.newContextName).Debug("name already exists in contexts")
		return true
	}

	return false
}
