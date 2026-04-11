<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div 
      class="modal-content" 
      @click.stop
      :style="contentStyle"
    >
      <div class="modal-header">
        <h3>
          <span class="version-tag">V{{ version.version_num }}</span>
          <span class="preset-name">{{ preset?.title }}</span>
        </h3>
        <div class="header-actions">
          <button 
            class="star-btn"
            :class="{ starred: version.is_starred }"
            @click="toggleStar"
          >
            <StarIcon class="w-5 h-5" />
          </button>
          <button class="close-btn" @click="$emit('close')">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </div>
      
      <div class="modal-body">
        <!-- 左侧：预览图 -->
        <div class="preview-section">
          <div class="main-preview">
            <img v-if="currentImage" :src="currentImage" :alt="`版本 V${version.version_num}`" />
            <div v-else class="no-preview">
              <PhotoIcon class="w-16 h-16 text-slate-600" />
              <span>暂无预览图</span>
            </div>
          </div>
          
          <!-- 预览图列表 -->
          <div v-if="version.previews?.length > 0" class="preview-list">
            <div 
              v-for="(preview, index) in version.previews" 
              :key="index"
              class="preview-thumb"
              :class="{ active: currentImageIndex === index }"
              @click="currentImageIndex = index"
            >
              <img :src="preview" :alt="`预览 ${index + 1}`" />
            </div>
          </div>
        </div>
        
        <!-- 右侧：版本信息 -->
        <div class="info-section">
          <!-- 时间信息 -->
          <div class="time-info">
            <div class="time-item">
              <ClockIcon class="w-4 h-4" />
              <span class="time-label">创建时间</span>
              <span class="time-value">{{ formatFullTime(version.created_at) }}</span>
            </div>
            <div v-if="timeAgo" class="time-ago">
              {{ timeAgo }}
            </div>
          </div>
          
          <!-- 变更统计 -->
          <div class="diff-stats-bar">
            <div class="diff-item added">
              <PlusIcon class="w-4 h-4" />
              <span class="diff-count">{{ addedCount }}</span>
              <span class="diff-label">新增</span>
            </div>
            <div class="diff-item removed">
              <MinusIcon class="w-4 h-4" />
              <span class="diff-count">{{ removedCount }}</span>
              <span class="diff-label">删除</span>
            </div>
            <div v-if="hasParamChanges" class="diff-item changed">
              <AdjustmentsHorizontalIcon class="w-4 h-4" />
              <span class="diff-label">参数变更</span>
            </div>
          </div>
          
          <!-- 详细变更列表 -->
          <div v-if="diffDetails.length > 0" class="diff-details">
            <div class="section-title">变更详情</div>
            <div class="diff-list">
              <div 
                v-for="(diff, index) in diffDetails" 
                :key="index"
                class="diff-row"
                :class="diff.type"
              >
                <span class="diff-icon">
                  <PlusIcon v-if="diff.type === 'added'" class="w-3 h-3" />
                  <MinusIcon v-else-if="diff.type === 'removed'" class="w-3 h-3" />
                  <ArrowPathIcon v-else class="w-3 h-3" />
                </span>
                <span class="diff-text">{{ diff.text }}</span>
              </div>
            </div>
          </div>
          
          <!-- 提示词 -->
          <div class="prompt-section">
            <div class="section-title">正向提示词</div>
            <div class="prompt-box positive">
              {{ snapshot.pos_text || '无' }}
            </div>
          </div>
          
          <div class="prompt-section">
            <div class="section-title">负向提示词</div>
            <div class="prompt-box negative">
              {{ snapshot.neg_text || '无' }}
            </div>
          </div>
          
          <!-- 参数 -->
          <div v-if="snapshot.params" class="params-section">
            <div class="section-title">生成参数</div>
            <div class="params-grid">
              <div class="param-item">
                <span class="param-label">采样步数</span>
                <span class="param-value">{{ snapshot.params.steps || 30 }}</span>
              </div>
              <div class="param-item">
                <span class="param-label">CFG Scale</span>
                <span class="param-value">{{ snapshot.params.cfg || 7 }}</span>
              </div>
              <div class="param-item">
                <span class="param-label">采样器</span>
                <span class="param-value">{{ snapshot.params.sampler || 'DPM++ 2M Karras' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="btn-secondary" @click="$emit('rollback', version)">
          <ArrowUturnLeftIcon class="w-4 h-4" />
          回滚到此版本
        </button>
        <div class="spacer"></div>
        <button class="btn-primary" @click="$emit('use', version)">
          <PlayIcon class="w-4 h-4" />
          使用该版本
        </button>
      </div>
      
      <!-- 高度调整手柄 -->
      <div 
        class="resize-handle"
        @mousedown="startResize"
        title="拖拽调整高度"
      >
        <div class="resize-indicator"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import dayjs from 'dayjs'
import {
  XMarkIcon,
  StarIcon,
  PhotoIcon,
  ClockIcon,
  PlusIcon,
  MinusIcon,
  AdjustmentsHorizontalIcon,
  ArrowPathIcon,
  ArrowUturnLeftIcon,
  PlayIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  version: {
    type: Object,
    required: true,
  },
  preset: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['close', 'rollback', 'use', 'toggle-star'])

const currentImageIndex = ref(0)

// 高度调整相关
const modalHeight = ref(0) // 0 表示使用默认高度
const isResizing = ref(false)
const startY = ref(0)
const startHeight = ref(0)

const contentStyle = computed(() => {
  if (modalHeight.value > 0) {
    return { height: `${modalHeight.value}px`, maxHeight: 'none' }
  }
  return {}
})

function startResize(e) {
  isResizing.value = true
  startY.value = e.clientY
  const content = document.querySelector('.modal-content')
  if (content) {
    startHeight.value = content.offsetHeight
  }
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
  e.preventDefault()
}

function onResize(e) {
  if (!isResizing.value) return
  const delta = e.clientY - startY.value
  const newHeight = startHeight.value + delta
  // 限制最小和最大高度
  modalHeight.value = Math.max(400, Math.min(newHeight, window.innerHeight - 40))
}

function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
}

onUnmounted(() => {
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
})

const snapshot = computed(() => {
  return props.version.snapshot || {}
})

const currentImage = computed(() => {
  if (props.version.previews?.length > 0) {
    return props.version.previews[currentImageIndex.value]
  }
  return props.version.thumbnail_path
})

const timeAgo = computed(() => {
  if (!props.version.created_at) return ''
  return dayjs(props.version.created_at).fromNow()
})

const addedCount = computed(() => {
  const stats = props.version.diff_stats || '+0/-0'
  const match = stats.match(/\+(\d+)/)
  return match ? parseInt(match[1]) : 0
})

const removedCount = computed(() => {
  const stats = props.version.diff_stats || '+0/-0'
  const match = stats.match(/-(\d+)/)
  return match ? parseInt(match[1]) : 0
})

const hasParamChanges = computed(() => {
  // 简化判断，实际应该对比前后版本
  return false
})

const diffDetails = computed(() => {
  // 模拟变更详情，实际应该从 snapshot 对比生成
  const details = []
  
  if (addedCount.value > 0) {
    // 模拟添加的词
    details.push(
      { type: 'added', text: 'sitting' },
      { type: 'added', text: 'classroom' }
    )
  }
  
  if (removedCount.value > 0) {
    // 模拟删除的词
    details.push(
      { type: 'removed', text: 'standing' },
      { type: 'removed', text: 'outdoor' }
    )
  }
  
  return details
})

function formatFullTime(timeStr) {
  if (!timeStr) return ''
  return dayjs(timeStr).format('YYYY-MM-DD HH:mm:ss')
}

function toggleStar() {
  emit('toggle-star', props.version)
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.modal-content {
  width: 100%;
  max-width: 1000px;
  max-height: 90vh;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* 高度调整手柄 */
.resize-handle {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 12px;
  cursor: ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  transition: background-color 0.2s;
  z-index: 10;
}

.resize-handle:hover {
  background-color: rgba(14, 165, 233, 0.1);
}

.resize-indicator {
  width: 60px;
  height: 4px;
  background-color: #475569;
  border-radius: 2px;
  transition: background-color 0.2s;
}

.resize-handle:hover .resize-indicator {
  background-color: #0ea5e9;
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
  gap: 12px;
  margin: 0;
}

.version-tag {
  font-size: 18px;
  font-weight: 700;
  color: white;
  padding: 4px 12px;
  background: linear-gradient(135deg, #7c3aed, #a78bfa);
  border-radius: 8px;
}

.preset-name {
  font-size: 16px;
  font-weight: 500;
  color: #e2e8f0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.star-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background-color: transparent;
  border: none;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.star-btn:hover,
.star-btn.starred {
  color: #f59e0b;
  background-color: rgba(245, 158, 11, 0.1);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
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
  grid-template-columns: 1fr 1.2fr;
  gap: 0;
}

/* 预览图区域 */
.preview-section {
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
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
}

.no-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #64748b;
}

.preview-list {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 4px;
}

.preview-thumb {
  flex-shrink: 0;
  width: 60px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  opacity: 0.6;
  transition: all 0.2s;
}

.preview-thumb:hover {
  opacity: 0.8;
}

.preview-thumb.active {
  border-color: #0ea5e9;
  opacity: 1;
}

.preview-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 信息区域 */
.info-section {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background-color: #1e293b;
  border-radius: 8px;
}

.time-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #94a3b8;
}

.time-label {
  font-size: 13px;
}

.time-value {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
  font-family: monospace;
}

.time-ago {
  font-size: 12px;
  color: #64748b;
  padding-left: 26px;
}

/* 变更统计 */
.diff-stats-bar {
  display: flex;
  gap: 12px;
}

.diff-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px;
  border-radius: 8px;
}

.diff-item.added {
  background-color: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.diff-item.removed {
  background-color: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.diff-item.changed {
  background-color: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.diff-count {
  font-size: 24px;
  font-weight: 700;
}

.diff-label {
  font-size: 12px;
}

/* 变更详情 */
.diff-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.diff-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.diff-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 13px;
}

.diff-row.added {
  background-color: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.diff-row.removed {
  background-color: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.diff-row.changed {
  background-color: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.diff-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
}

/* 提示词 */
.prompt-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.prompt-box {
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

.prompt-box.positive {
  border-left: 3px solid #22c55e;
}

.prompt-box.negative {
  border-left: 3px solid #ef4444;
}

/* 参数 */
.params-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.params-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.param-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 10px;
  background-color: #1e293b;
  border-radius: 6px;
}

.param-label {
  font-size: 11px;
  color: #64748b;
}

.param-value {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
}

/* 底部 */
.modal-footer {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid #334155;
  gap: 12px;
}

.spacer {
  flex: 1;
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

@media (max-width: 768px) {
  .modal-body {
    grid-template-columns: 1fr;
  }
  
  .preview-section {
    border-right: none;
    border-bottom: 1px solid #334155;
  }
  
  .params-grid {
    grid-template-columns: 1fr;
  }
}
</style>
