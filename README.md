## linkflow
存放各种流式处理函数

## Tee函数，一拖二(一个chan 生成两个)
伪代码如下
```go
 	voice := make(chan interface{})
	out1, out2 := linkflow.Tee(context.Background(), voice)
	var wg sync.WaitGroup

	wg.Add(2)
	defer wg.Wait()

	go func() {
		for k := range [100]int{} {
			voice <- k
		}

		close(voice)

	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case d, ok := <-out1:
				if !ok {
					return
				}
				fmt.Printf("out1 :%d\n", d.(int))
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case d, ok := <-out2:
				if !ok {
					return
				}
				fmt.Printf("out2 :%d\n", d.(int))
			}
		}
	}()

```
