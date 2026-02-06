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
    systemPrompt: ``,
    userPromptTemplate: '{{input}}',
    temperature: 0.7,
    responseFormat: 'json',
  },
  optimize: {
    id: 'optimize',
    name: '优化提示词',
    description: '优化提示词的质量和表达',
    systemPrompt: ``,
    userPromptTemplate: '{{input}}',
    temperature: 0.8,
    responseFormat: 'json',
  },
  translate: {
    id: 'translate',
    name: '翻译提示词',
    description: '将中文提示词翻译为英文',
    systemPrompt: ``,
    userPromptTemplate: '{{input}}',
    temperature: 0.5,
    responseFormat: 'json',
  },
  analyze: {
    id: 'analyze',
    name: '分析提示词',
    description: '分析提示词的结构和效果',
    systemPrompt: ``,
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
      api_key: provider.apiKey,
      endpoint: provider.baseUrl,
      model: provider.model,
    }
  }
  
  // 调用AI - 根据当前模式调用不同功能
  async function callAI(input, mode = null) {
    const promptMode = mode || currentPromptId.value
    const config = buildAIConfig()
    
    if (!config) {
      throw new Error('请先配置AI提供商')
    }
    
    isLoading.value = true
    lastResult.value = null
    
    try {
      let response
      
      switch (promptMode) {
        case 'explode':
          response = await ExplodePrompt({
            prompt: input,
            config: config,
          })
          break
          
        case 'optimize':
          response = await OptimizePrompt({
            prompt: input,
            config: config,
          })
          break
          
        case 'translate':
          response = await TranslatePrompt({
            prompt: input,
            config: config,
          })
          break
          
        case 'analyze':
          response = await AnalyzePrompt({
            prompt: input,
            config: config,
          })
          break
          
        default:
          // 使用通用接口
          response = await ProcessAI({
            mode: promptMode,
            prompt: input,
            config: config,
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
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await ExplodePrompt({
        prompt: prompt,
        categories: categories,
        category_map: categoryMap,
        config: config,
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
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await OptimizePrompt({
        prompt: prompt,
        config: config,
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
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await TranslatePrompt({
        prompt: prompt,
        config: config,
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
    const config = buildAIConfig()
    
    isLoading.value = true
    try {
      const response = await AnalyzePrompt({
        prompt: prompt,
        config: config,
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
    
    // Actions
    init,
    setCurrentProvider,
    setCurrentPrompt,
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
  }
})
