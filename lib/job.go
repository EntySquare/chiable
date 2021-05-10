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

	sectorDataDirHostType := corev1.HostPathDirectoryOrCreate
	jobRestartPolicy := corev1.RestartPolicyNever

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
					RestartPolicy: jobRestartPolicy,
					Volumes: []corev1.Volume{
						{
							Name: "chiadatadir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/root/rplots",
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
									MountPath: "/root/rplots",
								},
							},
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{". ./activate && chia init && chia plots create -f " + farmerKey + " -p " + poolKey + " -d /root/rplots/plots -t /root/rplots/cache -b 10000"},
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
							},
						},
					},
				},
			},
		},
	}

	return job
}
