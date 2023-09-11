package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestGetCacheInvalid(t *testing.T) {
	duration := time.Second
	cache := NewCache(duration)

	cases := []struct {
		key  string
		data []byte
	}{
		{
			key:  "",
			data: []byte("test data"),
		},
		{
			key:  "www.testwebsite.com",
			data: []byte(""),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("running test case %d", i), func(t *testing.T) {
			err := cache.Add(&c.key, c.data)
			if err == nil {
				t.Errorf("expected error when setting cache")
			}
		})
	}
}

func TestGetCacheValid(t *testing.T) {
	duration := time.Second
	cache := NewCache(duration)

	cases := []struct {
		key  string
		data []byte
	}{
		{
			key:  "www.testwebsite.com",
			data: []byte("test data"),
		},
		{
			key:  "www.testwebsite2.com",
			data: []byte("test data 2"),
		},
		{
			key:  "www.testwebsite3.com",
			data: []byte("x"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("running test case %d", i), func(t *testing.T) {
			err := cache.Add(&c.key, c.data)
			if err != nil {
				t.Errorf("error when setting cache: %s", err)
			}
			data, ok := cache.Get(&c.key)
			if !ok {
				t.Errorf("expected data to be present in cache")
			}
			if string(data) != string(c.data) {
				t.Errorf("expected data to be %s, got %s", string(c.data), string(data))
			}
		})
	}
}

func TestCacheReap(t *testing.T) {
	duration := time.Millisecond * 200
	cache := NewCache(duration)

	key := "www.testwebsite.com"
	data := []byte("test data")

	err := cache.Add(&key, data)
	if err != nil {
		t.Errorf("error when setting cache: %s", err)
	}

	time.Sleep(duration * 2)

	_, ok := cache.Get(&key)
	if ok {
		t.Errorf("expected data to be reaped from cache")
	}
}
