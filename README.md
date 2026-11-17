# JAVGet

JAVGet 是基于 [Satoing/NASSAV](https://github.com/Satoing/NASSAV) 修改的影片下载、刮削与整理工具，可在 Windows 和 Linux 使用，另有实验性的网页浏览与任务入口。项目保留上游 MIT 版权声明，详见 [LICENSE](LICENSE)。

下面先介绍命令行下载；网页服务是可选组件。下载源可能变化或失效，请在有权使用的范围内使用。

## 准备工作

- Git、Python 3.11 或更新版和 FFmpeg。只使用命令行时不需要安装 Go 或 Node.js。
- 从 [m3u8-Downloader-Go](https://github.com/Greyh4t/m3u8-Downloader-Go) 按操作系统和 CPU 架构获取可执行文件：Windows 文件放在 `tools/m3u8-Downloader-Go.exe`，Linux 文件放在 `tools/m3u8-Downloader-Go`。
- 仓库没有打包上述第三方程序，请从可信来源获取并遵守其许可。
- 网页功能还需要 Go 1.22.6+、Node.js 和 npm。Go 后端依赖 SQLite 的 CGO 驱动，编译机器须有 C 编译器（`gcc`）且启用 CGO；Windows 可使用提供 gcc 的工具链，并将其加入 PATH。

## Windows 部署与命令行使用

以下在 PowerShell 中操作；将仓库下载到你希望的位置后，进入项目根目录：

```powershell
git clone https://github.com/Mark-stimulate/JAVGet.git
cd JAVGet
py -3 -m venv .venv
.\.venv\Scripts\Activate.ps1
python -m pip install -r requirements.txt
Copy-Item cfg\configs.json.example cfg\configs.json
```

仓库为私有时，克隆需要已获授权的 GitHub 账号。若 PowerShell 阻止激活脚本，可按 Windows 的本地策略处理；也可以使用虚拟环境里的 Python 执行命令。

放好 `tools\m3u8-Downloader-Go.exe`。FFmpeg 可放在 `tools\ffmpeg.exe`，或安装后确保 `ffmpeg` 命令在 PATH 中可用。编辑 `cfg\configs.json`：例如 `"SavePath": "D:/JAVGet/media"`，Windows 的 JSON 路径可以使用正斜杠。按需填写 `Proxy`，不使用代理则填空字符串。

在项目根目录、已激活虚拟环境的 PowerShell 中运行：

```powershell
python main.py ABC-123
```

将 `ABC-123` 替换为目标编号；如需忽略已有下载记录，可运行 `python main.py ABC-123 -f`。下载后的默认目录为项目下的 `media`，也可以在配置中修改。

## Linux 部署与命令行使用

下面使用 Bash，假定已安装 Git、Python、venv 和 FFmpeg：

```bash
git clone https://github.com/Mark-stimulate/JAVGet.git
cd JAVGet
python3 -m venv .venv
source .venv/bin/activate
python -m pip install -r requirements.txt
cp cfg/configs.json.example cfg/configs.json
```

把对应架构的 m3u8-Downloader-Go 可执行文件放到 `tools/m3u8-Downloader-Go`，并赋予执行权限：

```bash
chmod +x tools/m3u8-Downloader-Go
```

编辑 `cfg/configs.json`，将 `SavePath` 设置为用户有写入权限的目录，如 `/home/用户名/Videos/JAVGet`；按需填写 `Proxy`。在项目根目录运行：

```bash
python main.py ABC-123
# 如需忽略已有记录：
python main.py ABC-123 -f
```

## 配置与整理

配置中的 `LogPath`、`SavePath`、`DBPath` 和 `QueuePath` 可以使用绝对路径；相对路径以运行命令时的工作目录为准，因此建议始终在项目根目录运行 Python 命令。`Downloader` 中的 `weight` 越高越先尝试，设为 0 表示禁用；`ScraperDomain` 是元数据来源列表。不要把实际配置、日志、数据库或队列提交到公开仓库。

对已有影片生成元数据时，将文件按 `SavePath/ABC-123/ABC-123.mp4` 的形式整理后，在项目根目录运行：

```bash
python metadata.py
```

Linux 可将编号逐行写入 `db/download_queue.txt`，在项目根目录执行 `bash cron_task.sh` 处理队列末尾的一条；该脚本依赖 Bash、sed 和 `python3`，Windows 不直接支持它。现有队列机制适用于单任务场景，异常中断后 `work` 标记可能需要人工复位。

## 可选：网页浏览与添加任务

网页功能仍为实验性，默认仅供本机使用。先完成上面的命令行部署。将 `cfg/configs.json` 的 `DBPath` 保持为默认 `./db/downloaded.db`（当前 Go 后端使用此固定路径），并确保媒体目录已存在。若尚未执行过下载，可在项目根目录运行以下命令初始化数据库，避免首次添加任务时报错：

```bash
python -c "import os; from src import data; from src.comm import downloaded_path, save_path; os.makedirs(os.path.dirname(downloaded_path) or '.', exist_ok=True); os.makedirs(save_path, exist_ok=True); data.initialize_db(downloaded_path, 'MissAV')"
```

启动后端的终端也应先激活同一个 Python 虚拟环境，使后台启动的下载任务能找到依赖。媒体目录必须与配置中的 `SavePath` 指向同一位置。

### Windows：启动后端

从项目根目录、已激活虚拟环境的 PowerShell 中运行。将示例媒体路径替换成你自己的**绝对路径**：

```powershell
$env:CGO_ENABLED = '1'
$env:JAVGET_MEDIA_PATH = 'D:\JAVGet\media'
$env:JAVGET_API_KEY = python -c "import secrets; print(secrets.token_hex(32))"
Write-Host "网页访问密钥：$env:JAVGET_API_KEY"
cd backend
go build -o javget-server.exe .
.\javget-server.exe
```

`go build` 找不到 gcc 时，先安装带 gcc 的 Windows C 编译工具链并加入 PATH。该命令生成的密钥只存在于当前终端进程；关闭终端后重新启动会生成新密钥。

### Linux：启动后端

从项目根目录、已激活虚拟环境的终端中运行。将示例媒体路径替换成自己的**绝对路径**：

```bash
export CGO_ENABLED=1
export JAVGET_MEDIA_PATH="$HOME/Videos/JAVGet"
export JAVGET_API_KEY="$(python -c 'import secrets; print(secrets.token_hex(32))')"
printf '网页访问密钥：%s\n' "$JAVGET_API_KEY"
cd backend
go build -o javget-server .
./javget-server
```

在后端启动后，保持这个终端运行。默认监听 `127.0.0.1:31471`，`JAVGET_BIND` 可修改监听地址。服务采用共享密钥和 HTTP，尚未经过完整安全审计；不要直接暴露公网。

### 启动前端

打开**另一个终端**，进入项目的 `frontend` 目录：

```bash
npm ci
npm run dev
```

浏览器打开终端显示的地址，默认是 `http://localhost:5173`，输入后端终端显示的密钥。页面可查看已有项目并输入编号提交任务；提交成功仅表示请求已接受，下载进度和结果请查看后端及 Python 日志。

前后端默认地址分别为 `http://localhost:5173` 和 `http://localhost:31471`。修改前端地址时，启动后端前还需设置 `JAVGET_FRONTEND_ORIGIN` 与前端页面的完整来源一致；修改后端地址时可在启动前端前设置 `VITE_API_BASE`。两端使用相同主机名，避免浏览器的文件会话 cookie 不匹配。

## Docker

仓库中的 Dockerfile 仅构建命令行下载器，并非完整网页服务。镜像工作目录为 `/JAVGet`，实际配置文件与媒体、数据库目录需要在运行时挂载。使用前请按自己的环境配置并测试挂载。

## 已知限制

- 下载源依赖网站页面结构，运行状态未经持续保证；有些来源或元数据可能需要额外网络配置。
- 网页服务、Windows 桌面部署和 Docker 尚未进行完整端到端验证；默认网页接口只支持固定数据库路径。
- 共享密钥认证较简化；日志和下载目录可能包含隐私数据，请妥善保存。
