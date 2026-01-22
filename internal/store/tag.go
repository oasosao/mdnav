package store

import "mdnav/internal/conf"

// GetConfigTags 获取标签数据
func GetConfigTags() map[string]string {

	tagsMap := make(map[string]string)

	for k, v := range conf.Config().GetStringMapString("tags") {
		tagsMap[v] = k
	}

	return tagsMap
}
