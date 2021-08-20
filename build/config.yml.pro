database:
  # 腾讯云数据库
  crawler:
    host: 172.21.32.16
    user: root
    pass: Guduodatabase123!@#
    dbname: guduo
    port: 3306
  clean:
    host: 172.21.32.16
    user: root
    pass: Guduodatabase123!@#
    dbname: business
    port: 3306

  # 只有数据迁移用到，但是代码里写死了所以这里写上
  lolipop:
    host: 172.21.0.207
    user: root
    pass: aipheu6us0Ea
    dbname: lollipop
    port: 3306
  carl:
    host: 172.21.0.205
    user: root
    pass: aBoh3caer9oo
    dbname: carl
    port: 3306
# 实际上没在用
redis:
  host: 127.0.0.1
  pass:
  db: 0
  port: 6381

# 代理的IP地址-在 82.156.34.31 上
proxy:
  host: proxyv2.guduodata.com
  port: 80
  secret: ueqy8qfyrhwov9c8f5ksnkxlrnk8eo33