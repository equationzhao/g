package cached

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCacheMap(t *testing.T) {
	cacheMap := NewCacheMap[int, int](10)
	if cacheMap == nil {
		t.Error("NewCacheMap failed")
	}
}

func TestMap_Keys(t *testing.T) {
	cacheMap := NewCacheMap[int, int](10)
	cacheMap.Set(1, 2)
	cacheMap.Set(2, 3)
	keys := cacheMap.Keys()
	if len(keys) != 2 {
		t.Error("Keys failed")
	}
	assert.ElementsMatch(t, []int{1, 2}, keys)
}

func TestMap_Values(t *testing.T) {
	cacheMap := NewCacheMap[int, int](10)
	cacheMap.Set(1, 2)
	cacheMap.Set(2, 3)
	values := cacheMap.Values()
	if len(values) != 2 {
		t.Error("Values failed")
	}
	assert.ElementsMatch(t, []int{2, 3}, values)
}

func TestMap_Pairs(t *testing.T) {
	cacheMap := NewCacheMap[int, int](10)
	cacheMap.Set(1, 2)
	cacheMap.Set(2, 3)
	pairs := cacheMap.Pairs()
	if len(pairs) != 2 {
		t.Error("Pairs failed")
	}
	assert.ElementsMatch(t, []Pair[int, int]{{First: 1, Second: 2}, {First: 2, Second: 3}}, pairs)
}

func TestMakePair(t *testing.T) {
	pair := MakePair(1, 2)
	if pair.First != 1 || pair.Second != 2 {
		t.Error("MakePair failed")
	}
}

func TestPair(t *testing.T) {
	pair := Pair[int, int]{}
	if pair.First != 0 || pair.Second != 0 {
		t.Error("Pair failed")
	}

	pair.Set(1, 2)
	if pair.First != 1 || pair.Second != 2 {
		t.Error("Pair Set failed")
	}

	if pair.Key() != 1 {
		t.Error("Pair Key failed")
	}

	if pair.Value() != 2 {
		t.Error("Pair Value failed")
	}
}
