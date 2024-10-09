# 一个简单的图床

goland + mysql
运行 main.go以启动项目。默认端口为8000

*接口文档没写完 2点半了 该睡觉了

## 接口

#### 上传(upload)

​	接口请求地址：/upload

​	功能说明：上传图片（大小不大于100MB，支持同时上传多个图片）

​	请求方法：POST

​	请求参数：数量任意 

​						参数名：任意参数名

​						单个文件大小不大于100MB

#### 删除(delete)

​	接口请求地址：/delete

​	功能说明：删除图片

​	请求方法：DELETE

​	请求参数：参数名：md5	要删除的图片的md5值

#### 查询(inquire)

​	接口请求地址：/inquire

​	功能说明：获取指定范围内的图片数据

​	请求方法：GET

​	请求参数：参数名：start	查询数据范围的起始点

​						参数名：end	查询数据范围的结束点

​	响应参数：一个数组。每一个包含四个数据：

​						参数名：ImageName	图片名

​						参数名：IImageURL	图片的实际储存位置

​						参数名：IMD5Hash	图片的MD5值

​						参数名：ICreatedAt	图片的上传时间

#### 获取(fetch)

​	接口请求地址：/fetch

​	功能说明：获取图片

​	请求方法：GET

​	请求参数：参数名：md5	要获取的图片的md5值

​	响应参数：查询的图片

#### 通过API获取(image_api)

​	URL：/image/api/ + 需要获取的图片的md5值

​	浏览器直接访问或调用

## 数据库

在 SQL_Operate/config.go 中修改数据库的连接信息

该项目会在数据库下创建一个表image_info，包含image_name，md5_hash，created_at三个列。主键为md5_hash

## 文件储存目录

在 SQL_Operate/config.go 中修改文件的储存目录
