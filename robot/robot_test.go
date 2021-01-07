package robot

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	login()
}

func TestDiffLocalServerTime(t *testing.T) {
	fmt.Println(diffLocalServerTime())
}

func TestGetGetTriggerTime(t *testing.T) {
	fmt.Println(getTriggerTime(ac.Skus[0], 20))
}

func TestNervousBlockWait(t *testing.T) {
	nervousBlockWait(1000)
}

//func TestGetSubmitOrderPostData(t *testing.T) {
//	getSubmitOrderPostData(ac.Sku)
//}

func TestName(t *testing.T) {
	fmt.Println(fakeKillUrl("xxxxfakeKillUrl"))
}

func fakeKillUrl(skuId string) string {
	ch := make(chan string)

	go func(ch chan string) {

		type Goon struct {
			value bool
		}

		giveUpLoopLeftTimes := 100
		goon := Goon{value: true}
		for {

			if !goon.value || giveUpLoopLeftTimes == 0 {
				break
			}

			fmt.Println("一个协程")
			go func(g *Goon, ch chan string) {
				//偷梁换柱。（其实也可以封装一些func形式的，还能复用这个逻辑，还能测试）
				fakeRequestTime := rand.Intn(200) + 3000
				fmt.Println("fake req开始", fakeRequestTime)

				time.Sleep(time.Duration(fakeRequestTime) * time.Millisecond)
				url := "test" + skuId
				if len(url) > 0 {
					ch <- url
					//不要惊讶这个日志，一个不是最小的fakeRequestTime返回了。因为每个fakeRequestTime的开始时间不一样。所以，你懂的，
					fmt.Println("一个线程返回了", fakeRequestTime)
					g.value = false
				}
			}(&goon, ch)

			time.Sleep(20 * time.Millisecond)
			giveUpLoopLeftTimes--
		}

	}(ch)

	killUrl := <-ch

	return killUrl
}

func TestWg(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	for index := range ac.Skus {
		//注意，用索引的方式，可以获取到值，而不是值copy.可以配合动态监听配置文件，修改一些内容
		sku := ac.Skus[index]
		fmt.Println("开启doSku id=", sku.Id)
		go func() {
			time.Sleep(10 * time.Second)
			fmt.Println("子线程结束")
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("结束")

}
