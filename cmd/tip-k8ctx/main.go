package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/florisvanmaldegem/tip-k8ctx/pkg/types"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	// Load flags
	kubeConfigLocation := flag.String("kubeconfig", "", "specifies the location of the kubeconfig. Defaults to $KUBECONFIG and then to '~/.kube/config'")
	newContextLocation := flag.String("context", "", "specifies the kubeconfig to be added. Default to './config.yaml'")
	newContextName := flag.String("name", "", "specifies the name to give to the context")
	logLevel := flag.Int("verbosity", 1, "set the log level")
	flag.Parse()

	// Set log level
	if *logLevel < 0 || *logLevel > 6 {
		*logLevel = 0
	}

	log.SetLevel(log.Level(*logLevel))

	// Set kubeconfig
	if *kubeConfigLocation == "" {
		log.Info("no location set for kubeconfig, trying KUBECONFIG environment variable")
		*kubeConfigLocation = os.Getenv("KUBECONFIG")
	}

	if *kubeConfigLocation == "" {
		log.Info("no location set for kubeconfig, trying default kubeconfig location")
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.WithError(err).Fatal("could not find home directory")
		}

		*kubeConfigLocation = fmt.Sprintf("%s/.kube/config", homedir)
	}

	kubeconfig, err := os.ReadFile(*kubeConfigLocation)
	if err != nil {
		displayHelp()
		log.WithError(err).Fatal("could not load kubeconfig file")
	}

	// set the new kubeconfig context
	if *newContextLocation == "" {
		*newContextLocation = "./config.yaml"
	}

	newContext, err := os.ReadFile(*newContextLocation)
	if err != nil {
		displayHelp()
		log.WithError(err).Fatal("could not load new context file")
	}

	// set the name of the new context
	if *newContextName == "" {
		log.Fatal("Please specify a name for the new context")
	}

	newKubeConfig := insertNewContext(*newContextName, kubeconfig, newContext)
	os.WriteFile(*kubeConfigLocation, newKubeConfig, os.FileMode(0666))
}

func displayHelp() {
	flag.CommandLine.Usage()
}

func insertNewContext(name string, currentConfigData, newConfigData []byte) []byte {
	log.Infof("inserting context '%s'", name)

	currentConfig := types.KubeConfig{}
	err := yaml.Unmarshal(currentConfigData, &currentConfig)
	if err != nil {
		log.WithError(err).Fatal("could not unmarshal kubeconfig")
	}

	configToAdd := types.KubeConfig{}
	err = yaml.Unmarshal(newConfigData, &configToAdd)
	if err != nil {
		log.WithError(err).Fatal("could not unmarshal new config")
	}

	if checkIfNameAlreadyExists(name, currentConfig) {
		log.WithField("name", name).Info("name already exists in kubeconfig")
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Do you want to overwrite the current context named '%s'", name),
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			log.Info("did not overwrite config")
			return currentConfigData
		}

		overwriteConfigByName(name, &currentConfig, &configToAdd)
	} else {
		appendConfigByName(name, &currentConfig, &configToAdd)
	}

	newData, err := yaml.Marshal(currentConfig)
	if err != nil {
		log.WithError(err).Fatal("could not marshal new config")
	}

	return newData
}

func checkIfNameAlreadyExists(name string, kubeconfig types.KubeConfig) bool {
	for _, cluster := range kubeconfig.Clusters {
		if cluster.Name == name {
			return true
		}
	}

	for _, user := range kubeconfig.Users {
		if user.Name == name {
			return true
		}
	}

	for _, context := range kubeconfig.Contexts {
		if context.Name == name {
			return true
		}
	}

	return false
}

func renameConfigName(name string, newConfig *types.KubeConfig) {
	newConfig.Clusters[0].Name = name
	newConfig.Users[0].Name = name
	newConfig.Contexts[0].Name = name
	newConfig.Contexts[0].Context.Cluster = name
	newConfig.Contexts[0].Context.User = name
}

func overwriteConfigByName(name string, currentConfig, newConfig *types.KubeConfig) {
	changedCluster := false
	changedUser := false
	changedContext := false

	for _, cluster := range currentConfig.Clusters {
		if cluster.Name == name {
			cluster.Cluster = newConfig.Clusters[0].Cluster
			changedCluster = true
		}
	}

	for _, user := range currentConfig.Users {
		if user.Name == name {
			user.User = newConfig.Users[0].User
			changedUser = true
		}
	}

	for _, context := range currentConfig.Contexts {
		if context.Name == name {
			context.Context = newConfig.Contexts[0].Context
			changedContext = true
		}
	}

	if !changedCluster || !changedUser || !changedContext {
		log.Fatalf("could not find a cluster, user or context with the name '%s'", name)
	}
}

func appendConfigByName(name string, currentConfig, newConfig *types.KubeConfig) {
	renameConfigName(name, newConfig)
	currentConfig.Clusters = append(currentConfig.Clusters, newConfig.Clusters[0])
	currentConfig.Users = append(currentConfig.Users, newConfig.Users[0])
	currentConfig.Contexts = append(currentConfig.Contexts, newConfig.Contexts[0])
}
