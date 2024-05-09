package runtime_test

import (
	"testing"

	minik8s_apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	minik8s_runtime "github.com/MiniK8s-SE3356/minik8s/pkg/runtime"
	minik8s_types "github.com/MiniK8s-SE3356/minik8s/pkg/types"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	test_pod := &minik8s_apiobject.Pod{
		Basic: minik8s_apiobject.Basic{
			APIVersion: "v1",
			Kind:       "Pod",
			Metadata: minik8s_apiobject.Metadata{
				UUID:      uuid.New().String(),
				Name:      "test-pod",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		Spec: minik8s_apiobject.PodSpec{
			NodeName: "node1",
			Containers: []minik8s_types.Container{
				{
					Name:  "nginx-container",
					Image: "nginx:latest",
					Ports: []minik8s_types.ContainerPort{
						{
							HostPort:      80,
							ContainerPort: 80,
						},
					},
				},
				{
					Name:  "redis-container",
					Image: "redis:latest",
					Ports: []minik8s_types.ContainerPort{
						{
							HostPort:      6379,
							ContainerPort: 6379,
						},
					},
				},
			},
		},
		Status: minik8s_types.PodStatus{
			Phase: minik8s_types.PodPending,
		},
	}

	runtimeManager := &minik8s_runtime.RuntimeManager{}
	_, err := runtimeManager.CreatePod(test_pod)
	if err != nil {
		println("Error creating pod")
		panic(err)
	}
}
