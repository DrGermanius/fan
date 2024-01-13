package fan

import (
	"fmt"
	"reflect"
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
func TestInCanDoSendToOldChan(t *testing.T) {
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

	if !reflect.DeepEqual(results, []int{11, 21, 22, 31, 32, 33}) {
		fmt.Println(results)
		t.Fatal("Source code is lying")
	}
}

func TestIn_empty(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	ch3 := make(chan int, 3)
	res := In[int](&ch1, &ch2, &ch3)

	if len(res) > 0 {
		t.Fatal("Source code is lying")
	}
}
