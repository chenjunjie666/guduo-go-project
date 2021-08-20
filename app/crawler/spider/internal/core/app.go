// 爬虫的主app程序，主要作用为管理爬虫的各种资源
// 如爬虫实体，代理池,等等
//
package core
//
//import (
//	"sync"
//	"time"
//)
//
//// 爬虫的 app 实例
//var app *application
//
//const (
//	maxWorkingNum = 5 // 最大工作的爬虫数量
//	maxWaitingNum = 3 // 最大等待数量
//)
//
//// 初始化 app 实例
//func initApp() {
//	app = &application{
//		make(chan *CollectorObj, maxWaitingNum), // 等待队列,上限 3
//		make([]*CollectorObj, 0, 10),            // 工作队列
//		make(chan bool, maxWorkingNum),
//	}
//
//	go app.start()
//}
//
//// 获取 app 实例
//func GetApp() *application {
//	return app
//}
//
//// 爬虫的应用程序中心，针对爬虫的采集器以及爬虫所需的各种资源进行管理
//type application struct {
//	waitingPool chan *CollectorObj // 等待开始的爬虫队列
//	workingPool []*CollectorObj    // 正在工作中的爬虫队列
//	workingNum  chan bool          // 用于表示以及控制正在工作中的爬虫数量
//}
//
//// 创建一个新采集器后，调用此方法进入等待队列，直至采集器被放入工作队列
//func (app *application) waitingToWork(c *CollectorObj) {
//	//fmt.Println("---------------")
//	app.waitingPool <- c
//	for !c.IsCanStart() {
//		//fmt.Println("---------------")
//		time.Sleep(time.Millisecond * 100)
//	}
//}
//
//// app 启动方法
//func (app *application) start() {
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	go app.checkingWorkNum()
//	// 检测工作中的爬虫队列
//	go app.releaseWorkingPool()
//
//	wg.Wait()
//}
//
//// 检测工作中的爬虫队列，对完成的爬虫执行释放操作
//func (app *application) releaseWorkingPool() {
//	for {
//		dec := 0
//		var n []*CollectorObj // 保存仍然进行中的 collector 的中间变量
//		for _, c := range app.workingPool {
//			if c.IsDone() || c.IsTimeOut() {
//				// 如果爬虫已经完成，或者超时了，计数+1
//				dec++
//			} else {
//				// 否则就是正在进行中的，放入变量 n
//				n = append(n, c)
//			}
//		}
//		// 用 n 对 workingPool 重新赋值，目的是去掉 workingPool 中的已完成的 collector
//		app.workingPool = n
//
//		// 将从计数池中弹出 dec 个，dec 就是已经完成/超时的 collector 数量
//		for i := 0; i < dec; i++ {
//			<-app.workingNum
//		}
//		time.Sleep(time.Millisecond * 500) // 每次检查间隔 500ms
//	}
//
//}
//
//// 检测如果 workingNum 长度不足上限值
//// 从 waitingPool 中取出等待的 collector，将之切换为启动状态，然后放入 workingPool
//// 如果 waitingPool 为空的话，那么这里就会阻塞，等待新建 collector 的时候
//// 通过 releaseWorkingPool 方法加入新的
//func (app *application) checkingWorkNum() {
//	for len(app.workingNum) < maxWorkingNum {
//		//fmt.Println("+++++++++++++")
//		c := <-app.waitingPool
//		c.CanStart()
//		app.workingPool = append(app.workingPool, c)
//		app.workingNum <- true
//		time.Sleep(time.Millisecond * 100)
//	}
//}
