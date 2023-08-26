# Intelligent customer service system based on go language
## The backend module has been added

I wrote a python version with a backend before, with a simple backend, but it is unstable, the performance is not good, and it is rubbish. Portal: https://github.com/haynb/chatgpt_Read_Local

#### The go language version is now updated!

The user-side Q&A is based on the chatgpt large model, and the vector database is based on qdrant.

### Official document:

chatgpt: https://platform.openai.com/docs/introduction

qdrant: https://qdrant.tech/documentation/

### Current progress:

This project currently uses the go language http package for database interaction.

In the future, when there is time, we will improve the deployment of docker-compose, and use the docker virtual network for database operations, which is both safe and practical.

### Deployment tutorial:

Install the database first:

```shell
docker pull qdrant/qdrant
docker run -p 6333:6333 \
     -v $(pwd)/qdrant_storage:/qdrant/storage \
     qdrant/qdrant
```

After the installation is complete, it can be used to clone this project:

```shell
git clone https://github.com/haynb/go_gpt_kefu.git
```

### Program configuration:

Before running please:

```go
go mod tidy
```

Then configure as follows:
qdrant.go under the db folder is the database configuration, please put

```go
var (
QdrantBase = "1**.**.2*.***"//database ip address
QdrantPort = "6333"//The open port of the database
id_file = "id.txt"
)
```

Modify the configuration for your own database. id_file is an auxiliary file for adding points, please read the source code for details.

gpt.go under the gpt folder is the gpt configuration

Please read the code by yourself and change your `api_baseurl` and `openai_key` to your own. The specific prompt can also `append` by itself.

main.go is the program entry point.

### at last:

```go
go run .
```

The interface of the program has been integrated in main.go.

Including: upload and automatically parse files, query file list, use chat.