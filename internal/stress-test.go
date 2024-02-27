package internal

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type SafeCounter struct {
	mu    sync.Mutex
	resul map[int]int
}

func (c *SafeCounter) Inc(key int) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	if _, found := c.resul[key]; !found {
		c.resul[key] = 1
	} else {
		c.resul[key]++
	}

	c.mu.Unlock()
}

func Execute(address string, requests int, concurrency int) {
	waitGroup := sync.WaitGroup{}
	waitGroup2 := sync.WaitGroup{}
	resul := SafeCounter{resul: make(map[int]int)}

	calculateData := calculate(requests, concurrency)

	start := time.Now()
	waitGroup.Add(len(calculateData))
	percent := 0
	for i := 0; i < len(calculateData); i++ {

		waitGroup2.Add(calculateData[i])
		for j := 0; j < calculateData[i]; j++ {
			go CallHttpGet(address, &resul, &waitGroup2)
		}

		waitGroup2.Wait()
		waitGroup.Done()
		percent = percent + ((calculateData[i] * 100) / requests)
		fmt.Printf("%d%% ", percent)
	}
	waitGroup.Wait()

	finish := time.Since(start)

	fmt.Printf("\n\nTime : %s", finish)
	for key, value := range resul.resul {
		fmt.Printf("\nCode %d -> %d", key, value)
	}

}

var ErrToMany = errors.New("stopped after 10 redirects")

func CallHttpGet(address string, resul *SafeCounter, waitGroup2 *sync.WaitGroup) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 20 {
				return ErrToMany
			}
			return nil
		},
	}

	req, err := client.Get(address)
	if err != nil {
		if errors.Is(err, ErrToMany) {
			resul.Inc(http.StatusTooManyRequests)
		} else {
			panic(err)
		}

	} else {
		resul.Inc(req.StatusCode)
	}

	waitGroup2.Done()
}

// Quebra a qtd de requisições em um array
// Ex : 5 -> requests 2 -> concurrency = [2,2,1]
func calculate(requests int, concurrency int) []int {

	quotient := requests / concurrency
	rest := requests % concurrency
	slice := make([]int, quotient+1)

	for i := 0; i <= quotient; i++ {
		if i == quotient {
			slice[i] = rest

		} else {
			slice[i] = concurrency
		}
	}

	if rest != 0 {
		return slice
	}
	return slice[:quotient]

}
