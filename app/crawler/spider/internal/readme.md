# 爬虫业务层代码目录

## 子目录说明
|目录名|用途说明|
|:---:|:---:|
|`collector`|爬虫业务的核心逻辑，旗下子目录以爬取指标为依据|
|`core`|spider服务的核心文件|
|`model`|数据库模型定义|
|`storage`|数据存储方法|
|`util`|工具方法|


## 业务代码规范说明：
以播放量指标项为例（播放量指标这里取名为play_count）：
1. 在`collector`下创建文件夹`play_count`
2. 在`play_count`下新建`entry.go`作为统一入口文件并建立`Run()`方法（**`Run`方法是这个包内唯一导出的方法名！**）
   并且entry.go内必须声明一个 `cost ModName = "播放量" // 指标中文名`
3. 根据播放量需要采集的源（对于播放量为：腾讯，芒果TV），在`collector/play_count/`下建立对应网站名的文件`tencent.go`
以及`mango.go`
4. 在`collector/common`以及`storage/`下分别建立`tecnent.go`以及`mango.go`。
**如果已经建立了文件，忽略这一步**
5. 如果第【4】步未跳过，新建文件内容可以根据已有的其他文件复制黏贴，并修改其中的内容。可参考已有文件
6. 对于`play_count/网站.go` 统一存在 `func {网站}Handle()` ，以及`func {网站}{指标名}()`方法。
比如`play_count/tencent.go`内必定存在`func tencentHandle()`以及`func tecnentPlayCount()`
7. 除【6】以外的内部方法，统一以`func {网站}xxxx()`命名，比如`func tencentFilterPlayCount()`
8. 完成编写后，需要在同级目录下编写测试用例，`{网站}_test.go`，方法名为`func Test{网站}{指标名}()`
   内容参考已有文件
   
## 爬虫指标任务分配 
详情见 [任务分配列表](./collector/readme.md)