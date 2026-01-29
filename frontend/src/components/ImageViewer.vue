<template>
  <div class="image-viewer-overlay" @click="$emit('close')">
    <div class="image-viewer-content" @click.stop>
      <!-- 工具栏 -->
      <div class="viewer-toolbar">
        <span class="image-counter">{{ currentIndex + 1 }} / {{ images.length }}</span>
        <div class="toolbar-actions">
          <button class="toolbar-btn" @click="zoomOut" title="缩小">
            <MagnifyingGlassMinusIcon class="w-5 h-5" />
          </button>
          <span class="zoom-level">{{ Math.round(scale * 100) }}%</span>
          <button class="toolbar-btn" @click="zoomIn" title="放大">
            <MagnifyingGlassPlusIcon class="w-5 h-5" />
          </button>
          <button class="toolbar-btn" @click="resetZoom" title="重置">
            <ArrowPathIcon class="w-5 h-5" />
          </button>
          <button class="toolbar-btn close" @click="$emit('close')" title="关闭">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </div>
      
      <!-- 图片显示区 -->
      <div class="viewer-main" @wheel="handleWheel">
        <button 
          v-if="images.length > 1"
          class="nav-btn prev"
          @click.stop="prevImage"
        >
          <ChevronLeftIcon class="w-8 h-8" />
        </button>
        
        <div 
          class="image-container"
          :style="containerStyle"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseUp"
        >
          <img 
            :src="currentImage" 
            :alt="`图片 ${currentIndex + 1}`"
            :style="imageStyle"
            @dragstart.prevent
          />
        </div>
        
        <button 
          v-if="images.length > 1"
          class="nav-btn next"
          @click.stop="nextImage"
        >
          <ChevronRightIcon class="w-8 h-8" />
        </button>
      </div>
      
      <!-- 缩略图列表 -->
      <div v-if="images.length > 1" class="viewer-thumbnails">
        <div 
          v-for="(img, index) in images"
          :key="index"
          class="thumb-item"
          :class="{ active: currentIndex === index }"
          @click="currentIndex = index"
        >
          <img :src="img" :alt="`缩略图 ${index + 1}`" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import {
  XMarkIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  MagnifyingGlassPlusIcon,
  MagnifyingGlassMinusIcon,
  ArrowPathIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  images: {
    type: Array,
    required: true,
  },
  initialIndex: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['close'])

const currentIndex = ref(props.initialIndex)
const scale = ref(1)
const translateX = ref(0)
const translateY = ref(0)
const isDragging = ref(false)
const dragStartX = ref(0)
const dragStartY = ref(0)

const currentImage = computed(() => {
  return props.images[currentIndex.value] || ''
})

const containerStyle = computed(() => ({
  cursor: isDragging.value ? 'grabbing' : 'grab',
}))

const imageStyle = computed(() => ({
  transform: `translate(${translateX.value}px, ${translateY.value}px) scale(${scale.value})`,
  transition: isDragging.value ? 'none' : 'transform 0.2s ease',
}))

// 监听初始索引变化
watch(() => props.initialIndex, (val) => {
  currentIndex.value = val
  resetZoom()
})

// 键盘导航
onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

function handleKeydown(e) {
  switch (e.key) {
    case 'Escape':
      emit('close')
      break
    case 'ArrowLeft':
      prevImage()
      break
    case 'ArrowRight':
      nextImage()
      break
    case '+':
    case '=':
      zoomIn()
      break
    case '-':
      zoomOut()
      break
    case '0':
      resetZoom()
      break
  }
}

function prevImage() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    resetZoom()
  }
}

function nextImage() {
  if (currentIndex.value < props.images.length - 1) {
    currentIndex.value++
    resetZoom()
  }
}

function zoomIn() {
  scale.value = Math.min(scale.value * 1.2, 5)
}

function zoomOut() {
  scale.value = Math.max(scale.value / 1.2, 0.2)
}

function resetZoom() {
  scale.value = 1
  translateX.value = 0
  translateY.value = 0
}

function handleWheel(e) {
  e.preventDefault()
  if (e.deltaY < 0) {
    zoomIn()
  } else {
    zoomOut()
  }
}

// 拖拽移动
function handleMouseDown(e) {
  if (scale.value <= 1) return
  isDragging.value = true
  dragStartX.value = e.clientX - translateX.value
  dragStartY.value = e.clientY - translateY.value
}

function handleMouseMove(e) {
  if (!isDragging.value) return
  translateX.value = e.clientX - dragStartX.value
  translateY.value = e.clientY - dragStartY.value
}

function handleMouseUp() {
  isDragging.value = false
}
</script>

<style scoped>
.image-viewer-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.image-viewer-content {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.viewer-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background-color: rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid #334155;
}

.image-counter {
  font-size: 14px;
  color: #94a3b8;
  font-weight: 500;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.zoom-level {
  font-size: 13px;
  color: #e2e8f0;
  min-width: 50px;
  text-align: center;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background-color: rgba(255, 255, 255, 0.1);
  border: none;
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.toolbar-btn:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.toolbar-btn.close:hover {
  background-color: rgba(239, 68, 68, 0.8);
}

.viewer-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.image-container {
  max-width: 100%;
  max-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-container img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  user-select: none;
}

.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.1);
  border: none;
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
  z-index: 10;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.nav-btn.prev {
  left: 20px;
}

.nav-btn.next {
  right: 20px;
}

.viewer-thumbnails {
  display: flex;
  gap: 8px;
  padding: 12px 20px;
  background-color: rgba(0, 0, 0, 0.5);
  border-top: 1px solid #334155;
  overflow-x: auto;
  justify-content: center;
}

.thumb-item {
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

.thumb-item:hover {
  opacity: 0.8;
}

.thumb-item.active {
  border-color: #0ea5e9;
  opacity: 1;
}

.thumb-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
