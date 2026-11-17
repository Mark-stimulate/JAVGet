# NASSAV（个人修改版）

基于 [Satoing/NASSAV](https://github.com/Satoing/NASSAV) 修改的 NAS 影音下载、元数据整理和网页浏览工具。保留原项目 MIT 版权声明，详见 [LICENSE](LICENSE)。本仓库是个人修改版，与上游项目没有从属或官方维护关系。

## 功能与状态

Python 下载器按配置选择来源并用 FFmpeg 整理文件；可抓取元数据并生成 NFO。Go 服务提供目录浏览、文件访问和添加任务接口；Vue 前端用于浏览和提交任务。下载来源可能随网站变更而失效。本版本仍需在实际 NAS 环境验证，不应视为生产安全方案。

## 初始配置

安装 Python 3.11+、FFmpeg 和 m3u8-Downloader-Go。将后者放在 `tools/m3u8-Downloader-Go`（Windows 为 `.exe`）；Windows 也可以将 `ffmpeg.exe` 放在 `tools/` 或加入 PATH。第三方工具按各自许可证获取。

将 `cfg/configs.json.example` 复制为 `cfg/configs.json`，设置媒体目录、元数据来源、下载源和代理，然后在项目根目录运行：

```sh
pip install -r requirements.txt
python main.py ABC-123
python metadata.py
```

目录、数据库、队列与日志属于本地数据，不应提交仓库。下载器和来源的有效性、合法使用权限需由部署者自行确认。

## 浏览服务（实验性）

在项目根目录先完成 Python 配置，创建媒体目录并初始化数据库。设置环境变量：`NASSAV_API_KEY` 为自行生成的随机十六进制密钥；`NASSAV_MEDIA_PATH` 为媒体目录绝对路径；`NASSAV_FRONTEND_ORIGIN` 默认为 `http://localhost:5173`；`NASSAV_BIND` 默认为 `127.0.0.1`。进入 `backend`，运行 `go build -o main` 后启动 `./main`。

前端进入 `frontend`，运行 `npm ci` 和 `npm run dev`。默认连接 `http://localhost:31471`，可通过构建环境变量 `VITE_API_BASE` 修改。打开网页时输入同一密钥。前后端应使用同一主机名（例如均使用 localhost），以便图片请求携带浏览器会话 cookie。

密钥和 cookie 是共享口令，HTTP 下不加密；默认只监听本机。不要直接开放到公网。远程访问应使用经过验证的加密连接和访问控制。本修改尚未进行完整安全审计。

## Docker

Dockerfile 构建下载器镜像而非完整网页服务，并从上游构建 m3u8 下载工具。构建时不包含实际配置和本地数据；运行时需挂载配置、数据库和媒体目录。请根据 NAS 路径调整 `cfg/configs.json` 与挂载路径，并先在测试环境验证。

## 已知限制

- 下载任务的 `work` 标记仍是简单文件标记，不是可靠的进程锁，异常中断后可能需要人工复位，不适合高并发。
- 定时任务和来源解析仍依赖原项目逻辑；网页前后端没有端到端测试。
- 网页登录使用简化的共享密钥机制，文件 cookie 仅供同站点使用，尚不支持完善的用户会话、退出与权限管理。
- 现有示例配置不代表来源当前可用；截图和第三方编译产物不包含在仓库中。
