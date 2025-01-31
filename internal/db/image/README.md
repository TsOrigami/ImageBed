# DB/Image

数据库操作-图片相关操作

### savePic

包含函数**SaveInfoToSql**，传入图片名，用户名，哈希值（string），创建时间（time.Time，[]uint8），将其数据保存至数据库。数据库相关信息从db/sql/connect.go中获取。正常无返回值，错误返回error。

### queryPic

包含函数**GetInfoQuery**，传入起始位置，结束位置，哈希值（string），从数据库查询范围内符合条件的图片的uuid并返回。正常返回string数组（[]string），错误返回error。

### inquirePic

包含函数**GetInfoByUUID**，传入UUID（string），从数据库查询范围内符合条件的图片数据并返回。正常返回PicInfo类型（自定义结构体，包含UUID，ImageName，Sha256Hash，CreatedAt string），错误返回error。

### deletePic

包含函数**DeleteInfoFromSQL**，传入用户名，UUID（string），从数据库查询范围内符合条件的图片的并删除（从正常表移到删除表）。正常无返回，错误返回error。