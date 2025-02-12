# HP-Lite 3.0内网穿透
> fork from https://gitee.com/HServer/hp-lite 2025.2.12

## HP Lite介绍
HP-Lite3.0是一个单机方案
我们采用的是数据转发实现 稳定性可靠性是有保证的即便是极端的环境只要能上网就能实现穿透。
我们支持TCP和UDP协议，针对 http/https ws/wss 协议做了大量的优化工作可以更加灵活的控制。让用户使用更佳舒服简单。

## 服务端
### 配置说明
> app.yml文件放在和二进制同目录下即可，不然会默认启动失败
```yaml
admin:
  username: 'admin' #后台账号
  password: '123456' #后台密码
  port: 9090 #管理后台监听的端口（TCP传输方式）

cmd:
  port: 6666 #控制指令端口，所有HP-lite 客户端需要连接这个端口（TCP传输方式）

tunnel:
  ip: '127.0.0.1' #隧道监听服务器外网的IP（记得改成你的服务器IP或者解析的域名也可以）
  port: 9091 #隧道传输数据端口，这个端口用来传输数据的，注意这个是UDP协议，如果是安全组设置记得UDP的放开
  open-domain: true #true 开启80，443端口域名转发（如果你的服务有宝塔或者nginx等，端口多半是被用了），false 关闭
acme:
  email: '232323@qq.com' #申请证书必须写一个邮箱可以随便写
  http-port: '5634' #证书验证会访问http接口，会通过80转发过来，所以这个端口不用暴露外网
```

### 服务端docker部署
```bash
docker run -d \
    --name hp \
    --net host \
    -v ./data:/app/data \
    -v ./app.yml:/app/app.yml \
    laoyutang/hp-lite-server:latest
```

## 客户端运行方式
### docker
```bash
# 通过 docker run 运行容器
sudo docker run --restart=always -d  -e server=xxx.com穿透服务:6666 -e deviceId=32位的设备ID registry.cn-shenzhen.aliyuncs.com/hserver/hp-lite:latest
# 通过 docker run 运行容器 ARM
sudo docker run --restart=always -d  -e server=xxx.com穿透服务:6666 -e deviceId=32位的设备ID registry.cn-shenzhen.aliyuncs.com/hserver/hp-lite:latest-arm64
```
### Linux或者win
```bash
chmod -R 777 ./hp-lite-amd64
./hp-lite-amd64 -server=xxx.com穿透服务:6666 -deviceId=32位的设备ID 
```

## 项目截图
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img.png"  />
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img_1.png"  />
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img_4.png"  />
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img_5.png"  />
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img_6.png"  />
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img_7.png"  />
<img src="https://gitee.com/HServer/hp-lite/raw/main/doc/img/img_8.png"  />

