## easychan
使用各种chan操作的函数

## Tee函数，一拖二(一个chan 生成两个)
使用者只管往一个chan里面写数据，从easychan里面生成的两个chan可以得到同样的数据
伪代码如下
```go
  // 这里支持任意类型chan
 	voice := make(chan []byte)
  var wg sync.WaitGroup
	out1, out2 := easychan.Tee(context.Background(), voice)

	wg.Add(2)
	defer wg.Wait()

	go func() {
    // 模拟生产者发音频
		for k := range [100]int{} {
			voice <- []byte{k}
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
				fmt.Printf("out1 :%d\n", d)
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
				fmt.Printf("out2 :%d\n", d)
			}
		}
	}()

```
