# Blog

## 前端
### 框架&依赖
1. 原生 JS
2. Tailwind CSS
3. Marked(highlight)
4. katex
5. remixicon

### 运行
配置 `public/static/basic.js` 文件，修改 `path` 为后端接口地址，
`cp ./basic.js.sample ./public/static/basic.js`
直接访问 `public/index.html`

## 后端
### 框架&依赖
1. Gin
2. SQLite

### 运行
配置 `.env` 文件，
```shell
cp env.sample .env
vim .env
```
然后执行
```shell
make run
```
或者使用 docker 运行：
```shell
make docker-run
```
详细查看`Makefile`文件