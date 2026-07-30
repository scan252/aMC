# aMC

鸣潮 **Mac 原生客户端** 抽卡记录助手（a Mac Client）。

从 Mac 版鸣潮的沙盒日志中自动提取唤取记录 URL，调用官方 API 抓取全部卡池历史，并保存到本地。

## 调研结论

GitHub 上同类工具大多面向 **Windows**，核心流程一致，但 Mac 版在日志路径与沙盒机制上有明显差异：

| 项目 | 平台 | 获取凭证方式 | 备注 |
|------|------|-------------|------|
| [ningnao/wuthering-waves-gacha-record](https://github.com/ningnao/wuthering-waves-gacha-record) | Windows | 读取 `Client.log` + XOR 解密 | Rust + egui，功能完整 |
| [BJY-STUDIO/wuwa-gacha-analyzer](https://github.com/BJY-STUDIO/wuwa-gacha-analyzer) | Windows | 手动粘贴 JSON / 抓包 | Python，分析能力强 |
| [dyar7474/WuWa_local_tracker](https://github.com/dyar7474/WuWa_local_tracker) | Windows | 自动扫描日志 | PowerShell，纯本地 |
| [GoneTone/wuthering-waves-convene-gacha-analyzer](https://github.com/GoneTone/wuthering-waves-convene-gacha-analyzer) | Windows | HTTPS 拦截 | 需管理员权限 |
| [Luzefiru macOS gist](https://gist.github.com/Luzefiru/b7a59e992a4db9c44cacd9178c5bb673) | **macOS** | grep `Client.log` | 仅提取 URL，无抓取 |

### Mac 版与 Windows 版的关键区别

1. **日志路径不同**
   - Windows: `{游戏目录}/Client/Saved/Logs/Client.log`
   - Mac 国服: `~/Library/Containers/com.kurogame.mingchao/Data/Library/Logs/Client/Client.log`
   - Mac 国际服: `~/Library/Containers/com.kurogame.wutheringwaves.global/Data/Library/Logs/Client/Client.log`

2. **沙盒机制**：Mac 原生客户端运行在 App Sandbox 内，日志不在游戏安装目录，而在 `Library/Containers` 下。

3. **日志可能加密**：新版客户端的 `Client.log` 可能经过 XOR 加密（魔数 `\xa5\xef\xa5`），需要解密后才能搜索 URL。

4. **API 相同**：抓取数据仍调用 Kuro 官方接口 `gmserver-api.aki-game2.com/net/gacha/record/query`，与 Windows 版一致。

## 功能

- 自动发现 Mac 鸣潮日志（国服 / 国际服）
- 支持加密日志解密
- 从日志提取唤取记录 URL
- 抓取全部 13 种卡池类型
- 本地 JSON 存储，支持增量合并（保留 API 已过期的历史记录）
- 命令行工具 `amc`

## 环境要求

- macOS（推荐 12+）
- Python 3.10+
- 已安装 Mac 原生鸣潮客户端

## 安装

```bash
# 使用 pip
pip install -e .

# 或使用 uv
uv pip install -e .
```

## 使用方式

### 1. 在游戏中生成日志

1. 启动鸣潮 Mac 版
2. 进入 **唤取 → 唤取记录** 页面
3. 等待页面加载完成（此时 URL 会写入 `Client.log`）

> 唤取 URL 约 **1 小时** 内有效，过期后需重新打开游戏内唤取记录页面。

### 2. 检查日志状态

```bash
amc status
```

### 3. 提取唤取 URL（可选）

```bash
amc url
# URL 会自动复制到剪贴板
```

### 4. 抓取全部抽卡记录

```bash
amc fetch
```

数据默认保存在 `~/.amc/data/{UID}/gacha_data.json`。

### 5. 导出数据

```bash
amc export <UID> -o wuwa_pulls.json
```

## 命令参考

| 命令 | 说明 |
|------|------|
| `amc status` | 检查 Mac 日志文件是否存在 |
| `amc url` | 从日志提取唤取 URL |
| `amc fetch` | 抓取全部卡池记录并保存 |
| `amc export <UID>` | 导出已保存的数据 |

### 常用参数

- `--log / -l`：指定自定义 `Client.log` 路径
- `--player / -p`：指定玩家 UID（多账号时）
- `--url / -u`：直接提供唤取 URL，跳过日志读取
- `--data-dir`：自定义数据存储目录

## 项目结构

```
amc/
├── cli.py           # 命令行入口
├── log_finder.py    # Mac 日志路径发现
├── log_parser.py    # 日志解密与 URL 提取
├── api_client.py    # 官方 API 调用
├── models.py        # 数据模型
└── storage.py       # 本地存储与合并
```

## 开发

```bash
pip install -e ".[dev]"
pytest
```

## 路线图

- [x] Mac 日志自动发现与 URL 提取
- [x] 全卡池 API 抓取
- [x] 本地数据存储与增量合并
- [ ] SwiftUI 原生图形界面
- [ ] 保底统计与可视化分析
- [ ] 多账号管理
- [ ] 菜单栏常驻助手

## 免责声明

本工具仅读取本地游戏日志并调用官方公开 API，不会收集或上传任何数据到第三方服务器。请遵守游戏相关服务条款。

## License

MIT
