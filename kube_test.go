package chiable

import (
	"chiable/lib"
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"log"
	"testing"
)

func TestK8sJob(t *testing.T) {

	var K8sClient = func() *kubernetes.Clientset {

		kclient, err := lib.GetK8sClientSet()
		if err != nil {
			log.Fatal(err)
		}
		return kclient

	}()

	sf := lib.StartPlotAffinity()
	limitList := corev1.ResourceList{}
	requestList := corev1.ResourceList{}
	limitList["cpu"] = resource.MustParse("2000m")
	requestList["cpu"] = resource.MustParse("2000m")
	limitList["memory"] = resource.MustParse("20Gi")
	requestList["memory"] = resource.MustParse("20Gi")
	jb := lib.GetJob("test-jb-"+rand.String(5), 1, 300, sf, limitList, requestList, "a6db5fdada01681be184e7499465df271a55db49c765cf67569bd35194922dca8df63f31ec7d1efd1825b84d86e5498d", "89f7d8e6d8a8887ab7aedc437f1582c4aaa41e73878b5857f061b8658b2bc0a02f163f2d31799e4ad4929e069ae69766")
	_, err := K8sClient.BatchV1().Jobs("default").Create(context.TODO(), jb, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}
