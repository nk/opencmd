package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/nk/opencmd/config"
)

type Command struct {
	Name        string // 名称
	Path        string // 完整路径
	OpencmdBase string // .opencmd目录所在的路径，一般为项目根目录

	Done   bool   // 是否已经执行
	Output string // 运行结束后的输出
	Err    error  // 运行结束后的错误
}

type CommandMeta struct {
	IsBin            bool
	Shebang          string
	HasExecPermition bool
}

type ExecParams struct {
	Name       string
	Args       []string
	WorkingDir string
}

func (cmd *Command) Run() error {
	commandMeta, err := cmd.GetCommandMeta()
	if err != nil {
		return err
	}

	return cmd.RunCommandWithMeta(commandMeta)
}

func (cmd *Command) RunCommandWithMeta(meta *CommandMeta) error {
	if meta.IsBin {
		// 二进制文件
		if !meta.HasExecPermition {
			err := fmt.Errorf(
				"You don't have permission to run this command: %v",
				cmd.Path,
			)
			return err
		}
		return cmd.RunCommandDirectly()
	} else {
		if meta.Shebang != "" {
			// 有shebang, 不管有没有可执行权限，都使用shebang执行
			// 如果有shebang同时有可执行权限，直接执行和使用shebang执行的结果一样
			return cmd.RunCommandViaShebang()
		} else {
			// 没有shebang，不管有没有可执行权限，都使用默认shell
			return cmd.RunCommandUseDefaultShell()
		}
	}
}

func (cmd *Command) GetCommandMeta() (*CommandMeta, error) {
	meta := &CommandMeta{}

	isBin, err := IsBinaryFile(cmd.Path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	shebang, err := GetShebang(cmd.Path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	hasExecPermition, err := HasExecPermition(cmd.Path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	meta.IsBin = isBin
	meta.Shebang = shebang
	meta.HasExecPermition = hasExecPermition

	// fmt.Printf("command meta: %+v\n", meta)
	return meta, nil
}

func (cmd *Command) RunCommandWithParams(params *ExecParams) error {

	execCmd := exec.Command(params.Name, params.Args...)
	execCmd.Dir = params.WorkingDir
	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr
	execCmd.Env = os.Environ()
	err := execCmd.Run()
	if err != nil {
		fmt.Println("run command [err]:", err)
	}
	return err
}

func (cmd *Command) RunCommandDirectly() error {
	perm, err := HasExecPermition(cmd.Path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !perm {
		err = fmt.Errorf(
			"You don't have permission to run this command: %v",
			cmd.Path,
		)
		return err
	}
	// fmt.Println("run command [direct]:", cmd.Path)
	// fixme: 只有指定了参数时才会把workingdir 设为OpencmdBase
	cmdParams := ExecParams{
		Name:       cmd.Path,
		Args:       []string{},
		WorkingDir: cmd.OpencmdBase,
	}
	return cmd.RunCommandWithParams(&cmdParams)
}

func (cmd *Command) RunCommandViaShebang() error {
	shebang, err := GetShebang(cmd.Path)
	if err != nil {
		return err
	}
	// fmt.Println("run command [shebang]:", cmd.Path, " ", shebang)
	if shebang == "" {
		err = fmt.Errorf("%v has no shebang", cmd.Path)
		return err
	}

	shebang_list := strings.Split(shebang, " ")
	shebang_list = append(shebang_list, cmd.Path)

	// fixme: 只有指定了参数时才会把workingdir 设为OpencmdBase
	cmdParams := ExecParams{
		Name:       shebang_list[0],
		Args:       shebang_list[1:],
		WorkingDir: cmd.OpencmdBase,
	}
	return cmd.RunCommandWithParams(&cmdParams)
}

func (cmd *Command) RunCommandUseDefaultShell() error {
	shell := config.DefaultConfig.DefaultShell
	if shell == "" {
		shell = os.Getenv("SHELL")
	}
	if shell == "" {
		shell = "/bin/bash"
	}

	// fmt.Println("run command [shell]:", shell, cmd.Path)
	// fixme: 只有指定了参数时才会把workingdir 设为OpencmdBase
	cmdParams := ExecParams{
		Name:       shell,
		Args:       []string{cmd.Path},
		WorkingDir: cmd.OpencmdBase,
	}
	return cmd.RunCommandWithParams(&cmdParams)
}

func GetShebang(file string) (string, error) {
	firstLine, err := ReadFirstLine(file)
	if err != nil {
		return "", err
	}

	hasShebang := strings.HasPrefix(firstLine, "#!")
	if !hasShebang {
		return "", nil
	}

	shebang := strings.TrimPrefix(firstLine, "#!")
	return shebang, nil
}

// read file's first non-blank line
func ReadFirstLine(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// read file's first non-blank line
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		trimedLine := strings.TrimSuffix(line, "\n")
		if trimedLine != "" {
			line = trimedLine
			break
		}
	}
	return line, nil
}

func HasExecPermition(path string) (bool, error) {
	// check file wheather is executable
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false, err
	}

	return fi.Mode().Perm()&0100 != 0, nil
}

func IsBinaryFile(file string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	return !isText(buf[:n]), nil
}

// check file is binary or text
func isText(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	isText := true
	for _, b := range data {
		if b == 0 {
			return false
		}
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			isText = false
			break
		}
	}

	return isText
}
