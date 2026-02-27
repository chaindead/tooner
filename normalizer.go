package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func isZero[T comparable](t T) bool {
	var zero T

	rT := reflect.TypeOf(t)
	if rT == nil {
		return t == zero
	}

	if rT.Kind() == reflect.Slice {
		return t == zero || reflect.ValueOf(t).Len() == 0
	}

	if rT.Kind() == reflect.Map {
		return t == zero || reflect.ValueOf(t).Len() == 0
	}

	return t == zero
}
func removeEmpty(obj any) any {
	jsonObj, _ := json.Marshal(obj)
	fmt.Printf("CALL removeEmpty %T: %s\n", obj, string(jsonObj))
	switch t := obj.(type) {
	case map[string]any:
		fmt.Printf("\tCALL SWICH MAP %T: %#v\n", t, t)
		for k, v := range obj.(map[string]any) {
			nV := removeEmpty(v)
			if isZero(nV) {
				fmt.Println("\t\tZERO MAP ", nV)
				delete(t, k)
				continue
			}

			t[k] = nV
		}

		if len(t) == 0 {
			fmt.Println("\t\tZERO MAP NIL")
			return nil
		}

		return t
	case []any:
		fmt.Printf("\tCALL SLICE %T: %#v\n", t, t)
		var nT []any
		for _, v := range t {
			nV := removeEmpty(v)
			if isZero(nV) {
				fmt.Println("\t\tZERO SLICE ", nV)
				continue
			}

			nT = append(nT, nV)
		}

		if len(t) == 0 {
			fmt.Println("ZERO SLICE NIL")
			return nil
		}

		return nT
	default:
	}

	return obj
}

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
				sV := removeEmpty(vv)
				if sV == nil {
					nv[kv] = "nil"
					continue
				}

				jVal, _ := json.Marshal(sV)
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
