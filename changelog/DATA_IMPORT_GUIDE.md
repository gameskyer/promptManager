# PromptMaster 数据导入导出指南

## 概述

PromptMaster 支持以下数据导入导出方式：

1. **自动种子数据** - 应用首次启动时自动导入默认数据
2. **JSON 导入** - 支持自定义 JSON 格式数据导入
3. **JSON 导出** - 导出原子词数据为 JSON 格式

## 自动种子数据

### 启动时自动导入

应用启动时会自动检查数据库状态，如果没有数据则会自动导入默认的种子数据：

- 7个一级分类（人物、场景、风格、质量、光照、道具、预设库）
- 5个二级分类（发型、眼睛、服装、姿势、表情）
- 40+ 常用原子词
- 3个示例预设

### 默认数据结构

```json
{
  "categories": [
    {"name": "人物", "type": "ATOM", "parent_id": 0},
    {"name": "发型", "type": "ATOM", "parent_id": 1}
  ],
  "atoms": [
    {
      "value": "masterpiece",
      "label": "杰作",
      "type": "Positive",
      "category": "质量",
      "synonyms": ["best quality"]
    }
  ],
  "presets": [
    {
      "title": "动漫女孩基础预设",
      "pos_text": "masterpiece, best quality, 1girl...",
      "neg_text": "low quality, bad anatomy...",
      "params": {"steps": 30, "cfg": 7.0}
    }
  ]
}
```

## 手动导入数据

### 通过后端 API

```go
// 获取 SeederHandler
seederHandler := handlers.NewSeederHandler(seeder)

// 导入默认数据
response := seederHandler.SeedAll()

// 从 JSON 导入自定义数据
response := seederHandler.ImportFromJSON(ImportFromJSONRequest{
    JSONData: `{"categories": [...], "atoms": [...]}`,
})
```

### 通过前端 Store

```javascript
import { useAtomStore } from './stores/atom'

const atomStore = useAtomStore()

// 批量导入原子词
const jsonData = JSON.stringify([
  { value: 'custom1', label: '自定义1', type: 'Positive', category_id: 1, synonyms: [] },
  { value: 'custom2', label: '自定义2', type: 'Positive', category_id: 1, synonyms: [] }
])
const importedCount = await atomStore.batchImport(jsonData)
```

## 数据导出

### 导出原子词

```javascript
import { useAtomStore } from './stores/atom'

const atomStore = useAtomStore()

// 导出所有原子词为 JSON
const jsonData = await atomStore.exportAtoms()
```

## 测试数据

### 运行测试

```bash
# 运行所有测试
go test ./backend/tests/... -v

# 运行特定测试
go test ./backend/tests/... -v -run TestSeeder
go test ./backend/tests/... -v -run TestIntegration
```

### 测试覆盖

- ✅ CategoryService - 分类 CRUD
- ✅ AtomService - 原子词 CRUD、搜索、批量导入
- ✅ PresetService - 预设 CRUD、软删除
- ✅ VersionService - 版本控制、对比
- ✅ Seeder - 数据导入
- ✅ Integration - 端到端工作流

## API 接口列表

### SeederHandler

| 方法 | 描述 |
|------|------|
| `SeedAll()` | 导入默认种子数据 |
| `ImportFromJSON(jsonData)` | 从 JSON 导入数据 |
| `GetDefaultSeedData()` | 获取默认种子数据 JSON |
| `GetSeedStatus()` | 获取数据库种子状态 |

### AtomHandler

| 方法 | 描述 |
|------|------|
| `CreateAtom(req)` | 创建原子词 |
| `GetAtomByID(id)` | 获取原子词详情 |
| `GetAtomsByCategory(req)` | 分页获取分类下的原子词 |
| `UpdateAtom(req)` | 更新原子词 |
| `DeleteAtom(id)` | 删除原子词 |
| `FindAtomsBySynonym(term)` | 搜索原子词 |
| `BatchImportAtoms(req)` | 批量导入原子词 |
| `ExportAtoms()` | 导出所有原子词 |

### CategoryHandler

| 方法 | 描述 |
|------|------|
| `CreateCategory(req)` | 创建分类 |
| `GetCategoryTree(type)` | 获取完整分类树 |
| `GetCategoriesByParent(parentID, type)` | 获取子分类 |
| `UpdateCategory(req)` | 更新分类 |
| `DeleteCategory(id)` | 删除分类 |
| `MoveCategory(req)` | 移动分类 |
| `ReorderCategories(req)` | 排序分类 |

### PresetHandler

| 方法 | 描述 |
|------|------|
| `CreatePreset(req)` | 创建预设（自动创建V1） |
| `GetPresets(req)` | 分页获取预设列表 |
| `GetPresetByID(id)` | 获取预设详情 |
| `UpdatePreset(req)` | 更新预设标题 |
| `SoftDeletePreset(id)` | 软删除预设 |
| `RestorePreset(id)` | 恢复预设 |
| `ForkPreset(req)` | Fork 预设 |

### VersionHandler

| 方法 | 描述 |
|------|------|
| `CreateVersion(req)` | 创建新版本 |
| `GetVersionHistory(presetID, limit)` | 获取版本历史 |
| `GetVersion(presetID, versionNum)` | 获取特定版本 |
| `StarVersion(req)` | 星标/取消星标版本 |
| `RollbackToVersion(presetID, versionNum)` | 回滚到版本 |
| `CompareVersions(presetID, v1, v2)` | 对比两个版本 |
| `DeleteVersion(versionID)` | 删除版本 |

## 数据库模型

### Category
- `id` - 主键
- `name` - 分类名称
- `parent_id` - 父分类ID（0表示一级分类）
- `type` - 类型（ATOM/PRESET）
- `sort_order` - 排序权重

### Atom
- `id` - 主键
- `value` - 英文值（唯一）
- `label` - 中文标签
- `type` - 类型（Positive/Negative）
- `category_id` - 所属分类ID
- `synonyms` - 同义词列表（JSON）
- `usage_count` - 使用次数

### Preset
- `id` - 主键
- `title` - 标题
- `current_version` - 当前版本号
- `is_deleted` - 软删除标记

### PresetVersion
- `id` - 主键
- `preset_id` - 所属预设ID
- `version_num` - 版本号
- `snapshot` - 完整快照（JSON）
- `diff_stats` - 变更统计（如 +2/-1）
- `is_starred` - 星标标记

## 数据文件位置

- Windows: `%APPDATA%\PromptMaster\`
- macOS: `~/Library/Application Support/PromptMaster/`
- Linux: `~/.config/PromptMaster/`

数据库文件：`promptmaster.db`
