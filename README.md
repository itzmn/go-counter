# 概览

go-counter主要用于接收请求报文，按照预设好的统计逻辑基于redis存储计算，生成结果json，接口返回

# 快速开始

## 本地启动
- redis需要准备好

1、预定义统计逻辑，见配置文件 ./config/statisticVars.json

```go
[
  {
    "filter": [
      {
        "func": "eq",
        "path": "name",
        "params": "zhangsan",
        "type": "string"
      }
    ],
    "function": "distinct",
    "dimensions": [
      {
        "path": "activityId",
        "type": "string"
      }
    ],
    "data": {
      "path": "user",
      "type": "string"
    },
    "window": {
      "type": "time",
      "size": 86400
    },
    "name": "user_cnt_per_activity_1d",
    "type": "int"
  }
]
```

2、修改配置文件 ./config/config.json 中redis的信息和http服务端口

3、运行 go-counter.go， http服务运行于19999端口

4、构造报文请求服务

```go
curl -d '{"activityId":"ac123", "amount":12, "name":"zhangsan", "user":"zhangsan", "timestamp":1733989840000, "requestId":"r123"}' http://localhost:19999/counter
```

5、得到预定义统计逻辑的结果并返回

```go
{"activityId":"ac123", "amount":12, "name":"zhangsan", "user":"zhangsan", "timestamp":1733989840000, "requestId":"fa6b7a82a55d421f4ef9649a22915d05","user_cnt_per_activity_1d":1}
```
得到统计结果
- user_cnt_per_activity_1d 去重用户数等于1
- requestId 该笔请求的流水


## docker启动

1、启动redis

```shell
docker run -d -p 16379:6379 --network jk-network --network-alias redis --name redis redis的镜像
```

由于编译好的项目需要连接redis，redis的容器需要和服务容器形成网络，构建了一个jk-network网络，也执行了redis的网络名称就是redis，在服务里面直接访问redis就可以访问到结果了

2、启动服务
```shell
docker run -d -p 19999:19999 -p 20000:20000 --network jk-network --network-alias counter-03 --name go-counter-v1.0.3 go-counter:v1.0.3
```
和redis是相同的网络，config.json中连接redis的地址就是redis

3、测试
```shell
curl -d '{"organization":"itzmn","activityId":"ac123", "amount":12, "name":"zhangsan", "user":"zhangsan", "timestamp":1733989840000, "requestId":"r123"}' http://localhost:19999/counter
{"organization":"itzmn","activityId":"ac123", "amount":12, "name":"zhangsan", "user":"zhangsan", "timestamp":1733989840000, "requestId":"dd2a73a44810cedebbc5798e6de3ed79","user_cnt_per_activity_1d":1,"amount_sum_per_activity_user_1d":12}
```

# 项目介绍

预设统计逻辑
```go
{
    "filter": [
        {
            "func": "eq",
            "path": "name",
            "params": "zhangsan"
        }
    ],
	"function":"count",
    "dimensions": [
        {
            "path": "activityId",
            "type": "string"
        }
    ],
	"data":{"path":"tokenId","type":"string"},
    "window": {
        "type": "time",
        "size": "86400"
    },
    "name": "cnt_per_activity_1d"
}



{"filter":[{"func":"eq","path":"name","params":"zhangsan","type":"string"}],"function":"count","dimensions":[{"path":"activityId","type":"string"}],"data":{"path":"requestId","type":"string"},"window":{"type":"time","size":"86400"},"name":"cnt_per_activity_1d"},
{"filter":[{"func":"eq","path":"name","params":"zhangsan","type":"string"}],"function":"sum","dimensions":[{"path":"activityId","type":"string"}],"data":{"path":"amount","type":"int"},"window":{"type":"time","size":"86400"},"name":"amount_sum_per_activity_1d"}
```

逻辑介绍
- name(必须)

  统计后的结果的名字
- function

  统计函数,支持 count,distinct,sum
- filter

  统计时候的过滤条件，可以不加
- dimensions

  统计维度
- window(必须)

  统计窗口，size 86400 数据为滑动窗口 秒
- data

  要统计的数据，字段

# 统计逻辑

从redis获取某个变量数据后，从当前时间向前滑动获取周期内的数据，进行汇总。 得到结果，同时更新redis内数据，
将当前笔请求更新到redis，且将超出时间周期的数据丢弃，不写入redis

# redis数据结构

value使用map结构，特定前缀counter: 加上维度真实值 作为 redis的key， value的map的key是统计变量名称，value是统计的值或者临时值

## 数据存储结构
counter:维度
- 统计变量名字 cnt_per_tokenId_1d
- 统计值 [{"ts": 1670000000, "val": 100}, {"ts": 1670000600, "val": 200}, ...]


### 各类函数数据存储格式
根据统计函数的类型区别，统计值存储不一样，由于是滑动窗口统计，需要切分一下时间槽，预设如下,当数据是第一次生成时候，ts以当前时间为准，后续的槽时间，按照时间+槽大小

|  统计窗口   | 槽大小  |
|  ----  | ----  |
| 0-10  | 1s |
| 11-60  | 2s |
| 60-1800  | 60s |
| 1801-21600  | 600s |
| 21601-86400  | 3600s 优先实现 | 


- count函数

  根据统计窗口大小，切分槽，存储每一个时间槽内的计数值，在计算时候，需要更新redis中槽的数据[{"ts":123,"val":2},{"ts":124,"val":4}]

[comment]: <> (- distinct函数，存储统计周期内的每个槽去重原始数据，记录最新一次出现时间[{"val":"zmn","ts":1234},{"val":"lisi","ts":1235}])

[comment]: <> (- sum函数 存储每一个时间槽内的统计汇总值，[{"ts":123,"val":2},{"ts":124,"val":4}])


