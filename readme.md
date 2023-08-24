# 基于go语言的智能客服系统

之前写了一个带后端的python版本，带简易后端，但是不稳定，性能不好，并且很垃圾。传送门：https://github.com/haynb/chatgpt_Read_Local

#### 现在更新了go语言版本！

用户端问答基于chatgpt大模型，向量数据库基于qdrant。

### 官方文档：

chatgpt：https://platform.openai.com/docs/introduction

qdrant：https://qdrant.tech/documentation/

### 目前进度：

本项目目前使用go语言http包进行数据库交互。

后续有时间将会完善为docker-compose搭建部署，使用docker虚拟网络进行数据库操作，兼具安全性与实用性。

### 部署教程：

先安装数据库：

```shell
docker pull qdrant/qdrant
docker run -p 6333:6333 \
    -v $(pwd)/qdrant_storage:/qdrant/storage \
    qdrant/qdrant
```

安装完成之后，即可使用，克隆本项目：

```shell
git clone https://github.com/haynb/go_gpt_kefu.git
```

### 程序配置：

在运行之前请先：

```go
go mod tidy
```

之后进行如下配置：
db文件夹下qdrant.go即为数据库配置，请把

```go
var (
	QdrantBase = "1**.**.2*.***"//数据库ip地址
	QdrantPort = "6333"//数据库的开放端口
	id_file    = "id.txt"
)
```

修改为你自己搭建的数据库的配置。id_file为方便添加points时使用的辅助文件，详情请阅读源码。

gpt文件夹下gpt.go即为gpt配置

请自行阅读代码，更改自己的`api_baseurl`和`openai_key`为自己的即可。具体的prompt也可以自行`append`。

main.go即为程序运行入口。

### 最后：

```go
go run .
```

