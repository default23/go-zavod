package cache

import (
	"sync"
)

/*
Дано:
	InMemoryCache - потоко-безопасная реализация Key-Value кэша, хранящая данные в оперативной памяти
Задача:
	1. Реализовать метод GetOrSet, предоставив следующие гарантии:
		- Значение каждого ключа будет вычислено ровно 1 раз
		- Конкурентные обращения к существующим ключам не блокируют друг друга
	2. Покрыть его тестами
*/

// ----------------------------------------------

type (
	Key   = string
	Value = string
)

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

// ----------------------------------------------

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
}

// (ВОПРОС) не правильнее было бы возвращать интерфейс `Cache` в фабричном методе
// для того, что бы на этапе компиляции проверять что структура
// правильно наследует интерфейс? Так получается что особого смысла от интерфейса нет
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

// GetOrSet возвращает значение ключа в случае его существования.
// Иначе, вычисляет значение ключа при помощи valueFn, сохраняет его в кэш и возвращает это значение.
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {

	if v, ok := cache.Get(key); ok {
		return v
	}

	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	value := valueFn()
	cache.data[key] = value
	return value
}
