# CDN Control 

**官方网站:** [https://cluckcdn.buzz](https://cluckcdn.buzz)<br>
**文档 (繁体中文):** [https://cluckcdn.buzz/docs/](https://cluckcdn.buzz/docs/)

-----

> English README: [README.md](README.md)

> 请将控制中心部署在单独的服务器上，以确保与所有 CDN 节点正确通信

> 推荐用于 Ubuntu/Debian 服务器（Centos 也不是不行）

> 我们仍处于测试阶段，欢迎您的意见 

### 1. 编译

> 请安装 Golang (>=1.15)

```bash
# Clone
git clone https://github.com/cluckcdn/control.git
cd control
# Downloads lib
go get github.com/prometheus/common/log
go get github.com/gin-contrib/sessions
go get github.com/gin-contrib/sessions/cookie
go get github.com/gin-gonic/gin
go get gopkg.in/yaml.v3
# Build
go build .
```

### 2. 运行！

```bash
rm -rf *.go
chmod 775 control
./control
```

### 3. 修改设置

#### 节点通讯: `/static/config.yaml`

您可以修改 "textToken" 来与其他节点通信，但请不要修改 "{ctrlServer}" （这是一个转义符） 

```yaml
control: {ctrlServer}
token: textToken
```

#### 管理员 / 节点列表: `/config.json`

您可以更改您的用户名和密码并添加更多管理员 

```json
{
    "admin": {
        "cluckbird": "123456",
        "Test": "123456"
    },
    "node": [
        {
            "ip": "192.168.48.138",
            "name": "TestNode"
        }
    ]
}
```

#### Vhost(网站): `vhost.json`

不建议手动修改此配置文件，您可以在管理员面板上更改

```json
[
    {
        "host": "testnode.com",
        "name": "TestWebSite",
        "proto": "https",
        "source": "172.217.31.227",
        "source_host": "www.google.com.hk",
        "text": "Test",
        "tls": false
    },
    {
        "host": "192.168.48.138",
        "key": "/node/tls/192.168.48.138.key",
        "name": "Test",
        "pen": "/node/tls/192.168.48.138.pen",
        "proto": "https",
        "source": "172.217.31.227",
        "source_host": "www.google.com.hk",
        "text": "Test",
        "tls": true
    }
]
```

### 4. 使用基于 Web 的管理面板 

> 暂时不支持英文

![](https://ci.cncn3.cn/0cde5a1ad3cec0ad54ad9100fc22786b.png)
![](https://ci.cncn3.cn/c2ab373ce3b7d57da3d8109d85460349.png)

#### 错误页面
![](https://ci.cncn3.cn/9738cba86de009a32dba739516a453ad.png)

-----

## 特别鸣谢

脑子抽了突然想自己写CDN的我自己