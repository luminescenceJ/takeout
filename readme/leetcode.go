package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	//rand.Seed(time.Now().UnixNano())

	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})

	// It must be a buffered channel.
	toStop := make(chan string, 1)

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop
		fmt.Println(stoppedBy)
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			for {
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					fmt.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	select {
	case <-time.After(time.Hour):
	}

}

//
//import (
//	"fmt"
//	"sort"
//)
//
//func main() {
//	// a := 0
//	// b := 0
//	// for {
//	//     n, _ := fmt.Scan(&a, &b)
//	//     if n == 0 {
//	//         break
//	//     } else {
//	//         fmt.Printf("%d\n", a + b)
//	//     }
//	// }
//	var n int
//	_, _ = fmt.Scan(&n)
//	arr := make([]int, n)
//	for i := 0; i < n; i++ {
//		_, _ = fmt.Scan(&arr[i])
//	}
//	sort.Ints(arr)
//	for i := 0; i < n; {
//		if i+1 >= len(arr) {
//			break
//		}
//		if arr[i] == arr[i+1] {
//			arr = append(arr[:i], arr[i+1:]...)
//		} else {
//			i++
//		}
//
//	}
//
//	for _, n := range arr {
//		fmt.Println(n)
//	}
//
//}
