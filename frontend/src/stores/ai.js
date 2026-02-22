import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  ExplodePrompt,
  OptimizePrompt,
  TranslatePrompt,
  AnalyzePrompt,
  ImportExtractedAtoms,
  ReverseImagePrompt,
  SaveAIConfig,
  GetAIConfig,
  ProcessAI,
} from '../lib/wailsjs/go/handlers/AIHandler'

// 默认AI提供商配置模板
const DEFAULT_PROVIDERS = [
  {
    id: 'openai',
    name: 'OpenAI',
    type: 'openai-compatible',
    baseUrl: 'https://api.openai.com/v1',
    apiKey: '',
    model: 'gpt-3.5-turbo',
    models: ['gpt-3.5-turbo', 'gpt-4', 'gpt-4-turbo'],
    headers: {},
    enabled: false,
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    type: 'openai-compatible',
    baseUrl: 'https://api.deepseek.com/v1',
    apiKey: '',
    model: 'deepseek-chat',
    models: ['deepseek-chat', 'deepseek-coder'],
    headers: {},
    enabled: false,
  },
  {
    id: 'kimi',
    name: 'Kimi (Moonshot)',
    type: 'openai-compatible',
    baseUrl: 'https://api.moonshot.cn/v1',
    apiKey: '',
    model: 'moonshot-v1-8k',
    models: ['moonshot-v1-8k', 'moonshot-v1-32k', 'moonshot-v1-128k'],
    headers: {},
    enabled: false,
  },
  {
    id: 'ollama',
    name: 'Ollama (本地)',
    type: 'ollama',
    baseUrl: 'http://localhost:11434',
    apiKey: '',
    model: 'llama2',
    models: ['llama2', 'mistral', 'codellama', 'vicuna'],
    headers: {},
    enabled: true,
  },
]

// 默认Prompt模板
const DEFAULT_PROMPTS = {
  explode: {
    id: 'explode',
    name: '拆解提示词',
    description: '将长段提示词拆解为原子词列表',
    systemPrompt: `你是一位专业的 AI 图像生成提示词拆解专家，擅长将复杂的提示词分解为结构化的原子词汇。

## 任务描述
将用户提供的提示词拆解为最小的语义单元（原子词），每个原子词代表一个独立的视觉概念或属性。

## 可用分类及ID
{{category_list}}

## 拆解规则

### 1. Value 字段（英文标识）
- 使用英文小写
- 多词使用下划线连接（如：blue_sky）
- 使用 Stable Diffusion 社区标准标签
- 避免使用生僻或模型不熟悉的词汇

### 2. Label 字段（中文说明）
- 提供简洁准确的中文翻译
- 保持与英文 value 的对应关系

### 3. Type 字段（类型标记）
- Positive: 正向提示词（期望出现的元素）
- Negative: 负向提示词（期望避免的元素）

### 4. Category 字段（分类ID）
- 必须是数字（ID），从可用分类中选择最合适的
- 可选分类：{{category_names}}
- 如果无法确定，使用 ID: 0

### 5. Synonyms 字段（近义词列表）
- 提供 1-3 个常用同义词或变体
- 使用英文，有助于提升召回率

## 输出示例

输入："一个女孩，蓝色眼睛，长发，在森林里，杰作品质"

输出：
{
  "atoms": [
    {"value": "1girl", "label": "一个女孩", "type": "Positive", "category": 2, "synonyms": ["one_girl", "single_girl"]},
    {"value": "blue_eyes", "label": "蓝眼睛", "type": "Positive", "category": 2, "synonyms": ["azure_eyes"]},
    {"value": "long_hair", "label": "长发", "type": "Positive", "category": 6, "synonyms": ["lengthy_hair"]},
    {"value": "forest", "label": "森林", "type": "Positive", "category": 4, "synonyms": ["woods", "jungle"]},
    {"value": "masterpiece", "label": "杰作", "type": "Positive", "category": 1, "synonyms": ["best_quality", "top_quality"]}
  ]
}

## 注意事项
- 必须严格返回 JSON 格式，不要包含 Markdown 代码块标记
- category 字段必须是数字 ID，不能是字符串
- 确保 JSON 格式合法，使用双引号`,
    userPromptTemplate: '{{input}}',
    temperature: 0.3,
    responseFormat: 'json',
  },
  optimize: {
    id: 'optimize',
    name: '优化提示词',
    description: '优化提示词的质量和表达',
    systemPrompt: `你是一位资深的 AI 图像生成提示词工程师，专注于优化提示词以获得最佳的图像生成效果。

## 任务描述
分析用户提供的提示词，进行专业优化，提升图像生成质量和一致性。

## 优化原则

### 1. 质量提升
- 添加必要的基础质量词（如：masterpiece, best quality, highres 等）
- 补充缺失的细节描述
- 使用更精确的专业术语

### 2. 结构优化
- 按重要性排序：质量词 → 主体 → 细节 → 风格 → 光照 → 背景
- 重要的描述词放在前面（AI 对前面的词权重更高）
- 相关概念分组排列

### 3. 去重与精简
- 去除语义重复的词汇
- 删除模糊的描述
- 合并相似的概念

### 4. 语法规范
- 使用英文逗号分隔
- 保持标签格式一致性
- 避免拼写错误

### 5. 风格一致性
- 确保风格描述与主体匹配
- 光照描述与场景协调
- 艺术风格词汇准确

## 输出格式

{
  "optimized": "优化后的完整提示词",
  "changes": [
    "修改说明1：具体改了什么",
    "修改说明2：为什么这样改"
  ],
  "suggestions": [
    "进一步优化的建议1",
    "可以添加的元素建议2",
    "风格调整建议3"
  ]
}

## 优化示例

输入："一个女孩在森林里，好看一点"

输出：
{
  "optimized": "masterpiece, best quality, 1girl, solo, detailed face, beautiful eyes, long hair, standing, forest, sunlight filtering through trees, dappled light, vibrant colors, detailed background",
  "changes": [
    "添加了质量词 'masterpiece, best quality' 提升整体画质",
    "将'一个女孩'标准化为'1girl, solo'标签",
    "细化了外貌描述：'detailed face, beautiful eyes, long hair'",
    "添加了姿态描述 'standing'",
    "丰富了场景细节：'sunlight filtering through trees, dappled light'",
    "将模糊的'好看一点'转化为具体的'vibrant colors, detailed background'"
  ],
  "suggestions": [
    "可以添加服装颜色描述以增强视觉效果",
    "考虑添加特定的艺术风格如'digital painting'或'anime style'",
    "可以指定光照方向如'from above'或'backlight'"
  ]
}`,
    userPromptTemplate: '{{input}}',
    temperature: 0.4,
    responseFormat: 'json',
  },
  translate: {
    id: 'translate',
    name: '翻译提示词',
    description: '将中文提示词翻译为英文',
    systemPrompt: `你是一位专业的 AI 图像生成提示词翻译专家，精通中文到英文的提示词翻译。

## 任务描述
将用户的中文提示词翻译成高质量、地道的英文提示词，适用于 Stable Diffusion、Midjourney 等 AI 图像生成模型。

## 翻译原则

### 1. 专业术语优先
- 使用 AI 绘画社区广泛认可的标准标签
- 优先使用 Danbooru 标签体系中的标准词汇
- 遵循 Stable Diffusion 模型的训练标签习惯

### 2. 格式规范
- 使用英文逗号分隔
- 标签使用小写（专有名词除外）
- 多词标签使用下划线连接

### 3. 语义准确性
- 保持原意的完整传达
- 文化特色词汇选择最贴切的英文表达
- 形容词和名词的搭配符合英文习惯

### 4. 结构保持
- 保持原提示词的概念顺序
- 质量词前置原则
- 主体 → 细节 → 场景 → 风格 → 质量的逻辑

### 5. 关键词提取
- 识别并提取核心视觉概念
- 按重要性排序关键词
- 包含风格、质量、光照等维度

## 常见翻译对照参考

| 中文概念 | 推荐英文 | 备选 |
|---------|---------|-----|
| 杰作品质 | masterpiece, best quality | ultra detailed |
| 一个女孩 | 1girl, solo | single girl |
| 精致面容 | detailed face, beautiful detailed eyes | delicate face |
| 长发 | long hair | lengthy hair |
| 阳光 | sunlight, sunshine | sunbeam |
| 森林 | forest, in the forest | woods |
| 赛博朋克 | cyberpunk, neon lights | sci-fi |
| 水墨风格 | ink wash painting, chinese ink style | sumi-e |

## 输出格式

{
  "translation": "完整的英文提示词，使用逗号分隔",
  "keywords": [
    "核心关键词1",
    "核心关键词2",
    "风格关键词",
    "质量关键词"
  ],
  "notes": "翻译说明和注意事项"
}

## 翻译示例

输入："一个穿着汉服的少女在樱花树下，古风，水墨画风格，高品质"

输出：
{
  "translation": "masterpiece, best quality, 1girl, solo, hanfu, traditional chinese clothes, cherry blossoms, under the tree, ancient style, ink wash painting, chinese ink style, soft colors, elegant, detailed background",
  "keywords": [
    "1girl",
    "hanfu",
    "cherry_blossoms",
    "ink_wash_painting",
    "masterpiece"
  ],
  "notes": "将'汉服'翻译为'hanfu'（Danbooru标准标签），'古风'拆分为'ancient style'和具体元素，'水墨画风格'提供两种表达方式供选择"
}`,
    userPromptTemplate: '{{input}}',
    temperature: 0.5,
    responseFormat: 'json',
  },
  analyze: {
    id: 'analyze',
    name: '分析提示词',
    description: '分析提示词的结构和效果',
    systemPrompt: `你是一位资深的 AI 图像生成提示词分析师，擅长深度分析提示词的构成、潜在问题和改进空间。

## 任务描述
对用户提供的提示词进行全面分析，识别结构特点、潜在问题，并提供专业的改进建议。

## 分析维度

### 1. Subject（主体分析）
- 识别图像的主要主体（人物、动物、物体、风景等）
- 分析主体的详细程度（是否有具体特征描述）
- 评估主体描述的清晰度和可识别性

### 2. Style（艺术风格分析）
- 识别指定的艺术风格（写实、动漫、油画、水彩等）
- 分析风格词汇的准确性和搭配合理性
- 评估风格与主体的匹配度

### 3. Quality（质量相关分析）
- 检查质量词汇的完整性和强度
- 识别分辨率、细节程度相关描述
- 评估质量词的位置（应前置）

### 4. Lighting（光照效果分析）
- 识别光照类型（自然光、人工光、特殊光效）
- 分析光照方向和强度描述
- 评估光照与场景的协调性

### 5. Other（其他要素分析）
- 构图相关词汇（视角、景深、焦距等）
- 色彩和氛围描述
- 背景和环境元素
- 情绪和氛围关键词

## 问题识别（Issues）
检查以下常见问题：
- 缺少质量基础词
- 描述过于模糊或抽象
- 词汇之间存在冲突
- 缺少必要的细节描述
- 标签顺序不合理
- 使用了模型不熟悉的生僻词
- 正负向提示词比例失衡

## 改进建议（Suggestions）
针对发现的问题提供：
- 具体的添加建议
- 替换为更优标签的建议
- 结构调整建议
- 可以探索的变体方向

## 输出格式

{
  "analysis": {
    "subject": "详细的主体分析...",
    "style": "风格分析...",
    "quality": "质量相关分析...",
    "lighting": "光照效果分析...",
    "other": "其他要素分析..."
  },
  "issues": [
    "问题1：缺少质量基础词，建议添加 'masterpiece, best quality'",
    "问题2：'好看'过于模糊，应具体描述喜欢的特征",
    "问题3：光照描述与夜晚场景冲突"
  ],
  "suggestions": [
    "建议1：将质量词移到提示词开头以提升权重",
    "建议2：添加具体的服装颜色描述",
    "建议3：考虑添加特定艺术家风格以获得独特效果",
    "建议4：可以尝试不同光照条件如'golden hour'或'volumetric lighting'"
  ]
}

## 分析示例

输入："一个女孩在森林里，好看一点，赛博朋克风格"

输出：
{
  "analysis": {
    "subject": "主体为单个人物（1girl），但缺少详细特征描述（如发型、服装、姿态等），仅描述了场景位置（森林）",
    "style": "指定了赛博朋克风格（cyberpunk），这是明确的风格方向，但缺少具体的表现元素如霓虹灯、高科技装备等",
    "quality": "完全缺少质量相关词汇，这将导致生成图像质量不稳定，建议添加 'masterpiece, best quality, highres'",
    "lighting": "未指定光照条件，赛博朋克风格通常配合特定的霓虹光照效果，建议补充",
    "other": "缺少视角（如from above, looking at viewer）、缺少背景细节、缺少色彩描述、氛围描述过于主观（'好看一点'）"
  },
  "issues": [
    "严重：缺少质量基础词，将显著影响输出质量",
    "模糊：'好看一点'过于主观，AI无法准确理解具体偏好",
    "不完整：赛博朋克风格需要更多配套元素如'neon lights', 'cityscape', 'holographic'等",
    "顺序：建议将质量词前置，主体描述紧随其后"
  ],
  "suggestions": [
    "添加质量词组：'masterpiece, best quality, ultra detailed, 8k uhd'",
    "具体化主体：添加发型、服装、配饰等细节描述",
    "丰富赛博朋克元素：'neon lights, holographic interface, futuristic city, glowing accessories'",
    "补充光照：'volumetric lighting, cyberpunk lighting, neon glow'",
    "添加色彩：指定主色调如'purple and cyan theme'",
    "考虑添加构图：如'dynamic pose, looking at viewer, depth of field'"
  ]
}`,
    userPromptTemplate: '{{input}}',
    temperature: 0.6,
    responseFormat: 'json',
  },
}

export const useAIStore = defineStore('ai', () => {
  // ========== State ==========
  
  // AI提供商列表
  const providers = ref([])
  
  // 当前选中的AI提供商ID
  const currentProviderId = ref('')
  
  // Prompt模板
  const prompts = ref({})
  
  // 当前选中的Prompt ID
  const currentPromptId = ref('explode')
  
  // 加载状态
  const isLoading = ref(false)
  
  // 上次调用的结果
  const lastResult = ref(null)
  
  // 调用历史
  const callHistory = ref([])

  // ========== Getters ==========
  
  // 当前AI提供商
  const currentProvider = computed(() => {
    return providers.value.find(p => p.id === currentProviderId.value) || providers.value[0]
  })
  
  // 当前Prompt模板
  const currentPrompt = computed(() => {
    return prompts.value[currentPromptId.value] || prompts.value['explode']
  })
  
  // 启用的提供商列表
  const enabledProviders = computed(() => {
    return providers.value.filter(p => p.enabled && p.apiKey)
  })
  
  // 可用的Prompt列表
  const availablePrompts = computed(() => {
    return Object.values(prompts.value)
  })
  
  // 是否已配置AI
  const isConfigured = computed(() => {
    const provider = currentProvider.value
    if (!provider) return false
    if (provider.type === 'ollama') return true // Ollama不需要API Key
    return !!provider.apiKey
  })

  // 当前可用模型列表（用于Ollama）
  const availableModels = computed(() => {
    const provider = currentProvider.value
    if (!provider) return []
    return provider.models || []
  })

  // ========== Actions ==========
  
  // 初始化 - 从localStorage加载并合并默认配置
  function init() {
    loadFromStorage()
    
    // 合并默认提供商（添加缺失的）
    const defaultProviders = JSON.parse(JSON.stringify(DEFAULT_PROVIDERS))
    const existingIds = new Set(providers.value.map(p => p.id))
    
    for (const defaultProvider of defaultProviders) {
      if (!existingIds.has(defaultProvider.id)) {
        providers.value.push(defaultProvider)
      }
    }
    
    // 合并默认Prompt模板（添加缺失的）
    const defaultPrompts = JSON.parse(JSON.stringify(DEFAULT_PROMPTS))
    if (!prompts.value) prompts.value = {}
    
    for (const [id, defaultPrompt] of Object.entries(defaultPrompts)) {
      if (!prompts.value[id]) {
        prompts.value[id] = defaultPrompt
      }
    }
    
    // 默认选中第一个启用的提供商
    if (!currentProviderId.value) {
      const enabled = enabledProviders.value
      currentProviderId.value = enabled.length > 0 ? enabled[0].id : providers.value[0]?.id || ''
    }
    
    // 保存合并后的配置
    saveToStorage()
  }
  
  // 从localStorage加载
  function loadFromStorage() {
    try {
      const savedProviders = localStorage.getItem('pm_ai_providers')
      const savedPrompts = localStorage.getItem('pm_ai_prompts')
      const savedCurrent = localStorage.getItem('pm_ai_current')
      const savedHistory = localStorage.getItem('pm_ai_history')
      
      if (savedProviders) providers.value = JSON.parse(savedProviders)
      if (savedPrompts) prompts.value = JSON.parse(savedPrompts)
      if (savedCurrent) {
        const { providerId, promptId } = JSON.parse(savedCurrent)
        currentProviderId.value = providerId
        currentPromptId.value = promptId
      }
      if (savedHistory) callHistory.value = JSON.parse(savedHistory)
    } catch (e) {
      console.error('加载AI配置失败:', e)
    }
  }
  
  // 保存到localStorage
  function saveToStorage() {
    try {
      localStorage.setItem('pm_ai_providers', JSON.stringify(providers.value))
      localStorage.setItem('pm_ai_prompts', JSON.stringify(prompts.value))
      localStorage.setItem('pm_ai_current', JSON.stringify({
        providerId: currentProviderId.value,
        promptId: currentPromptId.value,
      }))
      localStorage.setItem('pm_ai_history', JSON.stringify(callHistory.value.slice(-50)))
    } catch (e) {
      console.error('保存AI配置失败:', e)
    }
  }
  
  // 设置当前提供商
  function setCurrentProvider(providerId) {
    currentProviderId.value = providerId
    saveToStorage()
  }

  // 获取 Ollama 本地模型列表
  async function fetchOllamaModels(baseUrl = null) {
    const provider = currentProvider.value
    const url = baseUrl || provider?.baseUrl || 'http://localhost:11434'
    
    console.log('[Ollama] 正在获取模型列表:', url)
    
    try {
      // 尝试使用 Wails 的 HTTP 调用（如果可用）
      if (window.runtime && window.runtime.HTTPRequest) {
        console.log('[Ollama] 使用 Wails HTTP 调用')
        const response = await window.runtime.HTTPRequest('GET', `${url}/api/tags`)
        const data = JSON.parse(response)
        const models = data.models?.map(m => m.name) || []
        console.log('[Ollama] 获取成功:', models)
        return models
      }
      
      // 降级到普通 fetch（Wails 桌面应用通常没有 CORS 限制）
      console.log('[Ollama] 使用 fetch 调用')
      const response = await fetch(`${url}/api/tags`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      })
      
      console.log('[Ollama] 响应状态:', response.status, response.statusText)
      
      if (!response.ok) {
        throw new Error(`获取模型列表失败: ${response.status} ${response.statusText}`)
      }
      
      // 先获取原始文本，方便调试
      const rawText = await response.text()
      console.log('[Ollama] 原始响应:', rawText)
      
      // 解析 JSON
      let data
      try {
        data = JSON.parse(rawText)
      } catch (e) {
        console.error('[Ollama] JSON 解析失败:', e)
        throw new Error(`响应解析失败: ${rawText.slice(0, 200)}`)
      }
      
      console.log('[Ollama] 响应数据:', data)
      console.log('[Ollama] models 字段:', data.models)
      console.log('[Ollama] models 类型:', typeof data.models, Array.isArray(data.models))
      
      // Ollama返回的格式: { models: [{ name: 'model-name', ... }, ...] }
      const models = data.models?.map(m => m.name) || []
      console.log('[Ollama] 提取的模型列表:', models)
      return models
    } catch (error) {
      console.error('[Ollama] 获取模型列表失败:', error)
      console.error('[Ollama] 错误详情:', {
        message: error.message,
        name: error.name,
        stack: error.stack,
      })
      
      // 提供更友好的错误信息
      if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
        throw new Error('无法连接到 Ollama 服务，请检查：\n1. Ollama 是否已启动 (ollama serve)\n2. 地址是否正确\n3. 如果是浏览器访问，请设置 OLLAMA_ORIGINS=*')
      }
      
      throw error
    }
  }

  // 刷新当前 Ollama 提供商的模型列表
  async function refreshOllamaModels() {
    const provider = currentProvider.value
    if (!provider || provider.type !== 'ollama') {
      throw new Error('当前提供商不是Ollama')
    }
    
    const models = await fetchOllamaModels(provider.baseUrl)
    
    // 更新提供商的模型列表
    const index = providers.value.findIndex(p => p.id === provider.id)
    if (index !== -1) {
      providers.value[index].models = models
      // 如果当前选中的模型不在列表中，选择第一个
      if (models.length > 0 && !models.includes(providers.value[index].model)) {
        providers.value[index].model = models[0]
      }
      saveToStorage()
    }
    
    return models
  }

  // 设置当前模型
  function setCurrentModel(modelName) {
    const provider = currentProvider.value
    if (!provider) return
    
    const index = providers.value.findIndex(p => p.id === provider.id)
    if (index !== -1) {
      providers.value[index].model = modelName
      saveToStorage()
    }
  }
  
  // 设置当前Prompt
  function setCurrentPrompt(promptId) {
    currentPromptId.value = promptId
    saveToStorage()
  }
  
  // 添加自定义提供商
  function addProvider(provider) {
    const newProvider = {
      id: `custom_${Date.now()}`,
      name: provider.name,
      type: provider.type || 'openai-compatible',
      baseUrl: provider.baseUrl,
      apiKey: provider.apiKey || '',
      model: provider.model,
      models: provider.models || [provider.model],
      headers: provider.headers || {},
      enabled: true,
      isCustom: true,
    }
    providers.value.push(newProvider)
    saveToStorage()
    return newProvider.id
  }
  
  // 更新提供商
  function updateProvider(providerId, updates) {
    const index = providers.value.findIndex(p => p.id === providerId)
    if (index !== -1) {
      providers.value[index] = { ...providers.value[index], ...updates }
      saveToStorage()
    }
  }
  
  // 删除提供商
  function removeProvider(providerId) {
    providers.value = providers.value.filter(p => p.id !== providerId)
    if (currentProviderId.value === providerId) {
      currentProviderId.value = providers.value[0]?.id || ''
    }
    saveToStorage()
  }
  
  // 添加自定义Prompt模板
  function addPrompt(prompt) {
    const id = `custom_${Date.now()}`
    prompts.value[id] = {
      id,
      name: prompt.name,
      description: prompt.description,
      systemPrompt: prompt.systemPrompt,
      userPromptTemplate: prompt.userPromptTemplate || '{{input}}',
      temperature: prompt.temperature || 0.7,
      responseFormat: prompt.responseFormat || 'json',
      isCustom: true,
    }
    saveToStorage()
    return id
  }
  
  // 更新Prompt
  function updatePrompt(promptId, updates) {
    if (prompts.value[promptId]) {
      prompts.value[promptId] = { ...prompts.value[promptId], ...updates }
      saveToStorage()
    }
  }
  
  // 删除Prompt
  function removePrompt(promptId) {
    delete prompts.value[promptId]
    if (currentPromptId.value === promptId) {
      currentPromptId.value = 'explode'
    }
    saveToStorage()
  }
  
  // 重置为默认配置
  function resetToDefaults() {
    providers.value = JSON.parse(JSON.stringify(DEFAULT_PROVIDERS))
    prompts.value = JSON.parse(JSON.stringify(DEFAULT_PROMPTS))
    currentProviderId.value = providers.value[0]?.id || ''
    currentPromptId.value = 'explode'
    saveToStorage()
  }

  // 构建 AI 配置对象
  function buildAIConfig() {
    const provider = currentProvider.value
    if (!provider) return null

    return {
      provider: provider.id,
      provider_type: provider.type || 'openai-compatible',
      api_key: provider.apiKey,
      endpoint: provider.baseUrl,
      model: provider.model,
    }
  }
  
  // 调用AI - 根据当前模式调用不同功能
  async function callAI(input, mode = null) {
    const promptMode = mode || currentPromptId.value
    const promptTemplate = prompts.value[promptMode] || prompts.value['explode']
    const config = buildAIConfig()
    
    if (!config) {
      throw new Error('请先配置AI提供商')
    }
    
    isLoading.value = true
    lastResult.value = null
    
    try {
      let response
      
      // 构建模板参数
      const templateParams = {
        prompt: input,
        config: config,
        system_prompt: promptTemplate?.systemPrompt || '',
        user_prompt_template: promptTemplate?.userPromptTemplate || '{{input}}',
      }
      
      switch (promptMode) {
        case 'explode':
          response = await ExplodePrompt(templateParams)
          break
          
        case 'optimize':
          response = await OptimizePrompt(templateParams)
          break
          
        case 'translate':
          response = await TranslatePrompt(templateParams)
          break
          
        case 'analyze':
          response = await AnalyzePrompt(templateParams)
          break
          
        default:
          // 使用通用接口
          response = await ProcessAI({
            mode: promptMode,
            ...templateParams,
          })
      }
      
      if (!response.success) {
        throw new Error(response.error || 'AI调用失败')
      }
      
      // 记录历史
      callHistory.value.unshift({
        id: Date.now(),
        mode: promptMode,
        provider: currentProvider.value?.name,
        input: input.slice(0, 200),
        timestamp: new Date().toISOString(),
      })
      saveToStorage()
      
      lastResult.value = response.data
      return response.data
    } finally {
      isLoading.value = false
    }
  }

  // 拆解提示词
  // categories: 分类名称数组
  // categoryMap: 分类名称到ID的映射对象 {name: id}
  async function explodePrompt(prompt, categories = [], categoryMap = {}) {
    const promptTemplate = prompts.value['explode']
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await ExplodePrompt({
        prompt: prompt,
        categories: categories,
        category_map: categoryMap,
        config: config,
        system_prompt: promptTemplate?.systemPrompt || '',
        user_prompt_template: promptTemplate?.userPromptTemplate || '{{input}}',
      })
      
      if (!response.success) {
        throw new Error(response.error || '拆解失败')
      }
      
      return response.data
    } finally {
      isLoading.value = false
    }
  }

  // 导入拆解的原子词
  async function importExplodedAtoms(result, categoryId) {
    const response = await ImportExtractedAtoms({
      result: result,
      category_id: categoryId,
    })
    
    if (!response.success) {
      throw new Error(response.error || '导入失败')
    }
    
    return response.data
  }

  // 优化提示词
  async function optimizePrompt(prompt) {
    const promptTemplate = prompts.value['optimize']
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await OptimizePrompt({
        prompt: prompt,
        config: config,
        system_prompt: promptTemplate?.systemPrompt || '',
        user_prompt_template: promptTemplate?.userPromptTemplate || '{{input}}',
      })
      
      if (!response.success) {
        throw new Error(response.error || '优化失败')
      }
      
      return response.data
    } finally {
      isLoading.value = false
    }
  }

  // 翻译提示词
  async function translatePrompt(prompt) {
    const promptTemplate = prompts.value['translate']
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await TranslatePrompt({
        prompt: prompt,
        config: config,
        system_prompt: promptTemplate?.systemPrompt || '',
        user_prompt_template: promptTemplate?.userPromptTemplate || '{{input}}',
      })
      
      if (!response.success) {
        throw new Error(response.error || '翻译失败')
      }
      
      return response.data
    } finally {
      isLoading.value = false
    }
  }

  // 分析提示词
  async function analyzePrompt(prompt) {
    const promptTemplate = prompts.value['analyze']
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await AnalyzePrompt({
        prompt: prompt,
        config: config,
        system_prompt: promptTemplate?.systemPrompt || '',
        user_prompt_template: promptTemplate?.userPromptTemplate || '{{input}}',
      })
      
      if (!response.success) {
        throw new Error(response.error || '分析失败')
      }
      
      return response.data
    } finally {
      isLoading.value = false
    }
  }
  
  // 清空历史
  function clearHistory() {
    callHistory.value = []
    saveToStorage()
  }

  return {
    // State
    providers,
    currentProviderId,
    prompts,
    currentPromptId,
    isLoading,
    lastResult,
    callHistory,
    
    // Getters
    currentProvider,
    currentPrompt,
    enabledProviders,
    availablePrompts,
    isConfigured,
    availableModels,
    
    // Actions
    init,
    setCurrentProvider,
    setCurrentPrompt,
    setCurrentModel,
    addProvider,
    updateProvider,
    removeProvider,
    addPrompt,
    updatePrompt,
    removePrompt,
    resetToDefaults,
    callAI,
    explodePrompt,
    importExplodedAtoms,
    optimizePrompt,
    translatePrompt,
    analyzePrompt,
    clearHistory,
    fetchOllamaModels,
    refreshOllamaModels,
  }
})
