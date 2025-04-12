# 使用 Golang 1.20 的官方镜像作为基础镜像
FROM golang:1.20

# 将工作目录切换到 /app 目录下
WORKDIR /app

ENV ELASTICSEARCH_HOST="http://elasticsearch:9200"

# 将当前目录下的所有文件复制到 /app 目录下
COPY . /app

# 编译应用程序
RUN go build -o my_parser .

# 运行应用程序
CMD ["./my_parser"]
