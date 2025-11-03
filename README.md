## 项目启动
1. 在linux环境中（wsl，虚拟机都可以）配置好mysql、redis、kafka，根据实际情况修改config/config.yaml。

2. 下载nginx，修改resources/nginx-1.18.0/conf/nginx.conf文件，修改其中的前端所在位置。

3. 启动main.go