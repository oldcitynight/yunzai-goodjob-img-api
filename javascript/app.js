const express = require('express');
const fs = require('fs');
const exec = require('child_process').exec;

const lodash = require('lodash');
const request = require('sync-request');
const app = express();
const img_path = __dirname + '/goodjob-img/resources/';
const res_path = "https://gitee.com/SmallK111407/goodjob-img"
const img_dict = {};

const alia_json = JSON.parse(
    request(
        'GET',
        'https://gitee.com/SmallK111407/useless-plugin/raw/main/model/aliasData/alias.json'
    ).getBody('utf-8')
);

const isDir = (path) => {
    return fs.statSync(img_path + path).isDirectory();
};

const pickImg = ( alia ) => {
    point_name = alias_map[alia];
    var file_name
    while (file_name == undefined) {
        file_name = lodash.sample(img_dict[point_name])
    }
    let path = img_path + point_name + '/' + file_name;
    console.log('200 OK:' + path)
    return path;
};

const pickName = () => {
    let name = Object.keys(img_dict);
    var result;
    while (result == undefined) {
        result = lodash.sample(name);
    };
    return result;
};

const ReadPath = () => {
    const dirs = fs.readdirSync(img_path);
    dirs.forEach(dir => {
            if ( isDir(dir) ) { 
                img_dict[dir] = fs.readdirSync(img_path + dir);
            };
        }
    );
};

const update = () => {
    if (!fs.existsSync(__dirname + '/goodjob-img')) {
        exec('git clone' + res_path, (err, stdout, stderr) => {
            if (err) { console.log(err) } else {};
            if (stderr) { console.log(stderr) } else {};
            if (stdout) { console.log(stdout) } else {};
        });
    };

    exec('cd '+ __dirname +'/goodjob-img && git pull -f', (err, stdout, stderr) => {
        if (err) { console.log(err) } else {};
        if (stderr) { console.log(stderr) } else {};
        if (stdout) { console.log(stdout) } else {};
    });
};

ReadPath();

const get_alias_map = () => {
    let map = {};
    for (let key in alia_json) {
        if ( !Object.keys(img_dict).includes(key) ) {
            continue;
        };
        map[key] = key;
        alia_json[key].forEach( (value) => {
            map[value] = key;
            return;
        });
    };

    return map;
};

const alias_map = get_alias_map();

const routes = () => {
    return Object.keys(alias_map);
};

console.log('Folder Loaded:');

app.get('/*', (req, res) => {
    let route = routes();

    const _path = decodeURIComponent(req.path.slice(1));
    const result = _path.match(/^direct\/(.+)$/);

    if ( _path === '') {
        console.log('New Request at root')
        res.download(pickImg(pickName()));
        return;
    };

    if ( _path === 'HELPS' ) {
        res.send({
            '用法': '发送 GET 请求时会从词库随机抽取一个图片，如果需要指定某个人则可以 GET 对应地址',
            '请求方式': '发送 GET 请求获取任意图片，发送 GET 请求到对应地址获取某个人的图片',
            '别名处理': '以 https://gitee.com/SmallK111407/useless-plugin/blob/main/model/aliasData/alias.json 中的别名为准',
            '快速获得别名映射表': 'GET https://img-api.justrobot.dev/AliasMap',
            '示例': 'GET https://img-api.justrobot.dev/oldcitynight 来获得 oldcitynight 的随机图片',
            '更新频率': 'API 图库会在每天的 00:05 左右重启进行图库更新，耗时 10 秒以内',
            '速率限制': '图库有访问限制，单 IP 每秒限制 1 次，每 10 秒限制 20 次(包含无效访问), 超过任意限制均返回 429 错误',
            '图片类型': '图片类型可能为 png 或 gif , 如果图库有其他图片类型会原样提供',
        });
        return;
    };

    if ( _path == 'AliasMap' ) {
        res.send(alias_map);
        return;
    };

    if ( _path == 'direct' ) {
        res.sendFile(
            pickImg(pickName())
        );
        return;
    };

    if ( result ) {
        res.sendFile(
            pickImg(result[1])
        );
        return;
    }

    if ( route.includes(_path) ) {
        console.log('New Request at ' + _path);
        res.download(pickImg(_path));
        return;
    } else {
        console.log('New Request at Invalid Path ' + _path);
        res.sendStatus(404);
        return;
    };
});

update();

app.listen(10808, () => {
    console.log('Server Listening at 10808')
});
