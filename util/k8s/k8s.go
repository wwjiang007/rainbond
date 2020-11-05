package k8s

import (
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	rainbondv1alpha1 "github.com/goodrain/rainbond-operator/pkg/apis/rainbond/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/reference"
)

func init() {
	utilruntime.Must(rainbondv1alpha1.AddToScheme(scheme.Scheme))
}

// NewClientset -
func NewClientset(kubecfg string) (kubernetes.Interface, error) {
	c, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		logrus.Errorf("error reading kube config file: %s", err.Error())
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		logrus.Error("error creating kube api client", err.Error())
		return nil, err
	}
	return clientset, nil
}

// NewClientsetOrDie new clientset or die
// used for who just wants a kubernetes clientset
func NewClientsetOrDie(kubecfg string) kubernetes.Interface {
	restConfig, err := NewRestConfig(kubecfg)
	if err != nil {
		panic(err)
	}
	return kubernetes.NewForConfigOrDie(restConfig)
}

// NewRestConfig new rest config
func NewRestConfig(kubecfg string) (restConfig *rest.Config, err error) {
	if kubecfg == "" {
		return InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubecfg)
}

//NewRestClient new rest client
func NewRestClient(restConfig *rest.Config) (*rest.RESTClient, error) {
	return rest.RESTClientFor(restConfig)
}

// InClusterConfig in cluster config
func InClusterConfig() (*rest.Config, error) {
	// Work around https://github.com/kubernetes/kubernetes/issues/40973
	// See https://github.com/coreos/etcd-operator/issues/731#issuecomment-283804819
	if len(os.Getenv("KUBERNETES_SERVICE_HOST")) == 0 {
		addrs, err := net.LookupHost("kubernetes.default.svc")
		if err != nil {
			panic(err)
		}
		os.Setenv("KUBERNETES_SERVICE_HOST", addrs[0])
	}
	if len(os.Getenv("KUBERNETES_SERVICE_PORT")) == 0 {
		os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	}
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewRainbondFilteredSharedInformerFactory -
func NewRainbondFilteredSharedInformerFactory(clientset kubernetes.Interface) informers.SharedInformerFactory {
	return informers.NewFilteredSharedInformerFactory(
		clientset, 30*time.Second, corev1.NamespaceAll, func(options *metav1.ListOptions) {
			options.LabelSelector = "creator=Rainbond"
		},
	)
}

// ExtractLabels extracts the service information from the labels
func ExtractLabels(labels map[string]string) (string, string, string, string) {
	if labels == nil {
		return "", "", "", ""
	}
	return labels["tenant_id"], labels["service_id"], labels["version"], labels["creater_id"]
}

// ListEventsByPod -
type ListEventsByPod func(kubernetes.Interface, *corev1.Pod) *corev1.EventList

// DefListEventsByPod default implementatoin of ListEventsByPod
func DefListEventsByPod(clientset kubernetes.Interface, pod *corev1.Pod) *corev1.EventList {
	ref, err := reference.GetReference(scheme.Scheme, pod)
	if err != nil {
		logrus.Errorf("Unable to construct reference to '%#v': %v", pod, err)
		return nil
	}
	ref.Kind = ""
	if _, isMirrorPod := pod.Annotations[corev1.MirrorPodAnnotationKey]; isMirrorPod {
		ref.UID = types.UID(pod.Annotations[corev1.MirrorPodAnnotationKey])
	}
	events, _ := clientset.CoreV1().Events(pod.GetNamespace()).Search(scheme.Scheme, ref)
	return events
}
