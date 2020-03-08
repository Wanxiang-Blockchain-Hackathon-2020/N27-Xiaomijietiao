# 基于Taro的小借条小程序

### 如何运行
先讲一下如何将本项目起来。

``` bash
## clone
git clone https://github.com/despard/weyom

## install deps
npm install

## compile
npm run dev:weapp
```

然后把项目文件夹导入微信开发者工具即可预览。


### 前言
Taro是一个遵循React语法规范的多端开发的框架。自从其开源以来，就一直关注，从一开始0.x版本，到后来的1.0版本。时至今日，Taro已经发展成为一个成熟的框架，社区也逐渐完善。有一点值得点赞的是，Taro的维护人员处理问题很快，开发者很快能得到响应，专注认真的开源态度，我觉得是Taro发展好的原因之一。

安装Taro:https://nervjs.github.io/taro/docs/GETTING-STARTED.html

### 整体架构

https://www.processon.com/view/link/5dd8918ae4b052b7c58c6dfa#map
https://www.processon.com/view/link/5dce5eb9e4b0096e8c093998
后台架构：https://www.processon.com/diagraming/5ddb2b78e4b052b7c58e878a

### 业务流程

* 打借条：https://www.processon.com/view/link/5dd7fc2ee4b0f7888da4268c

* 还款提醒：

* 还款：

### 数据模型

* 借条
https://www.processon.com/view/link/5dd7f588e4b0bbcb8a695786
to be continue

### 后台接口

* 添加用户
* 用户借条查询
* 还款提醒


#### 目录结构
``` bash

├─ config
├─ client/dist            ##  编译目录
├─ qcloud/cloudFunctions  ##  云函数
└─ src
      ├─ app.js
      ├─ app.scss
      ├─ components       ##  公用组件
      ├─ index.html
      ├─ pages            ##  主页面
      │    ├─ index
      │    └─ person
      ├─ static           ##  静态资源
      │    ├─ icon
      │    └─ image
      ├─ subpages         ## 子包
      └─ utils.js         ## 工具函数
├─ package.json
├─ project.config.json
├─ .editorconfig
├─ .eslintrc
├─ .gitignore
├─ LICENSE
├─ README.md

```
