package lib

import corev1 "k8s.io/api/core/v1"

func StartPlotAffinity() corev1.NodeAffinity {
	return corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      "nodetype",
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{"plot"},
						},
					},
				},
			},
		},
	}
}

func StartTestAffinity() corev1.NodeAffinity {
	return corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      "nodetype",
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{"test"},
						},
					},
				},
			},
		},
	}
}
