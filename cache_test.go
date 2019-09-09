package cache_test

import (
	"github.com/default23/cache"
	"strconv"
	"sync"
	"testing"
)

func generateValueFor(counter int) cache.Value {
	return strconv.Itoa(counter) + "000"
}

func TestInMemoryCache_GetOrSet(t *testing.T) {

	c := cache.NewInMemoryCache()
	var wg sync.WaitGroup

	elementsToWrite := 10

	wg.Add(elementsToWrite)
	// конкурентная запись новых значений
	for i := 0; i < elementsToWrite; i++ {
		go func(counter int) {
			defer wg.Done()
			calls := 0 // счетчик вызовов генератора значений
			expected := generateValueFor(counter)

			actual := c.GetOrSet(strconv.Itoa(counter), func() cache.Value {
				calls++
				return expected
			})

			if actual != expected {
				t.Errorf("GetOrSet should return value, generated using the callback; expected: %s, got: %s", expected, actual)
			}

			if calls > 1 {
				t.Errorf("value generator should be triggered only once, but have been called %d times", calls)
			}
		}(i)
	}
	wg.Wait()

	// Все происходит внутри одного теста, для того что бы не наполнять кэш данными заново
	// примерно точно таким же алгоритмом как выше, будем использовать уже готовые данные для следующих тестов
	// В какой то мере это не совсем правильно, но для того что бы не плодить одинаковый код написал так>

	wg.Add(elementsToWrite)
	// Паралельная проверка на то, что все записанные значения присутствуют в кэше,
	// с правильным ключ-значением
	for i := 0; i < elementsToWrite; i++ {
		go func(counter int) {
			defer wg.Done()

			generatorCalls := 0
			key := strconv.Itoa(counter)
			expected := generateValueFor(counter)

			actual := c.GetOrSet(key, func() cache.Value {
				generatorCalls++
				return "----------"
			})

			if actual != expected {
				t.Errorf("unexpected value received, expected is: %s, but got: %s", expected, actual)
			}
			if generatorCalls > 0 {
				t.Errorf("")
			}
		}(i)
	}

	wg.Wait()
}
