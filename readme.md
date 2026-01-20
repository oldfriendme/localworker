# localworker

在 **自有本地 VPS** 上运行 **Worker**   

<br>

## 项目简介

**localworker** 是一个在 **自有 VPS / 本地服务器** 上运行的 **CF Worker**。 

---

<br>

## WebIDE 环境

![webide](img/webide.jpg?raw=true)

<br>

## 使用教程

cf官方的worker示例大多数都能正常运行，参考：https://developers.cloudflare.com/workers/examples/return-html/

也可以参考example目录的文件。[example目录](example)

<br>

### 注意

没有D1支持。没有R2支持，也没有Worker AI，Workflows支持。基本上CF的集成环境大多数都没有。不过runtime环境基本都有。


<br>

---

### 技术路线

前端**golang**(web IDE)，后端**rust**(worker runtime)，

并非官方这种的 `wrangler` 的依赖**JS dev开发环境**。也没有**nodejs**的**npm依赖**。
 

项目功能：

- 在 **自有 本地 VPS** 上运行 Worker
- 不依赖 Cloudflare 官方平台
- 提供 **Web IDE** + **Worker Runtime**

---

## 架构说明

```
┌──────────────┐
│   Web IDE    │  Golang
│ (开发面板)    │
└──────┬───────┘
       │
┌──────▼───────┐
│ Worker Runtime│  Rust
│ (CF Worker)  │
└──────────────┘
```

### 技术栈

- **前端 / 控制面板**
  - Golang
  - Web IDE（浏览器中直接开发）

- **后端 / Worker Runtime**
  - Rust
  - 高性能、可用于服务

---

<br>

## 与 CF / Wrangler 等的区别

| 特性| 官方平台 | Wrangler | localworker |
|----|----|----|----|
| 运行环境 | Cloudflare 平台| 开发机 (**npm依赖**) | 本地 VPS (**glibc依赖**) |
| 用途 | 生产 | 开发 / 测试 | **玩具** |
| 请求限制 | 有| 无限制 | **无限制** |
| KV 读写限制 | 有 | 无限制 | **无限制** |
| CPU Time | 严格限制 | 无限制 | **无限制** |
| 可控性 | 低 | 可控 | **可控** |

---

## 核心特性

✅ 本地部署  
✅ 不限量请求  
✅ 不限量 KV 读写  
✅ 不限量 CPU Time  
✅ Web IDE 在线开发  
✅ 支持直接运行与部署  

---

## 配置方法


### 下载runtime

去 https://github.com/cloudflare/workerd/releases 把rust版本的worker runtime环境下载下来。这个cf官方开源的worker裸核。

放到/var/worker/bin/workered，配置为appfile

### 测试
```bash
/var/worker/bin/workered --version
```

linux下的**glibc 版本2.35+** 

如果输出正常版本，说明可运行。如果报错环境依赖问题，可以使用docker/创建Debian 12+/Ubuntu 22.04+环境。在有依赖的环境中运行。

如果无法使用docker，可使用chroot/proot，创建虚拟环境

<br>

创建配置文件 `config.json`：

### https配置
```json
{
	"listen": "127.0.0.1:443",
	"enable_ssl": true,
	"ssl-config": {
		"crt": "/mnt/server.crt",
		"key": "/mnt/server.key"
	},
	"workdir": "/var/www",
	"appfile": "/var/worker/bin/workered",
	"user": "admin",
	"pass": "123456"
}
```

---

<br>

### http

```json
{
	"listen": "127.0.0.1:8080",
	"enable_ssl": false,
	"workdir": "/var/www",
	"appfile": "/var/worker/bin/workered",
	"user": "admin",
	"pass": "123456"
}
```

注意："workdir": "/var/www",指定的workdir目录，当前用户**必须可读可写**

windows环境的目录，需要加双斜杠 `"\\"`，如`"appfile": "D:\\worker\\workerd.exe"`,

---

<br>

#### 配置说明

| 字段 | 说明 |
|----|----|
| workdir | Worker 项目工作目录 |
| appfile | Worker Runtime 可执行文件 |
| user | 管理面板用户名 |
| pass | 管理面板密码 |

---

### 2.启动程序

```bash
localworker config.json
```

启动后服务将监听你配置的地址与端口。

---

### 3️.打开开发面板

在浏览器中访问：

```
https://你的地址/dev_page
```

---

### 4️.登录并开始开发

- 输入你在 `config.json` 中设置的用户名和密码
- 进入 Web IDE
- 即可：
  - 开发/运行 Worker
  - 部署到本地 Runtime

---

## 使用场景

- 自建 CF Worker 平替
- 高请求任务
- 无请求限制的 API 服务
- 私有 KV 存储
- 私有云部署

---

## 安全建议

- 请务必修改默认用户名和密码
- 生产环境建议启用 HTTPS
- 可通过防火墙限制 `/dev_page` 访问 IP

---

<br>

## 注意，KV实现与cf worker官方的KV api不一样。参考example目录的文件.