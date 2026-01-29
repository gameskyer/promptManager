# AI API 接入实现说明

## 实现功能

### 1. 后端 AI Service (`backend/services/ai_service.go`)

支持四种 AI 功能：

| 功能 | 方法 | 说明 |
|------|------|------|
| 拆解提示词 | `ExplodePrompt()` | 将长提示词拆解为原子词 |
| 优化提示词 | `OptimizePrompt()` | 优化提示词质量和结构 |
| 翻译提示词 | `TranslatePrompt()` | 中文翻译成英文 |
| 分析提示词 | `AnalyzePrompt()` | 分析提示词结构和效果 |

### 2. 后端 AI Handler (`backend/handlers/ai_handler.go`)

API 端点：

| 端点 | 请求类型 | 功能 |
|------|---------|------|
| `ExplodePrompt` | ExplodePromptRequest | 拆解提示词 |
| `OptimizePrompt` | OptimizePromptRequest | 优化提示词 |
| `TranslatePrompt` | TranslatePromptRequest | 翻译提示词 |
| `AnalyzePrompt` | AnalyzePromptRequest | 分析提示词 |
| `ProcessAI` | GenericAIRequest | 通用处理接口 |
| `ImportExtractedAtoms` | ImportExtractedRequest | 导入拆解的原子词 |

### 3. 前端 AI Store (`frontend/src/stores/ai.js`)

提供方法：

```javascript
// 拆解提示词
await aiStore.explodePrompt(prompt, categoryId)

// 优化提示词
await aiStore.optimizePrompt(prompt)

// 翻译提示词
await aiStore.translatePrompt(prompt)

// 分析提示词
await aiStore.analyzePrompt(prompt)

// 通用调用（根据当前模式）
await aiStore.callAI(input, mode)
```

### 4. 支持的 AI 提供商

- **OpenAI** (GPT-3.5/GPT-4)
- **DeepSeek**
- **Kimi (Moonshot)**
- **Ollama** (本地)

## 配置方法

### 在设置中配置 AI

1. 打开应用 → 点击设置（齿轮图标）
2. 选择 "AI 提供商" 标签
3. 点击已有提供商编辑，或点击 "添加" 新建
4. 填写：
   - 名称
   - 类型（OpenAI 兼容/Ollama）
   - API 地址
   - API Key（Ollama 可选）
   - 默认模型
5. 保存配置

### 使用 AI 功能

1. 点击顶部栏 "⚡ AI" 按钮
2. 选择功能模式：
   - 拆解提示词
   - 优化提示词
   - 翻译提示词
   - 分析提示词
3. 输入提示词内容
4. 点击 "开始处理"
5. 查看结果并导入（拆解功能）

## 数据结构

### 拆解结果 (ExplodeResult)
```json
{
  "atoms": [
    {
      "value": "masterpiece",
      "label": "杰作",
      "type": "Positive",
      "category": "质量",
      "synonyms": ["best quality"],
      "is_new": false
    }
  ],
  "raw_prompt": "原始提示词"
}
```

### 优化结果 (OptimizeResult)
```json
{
  "optimized": "优化后的提示词",
  "changes": ["修改说明1"],
  "suggestions": ["建议1"]
}
```

### 翻译结果 (TranslateResult)
```json
{
  "translation": "English translation",
  "keywords": ["keyword1"],
  "notes": "翻译说明"
}
```

### 分析结果 (AnalyzeResult)
```json
{
  "analysis": {
    "subject": "主体描述",
    "style": "风格描述",
    "quality": "质量描述",
    "lighting": "光照描述",
    "other": "其他"
  },
  "issues": ["问题1"],
  "suggestions": ["建议1"]
}
```

## 技术实现

### API 调用流程

```
Vue Component
    ↓
AI Store (ai.js) - callAI(mode, input)
    ↓
Wails Bridge
    ↓
Go Handler (ai_handler.go) - ProcessAI()
    ↓
Go Service (ai_service.go) - 根据 mode 分发
    ↓
AI Provider API (OpenAI/DeepSeek/Kimi/Ollama)
    ↓
返回 JSON 结果
```

### 本地模式（无 API Key）

当没有配置 API Key 时，系统使用基于规则的降级方案：
- **拆解**: 按逗号分割，检查现有原子词
- **优化**: 去重、规范化
- **翻译**: 返回原文（提示配置 AI）
- **分析**: 基础统计信息

## 测试

### 测试命令
```bash
go test ./backend/tests/... -v -run TestAI
```

### 手动测试

1. 配置 AI API（或使用 Ollama 本地模式）
2. 打开 AI 弹窗
3. 测试各功能模式
4. 验证结果格式

## 注意事项

1. **API Key 安全** - 存储在 localStorage，生产环境应使用系统密钥链
2. **超时设置** - API 调用超时 60 秒
3. **错误处理** - 失败时返回规则降级结果
4. **JSON 解析** - 支持从 Markdown 代码块中提取 JSON
