# bilibili_danmaku_v1.go 说明

## 概述
bilibili最新弹幕，采用 protobuf 传输弹幕数据  

bilibili_danmaku_v1.proto 为结构定义文件  
protoc 为 protobuf 的二进制文件  
bilibili_danmaku_v1.go 为经过编译后的go文件  


## 编译
1. 保证 $GOPATH/bin 在系统环境变量中，linux：`export PATH="$PATH:$GOPATH/bin"`  
2. 安装编译插件 `go get -u github.com/golang/protobuf/protoc-gen-go`
3. 安装protoc `go get -u github.com/golang/protobuf/proto`
4. 安装其他依赖 `go get google.golang.org/protobuf/reflect/protoreflect@v1.26.0`, 
   `go get google.golang.org/protobuf/reflect/protoreflect@v1.26.0`
5. 如果使用 go mod download 则 2-4 可跳过
6. 进入目录: `app/crawler/spider/internal/lib/bilibili/danmaku`
7. 执行： `protoc --go_out=. bilibili_danmaku_v1.proto` 获得编译后的go文件


## 使用 
见 `bilibili_danmaku_test.go` 文件