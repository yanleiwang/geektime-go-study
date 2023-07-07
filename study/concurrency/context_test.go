package concurrency

import (
	"context"
	"testing"
	"time"
)

type mykey struct{}

func TestContext_WithValue(t *testing.T) {
	// 一般是链路起点，或者调用的起点
	ctx := context.Background()
	// 在你不确定 context 该用啥的时候，用 TODO()
	//ctx := context.TODO()

	ctx = context.WithValue(ctx, mykey{}, "study-value")
	val := ctx.Value(mykey{}).(string)
	t.Log(val)
	newVal := ctx.Value("不存在的key")
	val, ok := newVal.(string)
	if !ok {
		t.Log("类型不对")
		return
	}
	t.Log(val)
}

// context包的经典用法----控制超时
func TestTimeoutExample(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	bsChan := make(chan struct{})
	go func() {
		// 模拟 耗时业务
		time.Sleep(2 * time.Second)
		// 业务结束，发送信号
		bsChan <- struct{}{}
	}()

	// 同时监听两个
	//channel，一个是正常业务结束的
	//channel，Done() 返回的。
	select {
	case <-ctx.Done():
		t.Log("超时了")
	case <-bsChan:
		t.Log("业务正常结束")
	}

}

func TestName(t *testing.T) {

}
