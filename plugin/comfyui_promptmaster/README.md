# PromptMaster for ComfyUI

ComfyUI 插件，集成 PromptMaster 的提示词搜索功能，支持在输入框中实时搜索中文/英文提示词。

## 功能特点

- 🔍 **实时搜索**：在字符串输入框中输入中文（如"裙"），自动弹出建议列表
- 🌐 **中英双语**：支持中文标签和英文 value 搜索
- 📂 **分类筛选**：按分类筛选提示词
- ⚡ **快捷选择**：点击建议项即可填充到输入框
- 📦 **直接读取数据库**：无需 HTTP 服务，直接读取 PromptMaster SQLite 数据库

## 安装方法

### 方法一：自动安装（推荐）

```bash
# 进入插件目录
cd plugin/comfyui_promptmaster

# 运行安装脚本
python install.py
```

### 方法二：手动安装

1. 将整个 `comfyui_promptmaster` 文件夹复制到 ComfyUI 的 `custom_nodes` 目录：
   ```
   ComfyUI/
   └── custom_nodes/
       └── comfyui_promptmaster/
           ├── __init__.py
           ├── database.py
           ├── prompt_search_node.py
           ├── web/
           │   └── prompt_search.js
           └── README.md
   ```

2. 重启 ComfyUI

## 使用方法

### 1. 使用 Prompt Search 节点

1. 在 ComfyUI 中右键添加节点：`PromptMaster` → `🔍 Prompt Search`
2. 在 `search_query` 输入框中输入中文（如"裙"）
3. 从下拉列表中选择合适的提示词
4. 选择的内容会自动填充到 `selected_prompt` 输出

### 2. 在 CLIP Text Encode 中使用

安装插件后，所有字符串输入框都会增强：

1. 在 CLIP Text Encode 节点的文本框中输入逗号分隔的提示词
2. 输入中文时（如"红"），会自动弹出建议列表
3. 选择后会自动替换当前输入的词

### 3. 示例工作流

```
[🔍 Prompt Search] ---> [CLIP Text Encode] ---> [KSampler]
       |                                               ^
       v                                               |
[选择 "red dress"]                          [生成图片]
```

## 节点说明

### 🔍 Prompt Search（基础搜索节点）

| 参数 | 类型 | 说明 |
|------|------|------|
| search_query | STRING | 搜索关键词（支持中文/英文） |
| category_filter | COMBO | 按分类筛选（全部/人物/服装/场景...） |
| max_results | INT | 最大结果数（1-50） |
| selected_prompt | STRING | 选择的提示词（输出） |

输出：`prompt` → 可直接连接到 CLIP Text Encode

### 🔍 Prompt Search Advanced（高级搜索节点）

额外功能：
- **weight**：权重调整（0.0-2.0）
- **prefix**：添加前缀（如"best quality"）
- **suffix**：添加后缀（如"masterpiece"）

输出：
- `positive_prompt` → 带权重的完整提示词
- `raw_value` → 原始英文 value

## 数据库配置

插件会自动搜索数据库文件，搜索顺序：

1. 插件同级目录：`../promptmaster.db`
2. ComfyUI 根目录：`./promptmaster.db`
3. 用户目录：`~/.promptmaster/promptmaster.db`

如果需要手动指定数据库路径，修改 `database.py`：

```python
def __init__(self, db_path: str = None):
    if db_path is None:
        db_path = "/your/custom/path/promptmaster.db"
```

## 文件结构

```
comfyui_promptmaster/
├── __init__.py              # 插件入口
├── database.py              # SQLite 数据库访问
├── prompt_search_node.py    # ComfyUI 节点定义
├── web/
│   └── prompt_search.js     # 前端搜索 UI 扩展
├── routes.py                # API 路由（可选）
├── install.py               # 安装脚本
└── README.md                # 本文件
```

## 技术说明

- **无需 HTTP 服务**：插件直接读取 SQLite 数据库，不需要启动 PromptMaster 的 HTTP 服务
- **实时搜索**：使用防抖技术（200ms），避免频繁查询
- **结果缓存**：最近 50 个搜索结果会被缓存
- **分类支持**：支持 PromptMaster 的分类系统

## 兼容性

- ✅ ComfyUI 最新版本
- ✅ 所有使用 STRING 类型输入的节点
- ✅ CLIP Text Encode
- ✅ 自定义提示词节点

## 故障排除

### 搜索无结果

1. 检查数据库文件是否存在
2. 查看 ComfyUI 控制台是否有数据库路径日志
3. 确认数据库中有原子词数据

### 界面不显示建议框

1. 检查浏览器控制台是否有 JavaScript 错误
2. 确保 JavaScript 文件正确加载
3. 尝试刷新页面或重启 ComfyUI

### 数据库连接失败

修改 `database.py` 中的 `_find_database` 方法，添加你的数据库路径：

```python
possible_paths = [
    "/path/to/your/promptmaster.db",  # 添加你的路径
    # ... 其他路径
]
```

## 更新日志

### v1.0.0
- 初始版本
- 支持中文/英文搜索
- 支持分类筛选
- 支持权重调整

## 许可证

与 PromptMaster 主项目相同
