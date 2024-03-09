package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("value1"))
	if val, _ := lru.Get("key1"); string(val.(String)) != "value1" {
		t.Error("not matched")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("should not exist key")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2 := "key1", "key2"
	v1, v2 := String("val1"), String("val2")
	length := int64(len(k1) + v1.Len())
	lru := New(length, nil)
	lru.Add(k1, v1)
	lru.Add(k2, v2)
	if _, ok := lru.Get(k1); ok {
		t.Error("k1 should not exist")
	}
	if _, ok := lru.Get(k2); !ok {
		t.Error("k2 should exist")
	}
	fmt.Println(lru.Len())
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Get("k2")
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k3"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
