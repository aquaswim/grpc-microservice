package methodmetamap

import (
	"fmt"
	"strings"
)

type MethodMeta struct {
	NeedAuth bool
}

type MethodMetaMap map[string]MethodMeta

func (m MethodMetaMap) Get(key string) (*MethodMeta, bool) {
	val, ok := m[key]
	if !ok {
		return nil, ok
	}
	return &val, ok
}

func (m MethodMetaMap) String() string {
	sb := strings.Builder{}
	for key, value := range m {
		sb.WriteString(fmt.Sprintf("%s: {%+v}\n", key, value))
	}
	return sb.String()
}
