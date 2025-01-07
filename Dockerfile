#依赖镜像
FROM golang:1.19-alpine3.18

#作者信息
MAINTAINER "zhangmengnan"

# 配置模块代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct


RUN mkdir -p /opt/code
#工作目录
WORKDIR /opt/code
# 将当前目录内容拷贝到容器的/opt目录
ADD .  /opt/code/

#在Docker工作目录下执行命令，编译服务
RUN /opt/code/bin/build.sh

RUN cp -rav /opt/code/dist /opt/go-counter

WORKDIR /opt/go-counter

#声明服务端口
EXPOSE 19999 20000

#执行项目的命令，启动服务
CMD ["./bin/docker-start.sh"]