# 前端与后端 API 集成说明

## 集成状态

✅ **已完成** - 前端已切换到后端 API 数据

## 变更内容

### 1. Store 层更新

所有 Store 已从模拟数据切换到后端 API：

| Store | 状态 | 连接的后端 Handler |
|-------|------|-------------------|
| `atom.js` | ✅ 已切换 | AtomHandler |
| `category.js` | ✅ 已切换 | CategoryHandler |
| `preset.js` | ✅ 已切换 | PresetHandler |
| `version.js` | ✅ 已切换 | VersionHandler |

### 2. 数据流

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Vue 组件   │────▶│   Pinia    │────▶│  Wails API  │
│             │     │   Store    │     │  (Go后端)   │
└─────────────┘     └─────────────┘     └─────────────┘
    MainContent           atom.js        AtomHandler
    SideMenu              category.js    CategoryHandler
    TopBar                preset.js      PresetHandler
    Timeline              version.js     VersionHandler
    PresetList
```

### 3. API 调用示例

#### Atom Store
```javascript
// 获取原子词
await atomStore.fetchAtoms(categoryId)

// 创建原子词
await atomStore.createAtom({
  value: 'long hair',
  label: '长发',
  type: 'Positive',
  category_id: 5,
  synonyms: ['lengthy hair']
})

// 搜索
await atomStore.searchAtoms('long hair')
```

#### Category Store
```javascript
// 获取分类树
await categoryStore.fetchCategories()

// 创建分类
await categoryStore.createCategory('新分类', 'ATOM', 0)
```

#### Preset Store
```javascript
// 获取预设
await presetStore.fetchPresets()

// 创建预设
await presetStore.createPreset(
  '标题',
  '正向提示词',
  '负向提示词',
  [],
  { steps: 30, cfg: 7 },
  []
)
```

#### Version Store
```javascript
// 获取版本历史
await versionStore.fetchVersions(presetId)

// 星标版本
await versionStore.starVersion(versionId, true)

// 回滚
await versionStore.rollbackToVersion(presetId, versionNum)
```

## 开发说明

### 启动开发服务器

```bash
wails dev
```

### 构建生产版本

```bash
wails build
```

### 数据库位置

- 开发模式：`./promptmaster.db`（项目根目录）
- 首次启动会自动导入种子数据

### 调试

浏览器开发者工具 (F12) → Console 查看 API 调用日志：

```javascript
// 检查原子词数据
const atomStore = useAtomStore()
console.log(atomStore.atoms)

// 检查分类数据
const categoryStore = useCategoryStore()
console.log(categoryStore.categories)
```

## 常见问题

### 1. 数据不显示
- 检查后端是否正确启动
- 检查数据库文件是否存在
- 查看浏览器控制台错误日志

### 2. API 调用失败
- 确保 Wails 绑定已重新生成：`wails generate module`
- 检查后端 Handler 是否正确注册在 `main.go` 中

### 3. 数据库为空
- 删除 `promptmaster.db` 重新启动应用
- 应用会自动导入种子数据

## 文件变更

### Store 文件
- `frontend/src/stores/atom.js` - 更新为连接后端 API
- `frontend/src/stores/category.js` - 更新为连接后端 API
- `frontend/src/stores/preset.js` - 更新为连接后端 API
- `frontend/src/stores/version.js` - 更新为连接后端 API

### 自动生成的绑定
- `frontend/src/lib/wailsjs/go/handlers/*.js` - Wails 自动生成

## 下一步

1. ✅ 前端连接后端 API
2. ✅ 数据库自动导入种子数据
3. ✅ 测试用例覆盖
4. 🔄 可以开始添加更多功能，如：
   - 图片上传功能
   - AI 实际调用
   - 导入/导出 UI
