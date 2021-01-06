参考python版本jd-assistant的实现逻辑。抽取出来可用逻辑，进行go语言重写。

## 预定目标
代码按照预期实现，并且业务上抢到一瓶茅台

## 技术选型
分析jd-asssistant python版本抽取，所需要的业务组件。任务拆分如下：

### jd-sdk
1. python requests库替代。换成异步io的，支持并发
   https://github.com/asmcos/requests
2. 重写所有，抢购遇到的请求接口，封装成sdk

### 文件读写
1. 实现图片的随便读写。指定位置，写，读取，删除图片。
2. 实现，内嵌db存储，或者内存数据库，后者文件数据库的存储
   https://github.com/spf13/afero

### 验证码打开
1. 调用windows，linux和mac版本的图片打开操作。
2. 支持将图片，发送邮箱，发送给某个人微信企业号
   https://github.com/jordan-wright/email
   https://github.com/h2non/bimg
   发邮件貌似不可行。京东手机端估计只能摄像头扫描才能验证。
   打开图片
   https://blog.csdn.net/qq_34857250/article/details/100569734
   https://blog.csdn.net/youngwhz1/article/details/88662172

### 定时器框架
1. 实现高精度的定时任务处理
2. 京东服务器时间同步，进行提升。
   https://juejin.cn/post/6844903901418749960
   高精度的算法，实现看来不行
   https://club.perfma.com/article/661495
   看来轮训还是靠谱。

### 融合所有的模块
1. 增加构建工具，自动化打包，自动化部署。
2. 融合所有模块，跑通预定逻辑
   配置文件库
   https://github.com/spf13/viper
   依赖管理
   https://github.com/Golang/dep [已放弃维护，推荐使用go mod]
   日志库
   https://github.com/sirupsen/logrus

## 开发记录
### idea开发环境
https://segmentfault.com/a/1190000019989029

### go项目结构规范
https://www.jianshu.com/p/4726b9ac5fb1

### 开源库参考
https://github.com/avelino/awesome-go

### go语言的工具库，和开发规范
https://github.com/xxjwxc/uber_go_guide_cn

### 入门基本语法
https://golang.org/doc/
https://geektutu.com/post/quick-golang.html
https://www.jianshu.com/p/fe16d3762043

1. func main：main 函数是整个程序的入口，main 函数所在的包名也必须为 main。
- 1. 问题： 该版本的 %1 与你运行的 Windows 版本不兼容。请查看计算机的系统信息，然后联系软件发布者
- go Main file has non-main package or doesn't contain main function
2. 变量声明。特殊的变量声明方式。后置声明var a int = 1 ，加一种特殊写法msg := "Hello World!"
3. 简单类型。
- 空值nil。
- 整形类型int（取决于操作系统），int8，int16
- 浮点float32，float64
- 字节byte（等价于uint8）
- 字符串string。 go语言中，字符串使用utf8编码。以 byte 数组形式保存的。打印长度和输出比较特殊
- 布尔 boolean（true、false）
```
var a int8 = 10
var c1 byte = 'a'
var b float32 = 12.2
var msg = "Hello World"
ok := false
```
- 数组： var i [5][5]int //二维数组。 带初始化var arr = [5]int{1, 2, 3, 4, 5}或 arr := [5]int{1, 2, 3, 4, 5}
  切片的的方式来，对数组进行，合并和拆分操作，使用了数组作为底层的数据结构
- 字典（map）和其他语言没啥区别
- 指针，类型定义时使用符号\*，对一个已经存在的变量，使用 & 获取该变量的地址。指针通常在函数传递参数，或者给某个类型定义新的方法时使用。Go 语言中，**参数是按值传递**的，如果不使用指针，函数内部将会拷贝一份参数的副本，对参数的修改并不会影响到外部变量的值

4. 流程控制
- 条件控制 if else .switch 比较特殊
- 循环 for

5. 函数
```
// 多个参数，多个返回值
func funcName(param1 Type1, param2 Type2, ...) (return1 Type3, ...) {
    // body
}
```
6. 错误处理
- error errors.New("返回自定义错误"),可以利用多返回值的特性，来返回错误
- 异常处理，defer 和recover机制来操作
7. 结构体struct.类似java的class。
- 注意实现方法，和实现函数的区别
8. 接口
- Go 语言中，并不需要显式地声明实现了哪一个接口，只需要直接实现该接口对应的方法即可。实例化 Student后，强制类型转换为接口类型 Person。
- 如何确保某个类型实现了某个接口的所有方法呢？一般可以使用下面的方法进行检测，如果实现不完整，编译期将会报错。var _ Person = (*Student)(nil)
9. 空接口。定义了一个没有任何方法的空接口，那么这个接口可以表示任意类型

```
func main() {
	m := make(map[string]interface{})
	m["name"] = "Tom"
	m["age"] = 18
	m["scores"] = [3]int{98, 99, 85}
	fmt.Println(m) // map[age:18 name:Tom scores:[98 99 85]]
}
```
7 并发.协程(goroutine)并发（
>协程，又称微线程。在一个子程序中中断，去执行其他子程序，不是函数调用，有点类似CPU的中断。
1. **优势**没有线程切换的开销
2. 不需要多线程的锁机制，只有一个线程，也不存在同时写变量冲突，在协程中控制共享资源不加锁，只需要判断状态就好了）

- sync  sync.WaitGroup；wg.Add(1)/wg.Done()，go xxx；sync.Wait()无阻塞并发，协程间不需要等待
- channel 阻塞并发，阻塞协程
8. 单测
- 新建xxxfunc_test.go。直接编写测试函数即可
- 运行 go test，将自动运行当前 package 下的所有测试用例。 go test -v 显示详细
9. 包
- 一般来说，一个文件夹可以作为 package，同一个 package 内部变量、类型、方法等定义可以相互看到。
- Go 语言也有 Public 和 Private粒度是包。如果类型/接口/方法/函数/字段的首字母大写，则是 Public 的，对其他 package 可见，如果首字母小写，则是 Private 的，对其他 package 不可见
10. 模块Modules
- Go Modules 是 Go 1.11 版本之后引入的，Go 1.11 之前使用 $GOPATH 机制
- Go Modules 在 1.13 版本仍是可选使用的，环境变量 GO111MODULE 的值默认为 AUTO，强制使用 Go Modules 进行依赖管理，可以将 GO111MODULE 设置为 ON。
- 1. go mod init example 2. go run .，将会自动触发第三方包 rsc.io/quote的下载，具体的版本信息也记录在了go.mod
- import 模块名/子目录名 的方式，来在相同的模块不同的包下调用

### 面向对象
https://xiaowing.github.io/post/20170816_is_go_object_oriented_cn/
https://flaviocopes.com/golang-is-go-object-oriented/

go语言是面向对象的。
- 要结构体,不要类
- 由于GO语言中没有继承这一概念，用组合代替继承。当我们在定义一个结构体时，我们可以追加类型为另一个结构体的匿名字段。这样一来，我们定义的这个结构体也就同时拥有了另一个结构体的所有字段以及方法。这种技法被称之为Struct Embedding。
- 接口是隐式实现的
- GO语言同时又允许函数脱离于对象而独立存在
### readme，任务拆分，分支建立。github同步

问题：模块下载不下来
设置go的代理，一个全球代理，为go模块而生！！！
https://goproxy.io/zh/

### 环境变量问题
问题：包引入，运行不起来。各种报错
https://studygolang.com/articles/15790
关于环境变量：
1. GOROOT：C:\app\go
2. Path追加 %GOROOT%/bin
3. GOPATH.go和其他语言不一样十分依赖于工作目录，即GOPATH所指向的目录。
   go的这种模式决定了你不能按照版本控制工具来作为代码的根目录，也不能随意的将某一个测试项目建立到随意的路径下，因为工作目录必须都在GOPATH所指向的路径中。

env的参数解释
https://wiki.jikexueyuan.com/project/go-command-tutorial/0.14.html
![](leanote://file/getImage?fileId=5fecca1104c0f95768000007)

### 包管理
https://zhuanlan.zhihu.com/p/60703832
1.  1.12之后支持
2.  GO111MODULE=auto
3.  $GOPATH/src ，并且cd进去，新建文件，当前目录下，命令行运行 go mod init + 模块名称 初始化模块
4.  直接 go run hello.go。go 会自动查找代码中的包，下载依赖包，并且把具体的依赖关系和版本写入到go.mod和go.sum文件中。

### 编写一个go项目，必须知道的gomoudle知识
https://blog.csdn.net/weixin_44676081/article/details/107279746
https://www.jianshu.com/p/c666ebdb462b
1. GO111MODULE="off"
   需要配置，$GOPATH的位置，指向我们的项目路径，
   在$GOPATH/src下面写代码。写模块，写main
   模块之间相互调用。import(直接文件夹(文件夹是就是moudle))。go语言会默认使用文件夹名作为包的导入。可以理解为import导入的不是包，而是文件夹。所以，一般情况下，包名和文件夹名字一样。

2. GO111MODULE="on"
3.

！！！**默认情况下，import的路径，文件夹路径**

### 包管理
1. 基于GOPATH和Vendor的构建方式
   go mod与之前的利用vendor特性的依赖管理工具的不同点在于，go mod 更类似于maven这种本地缓存库的管理方式,不论你有多少个工程，只要你引用的依赖的版本是一致的，那么在本地就只会有一份依赖文件的存在。而vendor即使依赖的版本是相同的，但如果在不同的工程中进行了引用，也会在工程目录下的vendor产生一份依赖文件。

2. Golang在1.11版本中引入了go mod机制,在统一的位置对依赖进行管理
   主要是通过GOPATH/pkg/mod下的缓存包来对工程进行构建。 可以通过GO111MODULE来控制是否启用，GO111MODULE有一下三种类型。
   on 所有的构建，都使用Module机制
   off 所有的构建，都不使用Module机制，而是使用GOPATH和Vendor
   auto 在GOPATH下的工程，不使用Module机制，不在GOPATH下的工程使用

将需要进行版本管理的代码从GOPATH路径下移出
在项目的根目录下使用命令go mod init projectName
在该目录下执行go build main.go

如果工程中存在go.mod文件，编译时是从GOPATH/pkg/mod下查找依赖。
如果主动使用export GO111MODULE=off命令不使用Module机制,进行编译就会从GOPATH/src下查找依赖。会产生以下输出。(编译失败是由于相应目录下无依赖文件)

建议在工程的根目录下执行go mod init projectName命令。在执行go mod化之后，所有的引用都不再是以GOPATH为相对路径的引用了，而是变成了以go.mod中初始化的项目名为起始的引用。

### go项目的目录结构，和import
https://www.bilibili.com/video/BV1Qk4y1R74q?p=5
![](leanote://file/getImage?fileId=5fed972404c0f9576800000b)
![](leanote://file/getImage?fileId=5fed973d04c0f9576800000c)
![](leanote://file/getImage?fileId=5fed96ec04c0f9576800000a)
go语言，使用goroot来表达go安装的位置
使用gopath/src，表示代码的位置
$gopath/src/github.com/fanlushuai/jd-assint
组成，由，$gopath+src+域名+用户名+项目名，的结构来使用。

通常情况下的，import会在基于gopath的路径上寻找，模块的目录。
包名，是从$gopath+src之后开始算的。
go语言禁止循环导入包。

### 尝试写一个部分

####golang中如何在包中引用另外一个包的结构体？
https://segmentfault.com/q/1010000009368531
将room.go所在的package引入到msg.go中，然后在msg.go中使用的时候加上包名。

#### go代码的执行顺序
https://davidchan0519.github.io/2019/06/21/golang-main-init-sequence/
![](leanote://file/getImage?fileId=5feee36704c0f9576800000d)

#### go字符处理
https://segmentfault.com/a/1190000019439223

#### glob对象序列化
https://www.cnblogs.com/yinzhengjie2020/p/12735277.html

#### go语言 interface{}怎么转化为具体类型
v.{T}
#### go 语言空对象判断
if objectA== (structname{}){ // your code }

#### 并发策略模型  赢者为王，生产者消费者模型，发布订阅模型
https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-06-goroutine.html

#### go语言返回对象还是返回指针
性能的问题：返回值和返回指针，个人感觉没吊差别。
https://segmentfault.com/q/1010000019133280
如果是接口，返回指针。可以指向具体类型
https://blog.csdn.net/K346K346/article/details/91041276

### go语言3元运算符
不存在。可读性的问题

### 对象转map[string]interface{}
https://www.liwenzhou.com/posts/Go/struct2map/

### go 对象转化 map[string]string
反射

go的50个坑
https://juejin.cn/post/6844903816018542600

#### go配置文件的路径问题
https://eddycjy.gitbook.io/golang/di-1-ke-za-tan/golang-relatively-path

在线变成go
https://play.golang.org/p/M7bukF-4Pk6

#### go的时间处理问题
特殊的时间格式！！！

#### go性能测试time获取
t1 := time.Now()
elapsed := time.Since(t1)
#### 通道应该什么时候被关闭
https://learnku.com/go/t/23459/how-to-close-the-channel-gracefully
Go 中没有一个內建函数去检测管道是否已经被关闭。
一般原则上使用通道是不允许接收方关闭通道和 不能关闭一个有多个并发发送者的通道。 换而言之， 你只能在发送方的 goroutine 中关闭只有该发送方的通道。

#### 类型转换

#### go的构造器？
https://learnku.com/go/t/47475
在所有这些情况下，声明变量或调用 new，或使用不带任何显式值的复合文字，或调用 make- 在所有这些情况下，我们得到的值为零。
但是有时我们希望使用一些合理的非零值来初始化变量。这是构造函数最常见的用例

#### go sleep精度的问题
//todo 编写一个搜索关键字，记录的chrome插件

#### 参数传递问题，考虑一下chain通道
https://www.flysnow.org/2018/02/24/golang-function-parameters-passed-by-value.html
Go语言中所有的传参都是值传递（传值），都是一个副本，一个拷贝。因为拷贝的内容有时候是非引用类型（int、string、struct等这些），这样就在函数中就无法修改原内容数据；有的是引用类型（指针、map、slice、chan等这些），这样就可以修改原内容数据。

区别java。java没有指针。有对象。普通类型和String（不可变）传递值，其他对象类型传递对象的地址的值。来实现引用传递

#### golang随机time.sleep的Duration问题
http://xiaorui.cc/archives/3034
time.Sleep(time.Duration(num) * time.Second)

#### 优质博客
https://eddycjy.gitbook.io/golang/
https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/
https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-06-goroutine.html
