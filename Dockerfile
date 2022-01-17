# 用户服务
FROM golang:1.14-alpine as builder

# 设置go module 以及代理
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装 git
RUN apk --no-cache add git

# 设置工作目录
WORKDIR /app/blog-user-service

COPY . .

# 下载依赖
RUN go mod download

# 重新下载net依赖包
RUN rm -rf $GOPATH/pkg/mod/golang.org/x/net\@v0.0.0-20220114011407-0dd24b26b47d
RUN cd $GOPATH/pkg/mod/golang.org/x && git clone https://github.com/golang/net.git net\@v0.0.0-20220114011407-0dd24b26b47d
RUN cd $GOPATH/pkg/mod/golang.org/x/net\@v0.0.0-20220114011407-0dd24b26b47d && git checkout release-branch.go1.14

WORKDIR /app/blog-user-service

# 由于重新下载了net依赖包需要重新加载一次所需要的依赖
RUN go mod tidy

# 打包go程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-user-service

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装相关软件
RUN apk update && apk add --no-cache bash ca-certificates

# 和上个阶段一样设置工作目录
RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/blog-user-service/blog-user-service .

CMD ["./blog-user-service"]