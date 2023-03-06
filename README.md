<!-- PROJECT SHIELDS -->
[![LICENSE](https://img.shields.io/github/license/basefas/admin-go.svg?style=flat-square)](/LICENSE)
[![Releases](https://img.shields.io/github/release/basefas/admin-go/all.svg?style=flat-square)](https://github.com/basefas/admin-go/releases)
![GitHub Repo stars](https://img.shields.io/github/stars/basefas/admin-go?style=social)

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a>
    <img src="https://raw.githubusercontent.com/basefas/files/main/logo.svg" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">admin-go</h3>

  <p align="center">
     一个使用 Go 语言开发的管理系统
    <br />
  </p>
</div>

<!-- Introduction -->

## 简介

admin 使用 Go & Gin & Gorm 开发，包含常用后台使用的基本模块，同时提供完整功能的前端程序，可快速用于二次开发及功能扩展。


|           | url                                         | introduction                               |
|-----------|---------------------------------------------|--------------------------------------------|
| backend   | https://github.com/basefas/admin-go         | 使用 Go & Gin 开发的后台管理系统后端           |
| frontend  | https://github.com/basefas/react-antd-admin | 使用 react & vite & antd 开发的后台管理系统前端|


## 页面截图

### 登录页面

![Screen Shot](https://github.com/basefas/files/blob/main/login.png)

### 用户管理

![Screen Shot](https://github.com/basefas/files/blob/main/user.png)

### 分组管理

![Screen Shot](https://github.com/basefas/files/blob/main/group.png)

### 菜单管理

![Screen Shot](https://github.com/basefas/files/blob/main/menu.png)

### 角色及权限管理

![Screen Shot](https://github.com/basefas/files/blob/main/permission.png)


<!-- GETTING STARTED -->

## 快速开始

1. 克隆项目到本地

```
git clone https://github.com/basefas/admin-go
```

2. 安装依赖

```
go mod download
```

3. 运行

```
go run ./cmd/app/main.go
```

<!-- USE DOCKER -->

## 使用 Docker 部署

> 注：需要提前安装好 docker 和 docker-compose

1. 切换目录

```
cd  ./deploy/docker-compose
```

2. 使用 docker-compose 一键部署

```
docker-compose up -d
```

3. 可以根据需要修改该文件夹下的配置文件及镜像版本

<!-- LICENSE -->

## 版权声明

admin-go 基于 MIT 协议， 详情请参考 [license](LICENSE)。
