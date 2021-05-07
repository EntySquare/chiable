package lib

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func GetJob(jobName string, jobParallelism int32, deleteJobAfterFinishSec int32, nodeAffinity corev1.NodeAffinity,
	limitList corev1.ResourceList, requestList corev1.ResourceList, farmerKey string, poolKey string) *batchv1.Job {
	entyChiaImage := os.Getenv("ENTY_CHIA_IMAGE")
	entyServiceAccountName := os.Getenv("ENTY_SERVICE_ACCOUNT")
	entyNatsServer := os.Getenv("NATS_SERVER")

	sectorDataHostPath := os.Getenv("SECTOR_DATA_HOST_PATH")
	sectorDataDirHostType := corev1.HostPathDirectoryOrCreate
	jobRestartPolicy := corev1.RestartPolicyNever

	dirName := os.Getenv("ALL_DIR")

	//Dont Restart a failed job pod!!!
	zeroBackoffLimitIsRetryTimeForNeverRestartFailedJobPod := int32(3)

	jobLabelMaps := map[string]string{
		"app":   "enty-chia",
		"phase": "test",
	}

	priorityClassName := ""

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   jobName,
			Labels: jobLabelMaps,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &zeroBackoffLimitIsRetryTimeForNeverRestartFailedJobPod,
			TTLSecondsAfterFinished: &deleteJobAfterFinishSec,
			Parallelism:             &jobParallelism,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: jobLabelMaps,
				},
				Spec: corev1.PodSpec{
					PriorityClassName: priorityClassName,
					Affinity: &corev1.Affinity{
						NodeAffinity: &nodeAffinity,
					},
					ServiceAccountName: entyServiceAccountName,
					RestartPolicy:      jobRestartPolicy,
					Volumes: []corev1.Volume{
						{
							Name: "chiadatadir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: sectorDataHostPath,
									Type: &sectorDataDirHostType,
								},
							},
						},
						{
							Name: "chiatmpdir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: dirName,
									Type: &sectorDataDirHostType,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "enty-chia",
							Image: entyChiaImage,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "chiadatadir",
									MountPath: "/tmp/chia",
								},
								{
									Name:      "chiatmpdir",
									MountPath: "/var/tmp",
								},
							},
							Command: []string{". /chia-blockchain/activate", "chia plots create -f " + farmerKey + " -p " + poolKey + " -d /tmp/plots -t /tmp/cache -n 2 -r 2 -b 20000 &>plotting$i.log"},
							Resources: corev1.ResourceRequirements{
								Limits:   limitList,
								Requests: requestList,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "TMP_PATH",
									Value: "/tmp",
								},
								{
									Name:  "ENTY_K8S_CONFIG_IN_CLUSTER",
									Value: "true",
								},
								{
									Name: "JOB_NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "JOB_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name:  "NATS_SERVER",
									Value: entyNatsServer,
								},
								{
									Name:  "EVENTING",
									Value: "true",
								},
							},
						},
					},
				},
			},
		},
	}

	return job
}
