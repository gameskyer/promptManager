# PromptMaster 项目结构与功能说明

> AI绘画提示词管理软件 - 基于 Wails + Vue3 的桌面应用

## 📌 项目概述

PromptMaster 是一款专为 AI 绘画（Stable Diffusion/Midjourney 等）设计的提示词管理工具，提供原子词管理、预设模板、版本控制、AI 辅助等核心功能。

---

## 🏗️ 技术架构

| 层级 | 技术 | 版本 |
|------|------|------|
| 桌面框架 | [Wails v2](https://wails.io) | v2.x |
| 后端语言 | Go | 1.21+ |
| 前端框架 | Vue | 3.4+ |
| 构建工具 | Vite | 5.4+ |
| 样式框架 | TailwindCSS | 3.4+ |
| 状态管理 | Pinia | 2.1+ |
| 数据库 | SQLite | - |
| ORM | Gorm | - |
| 图标库 | @heroicons/vue | 2.1+ |

---

## 📁 目录结构

```
promptManager/
│
├── 📄 main.go                    # 应用入口文件
├── 📄 wails.json                 # Wails 配置文件
├── 📄 go.mod / go.sum            # Go 依赖管理
├── 📄 Makefile                   # 构建脚本
├── 📄 run.bat                    # Windows 运行脚本
│
├── 📁 backend/                   # Go 后端代码
│   │
│   ├── 📁 config/               # 应用配置
│   │   └── 📄 config.go         # 全局配置常量（数据库路径、版本号等）
│   │
│   ├── 📁 handlers/             # API 处理器（暴露给前端）
│   │   ├── 📄 atom_handler.go       # 原子词 CRUD 接口
│   │   ├── 📄 preset_handler.go     # 预设管理接口
│   │   ├── 📄 category_handler.go   # 分类管理接口
│   │   ├── 📄 version_handler.go    # 版本控制接口
│   │   ├── 📄 search_handler.go     # 搜索接口
│   │   ├── 📄 ai_handler.go         # AI 功能接口
│   │   ├── 📄 image_handler.go      # 图片管理接口
│   │   ├── 📄 backup_handler.go     # 数据备份接口
│   │   ├── 📄 batch_handler.go      # 批量操作接口
│   │   └── 📄 seeder_handler.go     # 数据导入接口
│   │
│   ├── 📁 models/               # 数据模型与数据库
│   │   ├── 📄 models.go         # 所有数据模型定义
│   │   └── 📄 db.go             # 数据库连接与初始化
│   │
│   ├── 📁 services/             # 业务逻辑服务层
│   │   ├── 📄 atom_service.go       # 原子词业务逻辑
│   │   ├── 📄 preset_service.go     # 预设业务逻辑
│   │   ├── 📄 category_service.go   # 分类业务逻辑
│   │   ├── 📄 version_service.go    # 版本管理逻辑
│   │   ├── 📄 search_service.go     # 搜索逻辑
│   │   ├── 📄 ai_service.go         # AI 服务（拆解/优化/翻译/分析）
│   │   ├── 📄 image_service.go      # 图片服务
│   │   ├── 📄 backup_service.go     # 备份服务
│   │   └── 📄 batch_service.go      # 批处理服务
│   │
│   ├── 📁 utils/                # 工具类
│   │   └── 📄 seeder.go         # 数据种子/初始数据导入
│   │
│   └── 📁 logger/               # 日志系统
│       └── 📄 logger.go         # AI 请求日志记录器
│
├── 📁 frontend/                  # Vue 前端代码
│   │
│   ├── 📄 index.html            # HTML 模板
│   ├── 📄 package.json          # NPM 依赖
│   ├── 📄 vite.config.js        # Vite 配置
│   ├── 📄 tailwind.config.js    # Tailwind 配置
│   └── 📄 postcss.config.js     # PostCSS 配置
│   │
│   └── 📁 src/                  # 源代码
│       │
│       ├── 📄 main.js           # 前端入口
│       ├── 📄 App.vue           # 根组件
│       ├── 📄 style.css         # 全局样式
│       │
│       ├── 📁 components/       # Vue 组件
│       │   ├── 📄 TopBar.vue            # 顶部工具栏
│       │   ├── 📄 SideMenu.vue          # 左侧分类菜单
│       │   ├── 📄 MainContent.vue       # 主内容区（原子词列表）
│       │   ├── 📄 Workbench.vue         # 右侧工作台
│       │   ├── 📄 Timeline.vue          # 底部时间轴
│       │   ├── 📄 AtomCard.vue          # 原子词卡片
│       │   ├── 📄 AtomDialog.vue        # 原子词编辑弹窗
│       │   ├── 📄 PresetCard.vue        # 预设卡片
│       │   ├── 📄 PresetList.vue        # 预设列表页面
│       │   ├── 📄 PresetDialog.vue      # 预设编辑弹窗
│       │   ├── 📄 PresetDetailModal.vue # 预设详情弹窗
│       │   ├── 📄 VersionDetailModal.vue# 版本详情弹窗
│       │   ├── 📄 CategoryTreeNode.vue  # 分类树组件
│       │   ├── 📄 CategoryManagement.vue# 分类管理页面
│       │   ├── 📄 CategoryDialog.vue    # 分类编辑弹窗
│       │   ├── 📄 AtomManagement.vue    # 原子词管理页面
│       │   ├── 📄 ImageUpload.vue       # 图片上传组件
│       │   ├── 📄 ImageViewer.vue       # 图片查看器
│       │   ├── 📄 AIModal.vue           # AI 功能弹窗
│       │   ├── 📄 SettingsModal.vue     # 设置弹窗
│       │   ├── 📄 BackupModal.vue       # 备份弹窗
│       │   ├── 📄 CompareModal.vue      # 版本对比弹窗
│       │   └── 📄 ...
│       │
│       ├── 📁 stores/           # Pinia 状态管理
│       │   ├── 📄 index.js      # Store 入口
│       │   ├── 📄 app.js        # 应用全局状态（主题、时间轴显示等）
│       │   ├── 📄 atom.js       # 原子词状态
│       │   ├── 📄 preset.js     # 预设状态
│       │   ├── 📄 category.js   # 分类状态
│       │   ├── 📄 version.js    # 版本状态
│       │   ├── 📄 ai.js         # AI 配置状态
│       │   ├── 📄 image.js      # 图片状态
│       │   └── 📄 backup.js     # 备份状态
│       │
│       ├── 📁 lib/wailsjs/      # Wails 自动生成的 Go 绑定
│       │   ├── 📁 go/handlers/  # Go Handler 的 JS 绑定
│       │   └── 📁 runtime/      # Wails 运行时
│       │
│       ├── 📁 utils/            # 前端工具函数
│       │   ├── 📄 helpers.js    # 辅助函数
│       │   └── 📄 wails.js      # Wails 调用封装
│       │
│       └── 📁 assets/           # 静态资源
│
├── 📁 data/                    # 应用数据目录
│   └── 📁 backups/             # 备份文件存储
│
├── 📁 images/                  # 图片存储目录（预览图等）
│
├── 📁 logs/                    # 日志文件目录
│   └── 📄 ai_requests.log      # AI 请求日志
│
├── 📁 plugin/                  # 插件目录（预留）
│
├── 📁 build/                   # 构建输出目录
│
└── 📁 changelog/               # 版本更新日志
```

---

## 🧩 核心功能模块

### 1. 原子词管理 (Atom Management)

**功能描述**：管理 AI 绘画提示词的最小单元

| 功能 | 说明 |
|------|------|
| 原子词 CRUD | 创建、读取、更新、删除原子词 |
| 分类归属 | 将原子词归类到不同分类 |
| 同义词管理 | 支持多个同义词，搜索时可关联查找 |
| 使用统计 | 记录使用次数和最后使用时间 |
| 批量导入 | 支持 JSON 格式批量导入 |
| 热度排序 | 按使用频率排序显示 |

**数据字段**：
- `value`: 英文提示词（如 "masterpiece"）
- `label`: 中文标签（如 "杰作"）
- `synonyms`: 同义词列表
- `type`: Positive（正向）/ Negative（负向）
- `category_id`: 所属分类
- `usage_count`: 使用次数

---

### 2. 预设管理 (Preset Management)

**功能描述**：管理完整的提示词组合模板

| 功能 | 说明 |
|------|------|
| 预设 CRUD | 创建、编辑、删除预设 |
| 版本控制 | 每个预设支持多版本快照 |
| 正向/负向分离 | 分别存储正向和负向提示词 |
| 关联原子词 | 记录预设包含的原子词 ID 列表 |
| 预览图管理 | 支持上传多张预览图 |
| 软删除 | 支持删除后恢复 |
| Fork 功能 | 基于现有版本创建新预设 |

**数据字段**：
- `title`: 预设标题
- `category_id`: 所属分类
- `current_version`: 当前版本号
- `is_deleted`: 软删除标记
- 版本快照包含：pos_text, neg_text, params, atom_ids, preview_paths

---

### 3. 分类管理 (Category Management)

**功能描述**：层级化组织原子词和预设

| 功能 | 说明 |
|------|------|
| 多级分类 | 支持父子层级结构 |
| 类型区分 | ATOM 类型（原子词）/ PRESET 类型（预设） |
| 排序控制 | 自定义分类显示顺序 |
| 树形展示 | 侧边栏树形菜单浏览 |

**默认分类**：
- 人物（二级：发型、眼睛、服装、姿势、表情）
- 场景
- 风格
- 质量
- 预设库

---

### 4. 版本控制 (Version Control)

**功能描述**：预设的版本快照与历史管理

| 功能 | 说明 |
|------|------|
| 自动保存 | 保存预设时自动创建新版本 |
| 版本对比 | 显示版本间的差异（+/- 统计） |
| 版本回滚 | 切换到历史版本 |
| 星标版本 | 标记重要版本，清理时保留 |
| 版本清理 | 自动清理旧版本，保留最近 N 个 |

**版本数据结构**：
```json
{
  "pos_text": "正向提示词文本",
  "neg_text": "负向提示词文本",
  "params": { "其他参数": "值" },
  "atom_ids": [1, 2, 3],
  "preview_paths": ["/images/xxx.jpg"]
}
```

---

### 5. AI 辅助功能 (AI Assistant)

**功能描述**：集成 AI 大模型辅助提示词处理

**支持的 AI 提供商**：
| 提供商 | 类型 | 说明 |
|--------|------|------|
| Ollama | 本地 | 本地部署的 LLM 服务 |
| OpenAI | 云端 | GPT-3.5/GPT-4 |
| DeepSeek | 云端 | DeepSeek API |
| Kimi | 云端 | Moonshot API |

**AI 功能**：

| 功能 | 说明 |
|------|------|
| 🔍 拆解 (Explode) | 将长提示词拆解为原子词列表，自动匹配分类 |
| ✨ 优化 (Optimize) | 优化提示词质量，添加质量词，调整顺序 |
| 🌐 翻译 (Translate) | 中英互译，提供关键词提取 |
| 📊 分析 (Analyze) | 分析提示词结构（主体/风格/质量/光照） |

**Fallback 机制**：未配置 AI 时使用基于规则的本地处理

---

### 6. 工作台 (Workbench)

**功能描述**：提示词组合构建与实时预览

| 功能 | 说明 |
|------|------|
| 拖拽添加 | 从原子词列表拖拽到工作台 |
| 排序调整 | 调整提示词顺序 |
| 权重标记 | 支持括号权重语法 (word:1.2) |
| 一键复制 | 复制完整提示词到剪贴板 |
| 正负分离 | 分别显示正向和负向提示词 |

---

### 7. 搜索功能 (Search)

**功能描述**：快速查找原子词和预设

| 功能 | 说明 |
|------|------|
| 全文搜索 | 搜索 value、label、synonyms |
| 分类过滤 | 按分类筛选 |
| 类型过滤 | 正向/负向筛选 |
| 热度排序 | 按使用频率排序 |

---

### 8. 数据备份与恢复 (Backup)

**功能描述**：数据安全与迁移

| 功能 | 说明 |
|------|------|
| 自动备份 | 定期自动备份数据库 |
| 手动备份 | 一键导出完整数据 |
| 数据恢复 | 从备份文件恢复 |
| 导出 JSON | 导出原子词为 JSON 格式 |
| 数据导入 | 从 JSON 导入原子词 |

---

### 9. 图片管理 (Image)

**功能描述**：预览图的上传和展示

| 功能 | 说明 |
|------|------|
| Base64 上传 | 前端 Base64 编码上传 |
| 本地存储 | 图片保存到 ./images 目录 |
| 静态服务 | 通过 /images/ 路径访问 |
| 缩略图 | 自动生成缩略图 |

---

## 🗄️ 数据模型详解

### Category（分类）
```go
type Category struct {
    ID        uint      // 主键
    Name      string    // 分类名称
    ParentID  uint      // 父分类ID（0为顶级）
    Type      string    // ATOM 或 PRESET
    SortOrder int       // 排序顺序
}
```

### Atom（原子词）
```go
type Atom struct {
    ID         uint         // 主键
    Value      string       // 英文值（唯一）
    Label      string       // 中文标签
    Synonyms   StringSlice  // 同义词列表（JSON）
    Type       string       // Positive / Negative
    CategoryID uint         // 分类ID
    UsageCount int          // 使用次数
    LastUsedAt *time.Time   // 最后使用时间
}
```

### Preset（预设）
```go
type Preset struct {
    ID             uint      // 主键
    Title          string    // 标题
    CategoryID     uint      // 分类ID
    CurrentVersion int       // 当前版本号
    IsDeleted      bool      // 软删除标记
    Versions       []PresetVersion // 关联版本
}
```

### PresetVersion（预设版本）
```go
type PresetVersion struct {
    ID            uint      // 主键
    PresetID      uint      // 所属预设
    VersionNum    int       // 版本号（V1=1）
    Snapshot      JSON      // 完整快照（JSON）
    ThumbnailPath string    // 缩略图路径
    IsStarred     bool      // 是否星标
    DiffStats     string    // 差异统计（如 "+2/-1"）
}
```

---

## 🚀 开发指南

### 环境要求
- Go 1.21+
- Node.js 18+
- Wails CLI

### 安装依赖
```bash
# 安装 Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装前端依赖
cd frontend && npm install
```

### 常用命令

```bash
# 开发模式（热重载）
wails dev

# 构建生产版本
wails build

# 构建 Windows 安装包
wails build -nsis

# 仅前端开发
cd frontend && npm run dev
```

### 添加新功能流程

1. **定义模型** (`backend/models/models.go`)
2. **实现服务** (`backend/services/xxx_service.go`)
3. **创建 Handler** (`backend/handlers/xxx_handler.go`)
4. **注册 Handler** (`main.go` 中的 `Bind`）
5. **生成绑定** (`wails dev` 自动生成)
6. **前端 Store** (`frontend/src/stores/xxx.js`)
7. **前端组件** (`frontend/src/components/`)

---

## 📦 构建输出

构建后的应用位于 `build/bin/` 目录：

| 平台 | 输出文件 |
|------|----------|
| Windows | `PromptMaster.exe` |
| macOS | `PromptMaster.app` |
| Linux | `PromptMaster` |

---

## 🔧 配置文件

### wails.json
Wails 应用配置，定义应用名称、前端命令、输出文件名等。

### 应用数据
- **数据库**: `./promptmaster.db` (SQLite)
- **配置**: `./data/`
- **图片**: `./images/`
- **日志**: `./logs/`

---

## 📝 日志系统

| 日志类型 | 位置 | 说明 |
|----------|------|------|
| AI 请求日志 | `./logs/ai_requests.log` | 记录所有 AI API 调用 |
| 应用日志 | 控制台 | Wails 运行时日志 |

---

## 🤝 开发建议

1. **后端开发**：在 `services` 层实现业务逻辑，`handlers` 只负责参数校验和调用转发
2. **前端开发**：使用 Pinia 管理状态，组件保持纯展示逻辑
3. **数据库**：使用 GORM 的 AutoMigrate 自动维护表结构
4. **错误处理**：后端返回统一格式 `{success, data, error}`
5. **图片存储**：使用相对路径 `./images/`，通过 Wails 中间件提供静态服务

---

*最后更新：2025-04-09*
