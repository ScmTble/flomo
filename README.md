<p align="center">
    flomo
    go编写的flomo-cli
<p/>


### 用法
1. 获取token，形如: Bearer xxxxx|xxxxxxxxxxxxxxxx
2. 安装
```bash
    go install github.com/scmtble/flomo@latest
```
3. 设置token / 登录
```bash
    # 设置token
    flomo token "token"
    
    # 登录
    flomo login "xxxxx" "xxxxxx"
```

4. 发送mono
```bash
    flomo new "Hello,World"
    
    # vim模式发送内容
    flomo vim
```
