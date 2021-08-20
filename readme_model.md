# 模型约定

## 说明
在全局的模型包（`pkg/model`）中定义了一些常用的模型字段，如`id`，`created_at`等，
以及数据库与golang的类型映射,代码不多，请详细阅读包内的文件注释


## 规范说明
为了更好地规范项目以及后续维护，对数据库模型定义做一些限制，
请务必保证模型定义能够遵循以下全部规范！


## 规范
1. 模型全部创建在自己服务内的`model`包下
2. 每一个表定义都需要建立一个名为`表名_model`的包（下简称表定义包）
3. 在表定义包中至少建立两个文件`field.go`，`model.go`，前者用于定义表字段，
后者用户定义模型的方法，定义字段时必须定义 json 的 tag
4. 在`field.go`中一定要定义当前表的外键字段
5. 在`model.go`中的模型结构体名固定为`Def`，通过实现`TableName`方法定义表名
6. 所有字段的数据类型，全部应该使用全局`model`包中的定义，类型映射定义在`pkg/model/type.go`

## 其他的规范说明
由于gorm的定义问题，通常select的时候很多字段是不需要的，这时候必须临时定义一个`struct`用于获取数据
建议将常用的一些结构放入表定义包下的`custom.go`

## 一个遵循规范的完整示例  
假设有表`Admin`定义如下  

|字段名|数据类型|其他信息|
|:---:|:---:|:---:|
|`id`|int(11)|主键|
|`name`|varchar(255)||
|`pwd`|text||
|`last_login_at`| int(11) | 最后登录时间 |
|`created_at`|int(11)|创建时自动更新|
|`created_at`|int(11)|更新数据时自动更新|
|`deleted_at`|int(11)|软删除，调用删除时自动更新|

文件`field.go`内容 
`app/your_service/internal/model/admin_model/field.go`

```go
package admin_model

import "guduo/pkg/model"

// 目前统一交fields，并且这个字段是不导出的！
// 这个结构体本质上包含了表的所有字段定义
type fields struct {
	model.FieldsWithSoftDelete
	Name
	Pwd
	LastLoginAt
}

// 下面都是以结构体形式，定义的表字段

// 定义admin的id在其他表中作为外键的定义
// 格式必须是：表名（驼峰形式）+Id
// 特别注意最后是 Id 而不是 ID （由于IDE的提示机制，ID的情况下会搜索不到这个字段）
// 这个字段是为了在别的模型定义中中如果使用 admin_id 作为外键的时候可以：
// type fields struct {
//    	model.Fields
//    	admin_model.AdminId // 这样来定义其他表的 admin_id 外键
//    	Money // 字段定义结构体
// }
type AdminId struct {
	AdminId model.ForeignKey `json:"admin_id"`
}


type Name struct {
	Name model.Varchar `json:"name"`
}

type Pwd struct {
	Pwd model.Text `json:"pwd"`
}

type LastLoginAt struct {
	LastLoginAt model.SecondTimeStamp `json:"last_login_at"`
}
```


文件`model.go`内容
`app/your_service/internal/model/admin_model/model.go`

```go
package admin_model

import (
	"gorm.io/gorm"
	"guduo/pkg/errors"
	"guduo/pkg/model"
)

// 这是一个最终当做表结构完整定义的结构体
// 所有基于模型的方法，Hook等，都在这上面去做
// 结构体名固定为 Def，因为短，全名是 definition 定义的意思
type Def struct {
	fields // field.go中定义的字段结构体，匿名引入
}

func (d Def) TableName() string {
	return "admin"
}

// 这里定义各种hook或者其他模型相关的方法，这部分跟gorm的文档说明是一样的
func (d *Def) BeforeCreate(tx *gorm.DB) error {
	// ....
	if d.Id == 0 {
		return errors.AppError("admin_model", "id不能为空")
	}
	return nil
}
```


文件`custom.go`内容
`app/your_service/internal/model/admin_model/custom.go`

```go
package admin_model

type AuthVerify struct {
	Name // field.go 中的 Name 结构
	Pwd  // field.go 中的 Pwd 结构
}
```

