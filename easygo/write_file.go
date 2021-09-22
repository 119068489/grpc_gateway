// 利用 beego 的 logs 模块实现写文件
package easygo

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/logs"
)

type FileWriter struct {
	Once  sync.Once
	Loger *logs.BeeLogger
	File  string
}

func NewFileWriter(file string) *FileWriter {
	p := &FileWriter{}
	p.Init(file)
	return p
}
func (fwSelf *FileWriter) Init(file string) {
	fwSelf.File = file
}
func (fwSelf *FileWriter) Write(text string, v ...interface{}) {
	fwSelf.Once.Do(fwSelf.OnceCall)
	text = fmt.Sprintf(text, v...)
	fwSelf.Loger.Write([]byte(text))

}
func (fwSelf *FileWriter) OnceCall() {
	fwSelf.Loger = logs.NewLogger()
	config := fmt.Sprintf(`{"filename":"%s","rotate":false,"perm":"777"}`, fwSelf.File)
	fwSelf.Loger.SetLogger("file", config) // SetLogger 这个函数一调用就马上产生文件,我不希望这样，所以在 Write 时才调用 OnceCall
	fwSelf.Loger.Async()
}

func (fwSelf *FileWriter) Flush() {
	if fwSelf.Loger != nil {
		fwSelf.Loger.Flush()
	}
}

//------------------------------------------------------

// 利用上面的类，定义出一个更方便的对外接口
func WriteFile(file string, text string, v ...interface{}) {
	var writer *FileWriter
	_MutexForFileWriter.Lock()
	wt, ok := _FileWriterMap.Load(file)
	if ok {
		writer = wt.(*FileWriter)
	} else {
		writer = NewFileWriter(file)
		_FileWriterMap.Store(file, writer)
	}
	_MutexForFileWriter.Unlock()

	writer.Write(text, v...)

}

var _FileWriterMap sync.Map
var _MutexForFileWriter Mutex
