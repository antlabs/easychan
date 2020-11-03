## linkflow
存放可以流式处理的函数

## Tee函数，一拖二(一个chan 生成两个)
伪代码如下
```go
 go func() {
        for _, v := range need {
            voice <- v
        }

        close(voice)

    }() 

    go func() {
        var got []int
        defer wg.Done()
        for {
            select {
            case d, ok := <-out1:
                if !ok {
                    assert.Equal(t, need, got)
                    return
                }
                got = append(got, d.(int))
            }
        }
    }() 

    go func() {
        defer wg.Done()
        var got []int
        for {
            select {
            case d, ok := <-out2:
                if !ok {
                    assert.Equal(t, need, got)
                    return
                }
                got = append(got, d.(int))
            }
        }
    }()

```