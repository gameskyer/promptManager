# PromptMaster 变更日志 - 2026-01-29

## 🎯 今日主要更新

### 1. 多AI配置系统（新增功能）

**新增文件：**
- `frontend/src/stores/ai.js` - AI配置和Prompt模板管理Store

**功能说明：**
- 支持多个AI提供商：OpenAI、DeepSeek、Kimi、Ollama
- 支持自定义添加任意OpenAI兼容API
- 支持自定义Prompt模板（System Prompt + User Prompt模板）
- 支持调节Temperature和响应格式
- 配置持久化到localStorage

**使用方式：**
```javascript
import { useAIStore } from './stores'
const aiStore = useAIStore()

// 添加自定义AI
aiStore.addProvider({
  name: 'Claude',
  type: 'openai-compatible',
  baseUrl: 'https://api.anthropic.com/v1',
  apiKey: 'sk-xxx',
  model: 'claude-3-opus'
})

// 添加自定义Prompt模板
aiStore.addPrompt({
  name: '风格转换',
  description: '将提示词转换为特定风格',
  systemPrompt: '你是一个风格转换专家...',
  userPromptTemplate: '将以下提示词转换为{{style}}风格：\n\n"{{input}}"',
  temperature: 0.8,
  responseFormat: 'json'
})
```

---

### 2. AI弹窗升级（AIModal.vue）

**变更内容：**
- 添加AI提供商选择下拉框
- 添加功能模式选择（拆解/优化/翻译/分析）
- 根据选择的Prompt动态显示结果格式
- 支持复制和导入结果

**界面布局：**
```
┌─────────────────────────────────┐
│  AI 智能助手           [配置 ⚙️] │
├─────────────────────────────────┤
│  AI提供商: [DeepSeek ▼]         │
│  功能模式: [拆解提示词 ▼]        │
│  ┌───────────────────────────┐  │
│  │ 💡 将提示词拆解为原子词... │  │
│  └───────────────────────────┘  │
│  ┌───────────────────────────┐  │
│  │ 输入提示词...              │  │
│  └───────────────────────────┘  │
│     [⚡ 开始处理]               │
└─────────────────────────────────┘
```

---

### 3. 设置面板扩展（SettingsModal.vue）

**新增标签页：**
- **AI 提供商** - 管理AI提供商配置
  - 显示已配置的AI列表
  - 支持添加/编辑/删除自定义AI
  - 重置为默认配置
  
- **Prompt 模板** - 管理LLM Prompt模板
  - 显示所有Prompt模板
  - 支持添加/编辑/删除自定义模板
  - 可编辑System Prompt和User Prompt模板

**界面截图：**
```
┌─────────────────────────────────┐
│  设置                      [×]  │
├──────────┬──────────┬──────────┤
│  常规    │ AI提供商 │Prompt模板│
├──────────┴──────────┴──────────┤
│                                 │
│  [已配置的 AI 提供商]    [+添加] │
│  ┌───────────────────────────┐  │
│  │ OpenAI              [⚙️]  │  │
│  │ 模型：gpt-3.5-turbo       │  │
│  │ 地址：api.openai.com/v1   │  │
│  └───────────────────────────┘  │
│  ┌───────────────────────────┐  │
│  │ DeepSeek ★          [⚙️]  │  │
│  │ 模型：deepseek-chat       │  │
│  │ 地址：api.deepseek.com/v1 │  │
│  └───────────────────────────┘  │
│                                 │
│  [重置为默认]                    │
│                                 │
└─────────────────────────────────┘
```

---

### 4. 移除TopBar保存按钮（TopBar.vue）

**变更内容：**
- 移除顶部栏的"保存"按钮
- 保留快捷键Ctrl+S功能
- 保留Workbench中的"保存预设"按钮

**原因：**
- 功能重复（Workbench已有保存预设按钮）
- 简化顶部栏界面

**变更前：**
```
[新建 ▼]  [⚡AI]  [⚙️设置]  [💾保存]
```

**变更后：**
```
[新建 ▼]  [⚡AI]  [⚙️设置]
```

---

### 5. 原子词自动选择分类（AtomDialog.vue + MainContent.vue + TopBar.vue）

**功能说明：**
- 新建原子词时自动选中当前分类
- 优先使用子分类，其次父分类

**代码变更：**
```javascript
// AtomDialog.vue
const props = defineProps({
  defaultCategoryId: {
    type: Number,
    default: 0,
  },
})

// 监听 defaultCategoryId 变化
watch(() => props.defaultCategoryId, (newId) => {
  if (!props.atom && newId > 0) {
    form.value.category_id = newId
  }
}, { immediate: true })
```

```javascript
// TopBar.vue / MainContent.vue
const currentSelectedCategoryId = computed(() => {
  return currentSubCategory.value?.id || currentCategory.value?.id || 0
})

<AtomDialog :default-category-id="currentSelectedCategoryId" />
```

---

### 6. 组合工作区分区显示（Workbench.vue）

**重大更新：**
- 原子词按类型分组显示（正向/负向）
- 正向提示词使用绿色标识
- 负向提示词使用红色标识
- 独立的拖拽排序（同类型内）
- 独立的复制按钮

**界面布局：**
```
┌─────────────────────────────┐
│      组合工作区        [🗑️] │
├─────────────────────────────┤
│ [+] 正向提示词 (3)          │
│ ┌─────────────────────────┐ │
│ │ ≡ long hair   长发 [↑↓×]│ │
│ │ ≡ blonde      金发 [↑↓×]│ │
│ │ ≡ smile       微笑 [↑↓×]│ │
│ └─────────────────────────┘ │
│ [-] 负向提示词 (2)          │
│ ┌─────────────────────────┐ │
│ │ ≡ worst quality [↑↓×] │ │
│ │ ≡ blurry        [↑↓×] │ │
│ └─────────────────────────┘ │
├─────────────────────────────┤
│ [+] 正向提示词预览    [复制]│
│ long hair, blonde, smile    │
├─────────────────────────────┤
│ [-] 负向提示词预览    [复制]│
│ worst quality, blurry       │
├─────────────────────────────┤
│ [保存预设]        [生成]    │
└─────────────────────────────┘
```

---

### 7. 组合工作区保存预设对话框（Workbench.vue）

**新增功能：**
- 点击"保存预设"弹出对话框
- 输入预设名称和描述
- 预览正向/负向提示词内容
- 保存到预设库

**对话框内容：**
```
┌─────────────────────────────────┐
│  保存预设                  [×]  │
├─────────────────────────────────┤
│  预设名称 *                     │
│  ┌───────────────────────────┐  │
│  │ my-style-v1               │  │
│  └───────────────────────────┘  │
│  描述                           │
│  ┌───────────────────────────┐  │
│  │ 我的自定义风格预设...      │  │
│  └───────────────────────────┘  │
│  正向提示词 (3个)                │
│  ┌───────────────────────────┐  │
│  │ long hair, blonde, smile  │  │
│  └───────────────────────────┘  │
│  负向提示词 (2个)                │
│  ┌───────────────────────────┐  │
│  │ worst quality, blurry     │  │
│  └───────────────────────────┘  │
├─────────────────────────────────┤
│          [取消]    [保存]       │
└─────────────────────────────────┘
```

---

### 8. 修复子分类显示问题（SideMenu.vue + category.js）

**问题：**
- 添加子分类后不显示
- 保存后被重置

**修复内容：**

**category.js:**
```javascript
// 修复前：每次fetch都会重置为模拟数据
async function fetchCategories() {
  categories.value = [...DEFAULT_CATEGORIES] // 覆盖新添加的分类
}

// 修复后：只在没有数据时加载
async function fetchCategories() {
  if (categories.value.length > 0) return
  categories.value = [...DEFAULT_CATEGORIES]
}
```

**SideMenu.vue:**
```javascript
// 使用直接方法替代 computed getter
function getCategoryChildren(parentId) {
  return categoryStore.categories.filter(c => c.parent_id === parentId)
}

// 保存后不再调用 fetchCategories（避免重置）
async function saveCategory(data) {
  await categoryStore.createCategory(data)
  // 不再调用 fetchCategories()
}
```

---

### 9. 修复原子词添加后显示问题（MainContent.vue + atom.js）

**问题：**
- 添加原子词后调用 `fetchAtoms()` 重置为模拟数据
- 新添加的原子词被过滤掉

**修复：**
```javascript
// 修改前
async function saveAtom(data) {
  await atomStore.createAtom(data)
  await atomStore.fetchAtoms(currentSubCategory.value.id) // 重置！
}

// 修改后
async function saveAtom(data) {
  await atomStore.createAtom(data)
  // 不再调用 fetchAtoms，新原子词已在列表中
}
```

---

### 10. 移除筛选按钮（MainContent.vue）

**变更：**
- 移除"筛选"按钮（功能未实现）
- 顶部搜索栏已满足基础筛选需求

---

### 11. 缓存清理指南（CLEAN_CACHE_GUIDE.md）

**新增文档：**
- Wails 缓存位置说明
- 标准清理 vs 深度清理
- 一键清理脚本
- 常见问题解决方案

---

## 📊 变更统计

| 类别 | 数量 |
|------|------|
| 新增文件 | 2 (ai.js, CLEAN_CACHE_GUIDE.md) |
| 修改文件 | 9 |
| 新增功能 | 6 |
| Bug修复 | 3 |
| 优化改进 | 2 |

---

## 🔄 文件变更列表

### 新增文件
```
frontend/src/stores/ai.js
CLEAN_CACHE_GUIDE.md
```

### 修改文件
```
frontend/src/stores/category.js
frontend/src/stores/index.js
frontend/src/stores/app.js
frontend/src/components/AIModal.vue
frontend/src/components/SettingsModal.vue
frontend/src/components/TopBar.vue
frontend/src/components/AtomDialog.vue
frontend/src/components/MainContent.vue
frontend/src/components/Workbench.vue
frontend/src/components/SideMenu.vue
frontend/src/App.vue
```

---

## ⚠️ 已知限制

### 模拟数据模式
- 新增的分类/原子词/预设**刷新页面后会丢失**
- 等待后端 API 对接后数据会持久化到数据库

### 待实现功能
- 真实的 AI API 调用（目前是模拟数据）
- 后端数据持久化
- 预设库详细功能

---

**更新日期:** 2026-01-29  
**版本:** v2.1.0-dev
