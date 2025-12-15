# UFProject
## 数据结构课设
### 题目:最小生成树问题
### 技术栈：
1.语言采用Golang 1.25

2.后端采用Gin+Gorm

3.数据库采用Sqlite

4.前端采用vue

### 运行项目
本项目有两种运行方法，一种是以本地直接run main.go；
另一种是run docker；启动成功后访问`http://localhost:8098/` 来显示页面

#### 法一
在本项目根目录下使用终端运行 `cd frontend`切换至前端目录，再运行 `npm run build`生成前端文件，再`cd ..`返回根目录，最后`go run cmd/server/main.go`启动服务

#### 法二
在本项目根目录下使用终端运行 `docker build -t mst-app .`将项目打包成docker后再运行 ` docker run -d -p 8098:8098 --name mst-demo mst-app`

### 注意事项
本项目会占用8098端口
