# CONFIGS

此次处理配置相关部分。getConfig.go中为解析信息相关代码。config.conf为配置文件

### config.conf：

##### [mysql] 数据库相关配置

字段：

​	host：主机

​	post：端口

​	user：用户名

​	password：密码

​	dbname：数据库名称

##### [redis] Redis相关配置

字段：

​	addr：主机

​	port：端口

​	password：密码

​	db：Redis表编号

​	remains：用户登录Token保存时间（仅用于api请求。单位：分钟）

##### [server] 其他配置

字段：

​	imagesPath：图片保存的路径

​	thumbnailPath：缩略图保存的路径