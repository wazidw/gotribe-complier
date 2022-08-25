# gotribe-complier

## 介绍
- 在线编译器
- 实现语言： golang + docker


使用前需先安装 docker，并事先下载运行代码的镜像
```
docker pull python:3
docker pull rust
docker pull golang
docker pull php:5.6
docker pull php:7
docker pull php:8
```
测试

```
curl -H "Content-Type: application/json" -X POST -d '{"lang":"golang","code":"package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"Hello, gotribe\")\n}","input":""}' http://127.0.0.1:9091/run

curl -H "Content-Type: application/json" -X POST -d '{"lang":"rust","code":"fn main() {\n    println!(\"Hello, gotribe!\");\n}","input":""}' http://127.0.0.1:9091/run

```
