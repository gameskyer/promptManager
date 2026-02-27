<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ preset.title }}</h3>
        <div class="header-actions">
          <button class="version-badge">V{{ preset.current_version || 1 }}</button>
          <button class="close-btn" @click="$emit('close')">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </div>
      
      <div class="modal-body">
        <!-- 图片预览区 -->
        <div class="image-preview-section">
          <div class="main-preview">
            <img v-if="currentImage" :src="currentImage" :alt="preset.title" />
            <div v-else class="no-image">
              <PhotoIcon class="w-16 h-16 text-slate-600" />
              <span>暂无预览图</span>
            </div>
          </div>
          
          <!-- 缩略图列表 -->
          <div v-if="preset.previews?.length > 0" class="thumbnail-list">
            <div 
              v-for="(preview, index) in preset.previews" 
              :key="index"
              class="thumb-item"
              :class="{ active: currentImageIndex === index }"
              @click="currentImageIndex = index"
            >
              <img :src="preview" :alt="`预览 ${index + 1}`" />
            </div>
          </div>
        </div>
        
        <!-- 提示词信息 -->
        <div class="prompt-info-section">
          <!-- 模型信息 -->
          <div v-if="preset.params?.model" class="info-block model-block">
            <div class="block-header">
              <CubeIcon class="w-4 h-4" />
              <span>模型 (Model)</span>
              <button class="copy-btn" @click="copyText(preset.params.model)">
                <ClipboardIcon class="w-3 h-3" />
              </button>
            </div>
            <div class="block-content">
              {{ preset.params.model }}
            </div>
          </div>
          
          <!-- LoRA 信息 -->
          <div v-if="preset.loras?.length > 0" class="info-block lora-block">
            <div class="block-header">
              <SquaresPlusIcon class="w-4 h-4" />
              <span>LoRA</span>
            </div>
            <div class="lora-table">
              <div class="lora-row header">
                <span class="lora-name-header">名称</span>
                <span class="lora-weight-header">权重</span>
              </div>
              <div 
                v-for="(lora, index) in preset.loras" 
                :key="index"
                class="lora-row"
              >
                <span class="lora-name">{{ lora.name }}</span>
                <span class="lora-weight">{{ lora.weight }}</span>
              </div>
            </div>
          </div>
          
          <div class="prompt-block">
            <div class="block-header">
              <span class="block-title positive">正向提示词</span>
              <button class="copy-btn" @click="copyText(preset.pos_text)">
                <ClipboardIcon class="w-4 h-4" />
                复制
              </button>
            </div>
            <div class="prompt-content positive">
              {{ preset.pos_text || '无' }}
            </div>
          </div>
          
          <div class="prompt-block">
            <div class="block-header">
              <span class="block-title negative">负向提示词</span>
              <button class="copy-btn" @click="copyText(preset.neg_text)">
                <ClipboardIcon class="w-4 h-4" />
                复制
              </button>
            </div>
            <div class="prompt-content negative">
              {{ preset.neg_text || '无' }}
            </div>
          </div>
          
          <!-- 参数信息 -->
          <div class="params-block">
            <div class="block-header">
              <span class="block-title">生成参数</span>
            </div>
            <div class="params-grid">
              <div class="param-row">
                <span class="param-name">采样步数</span>
                <span class="param-value">{{ preset.params?.steps || 30 }}</span>
              </div>
              <div class="param-row">
                <span class="param-name">CFG Scale</span>
                <span class="param-value">{{ preset.params?.cfg || 7 }}</span>
              </div>
              <div class="param-row">
                <span class="param-name">采样器</span>
                <span class="param-value">{{ preset.params?.sampler || 'DPM++ 2M Karras' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="btn-secondary" @click="$emit('edit', preset)">
          <PencilIcon class="w-4 h-4" />
          编辑
        </button>
        <button class="btn-secondary history-btn" @click="$emit('view-history', preset)">
          <ClockIcon class="w-4 h-4" />
          版本历史
        </button>
        <button class="btn-primary" @click="usePresetWithComfyUI">
          <ClipboardDocumentListIcon class="w-4 h-4" />
          使用该预设
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  XMarkIcon,
  PhotoIcon,
  ClipboardIcon,
  PencilIcon,
  CubeIcon,
  SquaresPlusIcon,
  ClockIcon,
  ClipboardDocumentListIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['close', 'edit', 'use', 'view-history'])

const currentImageIndex = ref(0)

const currentImage = computed(() => {
  if (props.preset.previews?.length > 0) {
    return props.preset.previews[currentImageIndex.value]
  }
  return props.preset.thumbnail
})

async function copyText(text) {
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
    showCopySuccess()
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// 使用该预设 - 复制为 ComfyUI 标准格式
async function usePresetWithComfyUI() {
  const posText = props.preset.pos_text || ''
  const negText = props.preset.neg_text || ''
  
  // ComfyUI/A1111 标准格式
  const comfyFormat = `${posText}

Negative prompt: ${negText}

Steps: ${props.preset.params?.steps || 30}, CFG scale: ${props.preset.params?.cfg || 7}, Sampler: ${props.preset.params?.sampler || 'DPM++ 2M Karras'}, Model: ${props.preset.params?.model || 'Unknown'}`
  
  try {
    await navigator.clipboard.writeText(comfyFormat)
    // 触发 use 事件
    emit('use', props.preset)
    // 关闭弹窗
    emit('close')
    // 显示成功提示
    showToast('已复制到剪贴板，可粘贴到 ComfyUI')
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// 显示提示
function showToast(message) {
  const toast = document.createElement('div')
  toast.className = 'copy-toast'
  toast.innerHTML = `
    <div class="toast-content">
      <span>${message}</span>
    </div>
  `
  document.body.appendChild(toast)
  
  // 动画进入
  requestAnimationFrame(() => {
    toast.classList.add('show')
  })
  
  // 2秒后移除
  setTimeout(() => {
    toast.classList.remove('show')
    setTimeout(() => toast.remove(), 300)
  }, 2000)
}
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
  max-width: 900px;
  max-height: 90vh;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
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
  font-size: 18px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-badge {
  padding: 4px 12px;
  background-color: #7c3aed;
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 13px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
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
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;
}

.image-preview-section {
  padding: 20px;
  border-right: 1px solid #334155;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.main-preview {
  flex: 1;
  min-height: 300px;
  background-color: #1e293b;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-preview img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.no-image {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #64748b;
}

.thumbnail-list {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 4px;
}

.thumb-item {
  flex-shrink: 0;
  width: 60px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.thumb-item:hover {
  transform: scale(1.05);
}

.thumb-item.active {
  border-color: #0ea5e9;
}

.thumb-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.prompt-info-section {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-block {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  border-radius: 8px;
}

.model-block {
  background-color: rgba(14, 165, 233, 0.1);
  border: 1px solid rgba(14, 165, 233, 0.2);
}

.lora-block {
  background-color: rgba(124, 58, 237, 0.1);
  border: 1px solid rgba(124, 58, 237, 0.2);
}

.block-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.block-header span {
  font-size: 12px;
  font-weight: 600;
  color: #e2e8f0;
}

.model-block .block-header {
  color: #0ea5e9;
}

.lora-block .block-header {
  color: #a78bfa;
}

.block-content {
  font-size: 13px;
  color: #e2e8f0;
  word-break: break-all;
  font-family: monospace;
}

.lora-table {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.lora-row {
  display: grid;
  grid-template-columns: 1fr 80px;
  gap: 12px;
  padding: 8px 10px;
  border-radius: 6px;
  font-size: 13px;
}

.lora-row.header {
  background-color: rgba(0, 0, 0, 0.2);
  font-weight: 600;
  color: #94a3b8;
  font-size: 11px;
  text-transform: uppercase;
}

.lora-row:not(.header) {
  background-color: rgba(255, 255, 255, 0.03);
  color: #e2e8f0;
}

.lora-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lora-weight {
  text-align: center;
  font-weight: 600;
  color: #a78bfa;
}

.prompt-block {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.block-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.block-title {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.block-title.positive {
  color: #22c55e;
}

.block-title.negative {
  color: #ef4444;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #334155;
  border: none;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background-color: #475569;
}

.prompt-content {
  padding: 12px;
  background-color: #1e293b;
  border-radius: 8px;
  font-size: 13px;
  line-height: 1.6;
  color: #cbd5e1;
  max-height: 120px;
  overflow-y: auto;
  word-break: break-word;
}

.prompt-content.positive {
  border-left: 3px solid #22c55e;
}

.prompt-content.negative {
  border-left: 3px solid #ef4444;
}

.params-block {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.params-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

.param-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background-color: #1e293b;
  border-radius: 6px;
}

.param-name {
  font-size: 12px;
  color: #64748b;
}

.param-value {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #334155;
}

.btn-secondary {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: 8px;
  background-color: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background-color: #334155;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: 8px;
  background-color: #0284c7;
  border: none;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  background-color: #0ea5e9;
}

.comfy-btn {
  background-color: #1e293b;
  border-color: #f59e0b;
  color: #f59e0b;
}

.comfy-btn:hover {
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

@media (max-width: 768px) {
  .modal-body {
    grid-template-columns: 1fr;
  }
  
  .image-preview-section {
    border-right: none;
    border-bottom: 1px solid #334155;
  }
}
</style>
