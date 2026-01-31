# 更新日志 - 2026-01-30

## 🎯 今日开发概览

| 类别 | 数量 | 说明 |
|------|------|------|
| 新增功能 | 4 | 工作区完善、管理页面、预设多分类支持 |
| 问题修复 | 1 | 导出功能 synonyms 解析错误 |
| 代码提交 | 4 | 详见 Git 提交记录 |

---

## ✨ 新增功能

### 1. 工作区功能完善

**提交:** `13c9054`

完成了工作区的核心交互功能：

- **拖拽排序** - 支持正向/负向原子词分别拖拽排序
- **删除功能** - 支持单个删除和一键清空
- **保存预设** - 完善保存对话框，支持填写标题和描述
- **提示词预览** - 实时显示组合后的正负向提示词
- **一键复制** - 支持分别复制正向/负向提示词

### 2. 原子词管理页面

**提交:** `13c9054`

全新的原子词管理界面（表格视图）：

| 功能 | 描述 |
|------|------|
| 表格展示 | 列表形式展示所有原子词，包含 ID、类型、英文、中文、分类、近义词、使用次数 |
| 分类筛选 | 支持按分类筛选 |
| 类型筛选 | 支持按正向/负向筛选 |
| 关键词搜索 | 支持搜索英文、中文、近义词 |
| 分页显示 | 支持 20/50/100 条每页 |
| 快捷操作 | 编辑、删除、导出、添加到工作区 |

**文件:** `frontend/src/components/AtomManagement.vue`

### 3. 分类管理页面

**提交:** `13c9054`

可视化的分类管理界面：

- **树形结构** - 清晰展示分类层级关系
- **拖拽排序** - 支持拖拽移动分类，调整父子关系
- **添加子分类** - 快速添加子分类
- **展开/折叠** - 支持展开/折叠全部节点
- **编辑删除** - 完整的分类生命周期管理

**文件:** 
- `frontend/src/components/CategoryManagement.vue`
- `frontend/src/components/CategoryTreeNode.vue`

### 4. 预设多分类支持 ⭐

**提交:** `86886e6`, `bb5c889`

预设系统现在支持多分类管理：

#### 后端改动
- **Preset 模型** - 新增 `CategoryID` 字段，与 Category 建立外键关系
- **API 更新** - 创建/更新/查询预设时支持分类参数
- **分类筛选** - `GetPresets` 支持按分类 ID 筛选

#### 前端改动
- **预设列表侧边栏** - 左侧新增分类树，支持：
  - 全部预设
  - 各预设分类（显示数量徽章）
  - 未分类
- **预设对话框** - 新增分类选择下拉框
- **侧边栏更新** - "预设库"分组显示分类列表

#### 使用方式
```
侧边栏 → 预设库
├── 全部预设 (12)
├── 动漫风格 (5)
├── 写实风格 (3)
└── 未分类 (4)
```

**注意:** 创建预设分类需要点击侧边栏"预设库"旁的 + 按钮，或在分类管理页面选择类型为 PRESET。

---

## 🐛 问题修复

### 修复导出功能错误

**提交:** `ec6eef6`

**问题:** 导出原子词时报错：`invalid character 'b' looking for beginning of value`

**原因:** 数据库中 `synonyms` 列数据格式不一致
- 早期数据：`'best quality'` (普通字符串)
- 新数据：`'["tag1", "tag2"]'` (JSON 数组)

**修复:** 更新 `StringSlice.Scan()` 方法，兼容两种格式

```go
// 处理逻辑
1. 以 '[' 开头的 → 按 JSON 数组解析
2. 普通字符串 → 作为单元素数组
3. 空值 → 转为空数组
```

**文件:** `backend/models/models.go`

---

## 📁 新增文件

```
frontend/src/components/
├── AtomManagement.vue      # 原子词管理页面
├── CategoryManagement.vue  # 分类管理页面
└── CategoryTreeNode.vue    # 分类树节点组件（递归）
```

## 📝 修改文件

### 后端
```
backend/models/models.go           # 添加 CategoryID 到 Preset
backend/services/preset_service.go # 支持分类参数
backend/handlers/preset_handler.go # API 支持分类筛选
backend/utils/seeder.go            # 修复 seeder 调用
```

### 前端
```
frontend/src/App.vue               # 添加管理页面路由
frontend/src/components/SideMenu.vue          # 添加管理入口、预设分类
frontend/src/components/PresetList.vue        # 添加分类侧边栏
frontend/src/components/PresetDialog.vue      # 添加分类选择
frontend/src/components/Workbench.vue         # 完善工作区功能
frontend/src/stores/preset.js                 # 更新 store 方法
```

---

## 🔄 Git 提交记录

```bash
# 查看今日提交
git log --oneline --since="2026-01-30" --until="2026-01-31"

bb5c889 docs: update FEATURES.md with preset multi-category support
86886e6 feat: add multi-category support for presets
ec6eef6 fix: fix StringSlice Scan to handle legacy string format in synonyms column
13c9054 feat: complete workbench functionality and add management pages
```

---

## 📊 项目统计更新

| 类别 | 更新前 | 更新后 |
|------|--------|--------|
| 后端文件 | 20+ | 20+ |
| 前端组件 | 15+ | 18+ |
| Store 模块 | 8 | 8 |
| 已实现功能 | 80% | 85% |

---

## 🎯 下一步建议

### 高优先级
- [ ] AI 拆解结果一键导入工作区
- [ ] AI 优化结果应用到预设

### 中优先级
- [ ] 分类管理页面支持筛选类型（原子词/预设）
- [ ] 快捷键支持（Ctrl+S 保存、Ctrl+F 搜索等）

### 低优先级
- [ ] 回收站功能
- [ ] 深色/浅色主题切换完善

---

## 📝 附注

- 所有修改已合并到 `main` 分支
- 编译命令: `wails build`
- 清理缓存: `make clean` 或手动删除 `build/`, `frontend/dist/`

---

*最后更新: 2026-01-30*
