Sniper       [![Build Status](https://drone.io/github.com/lubia/sniper/status.png)](https://drone.io/github.com/lubia/sniper/latest)
======
>Sniper is a powerful and high-performance http load tester writing in Golang. Basing on advantage of goroutine,achieving high concurrency,low memory,rice graphics display.  

##Experience
Pre-compiled executables
* [Darwin 64 bit](http://lubia-me.qiniudn.com/sniper_darwin_amd64)      
* [Darwin 32 bit](http://lubia-me.qiniudn.com/sniper_darwin_386)
* [Linux 64 bit](http://lubia-me.qiniudn.com/sniper_linux_amd64)
* [Linux 32 bit](http://lubia-me.qiniudn.com/sniper_linux_386)
* [FreeBSD 64 bit](http://lubia-me.qiniudn.com/sniper_freebsd_amd64)
* [FreeBSD 32 bit](http://lubia-me.qiniudn.com/sniper_freebsd_386)

##Features
- GET / POST
- Keep-alive
- Https
- Graphics display result
- Multi-target
- Large file support
- Cross-platform——Linux,FreeBSD,Darwin

##Compare
<table class="table table-bordered table-striped table-condensed">
   <tr>
      <td>tool </td>
      <td>language </td>
      <td>keep-alive </td>
      <td>https </td>
      <td>multi-target </td>
      <td>result-show </td>
      <td>proxy</td>
   </tr>
   <tr>
      <td>ab </td>
      <td>c </td>
      <td>NO </td>
      <td>YES </td>
      <td>NO </td>
      <td>html，standard output</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>siege </td>
      <td>c </td>
      <td>YES </td>
      <td>YES </td>
      <td>YES </td>
      <td>csv，standard output</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>http_load </td>
      <td>c </td>
      <td>NO </td>
      <td>YES </td>
      <td>YES </td>
      <td>standard output</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>webbench </td>
      <td>c </td>
      <td>NO </td>
      <td>YES </td>
      <td>NO </td>
      <td>standard output</td>
      <td>YES </td>
   </tr>
   <tr>
      <td>sniper</td>
      <td>go</td>
      <td>YES </td>
      <td>YES </td>
      <td>YES </td>
      <td>js+html5，standard output</td>
      <td>NO </td>
   </tr>
</table>


##Performance
- Memory usage less than Apache Benchmark（ab）
- Execution speed close to ab
- Large file support

![Alt text](http://lubia-me.qiniudn.com/compare.jpg)

To get the detail of this performance comparison,[click me](http://www.lubia.me/http-loader-compare)

##Graphics display
- Analyse each request and record it
- Output http connect time
- Output server processing time
- Output total time

Basing on [dygraphs](http://dygraphs.com/)and html5，show the detail of server's performance. Get 1000 samples from whole result,show the details of connect,processing,and server's response.

The chart below show the total time and connect time. Wait,how can golang get the connect time ? In a word,Sniper implements part of HTTP protocol stack, discard net/http package to get the details. Also improve the performance. 

![Alt text](http://lubia-me.qiniudn.com/sniper_2.JPG)

##Usage manual
###1. install Golang

Please reference  [Go install](https://github.com/astaxie/build-web-application-with-golang/blob/master/ebook/01.1.md) chapter of open-source Golang book "build-web-application-with-golang".

###2. install Sniper

    $ go get github.com/lubia/sniper
    $ go install github.com/lubia/sniper
    $ cp src/github.com/lubia/sniper/.sniperc ~

###3. Parameter declaration

####Example
GET

    $sniper -c 10 -n 100 http://www.google.com 

POST

    $sniper -c 10 -n 100 -p postData.txt http://www.google.com
    
####Parameter

#####Command line parameter

```
Usage: 
   sniper [options] http[s]://hostname[:port][/path]                 http or https，ip or domain support
   sniper [options] -f urls.txt                                      multi-target，format：each url per line
Options: 
   -c, --concurrent     concurrent users, default is 1.              
   -n, --requests       number of requests to perform.               
   -r, --repetitions    number of times to run the test.             
   -t, --time           testing time, 30 mean 30 seconds.            
   -R, --sniperc        specify an sniperc file to get config        
                        (default is $HOME/.sniperc).               
   -f, --urlfile        select a specific URLS file.                 
   -p, --post           select a specific file to POST.              
   -T, --content-type   set Content-Type in request                  
                        (default is text/plain).
   -V, --Version        print the version number.                    
   -h, --help           print this section.                          
   -C, --config         show the current config.                     
   -s, --plot           plot detail transactions' info               
                        (true | false,default set true,              
                        notice: set -t will not plot anyhow).

```


#####Config file parameter

    Notice：default get config file from $HOME/.sniperc，config file and command line parameter complement each other. CMD -R to specified config file, CMD -C to get default configuration.  

```
[protocol]
version = HTTP/1.1                            
#connection = keep-alive                      # to comment
connection = close
accept-encoding = gzip                        
user-agent = golang & sniper                  

[header]
#cookie = SSID=Abh_TYcDc6YSQh-GB              user-defined header

[process]
timeout = 30                                  socket timeout 
failures = 64                                 max failure，socket failure over it then application break

[Authenticate]
login = jeff:supersecret                      HTTP Authenticate

[ssl]
ssl-cert = /root/cert.pem                     ssl-cert file
ssl-key = /root/key.pem                       ssl-key file 
ssl-timeout = 30                              https timeout
```

#####Output

chart output to current directory "plot.html"

```
Transactions:                   1000 hits           total requests
Availability:                   100.00 %            percentage completion    
Elapsed time:                   0.15 secs           sniper elapsed time
Document length:               1162 Bytes           single response length
TotalTransfer:                  1.11 MB             total data transfer
Transaction rate:            6625.60 trans/sec      transactions per second 
Throughput:                     7.34 MB/sec         throughput 
Successful:                     1000 hits           result code not 200 also successful
Failed:                           0 hits            socket errors 
TransactionTime:               1.495 ms(mean)       each request total time (average)
ConnectionTime:                0.596 ms(mean)       connect time(average，tcp handshake)
ProcessTime:                   0.900 ms(mean)       TransactionTime = ConnectionTime + ProcessTime
StateCode:                    1000(code 200)        the number of result code is 200
```
##About 
####Twin projects

- [gohttpbench](https://github.com/parkghost/gohttpbench)
- [vegeta](https://github.com/tsenart/vegeta)

####Author

Lubia Yang,programmer in finance

Blog：[Program Design](http://www.lubia.me)

Contact：yanyuan2046 at 126.com

####Licence
[Apache License, Version 2.0.](http://www.apache.org/licenses/LICENSE-2.0.html)

###中文文档
[点击此处](https://github.com/lubia/sniper/blob/master/README_CN.md)
