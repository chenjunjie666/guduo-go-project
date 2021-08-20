// 定义数据库常用的数据类型
package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

// 根据其特殊用途所定义的类型
// 比如时间戳在数据库内定义为int，这里使用 uint
// +------------------------------------------------------------+
// |只有极度通用且普遍的特殊类型才应该在这里定义！！！！！！！！！！！！！！！|
// |只有极度通用且普遍的特殊类型才应该在这里定义！！！！！！！！！！！！！！！|
// |只有极度通用且普遍的特殊类型才应该在这里定义！！！！！！！！！！！！！！！|
// +------------------------------------------------------------+

type PrimaryKey = uint64    // 主键
type ForeignKey = uint64    // 外键
type SecondTimeStamp = uint // 秒级时间戳, 这里不定义uint64是因为扩展包对秒级时间戳的定义就是uint

// Mysql数据库类型与golang类型的映射关系

type Varchar = string     // varchar 类型
type Int = int64          // int,bigint 类型
type Decimal = float64    // decimal 类型
type Text = string        // text 类型
type DateTime = time.Time // datetime 类型
type Tinyint = int8       // tinyint 类型
type Float = float64	  // float 类型







// BitBool is an implementation of a bool for the MySQL type BIT(1).
// This type allows you to avoid wasting an entire byte for MySQL's boolean type TINYINT.
type BitBool bool

// Value implements the driver.Valuer interface,
// and turns the BitBool into a bitfield (BIT(1)) for MySQL storage.
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

// Scan implements the sql.Scanner interface,
// and turns the bitfield incoming from MySQL into a BitBool
func (b *BitBool) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}
	*b = v[0] == 1
	return nil
}