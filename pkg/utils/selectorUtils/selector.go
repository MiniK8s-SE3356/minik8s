package selectorUtils

import (
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/persistVolume"
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
	result := []string{}
	for name, pod_item := range *pods {
		if canSelect(&pod_item.Metadata.Labels, &sel.MatchLabels) {
			result = append(result, name)
		}
	}
	return result
}

func SelectPVNameList(sel *selectordef.Selector, pv_list *(map[string]persistVolume.PersistVolume)) []string {
	result := []string{}
	for pv_name, pv_item := range *pv_list {
		if canSelect(&pv_item.Metadata.Labels, &sel.MatchLabels) {
			result = append(result, pv_name)
		}
	}
	return result
}

func IsStrInArray(traget string, arr *([]string)) bool {
	for _, v := range *arr {
		if v == traget {
			return true
		}
	}
	return false
}

func IsSameSet(arr1 *([]string), arr2 *([]string)) bool {
	if len(*arr1) != len(*arr2) {
		return false
	}
	for _, v1 := range *arr1 {
		ishave := false
		for _, v2 := range *arr2 {
			if v1 == v2 {
				ishave = true
			}
		}
		if !ishave {
			return false
		}
	}
	return true
}

func AddToMayNilMap(key string, value string, m *(map[string][]string)) {
	if v_list, e := (*m)[key]; e {
		v_list = append(v_list, value)
		(*m)[key] = v_list
	} else {
		new_list := []string{value}
		(*m)[key] = new_list
	}
}
