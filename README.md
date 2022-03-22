
Table of Contents
----------------------------
- [opencmd 是什么](#opencmd-是什么)
  - [# 为什么需要opencmd](#-为什么需要opencmd)
- [快速开始](#快速开始)
  - [安装](#安装)
  - [命令参考](#命令参考)

# opencmd 是什么
  opencmd帮你管理、运行各种零散的脚本的命令行工具


# 为什么需要opencmd
------------------------------
日常有很多功能单一的脚本需要执行，比如

- 只在本地开发环境用到的初始化和清理临时文件的脚本
- 项目手动打包部署到某一台测试机器
- 有些很少会用到，但是内容很复杂的命令行
- Makefile里不适合放到git里的命令
- 一些ugly但是确实快捷的自定义脚本
- ......

一般在开发环境下，不同的项目有不同的命令管理方式，比如放到Makefile里，package.json 的 scripts 字段等, 但这些命令一般会加上版本控制公开给项目的所有开发者，对上面列举的那临时的、测试用的、不完善的、有用但是用的很少的脚本，很多时候是加到一个 cheat-sheet 文件，需要时复制出来运行一次。

opencmd 试图解决这个问题，基本上可以把它的运行方式描述为:
1. 在项目根目录创建 ".opencmd/commands/" 目录， 并把上面提到的所有脚本放进去
2. 在项目的根目录或者任意子目录，用 opencmd list 可以列出所有命令，用 opencmd run [filename] 运行命令

把脚本放到单独的文件内，而不是像 Makefile 或者 package.json 那样集中在一个文件内，主要是因为两个原因:
1. 放到单独的文件不限制脚本的复杂度， 可以把任意长的脚本管理起来
2. 方便对脚本分组，对.gitignore来说更友好，从而把用到的脚本都放在项目内而不用管git的问题，比如可以在commands目录建立一个private目录，在.gitignore文件直接忽略掉整个目录。


# 快速开始
## 安装

```
curl "https://github.com/nk/opencmd/releases/download/v0.2.1/opencmd_0.2.1_linux_amd64.gz" |gunzip -c > /usr/local/bin/opencmd

chmod a+x /usr/local/bin/opencmd

echo 'alias oc="opencmd"' >> ~/.bashrc
echo 'alias ocr="opencmd run"' >> ~/.bashrc
```

## 命令参考
```
Usage:
  opencmd [command]

Available Commands:
  help        Help about any command
  list        show available commands and path
  run         run a command
  version     show version

Flags:
  -h, --help   help for opencmd

Use "opencmd [command] --help" for more information about a command.
```
