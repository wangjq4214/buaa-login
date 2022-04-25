# buaa-login

[![unit-test](https://github.com/wangjq4214/buaa-login/actions/workflows/test.yml/badge.svg)](https://github.com/wangjq4214/buaa-login/actions/workflows/test.yml)

提供在命令行中登录北航校园网的程序，编译为单一的可执行文件，免去安装依赖。

## 登录命令

执行下面命令可以登录北航校园网：

```bash
$ ./buaa-login login -u by2106100 -p xxx
```

## 开启一个守护进程

执行下面命令可以开启一个后台进程进行轮询，当出现断网情况时自动登录：

```bash
$ ./buaa-login daemon -u by2106100 -p xxx
```
