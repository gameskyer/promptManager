# ComfyUI 插件 - PromptMaster

## 概述

为 ComfyUI 开发的插件，实现实时提示词搜索功能，直接读取 PromptMaster 数据库。

## 功能特点

- 🔍 **实时搜索**：输入中文自动弹出建议列表（如输入"裙"显示"短裙"、"长裙"）
- 🌐 **中英双语**：支持中文标签和英文 value 搜索
- 📂 **分类筛选**：按分类筛选提示词
- ⚡ **快捷选择**：点击建议项即可填充到输入框
- 📦 **直接读取数据库**：无需 HTTP 服务，直接读取 SQLite 数据库

## 安装

```bash
cd plugin/comfyui_promptmaster
python install.py
```

或手动复制到 `ComfyUI/custom_nodes/comfyui_promptmaster`

## 使用方法

### 1. 使用 Prompt Search 节点

1. 添加节点：`PromptMaster` → `🔍 Prompt Search`
2. 在 `search_query` 输入框输入中文（如"裙"）
3. 从下拉列表选择提示词
4. 自动填充到 `selected_prompt` 输出

### 2. 在 CLIP Text Encode 中使用

安装后，所有字符串输入框自动增强：
- 输入中文时自动弹出建议列表
- 选择后自动替换当前输入的词

## 节点

| 节点 | 说明 |
|------|------|
| 🔍 Prompt Search | 基础搜索节点 |
| 🔍 Prompt Search (Advanced) | 高级搜索，支持权重调整 |

## 技术实现

- **后端**：直接读取 SQLite 数据库（`promptmaster.db`）
- **前端**：JavaScript 扩展，实时搜索建议
- **数据库**：搜索 `atoms` 表的 `value`、`label`、`synonyms` 字段

## 文件位置

```
plugin/comfyui_promptmaster/
├── __init__.py
├── database.py          # 数据库访问
├── prompt_search_node.py # 节点定义
├── web/
│   └── prompt_search.js # 前端扩展
└── README.md
```

## Git 提交

```
f7d11d6 feat: add ComfyUI plugin with real-time prompt search (Chinese/English)
```
