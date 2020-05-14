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
	
	Пустой кэш
1000 горутин
Каждая обращается к случайному ключу от 1 до 10
Когда все горутины отработали, valueFn должна быть вызвана ровно 10 раз

Тесты на конкурентность рекомендуем запускать с флагом -count 1000
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

	// лочим чтение и проверяем присутствует ли в кэше значение
	cache.dataMutex.RLock()
	if value, exists := cache.data[key]; !exists {
		cache.dataMutex.RUnlock()
		cache.dataMutex.Lock()

		// Лочим кэш на запись и проверем, что не было ничего записано паралельно
		if value, exists = cache.data[key]; !exists {
			// кэш все еще пуст, записываем значение и разлочиваем мапу для записи
			value = valueFn()
			cache.data[key] = value
			cache.dataMutex.Unlock()
			return value
		} else {
			// Кэш был записан паралельно, отдаем записанное значение
			cache.dataMutex.Unlock()
			return value
		}
	} else {
		// значение присутствует, отдаем его
		cache.dataMutex.RUnlock()
		return value
	}
}
