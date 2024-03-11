# 图库API源码

## API 说明

* 用法：发送 GET 请求时会从词库随机抽取一个图片，如果需要指定某个人则可以 GET 对应地址
* 浏览器使用：额外增加一个 direct 根目录，比如<br>
 `https://img-api.justrobot.dev/direct/` 相当于 `https://img-api.justrobot.dev/`
* 请求方式：发送 GET 请求获取任意图片，发送 GET 请求到对应地址获取某个人的图片
* 别名处理：以 `https://gitee.com/SmallK111407/useless-plugin/blob/main/model/aliasData/alias.json` 中的别名为准
* 快速获得别名映射表：`GET https://img-api.justrobot.dev/AliasMap`
* 示例：`GET https://img-api.justrobot.dev/oldcitynight` 来获得 oldcitynight 的随机图片
* 更新频率：API 图库会在每天的 00:05 左右重启进行图库更新，耗时 10 秒以内
* 速率限制：图库有访问速率限制：<br>单 IP 每秒限制 1 次，每 10 秒限制 20 次(包含无效访问), 超过任意限制均返回 429 错误
* 图片类型：图片类型可能为 png 或 gif , 如果图库有其他图片类型会原样提供

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

* nodejs

 ```python
 npm install pnpm -g
 pnpm i
 node app.js
 ```
 然后就搭建好了

## OPEN-SOURCED LICENSE
  * The project is **open-sourced** under the **[GNU General Public License Version 3](https://github.com/oldcitynight/yunzai-goodjob-img-api/blob/master/LICENSE)**.