package url

import (
	"fmt"
	"strings"
)

type Filter struct {
	ignoreUrls map[string]interface{}
}

func NewUrlFilter(ignoreUrls []string) *Filter {
	ignoreUrlsMap := make(map[string]interface{})

	for _, url := range ignoreUrls {
		path := strings.Split(strings.TrimPrefix(url, "/"), "/")

		for i, p := range path {
			tmpMap := &ignoreUrlsMap

			for j := 0; j < i; j++ {
				tmpMap = (*tmpMap)[path[j]].(*map[string]interface{})
			}

			if i == len(path)-1 {
				(*tmpMap)[p] = true
			} else {
				if _, ok := (*tmpMap)[p]; !ok {
					m := make(map[string]interface{})
					(*tmpMap)[p] = &m
				}
			}
		}
	}

	return &Filter{ignoreUrlsMap}
}

func (f *Filter) IsIgnore(path string) (res bool, err error) {
	m := f.ignoreUrls
	urlPathMap := &m

	for _, p := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		val := (*urlPathMap)[p]

		if val == nil {
			if v, ok := (*urlPathMap)["*"]; ok {
				val = v
			} else {
				break
			}
		}

		switch val.(type) {

		case bool:
			res = val.(bool)

			return

		case *map[string]interface{}:
			urlPathMap = val.(*map[string]interface{})

			continue

		default:
			err = fmt.Errorf("parse url error")

			return
		}
	}

	return
}
