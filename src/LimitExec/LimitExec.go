package limitexec

import (
	debughelper "CUGOj-Judger/src/DebugHelper"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

type ExecError struct {
	Info string
}

func (e ExecError) Error() string {
	return e.Info
}

//以限制时间，限制内存的方式执行一个命令
//timeLimit：时间限制 ms
//memoryLimit：空间限制 KB
//args：命令
//stdin：标准输入
//stdout：标准输出
//stderr：标准错误
//返回 时间消耗（ms）内存消耗（KB） 错误信息
func LimitExec(timeLimit, memoryLimit int64, args string, stdin io.Reader, stdout, stderr io.Writer) (time, memery int64, err error) {
	//构造资源限制串，加入用户需要命令
	var str = fmt.Sprintf("ulimit -t %d;ulimit -m %d;", (timeLimit+999)/1000, memoryLimit+32768) + "cd /test/workspace;" + args
	//构造cmd
	cmd := exec.Command("sh", "-c", str)
	//重定向输入输出
	if stdin != nil {
		cmd.Stdin = stdin
	}
	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = &bytes.Buffer{}
	}
	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = &bytes.Buffer{}
	}
	//运行cmd并获取其rusage
	err = cmd.Run()

	//debug
	debughelper.ShowInfo("标准输出信息")
	debughelper.ShowBuf(cmd.Stdout.(*bytes.Buffer).Bytes())

	debughelper.ShowInfo("标准错误信息")
	debughelper.ShowBuf(cmd.Stderr.(*bytes.Buffer).Bytes())

	if err != nil {
		debughelper.ShowError(err)
	}

	statu := (*cmd).ProcessState
	rusage := (*statu).SysUsage().(*syscall.Rusage)

	//计算使用时间和内存
	timeUse := (rusage.Utime.Sec+rusage.Stime.Sec)*1000 + (rusage.Utime.Usec+rusage.Stime.Usec)/1000
	memoryUse := rusage.Maxrss

	if timeUse > timeLimit {
		timeUse = timeLimit
	}
	if memoryUse > memoryLimit {
		memoryUse = memoryLimit
	}

	if timeUse >= timeLimit {
		return timeUse, memoryUse, ExecError{
			Info: "TLE",
		}
	}
	if memoryUse >= memoryLimit {
		return timeUse, memoryUse, ExecError{
			Info: "MLE",
		}
	}
	if err != nil {
		if err.Error() == "exit status 137" {
			timeUse = timeLimit
			return timeUse, memoryUse, ExecError{
				Info: "TLE",
			}
		} else if err.Error() == "exit status 141" {
			return timeUse, memoryUse, ExecError{
				Info: "OLE",
			}
		} else {
			return timeUse, memoryUse, ExecError{
				Info: "RE",
			}
		}
	}
	return timeUse, memoryUse, err
}

//自定义输出的LimitExec  Out
func LimitExecO(timeLimit, memoryLimit int64, args string, stdout *bytes.Buffer) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, stdout, nil)
}

//自定义错误输出的LimitExec  Err
func LimitExecE(timeLimit, memoryLimit int64, args string, stderr *bytes.Buffer) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, nil, stderr)
}

//自定义标准输出和错误输出的LimitExec Out/Err
func LimitExecOE(timeLimit, memoryLimit int64, args string, stdout, stderr *bytes.Buffer) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, stdout, stderr)
}

//不处理输入输出的LimitExec None
func LimitExecN(timeLimit, memoryLimit int64, args string) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, nil, nil)
}

//以限制时间，限制内存的方式执行两个交互式的命令
//timeLimit：时间限制 ms
//memoryLimit：空间限制 KB
//args：命令
//stdin：标准输入
//stdout：标准输出
//stderr：标准错误
//返回 时间消耗（ms）内存消耗（KB） 错误信息
func IntExec(timeLimit, memoryLimit int64, args, spj, spjargs string) (int64, int64, error) {
	//构造资源限制串，加入用户需要命令
	var str = fmt.Sprintf("ulimit -t %d;ulimit -m %d;", (timeLimit+999)/1000, memoryLimit+32768) + "cd /test/workspace;" + args
	var spjstr = fmt.Sprintf("ulimit -t %d;ulimit -m %d;", (timeLimit+999)/1000, memoryLimit+32768) + "cd /test/workspace;" + spj + " " + spjargs
	debughelper.ShowInfo("main命令：" + str)
	debughelper.ShowInfo("spj运行命令：" + spjstr)
	//构造cmd
	cmd := exec.Command("sh", "-c", str)
	//重定向输入输出

	spjcmd := exec.Command("sh", "-c", spjstr)

	// buf1 := bytes.Buffer{}
	// buf2 := bytes.Buffer{}

	// cmd.Stdout = &buf1
	// spjcmd.Stdin = &buf1

	// spjcmd.Stdout = &buf2
	// cmd.Stdin = &buf2
	/*
	   go run src/main.go gnu cpp11 intrund /code/CUGOj-Judger/test/workspace/main 10000 256 /code/CUGOj-Judger/test/test1 /code/CUGOj-Judger/test/workspace/spj
	*/

	// spjcmd.Stdin, cmd.Stdout = os.Pipe()
	// cmd.Stdin, spjcmd.Stdout = os.Pipe()

	spjcmd.Stdin, _ = cmd.StdoutPipe()
	cmd.Stdin, _ = spjcmd.StdoutPipe()

	spjcmd.Stderr = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	//运行cmd并获取其rusage
	spjcmd.Start()
	cmd.Start()

	err := cmd.Wait()

	statu := (*cmd).ProcessState
	rusage := (*statu).SysUsage().(*syscall.Rusage)

	//计算使用时间和内存
	timeUse := (rusage.Utime.Sec+rusage.Stime.Sec)*1000 + (rusage.Utime.Usec+rusage.Stime.Usec)/1000
	memoryUse := rusage.Maxrss

	if timeUse > timeLimit {
		timeUse = timeLimit
	}
	if memoryUse > memoryLimit {
		memoryUse = memoryLimit
	}

	if timeUse >= timeLimit {
		return timeUse, memoryUse, ExecError{
			Info: "TLE",
		}
	}
	if memoryUse >= memoryLimit {
		return timeUse, memoryUse, ExecError{
			Info: "MLE",
		}
	}
	if err != nil {
		debughelper.ShowError(err)
		if err.Error() == "exit status 137" {
			timeUse = timeLimit
			return timeUse, memoryUse, ExecError{
				Info: "TLE",
			}
		} else if err.Error() == "exit status 141" {
			return timeUse, memoryUse, ExecError{
				Info: "OLE",
			}
		} else {
			return timeUse, memoryUse, ExecError{
				Info: "RE",
			}
		}
	}

	debughelper.ShowInfo("测试代码stderr输出")
	debughelper.ShowBuf(cmd.Stderr.(*bytes.Buffer).Bytes())

	err = spjcmd.Wait()
	if err != nil {
		debughelper.ShowError(err)
		return timeUse, memoryUse, ExecError{
			Info: "SE",
		}
	}

	// errch := make(chan error, 2)

	// goCnt := &sync.WaitGroup{}

	// goCnt.Add(1)
	// go func(cmd1, cmd2 *exec.Cmd, ch chan error) {
	// 	cmdBuf, cmderr := cmd1.Output()
	// 	cmd1.Wait()
	// 	fmt.Println(cmdBuf)
	// 	debughelper.ShowInfo("cmd运行结束")
	// 	ch <- cmderr
	// 	debughelper.ShowInfo("cmd返回结果")
	// 	goCnt.Done()
	// }(cmd, spjcmd, errch)

	// ac := false

	// goCnt.Add(1)
	// go func(cmd1, cmd2 *exec.Cmd, ch chan error) {
	// 	cmd1.Wait()
	// 	debughelper.ShowInfo("spjcmd运行结束")
	// 	goCnt.Done()
	// 	if cha, err := spjcmd.Stderr.(*bytes.Buffer).ReadByte(); err != nil || cha == 'w' || cha == 'W' {
	// 		ch <- ExecError{Info: "WA"}
	// 	} else {
	// 		ac = true
	// 	}
	// }(spjcmd, cmd, errch)

	// err := <-errch

	// debughelper.ShowInfo("交互题输出-测试代码输入")
	// debughelper.ShowBuf(buf2.Bytes())

	// debughelper.ShowInfo("交互题输入-测试代码输出")
	// debughelper.ShowBuf(buf1.Bytes())

	debughelper.ShowInfo("交互题stderr输出")
	debughelper.ShowBuf(spjcmd.Stderr.(*bytes.Buffer).Bytes())

	var ac bool
	if cha, err := spjcmd.Stderr.(*bytes.Buffer).ReadByte(); err != nil || cha == 'w' || cha == 'W' {
		ac = false
	} else {
		ac = true
	}
	if !ac {
		return timeUse, memoryUse, ExecError{
			Info: "WA",
		}
	}
	return timeUse, memoryUse, nil
}
