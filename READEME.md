# 服务器框架

#支持中间件
#支持注册对象的方法为路由（自动路由）


### 存储路劲
* 数据存储到 datapath
* 配置文件存储在 confpath
* 日志文件存储在 logpath

### 创建目录
    src 源代码管理
    build 管理编译脚本的地方
    log 管理日志的地方
    main 主函数入口

### 需要的东西在src里面写(对象和对应方法)
    route来调用


### 登录
    token
    指定某些ip能访问（中间件）
    
    方案一：（不支持单用户操作，不支持修改密码后让登录的用户重新登录,不支持每个用户一个签名认证）用户登录之后服务器什么都不存， 信息全放在token中， 所有token都用同一个固定的签名
    方案二：（不支持单用户操作，支持修改密码后让登录的用户重新登录，支持每个用户一个签名认证）用户登录之后把用户信息存在内存中,当修改密码后修改内存中的信息。
            内存中的信息有效期为30分钟，当token带的用户信息在内存查不到，去数据库查询
            token采用jwt
    方案三：（支持单用户操作，支持修改密码后让登录的用户重新登录，支持每个用户一个签名认证）在第二方案上加一个 session

    状态
        其他用户登录
        修改密码

    用户表

    username(明文用户名) password（密码+salt的md5）salt（盐） ....(后期拓展(比如用户权限等)

    用户信息
        增删查改
            删和改操作 影响内存的用户信息 
    
设置 cookie

