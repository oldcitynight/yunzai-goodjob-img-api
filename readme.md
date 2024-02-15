# 图库API源码

## API 说明

* 请求方式: 发送普通 GET 请求到指定路径
* 返回方式: 图片文件(png/gif)
* 示例: GET https://img-api.justrobot.dev 可以随机抽取一张图片<br>
GET https://img-api.justrobot.dev/win11 可以从 win11 的图集中随机抽取一张图片
* 速率限制: 每个 IP 每秒 1 次限制, 每 10 秒限制 20 次访问(包括访问失败), <br> 
超过速率限制均会返回 429 , 前者无拉黑措施，后者会拉黑 10 秒禁止访问

## 搭建方法

* python:
  ```python
  pip install fastapi uvicorn
  ```
  然后下载脚本，放到无用图库上一级目录运行

* golang

1. 将二进制可执行文件放置于无用图库上一级目录下。
2. 运行可执行文件以启动图库API。

> 也可以选择编译，就是很麻烦，我就不多说了
