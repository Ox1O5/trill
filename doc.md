# v0.1
基础的server模块只包括
* 启动
* 停止
* 服务

# v0.2
封装链接
针对不同的业务定义专用API

# v0.3
路由
把所有请求封装到request模块里


# v0.4
全局控制

# v0.5 
因为tcp是面向字节流的,多个消息同时发送,接收端读取时可能无法区分两个消息的边界,需要在应用层提供封包拆包的方法
在golang的net包里定义了一个接口叫net.Conn,它被注释为
Conn是一个通用的面向流的网络连接

gc,我在写这个项目的时候用xx查看过程序的运行效率,然后发现xx可能是GC的原因