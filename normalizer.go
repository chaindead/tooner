package main

import (
	"encoding/json"
	"reflect"
	"strings"
)

func normalizeAny(obj any) any {
	switch tObj := obj.(type) {
	case map[string]any:
		for k, v := range tObj {
			tObj[k] = normalizeAny(v)
		}

		return tObj
	case []any:
		for i, v := range tObj {
			tObj[i] = normalizeAny(v)
		}
		return normalizeSlice(tObj)
	}

	return obj
}

func normalizeSlice(obj []any) any {
	newObj := make([]map[string]any, len(obj))
	for i, v := range obj {
		nv, ok := v.(map[string]any)
		if !ok {
			return obj
		}

		replacer := strings.NewReplacer(
			`"`, ``,
			`,`, `;`,
			`]`, `)`,
			`[`, `(`,
			`}`, `)`,
			`{`, `(`,
			`{`, `(`,
		)

		for kv, vv := range nv {
			switch vv.(type) {
			case []any, map[string]any:
				jVal, _ := json.Marshal(vv)
				nv[kv] = replacer.Replace(string(jVal))
			}
		}

		newObj[i] = nv
	}

	keys := make(map[string]int)
	types := make(map[string]reflect.Type)
	for _, m := range newObj {
		for k, v := range m {
			keys[k] += 1

			types[k] = reflect.TypeOf(v)
		}
	}

	for _, m := range newObj {
		for k := range keys {
			_, ok := m[k]
			if ok {
				continue
			}

			m[k] = reflect.Zero(types[k]).Interface()

		}
	}

	return newObj
}
