#!/usr/bin/python3
import random
import os
import uvicorn
import requests
import json
import mimetypes

from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse

# 全局变量
global img_dict, img_path, json_path, alia_dict, app
img_path = './goodjob-img/resources'
json_path = 'https://gitee.com/SmallK111407/useless-plugin/raw/main/model/aliasData/alias.json'
img_dict = {}
alia_dict = {}
app = FastAPI()

# 加载别名数据
def load_alias() -> dict:
    global json_path
    return json.loads(
        requests.get(json_path).text
    )

# 创建别名映射
def alias_map() -> None:
    _alias = load_alias()
    global alia_dict
    for _name, _alias in _alias.items():
        if not _name in list(img_dict.keys()):
            continue
        alia_dict[_name] = _name
        for _alia in _alias:
            alia_dict[_alia] = _name

# 用于动态构建 handler
class ServeHandler:

    def __init__(
            self,
            alia: str
    ) -> None:
        self.alia = alia

    def __call__(self) -> object:
        return pick_img(alia_dict[self.alia])

# 随机抓取图片
def pick_img(
        name: str
) -> object:
    global img_dict, img_path

    if name not in img_dict:
        raise HTTPException(
            status_code=404,
            detail="Name Not Found"
        )
    file_path = os.path.join(img_path, name, random.choice(img_dict[name]))
    file_types, _ = mimetypes.guess_type(file_path)

    return FileResponse(
        path=file_path,
        media_type=file_types,
        status_code=200
    )

# 加载图片路径
def load_path() -> None:
    global img_dict, img_path, alia_dict
    folder_list = os.listdir(img_path)
    img_dict = { _name: os.listdir( os.path.join(img_path, _name) ) for _name in folder_list }
    
    alias_map()
    alia_dict.update(
        {_name: _name for _name, _ in img_dict.items() if _name not in list(alia_dict.keys())}
    )
    
    # 注册路由节点
    [register(_alia) for _alia, _ in alia_dict.items()]

# 导入图片，注册路由
def load_point() -> None:
    # 我也不知道我为啥把这个扔这里了
    os.system('cd ./goodjob-img && git pull -f')
    # 加载、注册路由节点
    load_path()

# 随机一个名字
def random_name() -> str:
    global img_dict
    return random.choice(
        list(img_dict.keys())
    )

# 注册路由节点
def register(alia):
    global app
    app.add_api_route(
        f'/{alia}',
        ServeHandler(alia),
        methods=['GET']
    )

# '/' 节点的路由方法
def direct() -> object:
    return pick_img(
        random_name()
    )

# '/HELPS' 节点的路由方法
def help() -> dict:
    return {
        '用法': '发送 GET 请求时会从词库随机抽取一个图片，如果需要指定某个人则可以 GET 对应地址',
        '请求方式': '发送 GET 请求获取任意图片，发送 GET 请求到对应地址获取某个人的图片',
        '别名处理': '以 https://gitee.com/SmallK111407/useless-plugin/blob/main/model/aliasData/alias.json 中的别名为准',
        '快速获得别名映射表': 'GET https://img-api.justrobot.dev/AliasMap',
        '示例': 'GET https://img-api.justrobot.dev/win11 来获得 win11 的随机图片',
        '更新频率': 'API 图库会在每天的 00:05 左右重启进行图库更新，耗时 10 秒以内',
        '速率限制': '图库有访问限制，单 IP 每秒限制 1 次，每 10 秒限制 20 次(包含无效访问), 超过任意限制均返回 429 错误',
        '图片类型': '图片类型可能为 png 或 gif , 如果图库有其他图片类型会原样提供'
    }

def AliasMap() -> dict:
    return alia_dict

# 入口函数
def main() -> None:
    global app
    load_point()
    app.add_api_route('/HELPS', help, methods=['GET'])
    app.add_api_route('/AliasMap', AliasMap, methods=['GET'])
    app.add_api_route('/', direct, methods=['GET'])
    uvicorn.run(app, host='0.0.0.0', port=10808)


if __name__ == '__main__':
    main()
