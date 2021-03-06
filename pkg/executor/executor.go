package executor

import (
	"fmt"
	"os"

	"github.com/munnerz/kube-plex/pkg/kube-plex"
	ptjv1 "github.com/munnerz/kube-plex/pkg/apis/ptj/v1"

	"k8s.io/client-go/tools/cache"
)

func Run(controller kubeplex.Controller) error {
	env := os.Environ()

	args := os.Args
	kubeplex.RewriteArgs(args)

	dir, _ := os.Getwd()

	ptj := kubeplex.GeneratePlexTranscodeJob(args, env, dir)
	new_ptj, err := kubeplex.CreatePlexTranscodeJob(&ptj, controller.KubeClient)
	if err != nil {
		return err
	}

	controller.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) {
			updated := new.(*ptjv1.PlexTranscodeJob)

			if updated.ObjectMeta.Name != new_ptj.ObjectMeta.Name {
				return
			}

			fmt.Println(updated.ObjectMeta.Name, updated.Status.State)
			if updated.Status.State == ptjv1.PlexTranscodeStateCompleted {
				close(controller.Stop)
			}
		},
	})

	return nil
}
