# 微服务网关

所有服务统一入口。

请求报文：
```text
query{
    member{
    
    },
}
```


响应报文：



##部署说明

编译镜像：docker pull golang:1.13-alpine
运行环境：docker pull mszs/alpine:3.10

Docker service 
```
// docker镜像源
docker service create --name gateway --network cluster -p 39701:80 --replicas 1  -d mszs/gateway:latest
// 阿里云镜像源
docker service create --name gateway --network cluster -p 39701:80 --replicas 1  -d registry.cn-hangzhou.aliyuncs.com/mszs/gateway:latest

docker service update --force --update-parallelism 1 --update-delay 3s gateway
docker service update  --replicas 3  gateway


```

Docker service 
```
docker service create --name gateway --network cluster -p 39701:80  --replicas 1  -d hub.unionlive.com/gateway:latest
docker service update --force --update-parallelism 1 --update-delay 3s --image hub.unionlive.com/gateway:latest gateway
docker service update  --replicas 3  gateway

```

测试环境接口地址：http://211.152.57.29:39701/api/v2/graphql  

## Chanage Log 

待实现功能
- [Feature] 根据服


### v1.0.0 
    [Release 2019-09-11 ]
- [Feature] 创建项目 















































































































