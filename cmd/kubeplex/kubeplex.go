package main

import (
	"log"
	"os"

	"github.com/munnerz/kube-plex/pkg/kube-plex"
	kubeplexController "github.com/munnerz/kube-plex/pkg/controller"
	"github.com/munnerz/kube-plex/pkg/executor"
	"github.com/munnerz/kube-plex/pkg/worker"
)

func main() {
	log.Println("Initializing kubernetes API client.")
	kubeClient, err := kubeplex.NewKubeClient()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating kube-plex.")
	controller := kubeplex.NewController(kubeClient)

	switch os.Getenv("KUBEPLEX_ENV") {
	case "executor":
		log.Println("Running executor.")
		err = executor.Run(controller)
	case "worker":
		log.Println("Running worker.")
		err = worker.Run(controller)
	case "controller":
		log.Println("Running controller.")
		err = kubeplexController.Run(controller)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running kube-plex.")
	go controller.Run()


	log.Println("Waiting for stop signal.")
	<-controller.Stop
}
