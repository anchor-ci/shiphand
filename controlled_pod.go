package main

import (
	"errors"
	"fmt"
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var PODS_NAMESPACE = "default"

type ControlledPod struct {
	Pod *apiv1.Pod
	Id  string
}

func NewControlledPod(name, image string) (ControlledPod, error) {
	instance := ControlledPod{}

	// TODO: Uncouple the kube client stuff from here
	clientset := getKubernetesClient("./config")

	req := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"job-id": name,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "job",
					Image: image,
					// TODO: Make this a configurable timeout v
					Command: []string{"sleep", "1800"},
				},
			},
		},
	}

	pod, err := clientset.CoreV1().Pods(PODS_NAMESPACE).Create(req)

	instance.Pod = pod
	instance.Id = name

	return instance, err
}

// Blocks until pod is in a started state,
// and ready to accept commands. Returns an
// error if one occurs during pod startup
func (c *ControlledPod) WaitForStart() error {
	// TODO: Uncouple the kube client stuff from here
	clientset := getKubernetesClient("./config")

	api := clientset.CoreV1()
	options := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-id=%s", c.Id),
	}

	watcher, err := api.Pods(PODS_NAMESPACE).Watch(options)

	if err != nil {
		return err
	}

	ch := watcher.ResultChan()

	for event := range ch {
		if pod, ok := event.Object.(*apiv1.Pod); ok {
			if pod.Status.Phase == "Running" {
				return nil
			}
		} else {
			return errors.New("Failed to grab event")
		}
	}

	return errors.New("Timed out")
}

func (c *ControlledPod) RunCommand(command string) (int, error) {
	clientset := getKubernetesClient("./config")
	restClient := clientset.CoreV1().RESTClient()
	req := restClient.Post().
		Resource("pods").
		Name(c.Id).
		Namespace(PODS_NAMESPACE).
		SubResource("exec").
		Param("container", "job")

	req.VersionedParams(&apiv1.PodExecOptions{
		Container: "job",
		Command:   []string{"/bin/sh", "-c", command},
		Stdin:     false,
		Stdout:    true,
	}, scheme.ParameterCodec)

	restconf, cfgErr := clientcmd.BuildConfigFromFlags("", "./config")

	if cfgErr != nil {
		return 1, cfgErr
	}

	exec, err := remotecommand.NewSPDYExecutor(restconf, "POST", req.URL())
	if err != nil {
		return 1, err
	}

	opts := remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: os.Stdout,
		Stderr: nil,
		Tty:    false,
	}

	execErr := exec.Stream(opts)
	if execErr != nil {
		return 1, execErr
	}

	return 0, nil
}
