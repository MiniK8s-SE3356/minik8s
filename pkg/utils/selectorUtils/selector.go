package selectorUtils

import (
	selectordef "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
)

func canSelect(pod_labels *map[string]string, selector_labels *map[string]string) bool {
	// 对于筛选器中的每个标签做检查
	for k, sv := range *selector_labels {
		if pv, exist := (*pod_labels)[k]; !exist {
			// 如果pod标签组中不存在此标签，返回false
			return false
		} else {
			if pv != sv {
				return false
			}
		}

	}
	return true
}

func SelectPodNameList(sel *selectordef.Selector, pods *httpobject.HTTPResponse_GetAllPod) []string {
	if sel.MatchLabels == nil || len(sel.MatchLabels) == 0 {
		return []string{}
	}
	result := []string{}
	for name, pod_item := range *pods {
		if canSelect(&pod_item.Metadata.Labels, &sel.MatchLabels) {
			result = append(result, name)
		}
	}
	return result
}
