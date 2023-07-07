#### gRPC 服务注册与发现特征

gRPC 本身没有提供服务注册接口，也就是服务注册本身，是我们自己管的。

gRPC 只提供了服务发现的接口，也就是Resolver 和对应的 Builder 两个接口。

gRPC 服务发现的核心接口是:

+ ClientConn: 它是一个抽象，代表的是对一个服务(**应用维度的服务**)的连接，而不是一个 TCP 连接。当我们调用`grcp.Dial`的时候会返回ClientConn
+ resolver.Builder：负责创建 Resolver。gRPC 会维护一个 scheme-> Builder 的映射。
+ Resolver：它和服务进行绑定，一个服务一个 Resolver。它负责和注册中心交互，监听
  注册数据的变化

一般步骤：

+ 用户在初始化 gRPC 的时候指定 grpc.WithResolver选项，传入自定义的 Resolver。
+ 在 Dial 调用的时候传入服务标识符，一般形式是scheme:///service-name。scheme 代表的是如何通信。大多数时候，它就代表我们的注册中心。
+ gRPC 会根据 scheme 来找到我们注册的 Resolver，我们在 Resolver 里面更新可用的连接。  


#### 具体实现:
先看 grpc_resolver.go 学习如何实现 grpc 服务发现
参考链接 : https://grpc.io/docs/guides/custom-name-resolution/