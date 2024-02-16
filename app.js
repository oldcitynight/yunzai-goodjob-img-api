const express = require('express');
const fs = require('fs');
const exec = require('child_process').exec;

const app = express();
const img_path = __dirname + '/goodjob-img/resources/';
const img_dict = {};

const randint = (num) => {
    return Math.floor(Math.random() * num);
};

const routes = () => {
    return Object.keys(img_dict);
};

const isDir = (path) => {
    return fs.statSync(img_path + path).isDirectory();
};

const lenPath = (path) => {
    let conut = 0
    files = fs.readdirSync(path);
    files.forEach(file => {
        if ( !( file.endsWith('.png') || file.endsWith('.gif') ) ) { return };
        conut += 1;
    });
    return conut;
};

const pickImg = ( name ) => {
    let path = img_path + name + '/' + randint(img_dict[name]);
    if ( name === '茄子' ) { path += '.gif' } else { path += '.png' };
    console.log('200 OK:' + path)
    return path;
};

const pickName = () => {
    let name = routes();
    return name[randint(name.length)];
};

const ReadPath = () => {
    const files = fs.readdirSync(img_path);
    files.forEach(file => {
            if ( isDir(file) ) { img_dict[file] = lenPath(img_path + file) };
        }
    );
};

const update = () => {
    exec('cd '+ __dirname +'/goodjob-img && git pull -f', (err, stdout, stderr) => {
        if (err) { console.log(err) } else {};
        if (stderr) { console.log(stderr) } else {};
        if (stdout) { console.log(stdout) } else {};
    });
};

ReadPath();
console.log('Folder Loaded: \n' + routes());

app.get('/*', (req, res) => {
    let route = routes();
    const _path = decodeURIComponent(req.path.slice(1));
    if ( _path === '/') {
        console.log('New Request at root')
        res.sendFile( pickImg( pickName() ) );
        return;
    };

    if ( _path === '/HELPS' ) {
        res.send({
            '用法': '发送 GET 请求时会从词库随机抽取一个图片，如果需要指定某个人则可以 GET 对应地址',
            '请求方式': '发送 GET 请求获取任意图片，发送 GET 请求到对应地址获取某个人的图片',
            '示例': 'GET https://img-api.justrobot.dev/win11 来获得 win11 的随机图片',
            '更新频率': 'API 图库会在每天的 00:05 左右重启进行图库更新，耗时 10 秒以内',
            '速率限制': '图库有访问限制，单 IP 每秒限制 1 次，每 10 秒限制 20 次(包含无效访问), 超过任意限制均返回 429 错误',
        });
        return;
    };

    if ( route.includes(_path) ) {
        console.log('New Request at ' + _path);
        res.sendFile(pickImg(_path));
    } else {
        console.log('New Request at Invalid Path ' + _path);
        res.sendStatus(404);
    }
    return;
});

app.listen(10808, () => {
    console.log('Server Listening at 10808')
})