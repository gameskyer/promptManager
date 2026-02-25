<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>
          <BoltIcon class="w-5 h-5 text-amber-400" />
          AI 智能助手
        </h3>
        <div class="header-actions">
          <button 
            class="config-btn" 
            @click="$emit('open-settings')"
            title="配置 AI API"
          >
            <Cog6ToothIcon class="w-4 h-4" />
            配置
          </button>
          <button class="close-btn" @click="$emit('close')">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </div>
      
      <div class="modal-body">
        <!-- 左侧面板：输入和配置 -->
        <div class="left-panel">
          <!-- AI 和 Prompt 选择 -->
          <div class="selectors">
            <div class="selector">
              <label>AI 提供商</label>
              <select v-model="selectedProviderId" @change="onProviderChange">
                <option 
                  v-for="provider in aiStore.providers" 
                  :key="provider.id" 
                  :value="provider.id"
                  :disabled="!provider.enabled"
                >
                  {{ provider.name }} {{ !provider.enabled ? '(未配置)' : '' }}
                </option>
              </select>
            </div>
            
            <div class="selector">
              <label>功能模式</label>
              <select v-model="selectedPromptId" @change="onPromptChange">
                <option 
                  v-for="prompt in aiStore.availablePrompts" 
                  :key="prompt.id" 
                  :value="prompt.id"
                >
                  {{ prompt.name }}
                </option>
              </select>
            </div>
          </div>
          
          <!-- 模型选择（仅Ollama显示） -->
          <div v-if="isOllamaProvider" class="model-selector">
            <div class="selector">
              <label>选择模型</label>
              <div class="model-select-wrapper">
                <select 
                  v-model="selectedModel" 
                  @change="onModelChange"
                  :disabled="aiStore.availableModels.length === 0 || isLoadingModels"
                >
                  <option v-if="isLoadingModels" value="">加载中...</option>
                  <option v-else-if="aiStore.availableModels.length === 0" value="">暂无可用模型</option>
                  <option 
                    v-for="model in aiStore.availableModels" 
                    :key="model" 
                    :value="model"
                  >
                    {{ model }}
                  </option>
                </select>
                <button 
                  class="refresh-btn" 
                  @click="refreshModels"
                  :disabled="isLoadingModels"
                  title="刷新模型列表"
                >
                  <ArrowPathIcon v-if="!isLoadingModels" class="w-4 h-4" />
                  <div v-else class="spinner-small"></div>
                </button>
              </div>
            </div>
            <div v-if="modelError" class="model-error">
              <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0" />
              <pre class="error-text">{{ modelError }}</pre>
            </div>
          </div>
          
          <!-- 当前 Prompt 描述 -->
          <div v-if="currentPrompt" class="prompt-desc">
            <SparklesIcon class="w-4 h-4" />
            {{ currentPrompt.description }}
          </div>
          
          <!-- 输入区域 -->
          <div class="input-section">
            <label>输入提示词</label>
            <textarea
              v-model="inputPrompt"
              :placeholder="inputPlaceholder"
              rows="8"
            ></textarea>
            
            <!-- 错误提示 -->
            <div v-if="error" class="error-message">
              {{ error }}
            </div>
            
            <button 
              class="explode-btn"
              :disabled="!inputPrompt.trim() || aiStore.isLoading || !aiStore.isConfigured"
              @click="handleSubmit"
            >
              <SparklesIcon v-if="!aiStore.isLoading" class="w-4 h-4" />
              <div v-else class="spinner"></div>
              {{ aiStore.isLoading ? '处理中...' : '开始处理' }}
            </button>
            
            <div v-if="!aiStore.isConfigured" class="warning-message">
              <ExclamationTriangleIcon class="w-4 h-4" />
              请先在设置中配置 AI API
            </div>
          </div>
        </div>
        
        <!-- 右侧面板：结果展示 -->
        <div class="right-panel">
          <div v-if="!result" class="empty-result">
            <SparklesIcon class="w-12 h-12 text-slate-600" />
            <p>在左侧输入内容并开始处理</p>
            <p class="empty-hint">AI 处理结果将在此处显示</p>
          </div>
          
          <!-- 结果区域 -->
          <div v-else class="result-section">
            <div class="result-header">
              <span>处理结果</span>
              <div class="result-actions">
                <button class="action-btn" @click="copyResult">
                  <ClipboardDocumentIcon class="w-4 h-4" />
                  复制
                </button>
                <button v-if="canImport" class="action-btn import" @click="handleImport">
                  <ArrowDownTrayIcon class="w-4 h-4" />
                  导入
                </button>
              </div>
            </div>
            
            <!-- 拆解结果 -->
            <div v-if="result.atoms" class="atoms-result">
              <div class="atoms-header">
                <span>共 {{ result.atoms.length }} 个原子词</span>
                <div class="category-legend">
                  <span class="legend-item new">新词</span>
                  <span class="legend-item existing">已存在</span>
                </div>
              </div>
              <div
                v-for="(atom, index) in result.atoms"
                :key="index"
                class="atom-item"
                :class="{ 'is-new': atom.is_new }"
              >
                <div class="atom-main">
                  <div class="atom-info">
                    <div class="atom-value">{{ atom.value }}</div>
                    <div class="atom-label">{{ atom.label }}</div>
                  </div>
                  <div class="atom-meta">
                    <span v-if="atom.type" class="atom-type" :class="atom.type">{{ atom.type }}</span>
                  </div>
                </div>
                <div class="atom-category-row">
                  <label>分类:</label>
                  <!-- 可搜索的分类选择器 -->
                  <div class="searchable-select">
                    <input
                      v-model="atom.categorySearch"
                      type="text"
                      class="category-search-input"
                      placeholder="输入分类名称搜索..."
                      @focus="atom.showCategoryDropdown = true"
                      @input="filterCategories(atom)"
                    />
                    <div v-if="atom.showCategoryDropdown" class="category-dropdown">
                      <div
                        v-for="cat in getFilteredCategories(atom)"
                        :key="cat.id"
                        class="category-option"
                        :class="{ 'is-selected': atom.category === cat.id }"
                        @click="selectCategory(atom, cat.id)"
                      >
                        <span class="cat-name">{{ cat.name }}</span>
                        <span v-if="cat.parent_name" class="cat-parent">{{ cat.parent_name }}</span>
                      </div>
                      <div v-if="getFilteredCategories(atom).length === 0" class="no-category">
                        无匹配分类
                      </div>
                    </div>
                  </div>
                  <span class="category-name">{{ getCategoryNameById(atom.category) }}</span>
                </div>
                <div v-if="atom.synonyms?.length" class="atom-synonyms">
                  <span class="synonyms-label">近义词:</span>
                  <span class="synonyms-list">{{ atom.synonyms.join(', ') }}</span>
                </div>
              </div>
            </div>
            
            <!-- 优化结果 -->
            <div v-else-if="result.optimized" class="text-result">
              <div class="result-block">
                <label>优化后：</label>
                <div class="result-text">{{ result.optimized }}</div>
              </div>
              <div v-if="result.changes?.length" class="result-block">
                <label>修改说明：</label>
                <ul>
                  <li v-for="(change, i) in result.changes" :key="i">{{ change }}</li>
                </ul>
              </div>
              <button class="apply-btn" @click="applyToPreset">
                <CheckIcon class="w-4 h-4" />
                应用到当前预设
              </button>
            </div>
            
            <!-- 翻译结果 -->
            <div v-else-if="result.translation" class="text-result">
              <div class="result-block">
                <label>翻译结果：</label>
                <div class="result-text">{{ result.translation }}</div>
              </div>
              <div v-if="result.keywords?.length" class="result-block">
                <label>关键词：</label>
                <div class="keyword-tags">
                  <span v-for="kw in result.keywords" :key="kw" class="keyword-tag">{{ kw }}</span>
                </div>
              </div>
            </div>
            
            <!-- 分析结果 -->
            <div v-else-if="result.analysis" class="text-result">
              <div class="analysis-grid">
                <div v-for="(value, key) in result.analysis" :key="key" class="analysis-item">
                  <label>{{ getAnalysisLabel(key) }}：</label>
                  <span>{{ value }}</span>
                </div>
              </div>
              <div v-if="result.issues?.length" class="result-block warning">
                <label>⚠️ 问题：</label>
                <ul>
                  <li v-for="(issue, i) in result.issues" :key="i">{{ issue }}</li>
                </ul>
              </div>
              <div v-if="result.suggestions?.length" class="result-block">
                <label>💡 建议：</label>
                <ul>
                  <li v-for="(s, i) in result.suggestions" :key="i">{{ s }}</li>
                </ul>
              </div>
            </div>
            
            <!-- 通用结果 -->
            <div v-else-if="result.text" class="text-result">
              <pre>{{ result.text }}</pre>
            </div>
            
            <!-- 原始JSON -->
            <details class="raw-result">
              <summary>查看原始响应</summary>
              <pre>{{ JSON.stringify(result, null, 2) }}</pre>
            </details>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import {
  BoltIcon,
  XMarkIcon,
  SparklesIcon,
  ArrowDownTrayIcon,
  Cog6ToothIcon,
  ExclamationTriangleIcon,
  ClipboardDocumentIcon,
  CheckIcon,
  ArrowPathIcon,
} from '@heroicons/vue/24/outline'
import { useAIStore, useCategoryStore } from '../stores'

const aiStore = useAIStore()
const categoryStore = useCategoryStore()
const { currentProvider, currentPrompt, isConfigured } = storeToRefs(aiStore)

const emit = defineEmits(['close', 'import', 'open-settings', 'apply-to-preset'])

// 获取可用的原子词子分类列表（排除根分类）
const availableCategories = computed(() => {
  return categoryStore.atomCategories.filter(c => c.parent_id > 0)
})

// 获取分类名称列表（用于AI提示）
const categoryNames = computed(() => {
  return availableCategories.value.map(c => c.name)
})

// 构建分类名称到ID的映射
const categoryMap = computed(() => {
  const map = {}
  availableCategories.value.forEach(c => {
    map[c.name] = c.id
  })
  return map
})

// 根据分类ID获取分类名称
function getCategoryNameById(categoryId) {
  const category = availableCategories.value.find(c => c.id === categoryId)
  return category?.name || '未分类'
}

// ========== 可搜索分类选择器相关方法 ==========

// 获取带父分类名称的分类列表
const categoriesWithParent = computed(() => {
  return availableCategories.value.map(cat => {
    const parent = categoryStore.categories.find(c => c.id === cat.parent_id)
    return {
      ...cat,
      parent_name: parent?.name || ''
    }
  })
})

// 筛选分类（根据搜索关键词）
function getFilteredCategories(atom) {
  const searchTerm = (atom.categorySearch || '').toLowerCase().trim()
  if (!searchTerm) {
    return categoriesWithParent.value
  }
  return categoriesWithParent.value.filter(cat => 
    cat.name.toLowerCase().includes(searchTerm)
  )
}

// 选择分类
function selectCategory(atom, categoryId) {
  atom.category = categoryId
  atom.showCategoryDropdown = false
  // 设置搜索框显示为选中的分类名
  const cat = categoriesWithParent.value.find(c => c.id === categoryId)
  if (cat) {
    atom.categorySearch = cat.name
  }
}

// 输入时筛选
function filterCategories(atom) {
  atom.showCategoryDropdown = true
}

// 点击外部关闭下拉框
function handleClickOutside(event) {
  const dropdowns = document.querySelectorAll('.searchable-select')
  dropdowns.forEach((dropdown, index) => {
    if (!dropdown.contains(event.target)) {
      const atom = result.value?.atoms?.[index]
      if (atom) {
        atom.showCategoryDropdown = false
        // 如果没有选择分类，清空搜索框
        if (!atom.category) {
          atom.categorySearch = ''
        } else {
          // 恢复为已选分类的名称
          const cat = categoriesWithParent.value.find(c => c.id === atom.category)
          if (cat) {
            atom.categorySearch = cat.name
          }
        }
      }
    }
  })
}

// 本地状态
const inputPrompt = ref('')
const error = ref(null)
const result = ref(null)
const selectedProviderId = ref('')
const selectedPromptId = ref('')
const selectedModel = ref('')
const isLoadingModels = ref(false)
const modelError = ref(null)

// 计算属性：当前是否是Ollama提供商
const isOllamaProvider = computed(() => {
  return currentProvider.value?.type === 'ollama'
})

// 输入占位符
const inputPlaceholder = computed(() => {
  if (currentPrompt.value?.id === 'explode') {
    return '粘贴你的长段提示词，AI 会自动拆解为原子词...'
  } else if (currentPrompt.value?.id === 'optimize') {
    return '输入需要优化的提示词...'
  } else if (currentPrompt.value?.id === 'translate') {
    return '输入中文提示词，翻译成英文...'
  } else if (currentPrompt.value?.id === 'analyze') {
    return '输入提示词进行分析...'
  }
  return '输入内容...'
})

// 是否可以导入（有atoms结果时）
const canImport = computed(() => {
  return result.value?.atoms && result.value.atoms.length > 0
})

// 分析字段标签
function getAnalysisLabel(key) {
  const labels = {
    subject: '主体',
    style: '风格',
    quality: '质量',
    lighting: '光照',
    other: '其他',
  }
  return labels[key] || key
}

// 初始化
onMounted(async () => {
  aiStore.init()
  selectedProviderId.value = aiStore.currentProviderId
  selectedPromptId.value = aiStore.currentPromptId
  selectedModel.value = currentProvider.value?.model || ''
  
  // 确保分类数据已加载
  if (categoryStore.atomCategories.length === 0) {
    await categoryStore.fetchCategories()
  }
  
  // 如果是Ollama提供商，自动获取模型列表
  if (isOllamaProvider.value) {
    await loadOllamaModels()
  }
  
  // 添加点击外部关闭下拉框的事件监听
  document.addEventListener('click', handleClickOutside)
})

// 组件卸载时移除事件监听
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 监听store变化
watch(() => aiStore.currentProviderId, (val) => {
  selectedProviderId.value = val
})

watch(() => aiStore.currentPromptId, (val) => {
  selectedPromptId.value = val
})

// 切换提供商
async function onProviderChange() {
  aiStore.setCurrentProvider(selectedProviderId.value)
  selectedModel.value = currentProvider.value?.model || ''
  modelError.value = null
  
  // 如果切换到Ollama，自动获取模型列表
  if (isOllamaProvider.value) {
    await loadOllamaModels()
  }
}

// 切换模型
function onModelChange() {
  aiStore.setCurrentModel(selectedModel.value)
}

// 加载Ollama模型列表
async function loadOllamaModels() {
  isLoadingModels.value = true
  modelError.value = null
  try {
    await aiStore.refreshOllamaModels()
    selectedModel.value = currentProvider.value?.model || ''
  } catch (err) {
    modelError.value = err.message || '获取模型列表失败'
    console.error('获取Ollama模型列表失败:', err)
  } finally {
    isLoadingModels.value = false
  }
}

// 刷新模型列表
async function refreshModels() {
  await loadOllamaModels()
}

// 切换Prompt
function onPromptChange() {
  aiStore.setCurrentPrompt(selectedPromptId.value)
  result.value = null // 清空结果
}

// 提交处理
async function handleSubmit() {
  if (!inputPrompt.value.trim()) return
  
  error.value = null
  result.value = null
  
  try {
    // 如果是拆解模式，传入分类数据
    if (currentPrompt.value?.id === 'explode') {
      result.value = await aiStore.explodePrompt(
        inputPrompt.value, 
        categoryNames.value,
        categoryMap.value
      )
      
      // 初始化每个原子词的搜索状态
      if (result.value?.atoms) {
        result.value.atoms.forEach(atom => {
          const cat = categoriesWithParent.value.find(c => c.id === atom.category)
          atom.categorySearch = cat?.name || ''
          atom.showCategoryDropdown = false
        })
      }
    } else {
      result.value = await aiStore.callAI(inputPrompt.value)
    }
  } catch (err) {
    error.value = err.message || 'AI 调用失败'
    console.error('AI调用错误:', err)
  }
}

// 复制结果
function copyResult() {
  const text = result.value.optimized || 
               result.value.translation || 
               JSON.stringify(result.value, null, 2)
  navigator.clipboard.writeText(text).then(() => {
    // 可以显示一个toast提示
  })
}

// 导入原子词
function handleImport() {
  emit('import', result.value)
  emit('close')
}

// 应用到当前预设
function applyToPreset() {
  if (!result.value?.optimized) {
    alert('没有可应用的优化结果')
    return
  }
  emit('apply-to-preset', result.value)
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  width: 100%;
  max-width: 1200px;
  max-height: 85vh;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #334155;
  flex-shrink: 0;
}

.modal-header h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.config-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #94a3b8;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.config-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.close-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

/* 左右布局主体 */
.modal-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
}

/* 左侧面板 */
.left-panel {
  width: 45%;
  min-width: 380px;
  max-width: 480px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  border-right: 1px solid #334155;
}

/* 右侧面板 */
.right-panel {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #0b1120;
}

/* 空结果提示 */
.empty-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #64748b;
  gap: 12px;
}

.empty-result p {
  font-size: 15px;
  margin: 0;
}

.empty-result .empty-hint {
  font-size: 13px;
  color: #475569;
}

/* Selectors */
.selectors {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.selector {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.selector label {
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
}

.selector select {
  padding: 8px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 14px;
  cursor: pointer;
}

.selector select:focus {
  outline: none;
  border-color: #0ea5e9;
}

.selector select option:disabled {
  color: #64748b;
}

/* Model Selector */
.model-selector {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.model-select-wrapper {
  display: flex;
  gap: 8px;
}

.model-select-wrapper select {
  flex: 1;
}

.refresh-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background-color: #334155;
  color: #e2e8f0;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner-small {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.model-error {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12px;
  color: #ef4444;
  padding: 8px 12px;
  background-color: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 6px;
  margin-top: 8px;
}

.model-error .error-text {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.5;
}

.prompt-desc {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background-color: rgba(124, 58, 237, 0.1);
  border: 1px solid rgba(124, 58, 237, 0.3);
  border-radius: 8px;
  font-size: 13px;
  color: #a78bfa;
}

/* Input Section */
.input-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}

.input-section label {
  font-size: 13px;
  font-weight: 500;
  color: #94a3b8;
}

.input-section textarea {
  width: 100%;
  flex: 1;
  min-height: 120px;
  padding: 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 14px;
  resize: none;
}

.input-section textarea:focus {
  border-color: #0ea5e9;
  outline: none;
}

.input-section textarea::placeholder {
  color: #64748b;
}

.error-message {
  padding: 10px 12px;
  background-color: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 6px;
  color: #ef4444;
  font-size: 13px;
}

.warning-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background-color: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: 6px;
  color: #f59e0b;
  font-size: 13px;
}

.explode-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  background-color: #7c3aed;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: auto;
}

.explode-btn:hover:not(:disabled) {
  background-color: #8b5cf6;
}

.explode-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Result Section */
.result-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px solid #334155;
}

.result-header span {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.result-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #94a3b8;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.action-btn.import {
  background-color: #22c55e;
  border-color: #22c55e;
  color: white;
}

.action-btn.import:hover {
  background-color: #16a34a;
}

/* Atoms Result */
.atoms-result {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: calc(85vh - 200px);
  overflow-y: auto;
}

.atoms-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background-color: #1e293b;
  border-radius: 6px;
  font-size: 13px;
  color: #94a3b8;
}

.category-legend {
  display: flex;
  gap: 12px;
}

.legend-item {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
}

.legend-item.new {
  background-color: rgba(124, 58, 237, 0.2);
  color: #a78bfa;
}

.legend-item.existing {
  background-color: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.atom-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  transition: all 0.2s;
}

.atom-item:hover {
  border-color: #475569;
}

.atom-item.is-new {
  border-color: rgba(124, 58, 237, 0.5);
  background-color: rgba(124, 58, 237, 0.1);
}

.atom-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.atom-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.atom-value {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
}

.atom-label {
  font-size: 12px;
  color: #64748b;
}

.atom-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.atom-category-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.atom-category-row label {
  color: #64748b;
  white-space: nowrap;
}

.category-select {
  flex: 1;
  padding: 4px 8px;
  background-color: #0f172a;
  border: 1px solid #475569;
  border-radius: 4px;
  color: #e2e8f0;
  font-size: 12px;
  cursor: pointer;
}

.category-select:focus {
  border-color: #0ea5e9;
  outline: none;
}

.category-name {
  color: #94a3b8;
  font-size: 11px;
  padding: 2px 8px;
  background-color: rgba(148, 163, 184, 0.1);
  border-radius: 4px;
}

/* 可搜索分类选择器 */
.searchable-select {
  flex: 1;
  position: relative;
}

.category-search-input {
  width: 100%;
  padding: 4px 8px;
  background-color: #0f172a;
  border: 1px solid #475569;
  border-radius: 4px;
  color: #e2e8f0;
  font-size: 12px;
  outline: none;
  transition: all 0.2s;
}

.category-search-input:focus {
  border-color: #0ea5e9;
}

.category-search-input::placeholder {
  color: #64748b;
}

.category-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  max-height: 200px;
  overflow-y: auto;
  background-color: #1e293b;
  border: 1px solid #475569;
  border-radius: 4px;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.category-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.category-option:hover {
  background-color: #334155;
}

.category-option.is-selected {
  background-color: rgba(14, 165, 233, 0.2);
}

.category-option .cat-name {
  color: #e2e8f0;
  font-size: 12px;
}

.category-option .cat-parent {
  color: #64748b;
  font-size: 10px;
}

.no-category {
  padding: 12px;
  text-align: center;
  color: #64748b;
  font-size: 12px;
  white-space: nowrap;
}

.atom-synonyms {
  font-size: 11px;
  color: #94a3b8;
}

.synonyms-label {
  color: #64748b;
  margin-right: 4px;
}

.synonyms-list {
  color: #94a3b8;
}

.atom-type {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
}

.atom-type.Positive {
  color: #22c55e;
  background-color: rgba(34, 197, 94, 0.2);
}

.atom-type.Negative {
  color: #ef4444;
  background-color: rgba(239, 68, 68, 0.2);
}

/* Text Result */
.text-result {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-block {
  padding: 12px;
  background-color: #1e293b;
  border-radius: 8px;
}

.result-block.warning {
  background-color: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.result-block label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
  margin-bottom: 8px;
}

.result-text {
  font-size: 14px;
  color: #e2e8f0;
  line-height: 1.6;
}

.result-block ul {
  margin: 0;
  padding-left: 20px;
}

.result-block li {
  font-size: 13px;
  color: #94a3b8;
  margin-bottom: 4px;
}

.keyword-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.keyword-tag {
  font-size: 12px;
  color: #0ea5e9;
  padding: 4px 8px;
  background-color: rgba(14, 165, 233, 0.1);
  border-radius: 4px;
}

/* Analysis Grid */
.analysis-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.analysis-item {
  padding: 10px;
  background-color: #1e293b;
  border-radius: 6px;
}

.analysis-item label {
  font-size: 11px;
  color: #64748b;
}

.analysis-item span {
  display: block;
  font-size: 13px;
  color: #e2e8f0;
  margin-top: 4px;
}

/* Raw Result */
.raw-result {
  margin-top: 8px;
}

.raw-result summary {
  font-size: 12px;
  color: #64748b;
  cursor: pointer;
  padding: 8px 0;
}

.raw-result pre {
  padding: 12px;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 6px;
  font-size: 12px;
  color: #94a3b8;
  overflow-x: auto;
  max-height: 200px;
  overflow-y: auto;
}

/* Apply to Preset Button */
.apply-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  margin-top: 12px;
  background-color: #22c55e;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.apply-btn:hover {
  background-color: #16a34a;
}

/* 响应式适配 */
@media (max-width: 900px) {
  .modal-content {
    max-width: 95vw;
  }
  
  .modal-body {
    flex-direction: column;
  }
  
  .left-panel {
    width: 100%;
    max-width: none;
    max-height: 50vh;
    border-right: none;
    border-bottom: 1px solid #334155;
  }
  
  .right-panel {
    max-height: 35vh;
  }
}
</style>
