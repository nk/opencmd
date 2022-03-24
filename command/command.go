package command

import (
	"bufio"
	"fmt"
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

type ExecParams struct {
	Name       string
	Args       []string
	WorkingDir string
}

func (cmd *Command) Run() (string, error) {

	funcArray := []func() (string, error){
		cmd.RunCommandDirectly,
		cmd.RunCommandViaShebang,
		cmd.RunCommandUseDefaultShell,
	}

	for _, f := range funcArray {
		// 防止由于错误导致执行多次
		if cmd.Done {
			break
		}

		out, err := f()
		if err != nil {
			cmd.Output = out
			cmd.Err = err
			continue
		}
		cmd.Done = true
		return out, nil
	}
	return cmd.Output, cmd.Err
}

func (cmd *Command) RunCommandWithParams(params *ExecParams) (string, error) {

	execCmd := exec.Command(params.Name, params.Args...)
	execCmd.Dir = params.WorkingDir
	out, err := execCmd.CombinedOutput()
	return string(out), err
}

func (cmd *Command) RunCommandDirectly() (string, error) {
	perm, err := HasExecPermition(cmd.Path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if !perm {
		err = fmt.Errorf(
			"You don't have permission to run this command: %v",
			cmd.Path,
		)
		return "", err
	}

	// fixme: 只有指定了参数时才会把workingdir 设为OpencmdBase
	cmdParams := ExecParams{
		Name:       cmd.Path,
		Args:       []string{},
		WorkingDir: cmd.OpencmdBase,
	}
	return cmd.RunCommandWithParams(&cmdParams)
}

func (cmd *Command) RunCommandViaShebang() (string, error) {
	shebang, err := GetShebang(cmd.Path)
	if err != nil {
		return "", err
	}

	// fmt.Println("run command [shebang]:", cmd.Path, " ", shebang)

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

func (cmd *Command) RunCommandUseDefaultShell() (string, error) {
	shell := config.DefaultConfig.DefaultShell
	if shell == "" {
		shell = os.Getenv("SHELL")
	}
	if shell == "" {
		shell = "/bin/bash"
	}

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
		return "", fmt.Errorf("%v has no shebang", file)
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
