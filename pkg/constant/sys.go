package constant

import (
	"os"
	"path/filepath"
)

// 由于go最终会编译，所以工程目录的结构无意义，最终以build出的
// 可执行文件位置为基准，_appPath 实际上是可执行文件的位置，在工程目录中
// 可以认为是main.go所在的位置
// 所以需要特别注意，这个变量不是指当前文件的位置
var _appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var AppPath = _appPath
