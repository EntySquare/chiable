package core

import (
	"chiable/lib"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

// static strategy
type StaticStrategy struct {
	PoolKey   string
	FarmerKey string
	count     int64
	client    *kubernetes.Clientset
}

func NewStaticStrategy(farmerKey, poolKey string) *StaticStrategy {

	kclient, err := lib.GetK8sClientSet()
	if err != nil {
		panic(err)
	}

	s := &StaticStrategy{
		PoolKey:   poolKey,
		FarmerKey: farmerKey,
		client:    kclient,
	}
	return s
}

func (s *StaticStrategy) Run(n int64) error {
	s.count = n
	var i int64
	for {
		if i == n {
			return nil
		}
		// do plot
		s.plot()
		i++
		time.Sleep(time.Second)
	}
}

func (s *StaticStrategy) plot() {
	sf := lib.StartPlotAffinity()
	limitList := corev1.ResourceList{}
	requestList := corev1.ResourceList{}
	limitList["cpu"] = resource.MustParse("2000m")
	requestList["cpu"] = resource.MustParse("2000m")
	limitList["memory"] = resource.MustParse("20Gi")
	requestList["memory"] = resource.MustParse("20Gi")
	jbname := "entysquare-job-plot-farmer-" + s.FarmerKey + "-" + rand.String(5)
	fmt.Println("run job : " + jbname)
	jb := lib.GetJob(jbname, 1, 30000, sf, limitList, requestList, s.FarmerKey, s.PoolKey)
	_, err := s.client.BatchV1().Jobs("default").Create(context.TODO(), jb, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

// dynamic strategy
type DynamicStrategy struct {
}

func NewDynamicStrategy() *DynamicStrategy {
	return nil
}

func (s *DynamicStrategy) Run() error {
	return nil
}
