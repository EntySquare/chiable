package lib

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/scheduling/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
)

func getJob(jobName string, jobParallelism int32, deleteJobAfterFinishSec int32, nodeAffinity corev1.NodeAffinity,
	taskType string, sectorMinerIdStr string, sectorNumberStr string, proofTypeStr string, sectorSizeStr string, paramStr string, limitList corev1.ResourceList, requestList corev1.ResourceList) *batchv1.Job {

	notiIP := os.Getenv("CONTROLLER_NOTIFICATION")
	sectorDataHostPath := os.Getenv("SECTOR_DATA_HOST_PATH")
	lotusminerPath, ok := os.LookupEnv("LOTUSMINER_PATH")
	if ok && taskType == spec.RECEIVE_COPY {
		sectorDataHostPath = lotusminerPath
	}
	filtabSealerImage := os.Getenv("FILTAB_SEALER_IMAGE")
	filtabServiceAccountName := os.Getenv("FILTAB_SERVICE_ACCOUNT")
	filtabNatsServer := os.Getenv("NATS_SERVER")
	reserveGiBForCopySector := os.Getenv("RESERVE_GIB_FOR_COPY_SECTOR")

	filtabMinerIP := os.Getenv("MINER_IP")

	filtabAddPiecePortStr := os.Getenv("FILTAB_ADDPIECE_PORT")
	filtabAddPiecePort, err := strconv.ParseInt(filtabAddPiecePortStr, 10, 32)
	if err != nil {
		panic("Filtab log: FILTAB_ADDPIECE_PORT ParseInt Error")
	}

	filtabCopyPortStr := os.Getenv("FILTAB_COPY_PORT")
	filtabCopyPort, err := strconv.ParseInt(filtabCopyPortStr, 10, 32)
	if err != nil {
		panic("Filtab log: FILTAB_COPY_PORT ParseInt Error")
	}

	filecoinTmpDir := "/var/tmp"

	sectorDataDirHostType := corev1.HostPathDirectoryOrCreate
	jobRestartPolicy := corev1.RestartPolicyNever

	sectorDirName := os.Getenv("ALL_SECTORS_DIR")

	//Dont Restart a failed job pod!!!
	zeroBackoffLimitIsRetryTimeForNeverRestartFailedJobPod := int32(3)

	jobLabelMaps := map[string]string{
		"app":          "filtab-sealer",
		"phase":        "sealing",
		"minerid":      sectorMinerIdStr,
		"sectornumber": sectorNumberStr,
		"tasktype":     taskType,
	}

	priorityClassName := ""
	if taskType == spec.COMMIT_2 {
		priorityClassName = getPriorityClassName(sectorMinerIdStr, "0")
	} else if taskType == spec.COMMIT_1 {

		sn, err := strconv.ParseInt(sectorNumberStr, 10, 64)
		if err != nil {
			sn = 2
		}

		sn = sn - 2

		if sn < 0 {
			sn = 0
		}

		snstr := strconv.FormatInt(sn, 10)

		priorityClassName = getPriorityClassName(sectorMinerIdStr, snstr)

	} else if taskType == spec.PRE_COMMIT_2 {
		sn, err := strconv.ParseInt(sectorNumberStr, 10, 64)
		if err != nil {
			sn = 1
		}
		sn = sn - 1
		if sn < 0 {
			sn = 0
		}
		snstr := strconv.FormatInt(sn, 10)

		priorityClassName = getPriorityClassName(sectorMinerIdStr, snstr)
	} else {
		priorityClassName = getPriorityClassName(sectorMinerIdStr, sectorNumberStr)
	}

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
					ServiceAccountName: filtabServiceAccountName,
					RestartPolicy:      jobRestartPolicy,
					Volumes: []corev1.Volume{
						{
							Name: "sectordatadir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: sectorDataHostPath,
									Type: &sectorDataDirHostType,
								},
							},
						},
						{
							Name: "filecointmpdir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: filecoinTmpDir,
									Type: &sectorDataDirHostType,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "filtab-sealer",
							Image: filtabSealerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "copyport",
									ContainerPort: int32(filtabCopyPort),
								},
								{
									Name:          "addpieceport",
									ContainerPort: int32(filtabAddPiecePort),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "sectordatadir",
									MountPath: "/tmp/.lotusminer",
								},
								{
									Name:      "filecointmpdir",
									MountPath: "/var/tmp",
								},
							},
							Command: []string{"/bin/bash", "-c", "/filtabsealer"},
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
									Name:  "FILTAB_K8S_CONFIG_IN_CLUSTER",
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
									Name:  "SECTOR_DIR",
									Value: sectorDirName,
								},
								{
									Name:  "SECTOR_MINER_ID",
									Value: sectorMinerIdStr,
								},
								{
									Name:  "MINER_IP",
									Value: filtabMinerIP,
								},
								{
									Name:  "SECTOR_NUMBER",
									Value: sectorNumberStr,
								},
								{
									Name:  "TASK_TYPE",
									Value: taskType,
								},
								{
									Name:  "TASK_SECTOR_TYPE",
									Value: sectorSizeStr,
								},
								{
									Name:  "PROOF_TYPE",
									Value: proofTypeStr,
								},
								{
									Name:  "NATS_SERVER",
									Value: filtabNatsServer,
								},
								{
									Name:  "EVENTING",
									Value: "true",
								},
								{
									Name:  "PARAMS",
									Value: paramStr,
								},
								{
									Name:  "RESERVE_GIB_FOR_COPY_SECTOR",
									Value: reserveGiBForCopySector,
								},
								{
									Name:  "FIL_PROOFS_USE_GPU_TREE_BUILDER",
									Value: "1",
								},
								{
									Name:  "FIL_PROOFS_USE_GPU_COLUMN_BUILDER",
									Value: "1",
								},
								{
									Name:  "BELLMAN_CUSTOM_GPU",
									Value: "GeForce RTX 2080 Ti:4352",
								},
								{
									Name:  "RUST_BACKTRACE",
									Value: "1",
								},
								{
									Name:  "CONTROLLER_NOTIFICATION",
									Value: notiIP,
								},
								{
									Name:  "FIL_PROOFS_USE_MULTICORE_SDR",
									Value: "0",
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

func getPriorityClassName(sectorMinerIdStr string, sectorNumberStr string) string {
	return "priority" + "-" + sectorMinerIdStr + "-" + sectorNumberStr
}

func getPriorityClass(priorityClassName string, sectorNumberStr string) v1.PriorityClass {

	sectorNumber, err := strconv.ParseInt(sectorNumberStr, 10, 32)
	if err != nil {
		panic("Filtab log: getPriorityClass sectorNumberStr ParseInt Error")
	}

	sectorNumberInt32 := int32(sectorNumber)

	K8sMaxPriorityClassValue := int32(1000 * 1000 * 1000)

	modSectorNumber := sectorNumberInt32 % K8sMaxPriorityClassValue

	sectorNumberPriorityValue := K8sMaxPriorityClassValue - modSectorNumber

	priority := v1.PriorityClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: priorityClassName,
		},
		Value:         sectorNumberPriorityValue,
		GlobalDefault: false,
	}

	return priority
}
