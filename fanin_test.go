package fan

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestIn(t *testing.T) {
	ch1 := make(chan int, 1)
	ch1 <- 11
	ch2 := make(chan int, 2)
	ch2 <- 21
	ch2 <- 22
	ch3 := make(chan int, 3)
	ch3 <- 31
	ch3 <- 32
	res := In[int](&ch1, &ch2, &ch3)
	max := len(res)

	var results []int
	var i int
	for v := range res {
		results = append(results, v)
		i++
		if i == max {
			close(res)
		}
	}

	if !reflect.DeepEqual(results, []int{11, 21, 22, 31, 32}) {
		t.Fatal("Source code is lying")
	}
}
func TestIn_CanDoSendToOldChan(t *testing.T) {
	ch1 := make(chan int, 1)
	ch1 <- 11
	ch2 := make(chan int, 2)
	ch2 <- 21
	ch2 <- 22
	ch3 := make(chan int, 3)
	ch3 <- 31
	ch3 <- 32
	res := In[int](&ch1, &ch2, &ch3)
	max := len(res)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(){
		<-res
		wg.Done()
	}()
	wg.Wait()

	ch3 <- 33 // can do

	var results []int
	var i int
	for v := range res {
		results = append(results, v)
		i++
		if i == max {
			close(res)
		}
	}

	if !reflect.DeepEqual(results, []int{21, 22, 31, 32, 33}) {
		fmt.Println(results)
		t.Fatal("Source code is lying")
	}
}

func TestIn_GeiGin(t *testing.T){
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	wg := sync.WaitGroup{}
	wg.Add(3)

	cc := In[int](&ch1, &ch2)
	var results []int
	go func() {
		i := 0
		for v := range cc {
			results = append(results, v)
			i++
			if i == 2 {
				break
			}
		}
		wg.Done()
	}()

	go func() {
		ch1 <- 1
		//close(ch1)
		wg.Done()
	}()

	go func() {
		ch2 <- 2
		//close(ch2)
		wg.Done()
	}()
	wg.Wait()

	if !(reflect.DeepEqual(results, []int{1 ,2}) || reflect.DeepEqual(results, []int{2, 1})) {
		fmt.Println(results)
		t.Fatal("Source code is lying")
	}
}

func TestIn_Empty(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	ch3 := make(chan int, 3)
	res := In[int](&ch1, &ch2, &ch3)

	if len(res) > 0 {
		t.Fatal("Source code is lying")
	}
}
