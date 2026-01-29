<template>
  <div class="image-upload">
    <!-- 上传区域 -->
    <div
      class="upload-zone"
      :class="{ 'is-dragging': isDragging, 'has-images': images.length > 0 }"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="handleDrop"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        multiple
        style="display: none"
        @change="handleFileSelect"
      />
      
      <div v-if="images.length === 0" class="upload-placeholder">
        <PhotoIcon class="w-12 h-12 text-slate-600" />
        <p>点击或拖拽上传图片</p>
        <span>支持 JPG、PNG、WebP、GIF</span>
      </div>
      
      <!-- 图片预览网格 -->
      <div v-else class="image-grid">
        <div
          v-for="(image, index) in images"
          :key="image.id || index"
          class="image-item"
          :class="{ 'is-cover': index === coverIndex }"
        >
          <img :src="image.url" alt="预览" />
          
          <!-- 封面标记 -->
          <div v-if="index === coverIndex" class="cover-badge">
            <StarIcon class="w-3 h-3" />
            封面
          </div>
          
          <!-- 操作按钮 -->
          <div class="image-actions">
            <button
              v-if="index !== coverIndex"
              class="action-btn"
              title="设为封面"
              @click.stop="setCover(index)"
            >
              <StarIcon class="w-4 h-4" />
            </button>
            <button class="action-btn delete" title="删除" @click.stop="removeImage(index)">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
        
        <!-- 添加更多按钮 -->
        <div class="add-more" @click.stop="triggerFileInput">
          <PlusIcon class="w-8 h-8 text-slate-600" />
        </div>
      </div>
    </div>
    
    <!-- 提示信息 -->
    <div v-if="error" class="error-message">
      <ExclamationTriangleIcon class="w-4 h-4" />
      {{ error }}
    </div>
    
    <!-- 上传进度 -->
    <div v-if="uploading" class="upload-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <span class="progress-text">上传中... {{ progress }}%</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  PhotoIcon,
  PlusIcon,
  XMarkIcon,
  StarIcon,
  ExclamationTriangleIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
  maxCount: {
    type: Number,
    default: 10,
  },
  maxSize: {
    type: Number,
    default: 10 * 1024 * 1024, // 10MB
  },
})

const emit = defineEmits(['update:modelValue', 'update:coverIndex'])

const fileInput = ref(null)
const isDragging = ref(false)
const uploading = ref(false)
const progress = ref(0)
const error = ref(null)

// 内部图片列表
const images = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

// 封面索引
const coverIndex = ref(0)

// 触发文件选择
function triggerFileInput() {
  fileInput.value?.click()
}

// 处理文件选择
function handleFileSelect(event) {
  const files = Array.from(event.target.files)
  processFiles(files)
  // 清空 input 以便重复选择相同文件
  event.target.value = ''
}

// 处理拖拽
function handleDrop(event) {
  isDragging.value = false
  const files = Array.from(event.dataTransfer.files)
  processFiles(files)
}

// 处理文件
async function processFiles(files) {
  error.value = null
  
  // 检查总数限制
  if (images.value.length + files.length > props.maxCount) {
    error.value = `最多上传 ${props.maxCount} 张图片`
    return
  }
  
  // 过滤图片文件
  const imageFiles = files.filter(file => file.type.startsWith('image/'))
  
  if (imageFiles.length === 0) {
    error.value = '请选择图片文件'
    return
  }
  
  uploading.value = true
  progress.value = 0
  
  try {
    for (let i = 0; i < imageFiles.length; i++) {
      const file = imageFiles[i]
      
      // 检查文件大小
      if (file.size > props.maxSize) {
        error.value = `文件 ${file.name} 超过大小限制`
        continue
      }
      
      // 转换为 base64
      const base64 = await fileToBase64(file)
      
      // 添加到列表
      images.value.push({
        id: null, // 上传后从后端获取
        url: base64,
        file: file,
        name: file.name,
      })
      
      progress.value = Math.round(((i + 1) / imageFiles.length) * 100)
    }
  } catch (err) {
    error.value = err.message || '上传失败'
  } finally {
    uploading.value = false
  }
}

// 文件转 base64
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

// 设置封面
function setCover(index) {
  coverIndex.value = index
  emit('update:coverIndex', index)
}

// 删除图片
function removeImage(index) {
  images.value.splice(index, 1)
  
  // 如果删除的是封面，重置封面索引
  if (index === coverIndex.value) {
    coverIndex.value = 0
    emit('update:coverIndex', 0)
  } else if (index < coverIndex.value) {
    coverIndex.value--
    emit('update:coverIndex', coverIndex.value)
  }
}

// 暴露方法给父组件
defineExpose({
  getImages: () => images.value,
  getCoverIndex: () => coverIndex.value,
  clear: () => {
    images.value = []
    coverIndex.value = 0
  },
})
</script>

<style scoped>
.image-upload {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.upload-zone {
  border: 2px dashed #334155;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: #0f172a;
}

.upload-zone:hover,
.upload-zone.is-dragging {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.05);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px 20px;
  color: #64748b;
}

.upload-placeholder p {
  font-size: 14px;
  font-weight: 500;
}

.upload-placeholder span {
  font-size: 12px;
}

/* 图片网格 */
.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
}

.image-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.image-item.is-cover {
  border-color: #f59e0b;
}

.image-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-badge {
  position: absolute;
  top: 4px;
  left: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background-color: #f59e0b;
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: 4px;
}

.image-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.image-item:hover .image-actions {
  opacity: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.7);
  border: none;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: rgba(0, 0, 0, 0.9);
}

.action-btn.delete:hover {
  background-color: #ef4444;
}

/* 添加更多 */
.add-more {
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1;
  border: 2px dashed #334155;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-more:hover {
  border-color: #475569;
  background-color: rgba(255, 255, 255, 0.05);
}

/* 错误信息 */
.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background-color: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 6px;
  color: #ef4444;
  font-size: 13px;
}

/* 上传进度 */
.upload-progress {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.progress-bar {
  height: 4px;
  background-color: #1e293b;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background-color: #0ea5e9;
  transition: width 0.3s;
}

.progress-text {
  font-size: 12px;
  color: #64748b;
}
</style>
