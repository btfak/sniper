Sniper       [![Build Status](https://drone.io/github.com/lubia/sniper/status.png)](https://drone.io/github.com/lubia/sniper/latest)
======
>Sniper是一个功能强大、高性能的HTTP负载工具,采用Golang编写。利用协程并发优势，实现海量并发、

>超低内存占用、丰富图表展示。是测试、分析、优化服务端性能的绝佳助手！

##体验
提供以下可执行文件，可直接运行
* Mac OSX 64 bit      
* Mac OSX 32 bit
* Linux 64 bit
* Linux 32 bit

##功能
以实用为原则，实现以下功能
- GET / POST
- keep-alive模式
- https
- 图表展示结果
- 测试多个目标
- 支持大文件负载
- 跨平台，支持Linux,FreeBSD,Darwin

####对比同类工具
<table class="table table-bordered table-striped table-condensed">
   <tr>
      <td>工具 </td>
      <td>编写语言 </td>
      <td>keep-alive </td>
      <td>https </td>
      <td>多点测试 </td>
      <td>结果展示 </td>
      <td>代理</td>
   </tr>
   <tr>
      <td>ab </td>
      <td>c </td>
      <td>NO </td>
      <td>YES </td>
      <td>NO </td>
      <td>html，标准输出</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>siege </td>
      <td>c </td>
      <td>YES </td>
      <td>YES </td>
      <td>YES </td>
      <td>csv，标准输出</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>http_load </td>
      <td>c </td>
      <td>NO </td>
      <td>YES </td>
      <td>YES </td>
      <td>标准输出</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>webbench </td>
      <td>c </td>
      <td>NO </td>
      <td>YES </td>
      <td>NO </td>
      <td>标准输出</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>sniper</td>
      <td>go</td>
      <td>YES </td>
      <td>YES </td>
      <td>YES </td>
      <td>js+html5，标准输出</td>
      <td>NO </td>
   </tr>
</table>


##性能
- 内存占用低于Apache Benchmark（ab）等主流负载工具
- 执行速度接近ab，高并发时超过ab
- 支持10k以上并发
- 支持超大文件测试

![Alt text](http://lubia-me.qiniudn.com/cmp.png)

测试的详细情况，与各大负载测试工具的性能对比[在此](http://www.lubia.me/http-loader-compare)

##图表展示
- 统计分析每个请求
- 输出建立连接时间
- 输出服务端响应时间
- 输出总时间

基于[dygraphs](http://dygraphs.com/)与html5，详细展现服务端性能情况

从测试结果中等距采样约1000样本，详细展现连接建立，链路传输和服务端执行情况

下图展示了总时间和连接建立时间的对比

![Alt text](http://lubia-me.qiniudn.com/sniper_2.JPG)

##使用说明
###1. 安装Golang

请参考astaxie的开源Golang书籍《Go Web 编程》一书，[Go安装](https://github.com/astaxie/build-web-application-with-golang/blob/master/ebook/01.1.md)一节。

###2. 安装Sniper

    $ go get github.com/lubia/sniper
    $ cp src/github.com/lubia/sniper/.sniperc ~

###3. 参数说明

####示例
GET

    $sniper -c 10 -n 100 http://www.google.com 

POST

    $sniper -c 10 -n 100 -p postData.txt http://www.google.com
    
####参数

#####命令行参数

```
Usage: 
   sniper [options] http[s]://hostname[:port][/path]                 http或https，支持域名或ip
   sniper [options] -f urls.txt                                      测试多个服务端地址，文件格式：每个url一行
Options: 
   -c, --concurrent     concurrent users, default is 1.              并发数(默认为1)
   -n, --requests       number of requests to perform.               总请求数
   -r, --repetitions    number of times to run the test.             重复次数(n=c*r)
   -t, --time           testing time, 30 mean 30 seconds.            测试时间(单位秒)
   -R, --sniperc        specify an sniperc file to get config        配置文件地址(默认为$HOME/.sniperc)
                        (default is $HOME/.sniperc).               
   -f, --urlfile        select a specific URLS file.                 多个测试目标的url文件
   -p, --post           select a specific file to POST.              POST模式
   -T, --content-type   set Content-Type in request                  POST的数据类型(默认为text/plain)
                        (default is text/plain).
   -V, --Version        print the version number.                    打印sniper版本号
   -h, --help           print this section.                          输出帮助信息
   -C, --config         show the current config.                     输出当前配置文件的配置
   -s, --plot           plot detail transactions' info               是否输出html展示测试结果(默认为true) 
                        (true | false,default set true,              (注意:采用-t指定测试时间时,不会输出html)
                        notice: set -t will not plot anyhow).

```


#####配置文件参数

    说明：默认从$HOME/.sniperc读取配置文件，配置文件设置与命令行设置互为补充
    可通过命令行 -R 指定配置文件地址，-C 查看默认配置。

```
[protocol]
version = HTTP/1.1                            HTTP协议版本，1.1或1.0
#connection = keep-alive                      connection模式，# 符号作为注释
connection = close
accept-encoding = gzip                        
user-agent = golang & sniper                  

[header]
#cookie = SSID=Abh_TYcDc6YSQh-GB              自定义消息头，等号连接键值对

[process]
timeout = 30                                  socket超时时间 
failures = 64                                 最大失败次数，socket错误超过此值则程序退出

[Authenticate]
login = jeff:supersecret                      HTTP基本认证

[ssl]
ssl-cert = /root/cert.pem                     ssl-cert文件地址
ssl-key = /root/key.pem                       ssl-key文件地址
ssl-timeout = 30                              https超时
```

#####结果输出

图表输出到当前目录下plot.html

```
Transactions:                   1000 hits           总请求数
Availability:                   100.00 %            完成百分百    
Elapsed time:                   0.15 secs           sniper执行时间
Document length:               1162 Bytes           服务端单个返回长度
TotalTransfer:                  1.11 MB             总传输数据量
Transaction rate:            6625.60 trans/sec      每秒事务数 
Throughput:                     7.34 MB/sec         吞吐量 
Successful:                     1000 hits           成功次数(结果码不为200也是成功)
Failed:                           0 hits            失败次数(socket等链路错误) 
TransactionTime:               1.495 ms(mean)       单个请求总耗时(平均)
ConnectionTime:                0.596 ms(mean)       链路建立耗时(平均，tcp三次握手)
ProcessTime:                   0.900 ms(mean)       服务端执行时间+传输时间(TransactionTime = ConnectionTime + ProcessTime)
StateCode:                    1000(code 200)        结果码为200的数量
```
##关于
####作者

Lubia Yang,程序员

博客：[程式設計](http://www.lubia.me)

联络：yanyuan2046 at 126.com

寻找Golang or C 开发工作中,坐标北京or深圳

####Licence
[Apache License, Version 2.0.](http://www.apache.org/licenses/LICENSE-2.0.html)
