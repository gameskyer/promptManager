<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>
          <Cog6ToothIcon class="w-5 h-5" />
          设置
        </h3>
        <button class="close-btn" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      
      <!-- 设置标签页 -->
      <div class="setting-tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'general' }"
          @click="activeTab = 'general'"
        >
          常规
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'ai-providers' }"
          @click="activeTab = 'ai-providers'"
        >
          <CpuChipIcon class="w-4 h-4" />
          AI 提供商
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'ai-prompts' }"
          @click="activeTab = 'ai-prompts'"
        >
          <DocumentTextIcon class="w-4 h-4" />
          Prompt 模板
        </button>
      </div>
      
      <div class="modal-body">
        <!-- 常规设置 -->
        <template v-if="activeTab === 'general'">
          <!-- 主题设置 -->
          <div class="setting-section">
            <div class="section-title">外观</div>
            <div class="setting-item">
              <span class="item-label">主题</span>
              <select v-model="localTheme" class="item-select">
                <option value="dark">深色</option>
                <option value="light">浅色</option>
                <option value="auto">跟随系统</option>
              </select>
            </div>
          </div>
          
          <!-- 快捷键设置 -->
          <div class="setting-section">
            <div class="section-title">快捷键</div>
            <div class="setting-item">
              <span class="item-label">保存</span>
              <kbd class="item-shortcut">Ctrl + S</kbd>
            </div>
            <div class="setting-item">
              <span class="item-label">搜索</span>
              <kbd class="item-shortcut">Ctrl + F</kbd>
            </div>
          </div>
          
          <!-- 关于 -->
          <div class="setting-section">
            <div class="section-title">关于</div>
            <div class="about-info">
              <div class="app-name">PromptMaster</div>
              <div class="app-version">v2.1.0</div>
              <div class="app-desc">AI绘画提示词管理系统</div>
            </div>
          </div>
        </template>
        
        <!-- AI 提供商配置 -->
        <template v-if="activeTab === 'ai-providers'">
          <div class="setting-section">
            <div class="section-header">
              <div class="section-title">已配置的 AI 提供商 ({{ providers.length }})</div>
              <button class="add-btn" @click="showAddProvider = true">
                <PlusIcon class="w-4 h-4" />
                添加
              </button>
            </div>
            
            <!-- 提供商列表 -->
            <div class="provider-list">
              <div 
                v-for="provider in providers" 
                :key="provider.id"
                class="provider-card"
                :class="{ 
                  active: provider.id === currentProviderId,
                  disabled: !provider.enabled 
                }"
                @click="selectProvider(provider.id)"
              >
                <div class="provider-header">
                  <div class="provider-title">
                    <span class="provider-name">{{ provider.name }}</span>
                    <span v-if="provider.isCustom" class="custom-badge">自定义</span>
                    <span v-if="!provider.enabled" class="disabled-badge">未配置</span>
                  </div>
                  <div class="provider-actions">
                    <button 
                      class="icon-btn"
                      @click.stop="editProvider(provider)"
                      title="编辑"
                    >
                      <PencilIcon class="w-4 h-4" />
                    </button>
                    <button 
                      v-if="provider.isCustom"
                      class="icon-btn danger"
                      @click.stop="deleteProvider(provider.id)"
                      title="删除"
                    >
                      <TrashIcon class="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <div class="provider-info-row">
                  <span class="info-label">模型：</span>
                  <span class="info-value">{{ provider.model }}</span>
                </div>
                <div class="provider-info-row">
                  <span class="info-label">地址：</span>
                  <span class="info-value">{{ provider.baseUrl }}</span>
                </div>
              </div>
            </div>
            
            <!-- 编辑提供商弹窗 -->
            <div v-if="editingProvider" class="edit-modal">
              <div class="edit-content">
                <h4>{{ editingProvider.id ? '编辑' : '添加' }} AI 提供商</h4>
                
                <div class="form-group">
                  <label>名称</label>
                  <input v-model="editingProvider.name" placeholder="如：DeepSeek" />
                </div>
                
                <div class="form-group">
                  <label>类型</label>
                  <select v-model="editingProvider.type">
                    <option value="openai-compatible">OpenAI 兼容</option>
                    <option value="ollama">Ollama</option>
                  </select>
                </div>
                
                <div class="form-group">
                  <label>API 地址</label>
                  <input v-model="editingProvider.baseUrl" placeholder="https://api.xxx.com/v1" />
                </div>
                
                <div class="form-group">
                  <label>API Key {{ editingProvider.type === 'ollama' ? '(可选)' : '' }}</label>
                  <input 
                    v-model="editingProvider.apiKey" 
                    type="password"
                    placeholder="sk-..."
                  />
                </div>
                
                <div class="form-group">
                  <label>默认模型</label>
                  <input v-model="editingProvider.model" placeholder="模型名称" />
                </div>
                
                <div class="form-group">
                  <label>可用模型（逗号分隔）</label>
                  <input v-model="editingProvider.modelsText" placeholder="model1, model2, model3" />
                </div>
                
                <div class="form-actions">
                  <button class="btn-secondary" @click="editingProvider = null">取消</button>
                  <button class="btn-primary" @click="saveProvider">保存</button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 重置按钮 -->
          <div class="setting-actions">
            <button class="btn-secondary" @click="resetAIConfig">
              <ArrowPathIcon class="w-4 h-4" />
              重置为默认
            </button>
          </div>
        </template>
        
        <!-- Prompt 模板配置 -->
        <template v-if="activeTab === 'ai-prompts'">
          <div class="setting-section">
            <div class="section-header">
              <div class="section-title">Prompt 模板 ({{ availablePrompts.length }})</div>
              <button class="add-btn" @click="showAddPrompt = true">
                <PlusIcon class="w-4 h-4" />
                添加
              </button>
            </div>
            
            <!-- Prompt 列表 -->
            <div class="prompt-list">
              <div 
                v-for="prompt in availablePrompts" 
                :key="prompt.id"
                class="prompt-card"
                :class="{ active: prompt.id === currentPromptId }"
                @click="selectPrompt(prompt.id)"
              >
                <div class="prompt-header">
                  <div class="prompt-title">
                    <span class="prompt-name">{{ prompt.name }}</span>
                    <span v-if="prompt.isCustom" class="custom-badge">自定义</span>
                  </div>
                  <div class="prompt-actions">
                    <button 
                      class="icon-btn"
                      @click.stop="editPrompt(prompt)"
                      title="编辑"
                    >
                      <PencilIcon class="w-4 h-4" />
                    </button>
                    <button 
                      v-if="prompt.isCustom"
                      class="icon-btn danger"
                      @click.stop="deletePrompt(prompt.id)"
                      title="删除"
                    >
                      <TrashIcon class="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <div class="prompt-desc">{{ prompt.description }}</div>
                <div class="prompt-meta">
                  <span class="meta-item">
                    <span class="temp-icon">Temp:</span>
                    Temp: {{ prompt.temperature }}
                  </span>
                  <span class="meta-item">
                    <DocumentIcon class="w-3 h-3" />
                    {{ prompt.responseFormat || 'text' }}
                  </span>
                </div>
              </div>
            </div>
            
            <!-- 编辑 Prompt 弹窗 -->
            <div v-if="editingPrompt" class="edit-modal">
              <div class="edit-content">
                <h4>{{ editingPrompt.id ? '编辑' : '添加' }} Prompt 模板</h4>
                
                <div class="form-group">
                  <label>名称</label>
                  <input v-model="editingPrompt.name" placeholder="如：翻译提示词" />
                </div>
                
                <div class="form-group">
                  <label>描述</label>
                  <input v-model="editingPrompt.description" placeholder="简短描述这个模板的功能" />
                </div>
                
                <div class="form-group">
                  <label>System Prompt</label>
                  <textarea 
                    v-model="editingPrompt.systemPrompt" 
                    rows="6"
                    placeholder="系统提示词，定义AI的角色和任务..."
                  ></textarea>
                </div>
                
                <div class="form-group">
                  <label>User Prompt 模板</label>
                  <textarea 
                    v-model="editingPrompt.userPromptTemplate" 
                    rows="3"
                    placeholder="用户提示词模板，使用 {{input}} 作为输入占位符"
                  ></textarea>
                  <span class="hint">使用 {{input}} 作为用户输入的占位符</span>
                </div>
                
                <div class="form-row">
                  <div class="form-group half">
                    <label>Temperature</label>
                    <input 
                      v-model.number="editingPrompt.temperature" 
                      type="number"
                      min="0"
                      max="2"
                      step="0.1"
                    />
                  </div>
                  <div class="form-group half">
                    <label>响应格式</label>
                    <select v-model="editingPrompt.responseFormat">
                      <option value="json">JSON</option>
                      <option value="text">纯文本</option>
                    </select>
                  </div>
                </div>
                
                <div class="form-actions">
                  <button class="btn-secondary" @click="editingPrompt = null">取消</button>
                  <button class="btn-primary" @click="savePrompt">保存</button>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { 
  Cog6ToothIcon, 
  XMarkIcon, 
  CpuChipIcon, 
  DocumentTextIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
  ArrowPathIcon,
  DocumentIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, useAIStore } from '../stores'

const appStore = useAppStore()
const aiStore = useAIStore()
const { theme } = storeToRefs(appStore)
const { providers, availablePrompts, currentProviderId, currentPromptId } = storeToRefs(aiStore)

const localTheme = ref(theme.value)
const activeTab = ref('general')

// 编辑状态
const editingProvider = ref(null)
const editingPrompt = ref(null)
const showAddProvider = ref(false)
const showAddPrompt = ref(false)

// 监听主题变化
watch(localTheme, (newTheme) => {
  appStore.setTheme(newTheme)
})

// 确保AI Store已初始化
if (aiStore.providers.length === 0) {
  aiStore.init()
}

// 选择提供商
function selectProvider(id) {
  aiStore.setCurrentProvider(id)
}

// 编辑提供商
function editProvider(provider) {
  editingProvider.value = {
    ...provider,
    modelsText: provider.models?.join(', ') || provider.model,
  }
}

// 保存提供商
function saveProvider() {
  const providerData = {
    ...editingProvider.value,
    models: editingProvider.value.modelsText.split(',').map(s => s.trim()).filter(Boolean),
    enabled: true,
  }
  delete providerData.modelsText
  
  if (editingProvider.value.id && !editingProvider.value.isCustom) {
    // 更新内置提供商
    aiStore.updateProvider(editingProvider.value.id, providerData)
  } else if (editingProvider.value.id) {
    // 更新自定义提供商
    aiStore.updateProvider(editingProvider.value.id, providerData)
  } else {
    // 添加新提供商
    aiStore.addProvider(providerData)
  }
  
  editingProvider.value = null
  showAddProvider.value = false
}

// 删除提供商
function deleteProvider(id) {
  if (confirm('确定要删除这个 AI 提供商吗？')) {
    aiStore.removeProvider(id)
  }
}

// 选择 Prompt
function selectPrompt(id) {
  aiStore.setCurrentPrompt(id)
}

// 编辑 Prompt
function editPrompt(prompt) {
  editingPrompt.value = { ...prompt }
}

// 保存 Prompt
function savePrompt() {
  if (editingPrompt.value.id && !editingPrompt.value.isCustom) {
    // 更新内置模板
    aiStore.updatePrompt(editingPrompt.value.id, editingPrompt.value)
  } else if (editingPrompt.value.id) {
    // 更新自定义模板
    aiStore.updatePrompt(editingPrompt.value.id, editingPrompt.value)
  } else {
    // 添加新模板
    aiStore.addPrompt(editingPrompt.value)
  }
  
  editingPrompt.value = null
  showAddPrompt.value = false
}

// 删除 Prompt
function deletePrompt(id) {
  if (confirm('确定要删除这个 Prompt 模板吗？')) {
    aiStore.removePrompt(id)
  }
}

// 重置 AI 配置
function resetAIConfig() {
  if (confirm('确定要重置所有 AI 配置吗？这将恢复默认设置。')) {
    aiStore.resetToDefaults()
  }
}

// 监听添加按钮
watch(showAddProvider, (val) => {
  if (val) {
    editingProvider.value = {
      name: '',
      type: 'openai-compatible',
      baseUrl: '',
      apiKey: '',
      model: '',
      modelsText: '',
    }
  }
})

watch(showAddPrompt, (val) => {
  if (val) {
    editingPrompt.value = {
      name: '',
      description: '',
      systemPrompt: '',
      userPromptTemplate: '{{input}}',
      temperature: 0.7,
      responseFormat: 'json',
    }
  }
})

defineEmits(['close'])
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  width: 100%;
  max-width: 600px;
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
}

.modal-header h3 {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
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

.modal-body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.setting-tabs {
  display: flex;
  gap: 4px;
  padding: 0 20px;
  border-bottom: 1px solid #334155;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  color: #94a3b8;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: #e2e8f0;
}

.tab-btn.active {
  color: #0ea5e9;
  border-bottom-color: #0ea5e9;
}

.setting-section {
  margin-bottom: 24px;
}

.setting-section:last-child {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #22c55e;
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn:hover {
  background-color: #16a34a;
}

/* Provider List */
.provider-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.provider-card {
  padding: 14px;
  background-color: #1e293b;
  border: 2px solid transparent;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.provider-card:hover {
  border-color: #475569;
}

.provider-card.active {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.1);
}

.provider-card.disabled {
  opacity: 0.6;
}

.provider-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.provider-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.provider-name {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.custom-badge, .disabled-badge {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
}

.custom-badge {
  background-color: rgba(124, 58, 237, 0.2);
  color: #a78bfa;
}

.disabled-badge {
  background-color: rgba(100, 116, 139, 0.2);
  color: #64748b;
}

.provider-actions {
  display: flex;
  gap: 4px;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.icon-btn.danger:hover {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.provider-info-row {
  display: flex;
  gap: 6px;
  font-size: 12px;
}

.info-label {
  color: #64748b;
}

.info-value {
  color: #94a3b8;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Prompt List */
.prompt-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.prompt-card {
  padding: 14px;
  background-color: #1e293b;
  border: 2px solid transparent;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.prompt-card:hover {
  border-color: #475569;
}

.prompt-card.active {
  border-color: #7c3aed;
  background-color: rgba(124, 58, 237, 0.1);
}

.prompt-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.prompt-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.prompt-name {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.prompt-desc {
  font-size: 12px;
  color: #94a3b8;
  margin-bottom: 8px;
}

.prompt-meta {
  display: flex;
  gap: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #64748b;
}

/* Edit Modal */
.edit-modal {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
  padding: 20px;
}

.edit-content {
  width: 100%;
  max-width: 500px;
  max-height: 80vh;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 24px;
  overflow-y: auto;
}

.edit-content h4 {
  margin: 0 0 20px 0;
  font-size: 16px;
  color: #e2e8f0;
}

/* Form Styles */
.form-group {
  margin-bottom: 16px;
}

.form-group.half {
  flex: 1;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-group label {
  display: block;
  font-size: 13px;
  color: #94a3b8;
  margin-bottom: 6px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: #0ea5e9;
}

.form-group textarea {
  resize: vertical;
  font-family: inherit;
}

.form-group .hint {
  display: block;
  font-size: 11px;
  color: #64748b;
  margin-top: 4px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.btn-primary, .btn-secondary {
  flex: 1;
  padding: 12px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #0ea5e9;
  border: none;
  color: white;
}

.btn-primary:hover {
  background-color: #38bdf8;
}

.btn-secondary {
  background-color: transparent;
  border: 1px solid #334155;
  color: #94a3b8;
}

.btn-secondary:hover {
  background-color: #334155;
  color: #e2e8f0;
}

/* Setting Actions */
.setting-actions {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #334155;
}

/* General Settings */
.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #1e293b;
}

.setting-item:last-child {
  border-bottom: none;
}

.item-label {
  font-size: 14px;
  color: #e2e8f0;
}

.item-select {
  padding: 6px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 13px;
  cursor: pointer;
  outline: none;
}

.item-select:focus {
  border-color: #0ea5e9;
}

.item-shortcut {
  padding: 4px 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
  color: #94a3b8;
}

.about-info {
  text-align: center;
  padding: 20px 0;
}

.app-name {
  font-size: 18px;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 4px;
}

.app-version {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
}

.app-desc {
  font-size: 13px;
  color: #94a3b8;
}
</style>
