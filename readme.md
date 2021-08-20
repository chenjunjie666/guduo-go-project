# 应用总览  


## 概览
骨朵应用服务代码


## 基本规范
1. 多数情况下文件夹都应该有对应的`readme.md`文件来说明文件夹作用以及其一级子目录用途
2. 除程序启动阶段的检查以及初始化，其他代码中 <u>**不应该手动抛出任何形式的panic!**</u>
3. 错误处理：
    1. 默认应该将所有的error向上return
    2. 如果某些错误或者不符合预期的情况不需要向上传递，则必须记录日志，
       日至等级为`Error`或者`Warn`/`Warning`


## 根目录文件夹说明  
|目录名|用途说明|
|:---:|:---:|
|`app`|服务目录|
|`build`|构建应用相关的文件|
|`pkg`|应用程序公用代码|
|`test`|测试文件|

## 服务
当前应用服务

|目录名|服务名|状态|
|:---:|:---:|:---:|
|`app/crawler/sipder`|[爬虫服务-数据抓取](./app/crawler/spider/internal/readme.md)| DOING |
|`app/crawler/clean`|爬虫服务-数据清洗| DOING |


## 项目使用到的第三方包

|包名|名称|说明|
|:---:|:---:|:---:|
|`github.com/go-redis/redis/v8`|go-redis| redis |
|`github.com/gocolly/colly/v2`|colly爬虫框架| 程序主框架 |
|`github.com/pkg/errors`|错误处理包|获取错误栈|
|`github.com/sony/sonyflake`|随机ID包|项目随机UID的生成|
|`gorm.io/gorm`|go orm数据库框架|程序数据库操作|
|`github.com/sirupsen/logrus`|日志系统|用于处理系统的日志|



# 数据库模型
详情见根目录下的 [readme_model.md](./readme_model.md) 说明