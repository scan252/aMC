# aMC macOS 客户端架构设计

## 1. 目标

aMC 的最终形态是可直接安装和运行的 macOS 原生应用，用于获取、保存、分析和展示《鸣潮》抽卡记录。

设计优先级依次为：

1. **零用户依赖**：用户不需要安装 Python、Node.js、Homebrew 或其他运行环境。
2. **轻量**：安装包小、启动快、内存占用低。
3. **易于构建 GUI**：界面符合 macOS 使用习惯，并方便持续增加分析和可视化功能。
4. **本地优先**：用户数据默认只保存在本机。
5. **可维护**：抓取、存储、分析和界面彼此解耦，游戏接口变化时能够局部适配。

## 2. 技术选型

正式客户端采用以下原生技术：

| 领域 | 选择 | 原因 |
|---|---|---|
| 语言 | Swift | macOS 原生支持，无需附带运行时 |
| GUI | SwiftUI | 原生控件、主题和无障碍支持，适合快速迭代 |
| 并发 | Swift Concurrency | 使用 `async/await` 和 Actor 管理任务与共享状态 |
| 网络 | URLSession | 系统内置，不引入第三方 HTTP 库 |
| 数据解析 | Codable | 系统内置 JSON 编解码 |
| 数据库 | SQLite | macOS 自带，体积小且迁移行为可控 |
| 图表 | Swift Charts | 原生图表，无需 WebView 或前端运行时 |
| 凭证 | Keychain | 保存短期敏感凭证 |
| 设置 | UserDefaults | 保存非敏感应用设置 |

建议最低支持 **macOS 13 Ventura**，以直接使用 Swift Charts 和较完整的 SwiftUI API。

### 不采用的方案

- **Electron**：需要附带 Chromium 和 Node.js，安装包与内存占用不符合轻量目标。
- **Python 打包**：需要附带 Python 解释器，体积较大，GUI 和签名分发也更复杂。
- **Flutter**：需要附带 Flutter Engine，原生集成收益有限。
- **Tauri**：比 Electron 轻，但仍引入 Web 前端和 WebView 通信层。
- **Rust + SwiftUI**：性能足够但跨语言边界增加维护成本；当前任务主要受网络和文件 I/O 限制，纯 Swift 已足够。

如果未来明确需要 Windows 或 Linux 客户端，再评估将核心逻辑迁移为 Rust 库；当前不为尚未确定的跨平台需求增加复杂度。

## 3. 总体结构

```text
aMC.app
├── App
│   ├── App lifecycle
│   ├── Navigation
│   └── Dependency container
├── Features
│   ├── Dashboard
│   ├── History
│   ├── Analysis
│   ├── Accounts
│   ├── ImportExport
│   └── Settings
├── Core
│   ├── LogDiscovery
│   ├── LogParser
│   ├── GachaAPI
│   ├── GachaRepository
│   ├── AnalysisEngine
│   └── Models
├── Infrastructure
│   ├── SQLiteStore
│   ├── KeychainStore
│   ├── FileAccess
│   ├── AssetCache
│   └── Exporters
└── Resources
    ├── Localizable
    ├── PoolRules
    └── AppAssets
```

依赖方向保持单向：

```text
SwiftUI Features
       ↓
Core protocols and use cases
       ↓
Infrastructure implementations
```

界面不直接访问文件、网络或 SQLite。所有读写通过 Core 定义的协议完成，便于测试和替换实现。

## 4. 核心模块

### 4.1 LogDiscovery

职责：

- 检测国服与国际服日志路径
- 判断日志是否存在及最后更新时间
- 处理用户手动选择的日志文件
- 保存可恢复的文件访问授权

已知默认路径：

```text
~/Library/Containers/com.kurogame.mingchao/Data/Library/Logs/Client/Client.log
~/Library/Containers/com.kurogame.wutheringwaves.global/Data/Library/Logs/Client/Client.log
```

路径必须集中配置，不能散落在 View 或业务代码中。

### 4.2 LogParser

职责：

- 识别明文或加密日志
- 解密 XOR 格式的 `Client.log`
- 提取最近一次有效唤取 URL
- 解析国服、国际服、UID、服务器和请求凭证

解析器只接收 `Data` 或字符串，不自行读取文件，以便使用脱敏日志样本做单元测试。

### 4.3 GachaAPI

职责：

- 调用官方唤取记录接口
- 获取全部受支持卡池
- 设置超时、限速和有限次数重试
- 将服务端响应转换成内部模型
- 将“凭证过期”“网络错误”“接口变化”区分为明确错误类型

API 层不负责数据库合并，也不向 SwiftUI 直接返回展示文本。

### 4.4 GachaRepository

职责：

- 读取账号与抽卡记录
- 将新响应与本地历史合并
- 使用事务保证一次同步原子完成
- 创建备份与执行数据库迁移
- 为界面和分析引擎提供查询接口

同步流程：

```text
读取日志
  → 解析凭证
  → 请求各卡池
  → 校验响应
  → SQLite 事务合并
  → 更新同步元数据
  → 刷新界面
```

任何卡池失败时应保留既有数据。是否提交成功卡池的数据，需要在产品层明确展示为“部分同步成功”。

### 4.5 AnalysisEngine

分析引擎实现为不依赖 SwiftUI 和数据库的纯逻辑模块：

```swift
func analyze(
    records: [GachaRecord],
    rules: PoolRules
) -> AnalysisResult
```

计划输出：

- 当前 5★ / 4★ 保底进度
- 距离硬保底剩余抽数
- 每次高稀有度出货所用抽数
- 平均、最早、最晚出货抽数
- UP 命中与歪卡统计
- 各稀有度数量和比例
- 按日期、版本和卡池的趋势数据

保底规则、卡池分组和 UP 资源列表放入独立的 `PoolRules` 配置，不在页面代码中硬编码。规则必须带版本，以便游戏机制变化后仍可解释旧数据。

## 5. 数据设计

SQLite 数据库建议存放于：

```text
~/Library/Application Support/aMC/amc.sqlite3
```

主要数据表：

```text
accounts
- id
- player_id
- server_id
- server_area
- display_name
- created_at
- last_sync_at

gacha_records
- id
- account_id
- pool_type
- resource_id
- resource_name
- resource_type
- quality_level
- count
- pulled_at
- source_fingerprint
- occurrence_index

sync_history
- id
- account_id
- started_at
- finished_at
- status
- error_summary

schema_migrations
- version
- applied_at
```

### 去重原则

不能只用“时间 + 名称 + 稀有度”作为唯一键，因为十连中可能出现同一秒、同名、同稀有度的多条记录。

同步时应按完整字段生成 `source_fingerprint`，并为同一响应中的重复项分配稳定的 `occurrence_index`。合并采用多重集合计数，而不是简单的 Set 去重，以避免丢失真实抽卡记录。

### 数据迁移

- 每次数据库结构变化增加 migration 版本。
- 升级前创建数据库备份。
- migration 失败时回滚，不启动半迁移状态的数据库。
- 当前 Python MVP 的 JSON 通过一次性导入器迁移到 SQLite。

JSON、CSV 和 UIGF 仅作为交换格式，不作为正式客户端的主存储。

## 6. macOS 文件权限

读取另一个应用容器中的日志是正式客户端的主要平台风险。

首选分发方式为：

- GitHub Releases 或项目官网分发
- Developer ID 签名
- Apple Notarization 公证
- 不启用 App Sandbox

应用启动后先尝试自动发现日志。如果系统隐私策略阻止访问，则使用 `NSOpenPanel` 让用户选择一次 `Client.log`，并保存文件书签供后续访问。

这意味着“零依赖、双击可用”可以实现，但某些 macOS 版本上首次使用仍可能需要用户确认文件访问权限。应用必须提供清楚的引导，不能要求用户关闭 SIP 或执行高权限脚本。

如果未来计划提交 Mac App Store，则需要重新评估 App Sandbox 限制，不能假设商店版可以直接读取游戏容器。

## 7. 凭证与隐私

- 唤取 URL 和 `record_id` 不写入普通日志。
- 临时凭证如需缓存，存入 Keychain，并设置过期时间。
- 数据库只保存分析所需的账号信息与抽卡记录。
- 崩溃报告和遥测默认不启用。
- 导出、云同步或诊断日志必须由用户明确触发。
- UI 中展示 URL 时默认脱敏。

## 8. GUI 设计

建议采用 `NavigationSplitView`：

```text
侧边栏
├── 总览
├── 抽卡记录
├── 分析
├── 账号
└── 设置
```

第一版原生 GUI 包含：

- 日志与凭证状态
- 一键同步按钮和分卡池进度
- 同步失败的可操作提示
- 卡池筛选与抽卡记录列表
- 总抽数及 5★、4★、3★ 基础统计
- 数据导入、导出和备份

后续使用 Swift Charts 增加：

- 保底进度
- 五星出货抽数分布
- UP 命中率
- 月度与版本趋势
- 卡池对比

业务计算不得写在 View 中。View 只订阅 ViewModel 暴露的状态和用户操作。

## 9. 素材与缓存

为控制安装包大小：

- 不在应用内打包全部角色和武器大图。
- 使用系统字体，不附带大型中文字体。
- 仅内置应用图标、占位图和必要的小型配置。
- 角色与武器图片按需下载到 `~/Library/Caches/aMC/`。
- 缓存设置大小上限，并允许用户一键清理。
- 无网络或素材下载失败时，文本记录和分析仍应正常使用。

远程素材必须使用固定 HTTPS 来源，并提供资源版本和完整性校验。

## 10. 打包与发布

发布流水线：

```text
构建与测试
  → Xcode Archive
  → Developer ID 签名
  → Apple Notarization
  → Staple 公证票据
  → 生成 DMG
  → GitHub Release
```

用户安装流程：

1. 下载 DMG。
2. 将 aMC 拖入 Applications。
3. 双击运行。

正式版本不要求用户安装 Xcode、Python、Node.js、Homebrew 或命令行工具。

初期使用 GitHub Releases 提供更新，由应用检查新版本并引导用户下载。暂不引入 Sparkle，以维持零第三方依赖；若以后实现应用内自动更新，必须加入更新包签名验证和回滚机制。

## 11. 依赖策略

正式客户端默认只使用 Apple 系统框架。

新增第三方依赖前必须确认：

1. 系统框架无法合理实现该功能。
2. 对安装包和启动性能的影响可接受。
3. 许可证允许分发。
4. 项目仍在维护且安全风险可控。
5. 能够被固定版本并纳入供应链审查。

角色素材等内容资源不视为运行时依赖，但仍需确认版权与分发许可。

## 12. 实施阶段

### 阶段一：Swift 核心

- 建立 Xcode 工程和 Core 模块
- 迁移日志发现、解密和 URL 解析
- 迁移官方 API 客户端
- 建立 SQLite schema、migration 和 Repository
- 使用脱敏样本与模拟 API 响应覆盖测试

Python MVP 在此阶段继续作为行为参考和接口诊断工具。

### 阶段二：基础 GUI

- 完成首页、同步流程和记录列表
- 加入数据导入、导出和错误引导
- 验证国服、国际服以及不同 macOS 版本的文件权限

### 阶段三：分析与可视化

- 实现版本化卡池规则
- 实现保底、UP 和分布分析
- 使用 Swift Charts 构建图表

### 阶段四：产品化

- 多账号管理
- 菜单栏助手与日志变化监听
- 素材缓存
- 签名、公证、DMG 和发布流水线
- 更新检查、备份恢复和国际化

## 13. 测试要求

- LogParser：明文、加密、损坏、多账号和无 URL 日志样本。
- GachaAPI：成功、过期、限流、超时、部分卡池失败和响应结构变化。
- Repository：重复十连、增量合并、旧数据保留、事务回滚和 migration。
- AnalysisEngine：各类保底边界、UP 命中与规则版本。
- UI：首次使用、权限拒绝、空数据、同步中断和恢复。
- 发布物：Intel/Apple Silicon 支持策略、签名、公证和全新 Mac 安装验证。

所有测试样本必须移除真实 UID、服务器标识和唤取凭证。
