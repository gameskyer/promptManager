# PromptMaster - 项目开发文档

> 用于开发环境切换时的快速参考  
> **最后更新**: 2026-01-29

---

## 📋 项目概览

**PromptMaster** - AI绘画提示词管理系统

| 属性 | 值 |
|------|-----|
| 技术栈 | Wails v2.11.0 + Go 1.24 + Vue 3 + Tailwind CSS |
| 数据库 | SQLite (GORM) |
| 窗口尺寸 | 1200x750 (最小 900x600) |
| 开发命令 | `wails dev` |
| 构建命令 | `wails build` |

---

## 📁 项目结构

```
promptManager/
├── main.go                     # Wails 应用入口
├── wails.json                  # Wails 配置文件
├── go.mod / go.sum             # Go 依赖
├── frontend/                   # Vue 3 前端
│   ├── src/
│   │   ├── main.js             # 前端入口
│   │   ├── App.vue             # 根组件
│   │   ├── style.css           # 全局样式（含主题变量）
│   │   ├── components/         # Vue 组件
│   │   │   ├── TopBar.vue          # 顶部搜索栏
│   │   │   ├── SideMenu.vue        # 左侧菜单
│   │   │   ├── MainContent.vue     # 主内容区
│   │   │   ├── Workbench.vue       # 工作台
│   │   │   ├── AtomCard.vue        # 原子提示词卡片
│   │   │   ├── AtomDialog.vue      # 原子词 CRUD 弹窗
│   │   │   ├── CategoryDialog.vue  # 分类 CRUD 弹窗
│   │   │   ├── PresetCard.vue      # 预设卡片
│   │   │   ├── PresetDialog.vue    # 预设 CRUD 弹窗
│   │   │   ├── PresetList.vue      # 预设库列表
│   │   │   ├── PresetDetailModal.vue # 预设详情弹窗
│   │   │   ├── Timeline.vue        # 版本时间线（可拖拽调整高度）
│   │   │   ├── VersionDetailModal.vue # 版本详情弹窗
│   │   │   ├── CompareModal.vue    # 版本对比弹窗
│   │   │   ├── ImageViewer.vue     # 图片查看器
│   │   │   ├── AIModal.vue         # AI 优化弹窗
│   │   │   └── SettingsModal.vue   # 设置弹窗（主题/AI配置）
│   │   ├── stores/             # Pinia 状态管理
│   │   │   ├── index.js            # Store 导出
│   │   │   ├── app.js              # 应用状态（含主题、AI配置）
│   │   │   ├── atom.js             # 原子词状态
│   │   │   ├── category.js         # 分类状态
│   │   │   ├── preset.js           # 预设状态
│   │   │   └── version.js          # 版本历史状态
│   │   ├── utils/              # 工具函数
│   │   │   ├── helpers.js          # 通用辅助函数
│   │   │   └── wails.js            # Wails 运行时封装
│   │   └── lib/wailsjs/        # Wails 自动生成绑定
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
├── backend/                    # Go 后端
│   ├── handlers/               # HTTP/gRPC 处理器
│   │   ├── atom_handler.go
│   │   ├── category_handler.go
│   │   ├── preset_handler.go
│   │   ├── version_handler.go
│   │   ├── search_handler.go
│   │   └── ai_handler.go
│   ├── services/               # 业务逻辑层
│   │   ├── atom_service.go
│   │   ├── category_service.go
│   │   ├── preset_service.go
│   │   ├── version_service.go
│   │   ├── search_service.go
│   │   └── ai_service.go
│   ├── models/                 # 数据模型 & GORM
│   │   ├── db.go               # 数据库连接
│   │   └── models.go           # 实体模型
│   └── config/                 # 配置
│       └── config.go
└── build/                      # 构建输出
```

---

## ✅ 已实现功能

### 核心功能

| 功能 | 状态 | 说明 |
|------|------|------|
| 窗口配置 | ✅ | 1200x750，最小 900x600 |
| 数据库初始化 | ✅ | SQLite + GORM，自动迁移 |

### 原子词管理 (Atoms)

| 功能 | 状态 | 组件 |
|------|------|------|
| 原子词列表 | ✅ | `AtomCard.vue` |
| 创建原子词 | ✅ | `AtomDialog.vue` |
| 编辑原子词 | ✅ | `AtomDialog.vue` |
| 删除原子词 | ✅ | `AtomDialog.vue` |
| 分类筛选 | ✅ | `SideMenu.vue` |
| **实时搜索** | ✅ | `TopBar.vue` - 输入即搜索，支持英文/中文/近义词 |
| 卡片间距优化 | ✅ | 减小 grid 间距，自适应高度 |

**原子词字段：**
- 内容 (content)
- 分类 (category_id)
- 使用次数 (use_count) - 暂不显示
- 关联预览图

### 分类管理 (Categories)

| 功能 | 状态 | 组件 |
|------|------|------|
| 分类列表 | ✅ | `SideMenu.vue` |
| 创建分类 | ✅ | `CategoryDialog.vue` |
| 编辑分类 | ✅ | `CategoryDialog.vue` |
| 删除分类 | ✅ | `CategoryDialog.vue` |

### 预设库 (Presets)

| 功能 | 状态 | 组件 |
|------|------|------|
| 预设列表 (卡片布局) | ✅ | `PresetList.vue` |
| 预设详情弹窗 | ✅ | `PresetDetailModal.vue` |
| 创建预设 | ✅ | `PresetDialog.vue` |
| 编辑预设 | ✅ | `PresetDialog.vue` |
| 删除预设 | ✅ | `PresetDialog.vue` |
| 预设搜索 | ✅ | `TopBar.vue` |

**预设字段：**
- 标题、描述
- 正向提示词、负向提示词
- 封面图、预览图列表
- 模型输入
- LoRA 列表 (名称 + 权重，多行)
- 参数：seed, steps, cfg, sampler, width, height

### 版本历史 (Versions)

| 功能 | 状态 | 组件 |
|------|------|------|
| **版本时间线** | ✅ | `Timeline.vue` - 支持拖拽调整高度（默认400px） |
| **展开版本详情** | ✅ | 点击行展开显示正/负提示词、参数 |
| 版本详情弹窗 | ✅ | `VersionDetailModal.vue` |
| 版本对比 | ✅ | `CompareModal.vue` |
| 收藏版本 | ✅ | ⭐ 功能 |
| 回滚到版本 | ✅ | 基于旧版本创建新版本 |
| 基于版本 Fork | ✅ | 创建新预设 |

**版本数据结构：**
```javascript
{
  id, preset_id, version_num,
  diff_stats: '+1/-2',      // 变更统计
  is_starred: boolean,
  formattedTime: '30分钟前',
  snapshot: {
    pos_text, neg_text,
    model, loras: [{name, weight}],
    params: {seed, steps, cfg, ...}
  },
  previews: [] // 预览图 URL 列表
}
```

### 工作台 (Workbench)

| 功能 | 状态 | 说明 |
|------|------|------|
| 原子词拖拽到工作台 | ✅ | 拖放交互 |
| 提示词编辑 | ✅ | 正/负提示词文本框 |
| 参数设置 | ✅ | steps, cfg, seed, sampler |
| 一键复制 | ✅ | 复制完整提示词 |
| 发送到预设库 | ✅ | 保存为预设 |

### 图片处理

| 功能 | 状态 | 组件 |
|------|------|------|
| 图片查看器 | ✅ | `ImageViewer.vue` |
| 缩放/拖拽 | ✅ | 鼠标滚轮/拖拽 |
| 键盘导航 | ✅ | ← → 切换, ESC 关闭 |
| 封面设置 | ✅ | 预设中选择封面 |
| 预览图编辑 | ✅ | 预设中增删预览图 |

### 设置与主题

| 功能 | 状态 | 组件 |
|------|------|------|
| **设置弹窗** | ✅ | `SettingsModal.vue` - 齿轮图标打开 |
| **主题切换** | ✅ | 深色 / 浅色 / 跟随系统 |
| **AI 提供商管理** | ✅ | 添加/编辑/删除自定义 AI |
| **Prompt 模板管理** | ✅ | 自定义 System Prompt |

**主题使用方式：**
```javascript
appStore.setTheme('dark')   // 深色
appStore.setTheme('light')  // 浅色
appStore.setTheme('auto')   // 跟随系统
```

### AI 系统 (增强 2026-01-29)

| 功能 | 状态 | 组件/Store |
|------|------|------------|
| **多 AI 提供商支持** | ✅ | `ai.js` - OpenAI, DeepSeek, Kimi, Ollama |
| **自定义 AI 提供商** | ✅ | 添加任意 OpenAI 兼容 API |
| **Prompt 模板系统** | ✅ | 拆解/优化/翻译/分析 + 自定义模板 |
| **可编辑 System Prompt** | ✅ | 完全控制 AI 行为和输出格式 |
| **Temperature 调节** | ✅ | 控制 AI 输出随机性 |
| **响应格式选择** | ✅ | JSON / 纯文本 |
| **API 调用历史** | ✅ | 最近 50 条记录 |

**内置 AI 提供商：**
| 提供商 | ID | 类型 | 地址 |
|--------|-----|------|------|
| OpenAI | `openai` | openai-compatible | https://api.openai.com/v1 |
| DeepSeek | `deepseek` | openai-compatible | https://api.deepseek.com/v1 |
| Kimi | `kimi` | openai-compatible | https://api.moonshot.cn/v1 |
| Ollama | `ollama` | ollama | http://localhost:11434 |

**内置 Prompt 模板：**
| 模板 | ID | 功能 |
|------|-----|------|
| 拆解提示词 | `explode` | 将长提示词拆解为原子词列表 |
| 优化提示词 | `optimize` | 优化提示词质量和表达 |
| 翻译提示词 | `translate` | 中文翻译成英文 |
| 分析提示词 | `analyze` | 分析提示词结构和效果 |

**AI Store 使用示例：**
```javascript
import { useAIStore } from './stores'

const aiStore = useAIStore()

// 初始化（从 localStorage 加载配置）
aiStore.init()

// 调用 AI
const result = await aiStore.callAI('输入的提示词', {
  provider: aiStore.currentProvider,
  prompt: aiStore.currentPrompt
})

// 添加自定义 AI 提供商
aiStore.addProvider({
  name: 'Claude',
  type: 'openai-compatible',
  baseUrl: 'https://api.anthropic.com/v1',
  apiKey: 'sk-xxx',
  model: 'claude-3-opus'
})

// 添加自定义 Prompt 模板
aiStore.addPrompt({
  name: '风格转换',
  description: '将提示词转换为特定风格',
  systemPrompt: '你是一个风格转换专家...',
  userPromptTemplate: '请将以下提示词转换为{{style}}风格：\n\n"{{input}}"',
  temperature: 0.8,
  responseFormat: 'json'
})
```

### 其他功能

| 功能 | 状态 | 说明 |
|------|------|------|
| AI 优化弹窗 | 🟡 | 界面完成，配置完成，待接入真实 API |
| 紧凑模式 | ✅ | 窗口 < 1100px 自动切换 |
| 暗黑主题 | ✅ | Tailwind slate 配色 |
| **浅色主题** | ✅ | CSS 变量实现，实时切换 |

---

## 🚧 待实现功能

### 高优先级

| 功能 | 状态 | 说明 |
|------|------|------|
| **后端 API 对接** | 🚧 | 原子词、分类、预设、版本的 CRUD 接口 |
| **AI 功能完整实现** | 🚧 | 连接真实的 OpenAI / Ollama API |
| **版本快照恢复** | 🚧 | 点击"使用该版本"恢复工作台状态 |
| **真实版本对比** | 🚧 | 当前 diff 为模拟数据，需实现真实对比算法 |

### 中优先级

| 功能 | 状态 | 说明 |
|------|------|------|
| 导入/导出 | ⏳ | 预设的 JSON 导入导出 |
| 批量操作 | ⏳ | 批量删除、移动分类 |
| 拖拽排序 | ⏳ | 分类、原子词拖拽排序 |
| 最近使用 | ⏳ | 显示最近使用的原子词 |
| **快捷键系统** | ⏳ | Ctrl+S 保存, Ctrl+F 搜索 |

### 低优先级

| 功能 | 状态 | 说明 |
|------|------|------|
| 国际化 | ⏳ | 多语言支持 |
| 数据备份 | ⏳ | 数据库备份恢复 |
| 性能优化 | ⏳ | 大数据量虚拟滚动 |

---

## 🔧 开发注意事项

### 已知问题

1. **go-winio 导入大小写问题**
   - 解决方案：确保使用 `github.com/Microsoft/go-winio` (大写 M)

2. **Handler 绑定时机**
   - 必须在 `wails.Run` 之前初始化所有 Handler
   - 否则会导致运行时 nil pointer 错误

3. **Windows 命令分隔符**
   - 使用 `;` 而非 `&&`
   - 示例：`cd dir; command1; command2`

4. **前端语法注意**
   - Vue 模板中避免多行属性
   - Mustache 表达式 `{{ }}` 必须正确闭合

### 代理配置

```powershell
$env:GOPROXY = "https://goproxy.cn"
```

### 常用命令

```powershell
# 开发模式
wails dev

# 生产构建
wails build

# 清理构建
wails build -clean
```

---

## 📝 数据库模型

### 表结构

```go
// 原子词
type Atom struct {
    ID          uint
    Content     string
    CategoryID  uint
    UseCount    int
    ImagePath   string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 分类
type Category struct {
    ID       uint
    Name     string
    SortOrder int
    CreatedAt time.Time
}

// 预设
type Preset struct {
    ID              uint
    Title           string
    Description     string
    PosText         string
    NegText         string
    Model           string
    LoRAs           string // JSON 存储
    Params          string // JSON 存储
    ThumbnailPath   string
    Previews        string // JSON 存储 URL 列表
    CurrentVersion  int
    UseCount        int
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// 版本历史
type PresetVersion struct {
    ID          uint
    PresetID    uint
    VersionNum  int
    Snapshot    string // JSON 存储
    DiffStats   string // +n/-m
    IsStarred   bool
    Previews    string // JSON 存储
    CreatedAt   time.Time
}

// 预览图
type Preview struct {
    ID       uint
    FilePath string
    Type     string // "local" | "url"
    CreatedAt time.Time
}
```

---

## 🎨 UI 设计规范

### 颜色主题

```css
/* 深色主题（默认） */
--color-bg-dark: #0f172a;       /* slate-900 */
--color-bg-card: #1e293b;       /* slate-800 */
--color-bg-hover: #334155;      /* slate-700 */
--color-text-primary: #e2e8f0;   /* slate-200 */
--color-text-secondary: #94a3b8; /* slate-400 */

/* 浅色主题 */
--color-bg-light: #ffffff;
--color-bg-card-light: #f8fafc;
--color-text-light: #1e293b;

/* 状态色 */
--color-primary: #0284c7;       /* sky-600 */
--color-success: #22c55e;       /* green-500 */
--color-warning: #f59e0b;       /* amber-500 */
--color-danger: #ef4444;        /* red-500 */
--color-star: #f59e0b;          /* amber-500 */
```

### 组件尺寸

- 侧边栏：220px (180px compact)
- 顶部栏：60px
- 卡片间距：8px（优化后）
- 圆角：8px (卡片), 12px (弹窗)
- 版本时间线高度：可拖拽调整（默认 400px）

---

## 🔗 状态管理 (Pinia)

### Store 列表

| Store | 文件 | 用途 |
|-------|------|------|
| appStore | `app.js` | 应用级状态：搜索、主题、AI配置、时间线显示 |
| atomStore | `atom.js` | 原子词 CRUD、搜索过滤 |
| categoryStore | `category.js` | 分类 CRUD |
| presetStore | `preset.js` | 预设 CRUD |
| versionStore | `version.js` | 版本历史、mock 数据 |

### 配置持久化

| 配置项 | 存储位置 | 说明 |
|--------|----------|------|
| 主题设置 | localStorage | `theme: 'dark' \| 'light' \| 'auto'` |
| AI API 配置 | localStorage | OpenAI/Ollama 配置 |

---

## 💡 使用说明

### 配置 AI API
1. 点击顶部 ⚡ 按钮打开 AI 拆解
2. 点击右上角 "OpenAI" 或 "Ollama" 配置按钮
3. 在设置面板中选择提供商
4. 填写 API Key 和模型信息
5. 点击保存配置

### 切换主题
1. 点击右上角设置按钮（齿轮图标）
2. 在"常规"标签页选择主题
3. 支持实时切换，刷新后保持

### 搜索原子词
1. 在顶部搜索框输入关键词
2. 支持英文、中文、近义词搜索
3. 实时显示搜索结果（无需回车）

### 查看版本历史
1. 在预设库中点击预设卡片
2. 点击"版本历史"按钮
3. 底部面板显示版本时间线
4. 点击版本行可展开详情（显示提示词、参数）
5. 点击眼睛图标查看完整版本详情
6. 拖拽面板顶部可调整高度

---

*文档版本: v2.1.0*  
*更新日期: 2026-01-29*
