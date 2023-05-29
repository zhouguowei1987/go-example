package main

import (
	"fmt"
	"time"
)

func main() {
	//defer fmt.Println("!")
	//os.Exit(3)

	//arr := make([]int, 3, 5)
	//brr := append(arr, 8)
	//arr[1] = 2
	//brr[0] = 1
	//fmt.Println(arr)
	//fmt.Println(brr)
	//fmt.Println("==============")
	//brr = append(brr, 8)
	//brr = append(brr, 8)
	//brr[0] = 0
	//fmt.Println(arr)
	//fmt.Println(brr)

	//ch := make(chan struct{})
	//go func() {
	//	fmt.Println("start working")
	//	time.Sleep(time.Second)
	//	ch <- struct{}{}
	//}()
	//<-ch
	//fmt.Println("finished")

	//wg := sync.WaitGroup{}
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go func(wg *sync.WaitGroup, i int) {
	//		fmt.Printf("%d\n", i)
	//		wg.Done()
	//	}(&wg, i)
	//}
	//wg.Wait()

	//arr1 := make([]int, 5)
	//fmt.Println(arr1)
	//fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)
	//arr1 = append(arr1, 1)
	//fmt.Println(arr1)
	//fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)
	//for i := 0; i < len(arr1); i++ {
	//	fmt.Printf("%p\n", &arr1[i])
	//}

	//arr1 := make([]int, 5)
	//fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)
	//for i := 0; i < 5; i++ {
	//	arr1 = append(arr1, i)
	//	fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)
	//}
	//fmt.Println("====================")
	//for i := 0; i < 5; i++ {
	//	arr1 = append(arr1, i)
	//	fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)
	//}

	//arr1 := make([]int, 5)
	//fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)
	//for i := 0; i < 5; i++ {
	//	fmt.Printf("ptr = %p\n", &arr1[i])
	//}
	//arr1 = arr1[1:2]
	//fmt.Printf("len = %d , cap = %d, ptr = %p\n", len(arr1), cap(arr1), arr1)

	//sumPrint := fun1()
	//fmt.Println(sumPrint(1))
	//fmt.Println(sumPrint(1))

	//fmt.Println(fun3())

	data := make(map[int]int, 10)
	for i := 1; i <= 10; i++ {
		data[i] = i
	}
	for key, value := range data {
		go func() {
			fmt.Printf("key => %d, value => %d\n", key, value)
		}()
	}
	time.Sleep(time.Second * 5)
}

func fun1() func(int) int {
	sum1 := 0
	return func(val int) int {
		sum1 += val
		return sum1
	}
}

func fun2() (val int) {
	val = 10
	defer func() {
		val++
	}()
	return val
}

func fun3() int {
	val := 10
	defer func() {
		val++
		fmt.Println(val)
	}()
	return val
}
