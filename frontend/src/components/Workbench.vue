<template>
  <aside class="workbench">
    <!-- Header -->
    <div class="workbench-header">
      <h3 class="workbench-title">
        <SwatchIcon class="w-5 h-5" />
        组合工作区
      </h3>
      <div class="workbench-actions">
        <button 
          class="icon-btn"
          :disabled="selectedAtoms.length === 0"
          @click="clearAll"
          title="清空"
        >
          <TrashIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
    
    <!-- Version Info -->
    <div v-if="currentPreset" class="version-info">
      <div class="preset-name">{{ currentPreset.title }}</div>
      <div class="version-badge">
        <span class="version-label">当前版本</span>
        <span class="version-num">V{{ currentVersion?.version_num || currentPreset.current_version }}</span>
      </div>
      <button class="history-btn" @click="toggleTimeline">
        <ClockIcon class="w-4 h-4" />
        版本历史
      </button>
    </div>
    
    <!-- Empty State -->
    <div v-if="selectedAtoms.length === 0" class="empty-workbench">
      <SquaresPlusIcon class="w-12 h-12 text-slate-600" />
      <p>从左侧选择原子词</p>
      <span class="hint">点击原子词添加到此处</span>
    </div>
    
    <!-- Selected Atoms Lists -->
    <div v-else class="atoms-lists">
      <!-- Positive Atoms -->
      <div v-if="positiveAtoms.length > 0" class="atom-group positive">
        <div class="group-header">
          <PlusIcon class="w-4 h-4 text-green-400" />
          <span>正向提示词 ({{ positiveAtoms.length }})</span>
        </div>
        <div class="atoms-list">
          <div
            v-for="(atom, index) in positiveAtoms"
            :key="atom.id"
            class="workbench-item"
            :class="{ 'is-dragging': draggingIndex === index && draggingType === 'positive' }"
            draggable="true"
            @dragstart="handleDragStart(index, 'positive')"
            @dragover.prevent
            @drop="handleDrop(index, 'positive')"
            @dragend="handleDragEnd"
          >
            <div class="drag-handle">
              <Bars3Icon class="w-4 h-4" />
            </div>
            
            <div class="item-content">
              <div class="item-value">{{ atom.value }}</div>
              <div class="item-label">{{ atom.label }}</div>
            </div>
            
            <div class="item-actions">
              <button 
                class="action-btn up"
                :disabled="index === 0"
                @click="moveUp(atom.id)"
                title="上移"
              >
                <ChevronUpIcon class="w-4 h-4" />
              </button>
              <button 
                class="action-btn down"
                :disabled="index === positiveAtoms.length - 1"
                @click="moveDown(atom.id)"
                title="下移"
              >
                <ChevronDownIcon class="w-4 h-4" />
              </button>
              <button 
                class="action-btn remove"
                @click="removeAtom(atom.id)"
                title="移除"
              >
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Negative Atoms -->
      <div v-if="negativeAtoms.length > 0" class="atom-group negative">
        <div class="group-header">
          <MinusIcon class="w-4 h-4 text-red-400" />
          <span>负向提示词 ({{ negativeAtoms.length }})</span>
        </div>
        <div class="atoms-list">
          <div
            v-for="(atom, index) in negativeAtoms"
            :key="atom.id"
            class="workbench-item"
            :class="{ 'is-dragging': draggingIndex === index && draggingType === 'negative' }"
            draggable="true"
            @dragstart="handleDragStart(index, 'negative')"
            @dragover.prevent
            @drop="handleDrop(index, 'negative')"
            @dragend="handleDragEnd"
          >
            <div class="drag-handle">
              <Bars3Icon class="w-4 h-4" />
            </div>
            
            <div class="item-content">
              <div class="item-value">{{ atom.value }}</div>
              <div class="item-label">{{ atom.label }}</div>
            </div>
            
            <div class="item-actions">
              <button 
                class="action-btn up"
                :disabled="index === 0"
                @click="moveUp(atom.id)"
                title="上移"
              >
                <ChevronUpIcon class="w-4 h-4" />
              </button>
              <button 
                class="action-btn down"
                :disabled="index === negativeAtoms.length - 1"
                @click="moveDown(atom.id)"
                title="下移"
              >
                <ChevronDownIcon class="w-4 h-4" />
              </button>
              <button 
                class="action-btn remove"
                @click="removeAtom(atom.id)"
                title="移除"
              >
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Prompt Preview -->
    <div v-if="selectedAtoms.length > 0" class="prompt-preview">
      <!-- Positive Prompt -->
      <div class="prompt-section positive">
        <div class="preview-header">
          <div class="header-left">
            <PlusIcon class="w-4 h-4 text-green-400" />
            <span>正向提示词</span>
          </div>
          <button class="copy-btn" @click="copyPositive">
            <ClipboardIcon class="w-4 h-4" />
            {{ copiedPositive ? '已复制' : '复制' }}
          </button>
        </div>
        <div class="prompt-text positive">{{ positivePromptText }}</div>
      </div>
      
      <!-- Negative Prompt -->
      <div v-if="negativeAtoms.length > 0" class="prompt-section negative">
        <div class="preview-header">
          <div class="header-left">
            <MinusIcon class="w-4 h-4 text-red-400" />
            <span>负向提示词</span>
          </div>
          <button class="copy-btn" @click="copyNegative">
            <ClipboardIcon class="w-4 h-4" />
            {{ copiedNegative ? '已复制' : '复制' }}
          </button>
        </div>
        <div class="prompt-text negative">{{ negativePromptText }}</div>
      </div>
    </div>
    
    <!-- Bottom Actions -->
    <div v-if="selectedAtoms.length > 0" class="workbench-footer">
      <button class="btn-secondary" @click="showSaveDialog = true">
        <FolderArrowDownIcon class="w-4 h-4" />
        保存预设
      </button>
      <button class="btn-primary" @click="generatePrompt">
        <ClipboardDocumentListIcon class="w-4 h-4" />
        生成
      </button>
    </div>
    
    <!-- Save Preset Dialog -->
    <div v-if="showSaveDialog" class="modal-overlay" @click="showSaveDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>保存预设</h3>
          <button class="close-btn" @click="showSaveDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>预设名称 <span class="required">*</span></label>
            <input 
              v-model="presetForm.title" 
              type="text" 
              placeholder="如： my-style-v1"
            />
          </div>
          
          <div class="form-group">
            <label>描述</label>
            <textarea 
              v-model="presetForm.description" 
              rows="3"
              placeholder="描述这个预设的用途和特点..."
            ></textarea>
          </div>
          
          <div class="preview-section">
            <div class="preview-title">正向提示词 ({{ positiveAtoms.length }} 个)</div>
            <div class="preview-text">{{ positivePromptText || '无' }}</div>
          </div>
          
          <div class="preview-section" v-if="negativeAtoms.length > 0">
            <div class="preview-title">负向提示词 ({{ negativeAtoms.length }} 个)</div>
            <div class="preview-text">{{ negativePromptText }}</div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-secondary" @click="showSaveDialog = false">取消</button>
          <button class="btn-primary" @click="savePreset" :disabled="!presetForm.title.trim()">
            保存
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import {
  SwatchIcon,
  TrashIcon,
  ClockIcon,
  SquaresPlusIcon,
  Bars3Icon,
  ChevronUpIcon,
  ChevronDownIcon,
  XMarkIcon,
  ClipboardIcon,
  FolderArrowDownIcon,
  ClipboardDocumentListIcon,
  PlusIcon,
  MinusIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, usePresetStore, useVersionStore } from '../stores'

const appStore = useAppStore()
const presetStore = usePresetStore()
const versionStore = useVersionStore()

const { selectedAtoms, currentPreset, currentVersion } = storeToRefs(appStore)
const showSaveDialog = ref(false)
const showTimeline = ref(false)
const copiedPositive = ref(false)
const copiedNegative = ref(false)
const draggingIndex = ref(null)
const draggingType = ref(null)

// Preset form
const presetForm = ref({
  title: '',
  description: '',
})

// Watch for dialog open to reset form
watch(() => showSaveDialog.value, (val) => {
  if (val) {
    presetForm.value = {
      title: '',
      description: '',
    }
  }
})

// Group atoms by type
const positiveAtoms = computed(() => 
  selectedAtoms.value.filter(a => a.type === 'Positive')
)

const negativeAtoms = computed(() => 
  selectedAtoms.value.filter(a => a.type === 'Negative')
)

// Prompt texts
const positivePromptText = computed(() => {
  return positiveAtoms.value.map(a => a.value).join(', ')
})

const negativePromptText = computed(() => {
  return negativeAtoms.value.map(a => a.value).join(', ')
})

function clearAll() {
  appStore.clearSelection()
}

function removeAtom(atomId) {
  appStore.removeAtom(atomId)
}

function moveUp(atomId) {
  appStore.moveAtom(atomId, 'up')
}

function moveDown(atomId) {
  appStore.moveAtom(atomId, 'down')
}

function toggleTimeline() {
  appStore.toggleTimeline()
}

async function copyPositive() {
  try {
    await navigator.clipboard.writeText(positivePromptText.value)
    copiedPositive.value = true
    setTimeout(() => copiedPositive.value = false, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

async function copyNegative() {
  try {
    await navigator.clipboard.writeText(negativePromptText.value)
    copiedNegative.value = true
    setTimeout(() => copiedNegative.value = false, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// 生成 ComfyUI/A1111 标准格式并复制到剪贴板
async function generatePrompt() {
  const posText = positivePromptText.value || ''
  const negText = negativePromptText.value || ''
  
  // ComfyUI/A1111 标准格式
  const comfyFormat = `${posText}

Negative prompt: ${negText}

Steps: 30, CFG scale: 7, Sampler: DPM++ 2M Karras`
  
  try {
    await navigator.clipboard.writeText(comfyFormat)
    // 显示成功提示
    const toast = document.createElement('div')
    toast.className = 'copy-toast'
    toast.innerHTML = `
      <div class="toast-content">
        <span>已复制到剪贴板，可粘贴到 ComfyUI</span>
      </div>
    `
    document.body.appendChild(toast)
    requestAnimationFrame(() => toast.classList.add('show'))
    setTimeout(() => {
      toast.classList.remove('show')
      setTimeout(() => toast.remove(), 300)
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
    alert('复制失败')
  }
}

async function savePreset() {
  if (!presetForm.value.title.trim()) {
    alert('请输入预设名称')
    return
  }
  
  try {
    await presetStore.createPreset(
      presetForm.value.title,
      0, // categoryId - 默认分类
      positivePromptText.value,
      negativePromptText.value,
      [], // atoms
      { // params
        model: '',
        steps: 30,
        cfg: 7,
        sampler: 'DPM++ 2M Karras',
      },
      [], // loras
      [] // previews
    )
    
    alert('预设保存成功！')
    showSaveDialog.value = false
  } catch (err) {
    alert('保存失败：' + err.message)
  }
}

// Drag and drop
function handleDragStart(index, type) {
  draggingIndex.value = index
  draggingType.value = type
}

function handleDrop(targetIndex, targetType) {
  if (draggingIndex.value === null || draggingType.value !== targetType) return
  if (draggingIndex.value === targetIndex) return
  
  // Get the atoms array for the current type
  const atoms = draggingType.value === 'positive' 
    ? [...positiveAtoms.value] 
    : [...negativeAtoms.value]
  
  const [moved] = atoms.splice(draggingIndex.value, 1)
  atoms.splice(targetIndex, 0, moved)
  
  // Update the full selectedAtoms array while preserving order within types
  const otherAtoms = draggingType.value === 'positive'
    ? negativeAtoms.value
    : positiveAtoms.value
  
  appStore.selectedAtoms = draggingType.value === 'positive'
    ? [...atoms, ...otherAtoms]
    : [...otherAtoms, ...atoms]
}

function handleDragEnd() {
  draggingIndex.value = null
  draggingType.value = null
}

// Listen for save event
watch(() => showSaveDialog.value, (val) => {
  if (val) {
    // Show save dialog
  }
})

// Listen for global save event
if (typeof window !== 'undefined') {
  window.addEventListener('save-preset', () => {
    showSaveDialog.value = true
  })
}
</script>

<style scoped>
.workbench {
  width: 320px;
  background-color: #0f172a;
  border-left: 1px solid #1e293b;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.workbench-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #1e293b;
}

.workbench-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.workbench-actions {
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

.icon-btn:hover:not(:disabled) {
  background-color: #334155;
  color: #e2e8f0;
}

.icon-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.version-info {
  padding: 12px 16px;
  background-color: rgba(14, 165, 233, 0.1);
  border-bottom: 1px solid #1e293b;
}

.preset-name {
  font-size: 13px;
  font-weight: 500;
  color: #e2e8f0;
  margin-bottom: 4px;
}

.version-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.version-label {
  font-size: 11px;
  color: #64748b;
}

.version-num {
  font-size: 11px;
  font-weight: 600;
  color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.2);
  padding: 2px 8px;
  border-radius: 4px;
}

.history-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background-color: transparent;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #94a3b8;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.history-btn:hover {
  background-color: #1e293b;
  color: #e2e8f0;
}

.empty-workbench {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px 20px;
  color: #64748b;
}

.empty-workbench p {
  font-size: 14px;
  font-weight: 500;
}

.empty-workbench .hint {
  font-size: 12px;
  color: #475569;
}

.atoms-lists {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.atom-group {
  margin-bottom: 16px;
}

.atom-group:last-child {
  margin-bottom: 0;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 8px;
  padding: 0 4px;
}

.group-header .text-green-400 {
  color: #22c55e;
}

.group-header .text-red-400 {
  color: #ef4444;
}

.atoms-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.workbench-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  cursor: grab;
  transition: all 0.2s;
}

.workbench-item:hover {
  border-color: #475569;
}

.workbench-item.is-dragging {
  opacity: 0.5;
}

.drag-handle {
  display: flex;
  align-items: center;
  color: #475569;
  cursor: grab;
}

.item-content {
  flex: 1;
  min-width: 0;
}

.item-value {
  font-size: 13px;
  font-weight: 500;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-label {
  font-size: 11px;
  color: #64748b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-actions {
  display: flex;
  gap: 2px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover:not(:disabled) {
  background-color: #334155;
  color: #e2e8f0;
}

.action-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.action-btn.remove:hover {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.prompt-preview {
  padding: 12px;
  border-top: 1px solid #1e293b;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.prompt-section {
  background-color: #1e293b;
  border-radius: 8px;
  overflow: hidden;
}

.prompt-section.positive {
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.prompt-section.negative {
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background-color: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid #334155;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #94a3b8;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: transparent;
  border: 1px solid #334155;
  border-radius: 4px;
  color: #64748b;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.prompt-text {
  padding: 10px 12px;
  font-size: 12px;
  line-height: 1.5;
  color: #cbd5e1;
  word-break: break-word;
  max-height: 100px;
  overflow-y: auto;
}

.prompt-text.positive {
  color: #86efac;
}

.prompt-text.negative {
  color: #fca5a5;
}

.workbench-footer {
  display: flex;
  gap: 8px;
  padding: 12px;
  border-top: 1px solid #1e293b;
}

.btn-secondary,
.btn-primary {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-secondary {
  background-color: #1e293b;
  color: #e2e8f0;
  border: 1px solid #334155;
}

.btn-secondary:hover {
  background-color: #334155;
}

.btn-primary {
  background-color: #0284c7;
  color: white;
}

.btn-primary:hover {
  background-color: #0ea5e9;
}

.btn-comfy {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  background-color: #1e293b;
  border: 1px solid #f59e0b;
  color: #f59e0b;
}

.btn-comfy:hover {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: #fbbf24;
  color: #fbbf24;
}

/* 复制成功提示 */
.copy-toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%) translateY(-100px);
  z-index: 9999;
  opacity: 0;
  transition: all 0.3s ease;
}

.copy-toast.show {
  transform: translateX(-50%) translateY(0);
  opacity: 1;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background-color: #10b981;
  color: white;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

/* Modal Styles */
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
  max-width: 480px;
  max-height: 80vh;
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
  font-size: 16px;
  font-weight: 600;
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

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.modal-footer {
  display: flex;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid #334155;
}

.modal-footer .btn-secondary,
.modal-footer .btn-primary {
  flex: 1;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #94a3b8;
  margin-bottom: 6px;
}

.form-group .required {
  color: #ef4444;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: #0ea5e9;
}

.form-group textarea {
  resize: vertical;
  min-height: 60px;
}

.preview-section {
  margin-top: 16px;
  padding: 12px;
  background-color: #1e293b;
  border-radius: 8px;
}

.preview-title {
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
  margin-bottom: 8px;
}

.preview-text {
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
  word-break: break-word;
  max-height: 80px;
  overflow-y: auto;
}
</style>
