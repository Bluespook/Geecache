package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Error("not matched")
	}
}

func TestGet(t *testing.T) {
	loadCounter := make(map[string]int, len(DB))
	g := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slowDB] search key", key)
		if v, ok := DB[key]; ok {
			if _, ok2 := loadCounter[key]; !ok2 {
				loadCounter[key] = 0
			}
			loadCounter[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range DB {
		if ans, err := g.Get(k); err != nil || ans.String() != v {
			t.Error("load from db failed:  ", err)
		}
		if _, err := g.Get(k); err != nil || loadCounter[k] > 1 {
			t.Error("cache miss:", k)
		}
	}
	if _, err := g.Get("John"); err == nil {
		t.Fatal("should not get result from John")
	}
}
