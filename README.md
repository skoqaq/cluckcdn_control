# CDN Control 

**Official Website:** [https://cluckcdn.buzz](https://cluckcdn.buzz)<br>
**Documentation (Traditional Chinese):** [https://cluckcdn.buzz/docs/](https://cluckcdn.buzz/docs/)

-----

> 简体中文 README: [README_CN.md](README_CN.md)

> Please deploy the control centre on a separate server to ensure proper communication with all CDN nodes.

> Recommended for Ubuntu/Debian servers (Centos is also available)

> We are still in the testing stage and welcome your comments

### 1. Build

> Please install Golang

```bash
# Clone
git clone https://github.com/ArsFy/cluckcdn_control.git
cd cluckcdn_control
# Build
go build .
```

### 2. Run!

```bash
rm -rf *.go
chmod 775 control
./control
```

### 3. Change setting 

#### Node communication: `/static/config.yaml`

You can modify "textToken" to communicate with other nodes, but please do not modify {ctrlServer} (This is an escape character)

```yaml
control: {ctrlServer}
token: textToken
```

#### Admin / Node List: `/config.json`

You can change your username and password and add more admins

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

#### Vhost(WebSite): `vhost.json`

Manual modification of this configuration file is not recommended, you can change it on the web.

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

### 4. Using the web-based admin panel

> English is temporarily not supported

![](https://ci.cncn3.cn/0cde5a1ad3cec0ad54ad9100fc22786b.png)
![](https://ci.cncn3.cn/c2ab373ce3b7d57da3d8109d85460349.png)

#### Error Page
![](https://ci.cncn3.cn/9738cba86de009a32dba739516a453ad.png)

## Special thanks

Suddenly fucking want to write CDN's me