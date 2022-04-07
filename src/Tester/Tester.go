package Tester

import (
	"TestMachine/src/LimitExec"
	"bytes"
	"encoding/json"

	//"fmt"
	"io"
	"os"
	//"strconv"
)

type TestBuffer struct {
	buf []byte
	siz int
	len int
	off int
	out bool
}

func NewTestBuffer(length int) TestBuffer {
	return TestBuffer{
		buf: make([]byte, length),
		len: 0,
		siz: length,
		off: 0,
		out: false,
	}
}

func (t *TestBuffer) Peek() (byte, error) {
	if t.off == t.len {
		return 0, LimitExec.ExecError{Info: "EOF"}
	}
	return t.buf[t.off], nil
}

func (t *TestBuffer) Move() {
	if t.off < t.len {
		t.off++
	}
}

func (t *TestBuffer) GetByte() (byte, error) {
	if t.off == t.len {
		return 0, LimitExec.ExecError{Info: "EOF"}
	}
	t.off++
	return t.buf[t.off-1], nil
}

func (t *TestBuffer) Check(p []byte) bool {
	len_std := len(p)
	//fmt.Println(string(p))
	//fmt.Println(string(t.buf))
	for len_std > 0 && (p[len_std-1] == 0 || p[len_std-1] == '\n' || p[len_std-1] == ' ') {
		len_std--
	}
	off := 0
	for off < len_std {
		for off < len_std && p[off] == ' ' {
			off++
		}
		for {
			ch, err := t.Peek()
			if err != nil || ch != ' ' {
				break
			}
			t.Move()
		}
		if off == len_std {
			break
		}
		if p[off] == '\n' {
			ch, err := t.Peek()
			if err == nil && ch == '\n' {
				off++
				t.Move()
			} else {
				return false
			}
		} else {
			for off < len_std {
				if p[off] == ' ' || p[off] == '\n' {
					break
				}
				ch, err := t.Peek()

				//fmt.Println(ch)
				//fmt.Println(p[off])
				//fmt.Println(err)

				if err == nil && ch == p[off] {
					off++
					t.Move()
				} else {
					//fmt.Println("return false")
					return false
				}
			}
		}
	}
	return true
}

func (t *TestBuffer) Write(p []byte) (n int, err error) {
	plen := len(p)

	if t.len+plen > t.siz {
		t.out = true
		return 0, LimitExec.ExecError{Info: "OLE"}
	}
	tmp := t.buf[t.off:]
	copy(tmp, p)
	t.len += plen

	return plen, nil
}

type TestInfo struct {
	Statu   string
	Info    string
	RunTime int64
	Memory  int64
}

func (info TestInfo) ToBytes() []byte {
	tmp, err := json.Marshal(info)
	if err != nil {
		return []byte("{\"Statu\":\"005\",\"Info\":\"" + err.Error() + "\",\"RunTime\":-1,\"Memory\":-1}")
	}
	return tmp
}

func (info TestInfo) ToString() string {
	tmp := info.ToBytes()
	str := string(tmp)
	return str
}

func (info TestInfo) ToStdString() string {
	tmp := info.ToBytes()
	str := bytes.Buffer{}
	_ = json.Indent(&str, tmp, "", "	")
	return str.String()
}

type Tester interface {
	Compile() TestInfo
	Run(in, out string) TestInfo
	Spj(in, out, spj string) TestInfo
}

//
//
//
//
//
func CompileBase(cmd string, timeLimit, memoryLimit int64) TestInfo {
	stderr := bytes.Buffer{}
	timeUse, memoryUse, err := LimitExec.LimitExecE(timeLimit, memoryLimit, cmd, &stderr)
	var res TestInfo
	res.Statu = "006"
	if err != nil {
		if err.Error() == "MLE" {
			res.Info = "编译使用内存超限"
		} else if err.Error() == "TLE" {
			res.Info = "编译超时"
		} else if err.Error() == "RE" {
			res.Info = "编译出现错误" + stderr.String()
		} else {
			res.Info = err.Error() + stderr.String()
		}
		res.Statu = "007"
	} else {
		res.Info = "编译成功"
	}
	res.RunTime = timeUse
	res.Memory = memoryUse
	return res
}

func ReadFile(f io.Reader) *bytes.Buffer {
	resbuf := bytes.Buffer{}
	buf := make([]byte, 1024*1024)
	for {
		len, err := f.Read(buf)
		if len == 0 {
			break
		} else if err != nil {
			return nil
		} else {
			resbuf.Write(buf)
		}
	}
	return &resbuf
}

func RunBase(cmd, in, out string, timeLimit, memoryLimit int64) TestInfo {
	infile, err := os.OpenFile(in, os.O_RDONLY, 0444)
	if err != nil {
		return TestInfo{
			Statu: "009",
			Info:  "测试输入用例不存在：" + in,
		}
	}
	outfile, err := os.OpenFile(out, os.O_RDONLY, 0444)
	if err != nil {
		return TestInfo{
			Statu: "009",
			Info:  "测试输出用例不存在：" + out,
		}
	}

	stdoutbuf := ReadFile(outfile)
	if stdoutbuf == nil {
		return TestInfo{
			Statu: "016",
			Info:  "文件读取出错：" + out,
		}
	}
	stdout := NewTestBuffer(stdoutbuf.Len() + 1024*16)

	stderr := bytes.Buffer{}

	timeUse, memoryUse, err := LimitExec.LimitExec(timeLimit, memoryLimit, cmd, infile, &stdout, &stderr)
	res := TestInfo{}
	if err != nil {
		if err.Error() == "MLE" {
			res.Info = "内存超限"
			res.Statu = "013"
		} else if err.Error() == "TLE" {
			res.Info = "运行超时"
			res.Statu = "011"
		} else if err.Error() == "OLE" {
			res.Info = "输出超限"
			res.Statu = "015"
		} else {
			res.Statu = "012"
			res.Info = "运行错误"
		}
	} else {
		if stdout.out {
			res.Info = "输出超限"
			res.Statu = "015"
		} else {
			if stdout.Check(stdoutbuf.Bytes()) {
				res.Info = "运行通过"
				res.Statu = "010"
			} else {
				res.Info = "程序运行结果错误"
				res.Statu = "014"
			}
		}
	}
	res.RunTime = timeUse
	res.Memory = memoryUse
	return res
}
