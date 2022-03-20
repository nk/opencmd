package utils

import (
	"fmt"
	"opencmd/config"
	"os"
	"path/filepath"
)

// get current dir
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// get parent dir of dir
func GetParentDir(dir string) string {
	return filepath.Dir(dir)
}

func PrintCommands(commands []string) {
	for _, command := range commands {
		fileName := filepath.Base(command)
		fmt.Println(fileName, "\n   ", command)
	}
}

func GetAllParentDir(dir string) []string {
	dirList := []string{}

	for {
		dirList = append(dirList, dir)
		parentDir := GetParentDir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}
	return dirList
}

func GetAllCommands() []string {
	var AllCommands []string
	dirList := GetAllParentDir(GetCurrentDir())
	for _, dir := range dirList {
		commands, err := FindCommands(dir)
		if err != nil {
			fmt.Println(err)
			continue
		}
		AllCommands = append(AllCommands, commands...)
	}
	return AllCommands
}

func FindCommandByName(dir, cmd string) (string, error) {
	dirList := GetAllParentDir(dir)
	for _, dir := range dirList {
		command, err := FindCommandByNameInSingleDir(dir, cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if command != "" {
			return command, nil
		}
	}
	return "", fmt.Errorf("can not find command: %v", cmd)
}

func FindCommandByNameInSingleDir(dir, cmd string) (string, error) {
	files, err := FindCommands(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		// get file name
		fileName := filepath.Base(file)
		if fileName == cmd {
			return file, nil
		}
	}
	return "", nil
}

func FindCommands(dir string) ([]string, error) {

	CommandDirName := filepath.Join(
		dir,
		config.DefaultConfig.OpencmdDir,
		config.DefaultConfig.CommandDir,
	)
	pexist, perr := DirExists(CommandDirName)
	if perr != nil {
		panic(perr)
	}
	if !pexist {
		return nil, nil
	}

	// find all files
	var files []string
	err := filepath.Walk(CommandDirName, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// check directory exists
func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
