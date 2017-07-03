# go engine

## 目的

golang虽然编写程序非常方便，但一个完善的程序需要考虑很多方便，每次新写一个项目总
是需要去设置初始化一些共同的东西：信号，日志，GC等。本项目旨在让开发人员从这些琐
事中解脱出来，专注于自己的业务，进一步提高生产力。本项目适用于流式处理场景，类似
storm式的流式处理。经生产上验证，storm程序实际使用起来性能差强人意，而且依赖多，
搭建起来比较费劲，twitter已经放弃。

## 使用

本项目提供nsq数据源的支持，并且提供redis，mongo等工具，使用的时候只需要实现`flowprocess.Handler`
接口，详见项目`github.com/guyannanfei25/flowprocess`，将文件放在src/consumer文件
夹下面，并在creator_factory.go CreatorFactory函数中增加相应代码，case中的名字以及
json配置文件中consumer配置项名称与代码中handler名称一致。
