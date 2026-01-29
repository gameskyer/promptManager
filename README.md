# PromptMaster - AI绘画提示词管理客户端

一款基于 Go 1.24 和 Wails 开发的桌面级客户端，旨在为 AI 绘画爱好者提供结构化的提示词（Prompt）管理、智能组合、版本控制，以及基于 AI 大模型的自动拆解入库功能。

核心设计理念：**零摩擦创作流** —— 版本控制应隐形化，不打断用户的心流状态。

## 技术栈

| 层级       | 技术选型                          | 说明                        |
| ---------- | --------------------------------- | --------------------------- |
| 后端       | Go 1.24                           | 核心逻辑与本地服务          |
| 前端       | Vue 3 + Tailwind CSS              | 响应式界面                  |
| 本地数据库 | SQLite (GORM)                     | 主数据存储                  |
| 全文检索   | Bleve                             | 支持中英文、拼音模糊搜索    |
| AI 驱动    | OpenAI/DeepSeek API 或本地 Ollama | 语义拆解与优化              |
| 版本控制   | 自研差异算法 + 本地快照           | 无 Git 依赖的轻量化版本管理 |

## 功能特性

### 1. 原子词与近义词管理
- 正/负提示词显式区分
- JSON 存储近义词数组，搜索命中近义词即展示该原子词
- 使用统计：记录各原子词被使用的频次与最近使用时间

### 2. 组合预设与参数
- 手动参数配置：Checkpoint、LoRA 强度、采样参数（Sampler/Steps/CFG Scale）
- 多图预览：每个预设支持绑定最多 10 张本地预览图
- 版本绑定：每个版本快照可独立绑定预览图

### 3. 版本控制系统
- **自动保存**：点击保存按钮（Ctrl+S）时，静默创建新版本
- **无描述原则**：禁止要求用户输入版本变更说明，依赖自动 Diff 和缩略图识别
- **星标系统**：用户可对关键版本点击星标（⭐），星标版本置顶高亮
- **紧凑徽章**：`[+2/-1]` 表示该版本相比上一版增加了 2 个原子词，删除了 1 个
- **回滚**：双击历史版本，基于该状态创建新版本
- **Fork**：右键版本选择"基于此创建新预设"，生成独立分支

### 4. AI 智能拆解
- **文本拆解**：粘贴长段 Prompt，AI 分析拆解为原子词
- **入库逻辑**：
  - 库中已存在：以库为准，仅补充新中文解释至近义词列表
  - 库中不存在：识别为新词，引导用户确认分类后入库
- **版本感知**：拆解修改后的保存操作自动触发版本快照

### 5. 辅助增强功能
- **图像反推 (WD14/Tagger)**：拖拽上传已生成图片，自动反推原子词并建议入库
- **数据迁移**：支持导入 ComfyUI/SD WebUI 的 styles.csv、Civitai 模型页面抓取
- **A/B 测试模式**：选中两个版本，一键生成对比参数表

## 项目结构

```
promptmaster/
├── backend/                    # Go 后端代码
│   ├── config/                # 配置管理
│   ├── handlers/              # Wails API 处理器
│   ├── models/                # 数据库模型 (GORM)
│   ├── services/              # 业务逻辑层
│   └── utils/                 # 工具函数
├── frontend/                   # Vue 3 前端代码
│   ├── src/
│   │   ├── components/        # Vue 组件
│   │   ├── stores/            # Pinia 状态管理
│   │   ├── utils/             # 工具函数
│   │   ├── App.vue            # 根组件
│   │   ├── main.js            # 入口文件
│   │   └── style.css          # 全局样式
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
├── main.go                     # 应用入口
├── wails.json                  # Wails 配置
├── go.mod                      # Go 依赖
└── README.md                   # 项目说明
```

## 数据库设计

### 原子词表 (atoms)
| 字段          | 类型    | 说明                      |
| ------------- | ------- | ------------------------- |
| id            | INTEGER | 主键                      |
| value         | STRING  | 英文原词 (Unique Index)   |
| label         | STRING  | 中文展示标签              |
| synonyms      | TEXT    | JSON 字符串存储近义词数组 |
| type          | STRING  | Positive / Negative       |
| category_id   | INTEGER | 关联分类 ID               |
| usage_count   | INTEGER | 使用频次统计              |

### 分类表 (categories)
| 字段         | 类型    | 说明                              |
| ------------ | ------- | --------------------------------- |
| id           | INTEGER | 主键                              |
| name         | STRING  | 分类名称                          |
| parent_id    | INTEGER | 父类 ID（0 表示一级分类）         |
| type         | STRING  | ATOM (原子词) / PRESET (预设)     |
| sort_order   | INTEGER | 排序权重                          |

### 预设主表 (presets)
| 字段              | 类型      | 说明                       |
| ----------------- | --------- | -------------------------- |
| id                | INTEGER   | 主键                       |
| title             | STRING    | 预设标题                   |
| current_version   | INTEGER   | 当前最新版本号             |
| created_at        | TIMESTAMP | 创建时间                   |
| updated_at        | TIMESTAMP | 最后修改时间               |
| is_deleted        | BOOLEAN   | 软删除标记                 |

### 预设版本表 (preset_versions)
| 字段             | 类型      | 说明                       |
| ---------------- | --------- | -------------------------- |
| id               | INTEGER   | 主键                       |
| preset_id        | INTEGER   | 关联 presets 表            |
| version_num      | INTEGER   | 版本号，自增               |
| snapshot         | TEXT      | JSON 存储完整快照          |
| thumbnail_path   | STRING    | 预览图路径                 |
| is_starred       | BOOLEAN   | 用户星标标记               |
| created_at       | TIMESTAMP | 创建时间                   |
| diff_stats       | STRING    | 紧凑变更统计（如 `+2/-1`） |

## 界面布局

采用"顶部-左侧-中间-右侧抽屉-底部时间线"五区设计：

- **顶部栏**：全局搜索框（支持拼音/中英）、API 配置、快速保存（Ctrl+S）
- **左侧导航栏**：树形多级菜单
  - 一级：大类（提示词库、预设库、版本归档）
  - 二级：详细分类
- **中间内容区**：
  - 细化导航区：基于左侧选中的二级分类，展示子级细化分类
  - 原子词网格：紧凑卡片（含中文、英文、前 3 条近义词）
- **右侧组合工作区**：
  - 可伸缩列表：原子词条（原词、翻译、上移/下移/权重调节/移除）
  - 版本控制条：显示当前预设版本号、保存按钮、历史展开开关
- **底部时间线抽屉**：
  - 默认收起，点击右侧"版本历史"展开
  - 无文字描述设计：仅展示缩略图、版本号、时间戳、星标状态
  - 紧凑 Diff 徽章：如 `[+2/-1]`

## 开发环境设置

### 前提条件
- Go 1.24+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### 安装依赖

```bash
# 克隆项目
git clone <repo-url>
cd promptmaster

# 安装 Go 依赖
go mod tidy

# 安装前端依赖
cd frontend
npm install
cd ..
```

### 开发模式运行

```bash
# 开发模式（带热重载）
wails dev
```

### 构建应用

```bash
# 构建生产版本
wails build

# 构建不同平台
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

## 快捷键

| 快捷键      | 功能           |
| ----------- | -------------- |
| Ctrl+S      | 保存当前预设   |
| Ctrl+F      | 聚焦搜索框     |
| Ctrl+Z      | 撤销操作       |
| Delete      | 删除选中项     |
| Escape      | 关闭弹窗/面板  |

## 贡献指南

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'Add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 创建 Pull Request

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 致谢

- [Wails](https://wails.io/) - 跨平台桌面应用框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Tailwind CSS](https://tailwindcss.com/) - 实用优先的 CSS 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [Bleve](https://blevesearch.com/) - 全文搜索引擎
