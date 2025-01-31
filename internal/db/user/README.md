# DB/User

数据库操作-用户相关操作

### getLoginInfo

包含函数**GetLoginInfo**，传入用户名（登录账户）string，正常返回显示用户名，密码。错误返回error。

### getUsername

包含函数**GetUsername**，传入token string，正常返回用户名，错误返回error。该函数**会重置Redis用户信息的保留时间**。

### setToken

包含函数**SetUserToken**，传入用户名（登录账户），token string，正常无返回，错误返回error。该函数会在Redis中保存用户信息，并设置保留时间。