package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fvanmaldegem/tip-k8ctx/pkg/app"

	log "github.com/sirupsen/logrus"
)

var (
	kubeConfigLocation    = ""
	newKubeConfigLocation = ""
	newContextName        = ""
	forceOverwrite        = false
	verbose               = true
)

func main() {
	parseFlags()
	setLogLevel()

	app := app.New(kubeConfigLocation, newKubeConfigLocation, newContextName, forceOverwrite)
	app.Run()
}

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <context-name>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&kubeConfigLocation, "kubeconfig", "", "specify the location of your kubeconfig")
	flag.StringVar(&newKubeConfigLocation, "new-config", "", "specify the location of the kubeconfig you want to add")
	flag.BoolVar(&forceOverwrite, "f", false, "force the overwrite of a context")
	flag.BoolVar(&verbose, "v", false, "set the output to verbose")
	flag.Parse()

	newContextName = flag.Arg(0)
}

func setLogLevel() {
	log.SetLevel(log.InfoLevel)
	if verbose {
		log.SetLevel(log.TraceLevel)
	}
}
