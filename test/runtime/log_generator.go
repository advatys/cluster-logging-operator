package runtime

import (
	"fmt"
	"github.com/openshift/cluster-logging-operator/internal/utils"
	"time"

	"github.com/openshift/cluster-logging-operator/internal/runtime"

	corev1 "k8s.io/api/core/v1"
)

// NewLogGenerator creates a pod that will print `count` lines to stdout, waiting for
// `delay` between each line.  Lines are of the form "<timestamp> [n] `message`"
// where n is the number of lines output so far. Once done printing the pod will
// be idle but will not exit until deleted.
func NewLogGenerator(namespace, name string, count int, delay time.Duration, message string) *corev1.Pod {
	cmd := fmt.Sprintf(`i=0; while [ $i -lt %v ]; do echo "$(date) [$i]: %v"; i=$((i+1)); sleep %f; done; sleep infinity`, count, message, delay.Seconds())
	l := runtime.NewPod(namespace, "log-generator", corev1.Container{
		Name:    name,
		Image:   "quay.io/quay/busybox",
		Command: []string{"sh", "-c", cmd},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: utils.GetBool(false),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
			},
		},
	})
	l.Spec.RestartPolicy = corev1.RestartPolicyNever
	l.Spec.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot: utils.GetBool(true),
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}
	return l
}

// NewOneLineLogGenerator creates a pod that will print given lines to stdout.
// Once done printing the pod will be idle but will not exit until deleted.
func NewOneLineLogGenerator(namespace, containerName, message string) *corev1.Pod {
	cmd := fmt.Sprintf(`echo "%v"; sleep infinity`, message)
	l := runtime.NewPod(namespace, "log-generator", corev1.Container{
		Name:    containerName,
		Image:   "quay.io/quay/busybox",
		Command: []string{"sh", "-c", cmd}},
	)
	l.Spec.RestartPolicy = corev1.RestartPolicyNever
	return l
}
