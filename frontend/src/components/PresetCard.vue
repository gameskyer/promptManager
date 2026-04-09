<template>
  <div class="preset-card" :class="{ 'compact-mode': isCompact }">
    <!-- 封面图区域 - 非紧凑模式下显示 -->
    <div v-if="!isCompact" class="card-cover">
      <div class="cover-image-wrapper" @click="openImageViewer">
        <img v-if="preset.thumbnail" :src="preset.thumbnail" :alt="preset.title" />
        <div v-else class="cover-placeholder">
          <PhotoIcon class="w-12 h-12 text-slate-600" />
          <span>暂无封面</span>
        </div>
        <!-- 查看大图提示 -->
        <div v-if="preset.thumbnail" class="view-image-hint">
          <EyeIcon class="w-5 h-5" />
          <span>查看大图</span>
        </div>
      </div>
      
      <!-- 缩略图预览条（横向滚动） -->
      <div v-if="preset.previews?.length > 0" class="previews-strip">
        <div 
          v-for="(preview, index) in preset.previews" 
          :key="index"
          class="preview-thumb"
          :class="{ active: currentPreviewIndex === index }"
          @click="selectPreview(index)"
        >
          <img :src="preview" :alt="`预览 ${index + 1}`" />
        </div>
        <button class="add-preview-btn" @click.stop="addPreview" title="添加预览图">
          <PlusIcon class="w-4 h-4" />
        </button>
      </div>
      <button v-else class="add-preview-btn standalone" @click.stop="addPreview">
        <PlusIcon class="w-4 h-4" />
        <span>添加预览</span>
      </button>
      
      <!-- 版本标签 -->
      <div class="version-badge">
        <span class="version-label">V{{ preset.current_version || 1 }}</span>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="card-content">
      <div class="content-header">
        <h3 class="preset-title">{{ preset.title }}</h3>
        <!-- 紧凑模式下显示版本号 -->
        <span v-if="isCompact" class="version-tag">V{{ preset.current_version || 1 }}</span>
      </div>
      
      <!-- 模型和LoRA - 非紧凑模式下显示 -->
      <template v-if="!isCompact">
        <div v-if="preset.params?.model" class="model-info">
          <CubeIcon class="w-3 h-3" />
          <span class="model-name">{{ formatModelName(preset.params.model) }}</span>
        </div>
        
        <div v-if="preset.loras?.length > 0" class="lora-info">
          <div class="lora-header">
            <SquaresPlusIcon class="w-3 h-3" />
            <span>LoRA ({{ preset.loras.length }})</span>
          </div>
          <div class="lora-list">
            <span 
              v-for="(lora, index) in displayedLoras" 
              :key="index"
              class="lora-tag"
            >
              {{ lora.name }}:{{ lora.weight }}
            </span>
            <span v-if="hasMoreLoras" class="lora-more">+{{ remainingLoras }}</span>
          </div>
        </div>
      </template>
      
      <!-- 提示词预览 -->
      <div class="prompt-section">
        <div class="prompt-label">
          <span class="label-text">正向词</span>
          <button class="copy-btn" @click="copyPrompt('pos')">
            <ClipboardIcon class="w-3 h-3" />
          </button>
        </div>
        <p class="prompt-text positive">{{ truncatedPosText }}</p>
      </div>
      
      <div class="prompt-section">
        <div class="prompt-label">
          <span class="label-text">负向词</span>
          <button v-if="preset.neg_text" class="copy-btn" @click="copyPrompt('neg')">
            <ClipboardIcon class="w-3 h-3" />
          </button>
        </div>
        <p class="prompt-text negative">{{ preset.neg_text ? truncatedNegText : '无' }}</p>
      </div>
      
      <!-- 参数信息 - 非紧凑模式下显示 -->
      <div v-if="!isCompact && preset.params" class="params-bar">
        <span v-if="preset.params.steps" class="param-tag">
          <AdjustmentsHorizontalIcon class="w-3 h-3" />
          {{ preset.params.steps }}步
        </span>
        <span v-if="preset.params.cfg" class="param-tag">
          <ScaleIcon class="w-3 h-3" />
          CFG {{ preset.params.cfg }}
        </span>
        <span v-if="preset.params.sampler" class="param-tag">
          <SwatchIcon class="w-3 h-3" />
          {{ preset.params.sampler }}
        </span>
      </div>
    </div>
    
    <!-- 操作按钮 -->
    <div class="card-actions">
      <button class="action-btn view" @click="$emit('view', preset)" title="查看详情">
        <EyeIcon class="w-4 h-4" />
      </button>
      <button class="action-btn edit" @click="$emit('edit', preset)" title="编辑">
        <PencilIcon class="w-4 h-4" />
      </button>
      <button class="action-btn use" @click="$emit('use', preset)" title="使用该预设">
        <PlayIcon class="w-4 h-4" />
      </button>
      <button class="action-btn delete" @click="handleDelete" title="删除">
        <TrashIcon class="w-4 h-4" />
      </button>
    </div>
    
  </div>
  
  <!-- 图片查看器弹窗 - 使用 Teleport 传送到 body，避免被卡片样式影响 -->
  <Teleport to="body">
    <div 
      v-if="showImageViewer" 
      class="image-viewer-modal"
      @click="closeImageViewer"
    >
      <div class="image-viewer-content" @click.stop>
        <!-- 关闭按钮 -->
        <button class="viewer-close-btn" @click="closeImageViewer">
          <XMarkIcon class="w-6 h-6" />
        </button>
        
        <!-- 图片计数器 -->
        <div v-if="allImages.length > 1" class="image-counter">
          {{ viewerCurrentIndex + 1 }} / {{ allImages.length }}
        </div>
        
        <!-- 左切换按钮 -->
        <button 
          v-if="allImages.length > 1"
          class="nav-btn prev-btn"
          :class="{ disabled: viewerCurrentIndex === 0 }"
          @click.stop="prevImage"
        >
          <ChevronLeftIcon class="w-8 h-8" />
        </button>
        
        <!-- 主图片显示 -->
        <div class="viewer-image-wrapper">
          <img 
            :src="allImages[viewerCurrentIndex]" 
            :alt="`预览图 ${viewerCurrentIndex + 1}`"
            @click.stop
          />
        </div>
        
        <!-- 右切换按钮 -->
        <button 
          v-if="allImages.length > 1"
          class="nav-btn next-btn"
          :class="{ disabled: viewerCurrentIndex === allImages.length - 1 }"
          @click.stop="nextImage"
        >
          <ChevronRightIcon class="w-8 h-8" />
        </button>
        
        <!-- 底部缩略图列表 -->
        <div v-if="allImages.length > 1" class="viewer-thumbnails">
          <div 
            v-for="(img, index) in allImages" 
            :key="index"
            class="viewer-thumb"
            :class="{ active: viewerCurrentIndex === index }"
            @click.stop="viewerCurrentIndex = index"
          >
            <img :src="img" :alt="`缩略图 ${index + 1}`" />
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import {
  PhotoIcon,
  PlusIcon,
  ClipboardIcon,
  AdjustmentsHorizontalIcon,
  ScaleIcon,
  SwatchIcon,
  EyeIcon,
  PencilIcon,
  PlayIcon,
  TrashIcon,
  CubeIcon,
  SquaresPlusIcon,
  XMarkIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
})

// Debug: 监听 preset 变化
watch(() => props.preset, (newVal) => {
  console.log('[PresetCard] preset changed:', {
    id: newVal.id,
    title: newVal.title,
    thumbnail: newVal.thumbnail,
    previews: newVal.previews,
    previewsLength: newVal.previews?.length,
  })
}, { immediate: true })

const emit = defineEmits(['view', 'edit', 'use', 'delete', 'update-thumbnail'])

const currentPreviewIndex = ref(0)
const isCompact = ref(false)
const COMPACT_THRESHOLD = 1100 // 窗口宽度小于此值进入紧凑模式

// 图片查看器状态
const showImageViewer = ref(false)
const viewerCurrentIndex = ref(0)

const MAX_LORAS_DISPLAY = 2

// 检测窗口大小
function checkWindowSize() {
  isCompact.value = window.innerWidth < COMPACT_THRESHOLD
}

onMounted(() => {
  checkWindowSize()
  window.addEventListener('resize', checkWindowSize)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkWindowSize)
})

const displayedLoras = computed(() => {
  if (!props.preset.loras || props.preset.loras.length === 0) return []
  return props.preset.loras.slice(0, MAX_LORAS_DISPLAY)
})

const hasMoreLoras = computed(() => {
  return props.preset.loras && props.preset.loras.length > MAX_LORAS_DISPLAY
})

const remainingLoras = computed(() => {
  if (!props.preset.loras) return 0
  return props.preset.loras.length - MAX_LORAS_DISPLAY
})

// 紧凑模式下显示更少的文字
const truncatedPosText = computed(() => {
  const text = props.preset.pos_text || ''
  const maxLen = isCompact.value ? 150 : 100
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
})

const truncatedNegText = computed(() => {
  const text = props.preset.neg_text || ''
  const maxLen = isCompact.value ? 100 : 80
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
})

// 所有图片列表（包括缩略图和预览图）
const allImages = computed(() => {
  const images = []
  // 如果有 thumbnail，放在第一位
  if (props.preset.thumbnail) {
    images.push(props.preset.thumbnail)
  }
  // 添加所有预览图
  if (props.preset.previews?.length > 0) {
    props.preset.previews.forEach(preview => {
      // 避免重复添加 thumbnail
      if (preview !== props.preset.thumbnail) {
        images.push(preview)
      }
    })
  }
  return images
})

function formatModelName(modelPath) {
  if (!modelPath) return ''
  const parts = modelPath.split('/')
  const filename = parts[parts.length - 1]
  if (filename.length > 25) {
    return filename.slice(0, 22) + '...'
  }
  return filename
}

function selectPreview(index) {
  currentPreviewIndex.value = index
  emit('update-thumbnail', props.preset.id, props.preset.previews[index])
}

async function copyPrompt(type) {
  const text = type === 'pos' ? props.preset.pos_text : props.preset.neg_text
  try {
    await navigator.clipboard.writeText(text)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

function addPreview() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = (e) => {
    const file = e.target.files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        const newPreviews = [...(props.preset.previews || []), event.target.result]
        emit('update-thumbnail', props.preset.id, event.target.result, newPreviews)
      }
      reader.readAsDataURL(file)
    }
  }
  input.click()
}

function handleDelete() {
  if (confirm(`确定要删除预设 "${props.preset.title}" 吗？`)) {
    emit('delete', props.preset.id)
  }
}

// ========== 图片查看器功能 ==========

function openImageViewer() {
  if (!props.preset.thumbnail && !props.preset.previews?.length) {
    return
  }
  // 从当前选中的缩略图索引开始
  viewerCurrentIndex.value = currentPreviewIndex.value
  showImageViewer.value = true
  // 禁止背景滚动
  document.body.style.overflow = 'hidden'
  // 添加键盘事件监听
  document.addEventListener('keydown', handleKeydown)
}

function closeImageViewer() {
  showImageViewer.value = false
  // 恢复背景滚动
  document.body.style.overflow = ''
  // 移除键盘事件监听
  document.removeEventListener('keydown', handleKeydown)
}

function prevImage() {
  if (viewerCurrentIndex.value > 0) {
    viewerCurrentIndex.value--
  }
}

function nextImage() {
  if (viewerCurrentIndex.value < allImages.value.length - 1) {
    viewerCurrentIndex.value++
  }
}

function handleKeydown(e) {
  if (!showImageViewer.value) return
  
  switch (e.key) {
    case 'Escape':
      closeImageViewer()
      break
    case 'ArrowLeft':
      prevImage()
      break
    case 'ArrowRight':
      nextImage()
      break
  }
}
</script>

<style scoped>
.preset-card {
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  height: auto;
  min-height: 200px;
}

.preset-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  border-color: #475569;
}

/* 紧凑模式样式 */
.preset-card.compact-mode {
  display: grid;
  grid-template-columns: 1fr auto;
  grid-template-rows: auto auto;
}

.preset-card.compact-mode .card-content {
  grid-column: 1;
  grid-row: 1 / -1;
  padding: 14px;
  gap: 10px;
}

.preset-card.compact-mode .card-actions {
  grid-column: 2;
  grid-row: 1 / -1;
  flex-direction: column;
  border-top: none;
  border-left: 1px solid #334155;
  padding: 8px;
}

.preset-card.compact-mode .action-btn {
  padding: 10px;
}

.preset-card.compact-mode .content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.preset-card.compact-mode .version-tag {
  font-size: 11px;
  font-weight: 600;
  color: #a78bfa;
  padding: 2px 8px;
  background-color: rgba(124, 58, 237, 0.2);
  border-radius: 4px;
  white-space: nowrap;
}

/* 封面区域 */
.card-cover {
  position: relative;
  aspect-ratio: 16 / 10;
  background-color: #0f172a;
  overflow: hidden;
}

.cover-image-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  cursor: pointer;
  overflow: hidden;
}

.cover-image-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.cover-image-wrapper:hover img {
  transform: scale(1.05);
}

/* 查看大图提示 */
.view-image-hint {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  font-size: 13px;
  font-weight: 500;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.cover-image-wrapper:hover .view-image-hint {
  opacity: 1;
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #64748b;
}

.cover-placeholder span {
  font-size: 12px;
}

/* 缩略图预览条 */
.previews-strip {
  position: absolute;
  bottom: 8px;
  left: 8px;
  right: 8px;
  display: flex;
  gap: 6px;
  padding: 6px;
  background-color: rgba(0, 0, 0, 0.7);
  border-radius: 8px;
  overflow-x: auto;
  scrollbar-width: none;
}

.previews-strip::-webkit-scrollbar {
  display: none;
}

.preview-thumb {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.preview-thumb:hover {
  transform: scale(1.05);
}

.preview-thumb.active {
  border-color: #0ea5e9;
}

.preview-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 添加预览按钮 */
.add-preview-btn {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background-color: rgba(255, 255, 255, 0.1);
  border: 1px dashed rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s;
}

.add-preview-btn:hover {
  background-color: rgba(255, 255, 255, 0.2);
  color: #e2e8f0;
}

.add-preview-btn.standalone {
  position: absolute;
  bottom: 12px;
  right: 12px;
  width: auto;
  height: auto;
  padding: 8px 12px;
  gap: 6px;
  font-size: 12px;
}

/* 版本标签 */
.version-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 10px;
  background-color: rgba(124, 58, 237, 0.9);
  border-radius: 6px;
}

.version-label {
  font-size: 12px;
  font-weight: 700;
  color: white;
  font-variant-numeric: tabular-nums;
}

/* 内容区域 */
.card-content {
  flex: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 0;
  overflow: visible;
}

.preset-title {
  font-size: 15px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 模型信息 */
.model-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background-color: rgba(14, 165, 233, 0.1);
  border-radius: 6px;
  color: #0ea5e9;
  font-size: 12px;
}

.model-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* LoRA 信息 */
.lora-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px 10px;
  background-color: rgba(124, 58, 237, 0.1);
  border-radius: 6px;
}

.lora-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #a78bfa;
  font-weight: 500;
}

.lora-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.lora-tag {
  font-size: 11px;
  color: #e2e8f0;
  background-color: rgba(124, 58, 237, 0.2);
  padding: 3px 8px;
  border-radius: 4px;
}

.lora-more {
  font-size: 11px;
  color: #64748b;
  padding: 3px 6px;
}

/* 提示词预览 */
.prompt-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.prompt-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.label-text {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
}

.copy-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.prompt-text {
  font-size: 12px;
  line-height: 1.5;
  color: #94a3b8;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.prompt-text.positive {
  color: #22c55e;
}

.prompt-text.negative {
  color: #ef4444;
}

/* 参数栏 */
.params-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid #334155;
}

.param-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: #0f172a;
  border-radius: 4px;
  font-size: 11px;
  color: #94a3b8;
}

/* 操作按钮 */
.card-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1px;
  border-top: 1px solid #334155;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

.action-btn.view:hover {
  color: #0ea5e9;
}

.action-btn.edit:hover {
  color: #f59e0b;
}

.action-btn.use:hover {
  color: #22c55e;
}

.action-btn.delete:hover {
  color: #ef4444;
  background-color: rgba(239, 68, 68, 0.1);
}

/* ========== 图片查看器样式 ========== */
.image-viewer-modal {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.image-viewer-content {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

/* 关闭按钮 */
.viewer-close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  z-index: 10;
}

.viewer-close-btn:hover {
  background-color: rgba(255, 255, 255, 0.2);
  transform: scale(1.1);
}

/* 图片计数器 */
.image-counter {
  position: absolute;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 8px 16px;
  background-color: rgba(0, 0, 0, 0.6);
  border-radius: 20px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  z-index: 10;
}

/* 主图片区域 */
.viewer-image-wrapper {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 80px;
  box-sizing: border-box;
}

.viewer-image-wrapper img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
}

/* 导航按钮 */
.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 80px;
  background-color: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 8px;
}

.nav-btn:hover:not(.disabled) {
  background-color: rgba(255, 255, 255, 0.2);
}

.nav-btn.disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.prev-btn {
  left: 20px;
}

.next-btn {
  right: 20px;
}

/* 底部缩略图 */
.viewer-thumbnails {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
  padding: 12px;
  background-color: rgba(0, 0, 0, 0.6);
  border-radius: 12px;
  max-width: 80%;
  overflow-x: auto;
  scrollbar-width: none;
}

.viewer-thumbnails::-webkit-scrollbar {
  display: none;
}

.viewer-thumb {
  flex-shrink: 0;
  width: 60px;
  height: 60px;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  opacity: 0.6;
  transition: all 0.2s;
}

.viewer-thumb:hover {
  opacity: 0.9;
  transform: scale(1.05);
}

.viewer-thumb.active {
  border-color: #0ea5e9;
  opacity: 1;
}

.viewer-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
