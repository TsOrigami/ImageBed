# API

代码参考 /internal/handers/http

前后端通信预计使用websocket进行连接。http提供外部api

### 登录 login

接口：/login

方法：POST

参数：

​		Body参数

​		内容格式：application/x-www-form-urlencoded

​			account:  string;

​			salt:  string;

​			sign:  string;

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					token: string;

​					username: string;

*sign的计算方式：

账号(account)+盐(salt)+密码 组合之后的字符串 通过md5计算 得到32位sign

*token存在有效期。在/configs/config.conf中设置[redis]的remains字段

### 注册 register

接口：/register

方法：POST

参数：

​		Body参数

​		内容格式：application/x-www-form-urlencoded

​			account:  string;

​			user:  string;

​			password:  string;

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					Register: "success";

### 上传 upload

接口：/upload

方法：POST

参数：

​		Headers 参数

​			Authorization: Bearer + {{token}};

​		Body 参数

​		内容格式：multipart/form-data

​			{{参数名任意}}: file;

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					code: 200

​					msg: "文件上传成功"

*Headers的token：

​	为登录时返回的token

*Body的”参数名任意“：

​	该接口的读取为选取Body中的file类型，参数名不作影响

### 查询 query

接口：/query

方法：POST

参数：

​		Headers 参数

​			Authorization: Bearer + {{token}};

​		Body 参数

​		内容格式：application/x-www-form-urlencoded

​			start: int/string

​			end: int/string

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					code: 200

​					uuid: string[]

*Headers的token：

​	为登录时返回的token

*Body的”start“和"end"：

​	查询的起点和终点。起点不可大于终点。

*返回的"uuid"：

​	字符串数组。每一个字符串为一个uuid，用于调用和删除

### 删除 delete

接口：/delete

方法：POST

参数：

​		Headers 参数

​			Authorization: Bearer + {{token}};

​		Body 参数

​		内容格式：application/x-www-form-urlencoded

​			uuid: string

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					code: 200

​					msg: "删除成功: "+ {{uuid}}

*Headers的token：

​	为登录时返回的token

*Body和返回的”uuid“：

​	查询时获取的图片的uuid。

### 调用 invoke

接口：/invoke

方法：get

参数：

​		query 参数:

​			uuid: string

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容:

​				返回的图片

*Body的”uuid“：

​	查询时获取的图片的uuid。

### 调用缩略图 thumbnail

接口：/thumbnail

方法：get

参数：

​		query 参数:

​			uuid: string

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容:

​				返回的图片的缩略图

*Body的”uuid“：

​	查询时获取的图片的uuid。

### 获取图片总数 amount

接口：/amount

方法：GET

参数：

​		Headers 参数

​			Authorization: Bearer + {{token}};

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					code: 200

​					Count: int

*Headers的token：

​	为登录时返回的token

### 获取通信地址 prefix

接口：/prefix

方法：GET

返回响应：

​		成功(200)：

​			HTTP 状态码: 200

​			内容格式:  JSONapplication/json

​					Register: string

*实际没什么用 因为获取通信地址的前提是你能访问通。但是你能访问通了说明你有通信地址了。。。
