<template>
  <div class="lora-tag-cleaner">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <h2 class="title">
          <SparklesIcon class="w-5 h-5 text-amber-400" />
          LoRA 标签清洗工具
        </h2>
        <span class="subtitle">上传图片和对应的 TXT 标签文件进行清洗</span>
      </div>
      <div class="toolbar-actions">
        <button class="btn-secondary" @click="showUploadDialog = true">
          <FolderOpenIcon class="w-4 h-4" />
          上传文件
        </button>
        <button 
          v-if="filePairs.length > 0"
          class="btn-danger" 
          @click="cleanupFilePairs"
        >
          <TrashIcon class="w-4 h-4" />
          清空
        </button>
        <button 
          class="btn-primary" 
          @click="exportFiles"
          :disabled="filePairs.length === 0"
        >
          <ArchiveBoxIcon class="w-4 h-4" />
          导出全部
        </button>
      </div>
    </div>

    <!-- 主工作区 -->
    <div v-if="filePairs.length > 0" class="workspace">
      <!-- 左侧：文件列表和图片预览 -->
      <div class="left-panel">
        <div class="panel-header">
          <span class="panel-title">图片列表 ({{ filePairs.length }})</span>
        </div>
        <div class="file-list">
          <div
            v-for="(pair, index) in filePairs"
            :key="pair.id"
            class="file-item"
            :class="{ active: currentIndex === index, 'no-tags': !pair.hasTags }"
            @click="selectFile(index)"
          >
            <div class="file-thumb-wrapper">
              <img 
                v-if="!pair.imageError"
                :src="pair.preview" 
                class="file-thumb" 
                @error="handleImageError(pair)"
              />
              <div v-else class="file-thumb-error">
                <PhotoIcon class="w-6 h-6" />
              </div>
            </div>
            <div class="file-info">
              <span class="file-name">{{ pair.name }}</span>
              <span class="file-status">
                <TagIcon v-if="pair.hasTags" class="w-3 h-3" />
                <span v-else class="no-tag-indicator">无标签</span>
              </span>
            </div>
            <button 
              class="btn-icon remove-file"
              @click.stop="removeFilePair(index)"
              title="删除"
            >
              <XMarkIcon class="w-3 h-3" />
            </button>
          </div>
        </div>
      </div>

      <!-- 中间：TAG 列表编辑 -->
      <div class="center-panel">
        <div class="panel-header">
          <span class="panel-title">TAG 编辑</span>
          <div class="header-actions">
            <button class="btn-icon" @click="addEmptyTag" title="添加 TAG">
              <PlusIcon class="w-4 h-4" />
            </button>
            <button class="btn-icon danger" @click="clearAllTags" title="清空 TAG">
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
        
        <div v-if="currentPair" class="tag-editor">
          <div class="current-image">
            <img 
              v-if="!currentPair.imageError"
              :src="currentPair.preview" 
              @error="handleImageError(currentPair)"
            />
            <div v-else class="current-image-error">
              <PhotoIcon class="w-16 h-16" />
              <span>图片加载失败</span>
            </div>
          </div>
          
          <div class="tag-list-header">
            <span class="col-tag">TAG</span>
            <span class="col-trans">翻译</span>
            <span class="col-action">操作</span>
          </div>
          
          <draggable
            v-model="currentTags"
            item-key="id"
            handle=".drag-handle"
            class="tag-list"
          >
            <template #item="{ element, index }">
              <div class="tag-item">
                <div class="drag-handle">
                  <Bars3Icon class="w-4 h-4" />
                </div>
                <input
                  v-model="element.tag"
                  class="tag-input"
                  placeholder="输入 TAG"
                  @change="translateTag(index)"
                />
                <input
                  v-model="element.translation"
                  class="trans-input"
                  placeholder="翻译"
                  readonly
                />
                <button class="btn-icon danger" @click="removeTag(index)">
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </template>
          </draggable>
          
          <div v-if="currentTags.length === 0" class="empty-tags">
            <p>暂无 TAG</p>
            <button class="btn-link" @click="addEmptyTag">添加第一个 TAG</button>
          </div>
        </div>
        
        <div v-else class="empty-state">
          <PhotoIcon class="w-12 h-12" />
          <p>请选择左侧图片开始编辑</p>
        </div>
      </div>

      <!-- 右侧：TAG 管理和翻译 -->
      <div class="right-panel">
        <!-- 上半部分：所有 TAG -->
        <div class="top-section">
          <div class="panel-header">
            <span class="panel-title">所有 TAG ({{ allUniqueTags.length }})</span>
          </div>
          <div class="all-tags-list">
            <div
              v-for="item in allUniqueTags"
              :key="item.tag"
              class="tag-chip"
              :class="{ active: isTagInCurrent(item.tag) }"
              @click="addTagToCurrent(item.tag, item.translation)"
            >
              <span class="tag-text">{{ item.tag }}</span>
              <span class="tag-count">{{ item.count }}</span>
            </div>
            <div v-if="allUniqueTags.length === 0" class="empty-hint">
              暂无 TAG 数据
            </div>
          </div>
        </div>

        <!-- 下半部分：新增 TAG -->
        <div class="bottom-section">
          <div class="panel-header">
            <span class="panel-title">新增 TAG</span>
          </div>
          <div class="new-tag-form">
            <div class="form-group">
              <label>TAG 名称</label>
              <input
                v-model="newTag.tag"
                type="text"
                placeholder="如：solo"
                @keyup.enter="focusTranslation"
              />
            </div>
            <div class="form-group">
              <label>翻译</label>
              <input
                ref="transInput"
                v-model="newTag.translation"
                type="text"
                placeholder="如：单人"
                @keyup.enter="addNewTag"
              />
            </div>
            <div class="form-options">
              <label class="checkbox">
                <input v-model="newTag.applyToAll" type="checkbox" />
                <span>添加到所有图片</span>
              </label>
            </div>
            <button class="btn-primary w-full" @click="addNewTag">
              <PlusIcon class="w-4 h-4" />
              添加 TAG
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-workspace">
      <FolderOpenIcon class="w-16 h-16 text-slate-600" />
      <h3>开始标签清洗</h3>
      <p>上传图片和对应的 TXT 标签文件</p>
      <button class="btn-primary" @click="showUploadDialog = true">
        <FolderOpenIcon class="w-4 h-4" />
        选择文件
      </button>
      <div class="upload-hint">
        <p>提示：图片文件（.jpg/.png）和标签文件（.txt）需要同名</p>
        <p>例如：0001.jpg 对应 0001.txt</p>
      </div>
    </div>

    <!-- 上传对话框 -->
    <div v-if="showUploadDialog" class="modal-overlay" @click.self="showUploadDialog = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>上传文件</h3>
          <button class="btn-icon" @click="showUploadDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="modal-body">
          <div 
            class="upload-area"
            :class="{ dragging: isDragging }"
            @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="handleDrop"
            @click="triggerFileInput"
          >
            <input
              ref="fileInput"
              type="file"
              multiple
              accept=".jpg,.jpeg,.png,.webp,.txt"
              class="hidden"
              @change="handleFileSelect"
            />
            <CloudArrowUpIcon class="w-12 h-12 text-slate-500" />
            <p>点击或拖拽文件到此处</p>
            <span class="hint">支持 .jpg, .png, .webp 和对应的 .txt 文件</span>
          </div>
          
          <div v-if="uploadPreview.length > 0" class="upload-preview">
            <p class="preview-title">待上传文件 ({{ uploadPreview.length }})</p>
            <div class="preview-list">
              <div
                v-for="file in uploadPreview"
                :key="file.name"
                class="preview-item"
              >
                <PhotoIcon v-if="file.type.startsWith('image')" class="w-4 h-4" />
                <DocumentTextIcon v-else class="w-4 h-4" />
                <span class="preview-name">{{ file.name }}</span>
                <span class="preview-size">{{ formatSize(file.size) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="showUploadDialog = false">取消</button>
          <button 
            class="btn-primary" 
            @click="confirmUpload"
            :disabled="uploadPreview.length === 0"
          >
            确认上传 ({{ uploadPreview.filter(f => f.type.startsWith('image')).length }} 张图片)
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onUnmounted } from 'vue'
import draggable from 'vuedraggable'
import JSZip from 'jszip'
import {
  SparklesIcon,
  FolderOpenIcon,
  ArchiveBoxIcon,
  PlusIcon,
  TrashIcon,
  XMarkIcon,
  Bars3Icon,
  PhotoIcon,
  CloudArrowUpIcon,
  DocumentTextIcon,
  TagIcon,
} from '@heroicons/vue/24/outline'

// 文件对数据结构
class FilePair {
  constructor(id, name, imageFile, textFile = null) {
    this.id = id
    this.name = name
    this.imageFile = imageFile
    this.textFile = textFile
    this.preview = URL.createObjectURL(imageFile)
    this.hasTags = false
    this.tags = [] // { id, tag, translation }
    this.imageError = false
  }
}

// 状态
const filePairs = ref([])
const currentIndex = ref(-1)
const showUploadDialog = ref(false)
const isDragging = ref(false)
const uploadPreview = ref([])
const fileInput = ref(null)
const transInput = ref(null)

// 新增 TAG 表单
const newTag = ref({
  tag: '',
  translation: '',
  applyToAll: false
})

// 计算属性
const currentPair = computed(() => {
  if (currentIndex.value >= 0 && currentIndex.value < filePairs.value.length) {
    return filePairs.value[currentIndex.value]
  }
  return null
})

const currentTags = computed({
  get() {
    return currentPair.value?.tags || []
  },
  set(val) {
    if (currentPair.value) {
      currentPair.value.tags = val
    }
  }
})

// 所有唯一的 TAG 及其计数
const allUniqueTags = computed(() => {
  const tagMap = new Map()
  
  for (const pair of filePairs.value) {
    for (const item of pair.tags) {
      if (!tagMap.has(item.tag)) {
        tagMap.set(item.tag, { tag: item.tag, translation: item.translation, count: 0 })
      }
      tagMap.get(item.tag).count++
    }
  }
  
  return Array.from(tagMap.values()).sort((a, b) => b.count - a.count)
})

// 选择文件
function selectFile(index) {
  currentIndex.value = index
}

// 触发文件选择
function triggerFileInput() {
  fileInput.value?.click()
}

// 处理文件选择
function handleFileSelect(e) {
  const files = Array.from(e.target.files)
  addFilesToPreview(files)
}

// 处理拖拽
function handleDrop(e) {
  isDragging.value = false
  const files = Array.from(e.dataTransfer.files)
  addFilesToPreview(files)
}

// 添加文件到预览
function addFilesToPreview(files) {
  const validFiles = files.filter(f => {
    const ext = f.name.toLowerCase().split('.').pop()
    return ['jpg', 'jpeg', 'png', 'webp', 'txt'].includes(ext)
  })
  uploadPreview.value.push(...validFiles)
}

// 格式化文件大小
function formatSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

// 确认上传
async function confirmUpload() {
  const imageFiles = uploadPreview.value.filter(f => f.type.startsWith('image'))
  const textFiles = uploadPreview.value.filter(f => f.name.toLowerCase().endsWith('.txt'))
  
  // 创建文件名映射（使用小写进行不区分大小写的匹配）
  const textMap = new Map()
  for (const textFile of textFiles) {
    const baseName = textFile.name.replace(/\.txt$/i, '').toLowerCase()
    textMap.set(baseName, textFile)
  }
  
  // 统计信息
  let matchedCount = 0
  let unmatchedCount = 0
  
  // 创建文件对
  let id = filePairs.value.length
  for (const imageFile of imageFiles) {
    const baseName = imageFile.name.replace(/\.(jpg|jpeg|png|webp)$/i, '')
    const baseNameLower = baseName.toLowerCase()
    const textFile = textMap.get(baseNameLower) || null
    
    // 检查是否已存在相同名称的文件对（避免重复）
    const existingIndex = filePairs.value.findIndex(p => p.name.toLowerCase() === baseNameLower)
    
    if (existingIndex >= 0) {
      // 更新现有文件对的文本文件
      if (textFile) {
        const pair = filePairs.value[existingIndex]
        pair.textFile = textFile
        pair.hasTags = true
        try {
          const content = await readFile(textFile)
          pair.tags = parseTags(content)
          matchedCount++
        } catch (e) {
          console.error('Failed to parse tags:', e)
        }
      }
      continue
    }
    
    // 创建新文件对
    const pair = new FilePair(id++, baseName, imageFile, textFile)
    
    // 解析 TXT 文件
    if (textFile) {
      try {
        const content = await readFile(textFile)
        pair.tags = parseTags(content)
        pair.hasTags = true
        matchedCount++
      } catch (e) {
        console.error('Failed to parse tags:', e)
        pair.hasTags = false
        unmatchedCount++
      }
    } else {
      unmatchedCount++
    }
    
    filePairs.value.push(pair)
  }
  
  // 显示上传结果提示
  if (matchedCount > 0 && unmatchedCount > 0) {
    alert(`上传完成！\n✓ ${matchedCount} 个文件已匹配标签\n○ ${unmatchedCount} 个图片无对应标签文件`)
  } else if (matchedCount > 0) {
    alert(`上传完成！成功加载 ${matchedCount} 个图片及其标签`)
  } else if (unmatchedCount > 0) {
    alert(`已上传 ${unmatchedCount} 个图片（无标签文件）\n提示：可后续上传对应的 .txt 文件补充标签`)
  }
  
  // 自动选择第一个
  if (filePairs.value.length > 0 && currentIndex.value === -1) {
    currentIndex.value = 0
  }
  
  // 清空预览并关闭对话框
  uploadPreview.value = []
  showUploadDialog.value = false
  
  // 自动翻译所有 TAG
  await translateAllTags()
}

// 读取文件内容
function readFile(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = e => resolve(e.target.result)
    reader.onerror = reject
    reader.readAsText(file)
  })
}

// 解析 TAG
function parseTags(content) {
  // 支持逗号分隔的 TAG
  const tags = content.split(/[,，]/).map(t => t.trim()).filter(Boolean)
  return tags.map((tag, index) => ({
    id: `tag-${Date.now()}-${index}`,
    tag,
    translation: ''
  }))
}

// AI 翻译 TAG（模拟，实际应调用 AI 服务）
async function translateTag(index) {
  const tag = currentTags.value[index]
  if (!tag || !tag.tag || tag.translation) return
  
  // 先从缓存查找
  const cached = findTranslationInCache(tag.tag)
  if (cached) {
    tag.translation = cached
    return
  }
  
  // TODO: 调用实际的 AI 翻译服务
  // 这里使用简单的模拟翻译
  tag.translation = await mockTranslate(tag.tag)
}

// 翻译所有 TAG
async function translateAllTags() {
  const allTags = new Set()
  for (const pair of filePairs.value) {
    for (const tag of pair.tags) {
      if (!tag.translation) {
        allTags.add(tag.tag)
      }
    }
  }
  
  for (const tagText of allTags) {
    const translation = await mockTranslate(tagText)
    // 更新所有匹配的 TAG
    for (const pair of filePairs.value) {
      for (const tag of pair.tags) {
        if (tag.tag === tagText) {
          tag.translation = translation
        }
      }
    }
  }
}

// 模拟翻译（实际项目中应替换为真实的 AI 翻译 API）
const translationCache = new Map()
const commonTranslations = {
  'solo': '单人',
  '1girl': '1个女孩',
  '1boy': '1个男孩',
  'masterpiece': '杰作',
  'best quality': '最佳质量',
  'highres': '高分辨率',
  'ultra-detailed': '超详细',
  '8k uhd': '8K超高清',
  'long hair': '长发',
  'short hair': '短发',
  'blonde hair': '金发',
  'black hair': '黑发',
  'blue eyes': '蓝眼睛',
  'red eyes': '红眼睛',
  'green eyes': '绿眼睛',
  'smile': '微笑',
  'standing': '站立',
  'sitting': '坐着',
  'school uniform': '校服',
  'dress': '连衣裙',
  'outdoors': '户外',
  'indoors': '室内',
  'anime style': '动漫风格',
  'realistic': '写实风格',
}

function findTranslationInCache(tag) {
  const lowerTag = tag.toLowerCase()
  return translationCache.get(lowerTag) || commonTranslations[lowerTag]
}

async function mockTranslate(tag) {
  const lowerTag = tag.toLowerCase()
  
  // 检查缓存
  if (translationCache.has(lowerTag)) {
    return translationCache.get(lowerTag)
  }
  
  // 检查常用翻译
  if (commonTranslations[lowerTag]) {
    translationCache.set(lowerTag, commonTranslations[lowerTag])
    return commonTranslations[lowerTag]
  }
  
  // TODO: 这里应该调用实际的 AI 翻译 API
  // 现在返回空字符串，表示需要手动填写
  return ''
}

// TAG 操作
function addEmptyTag() {
  if (!currentPair.value) return
  currentTags.value.push({
    id: `tag-${Date.now()}-${Math.random()}`,
    tag: '',
    translation: ''
  })
}

function removeTag(index) {
  currentTags.value.splice(index, 1)
}

function clearAllTags() {
  if (!currentPair.value) return
  if (confirm('确定要清空当前图片的所有 TAG 吗？')) {
    currentTags.value = []
  }
}

function isTagInCurrent(tag) {
  return currentTags.value.some(t => t.tag === tag)
}

function addTagToCurrent(tag, translation) {
  if (!currentPair.value) return
  if (isTagInCurrent(tag)) return
  
  currentTags.value.push({
    id: `tag-${Date.now()}-${Math.random()}`,
    tag,
    translation: translation || ''
  })
}

// 新增 TAG
async function addNewTag() {
  const { tag, translation, applyToAll } = newTag.value
  
  if (!tag.trim()) {
    alert('请输入 TAG 名称')
    return
  }
  
  if (applyToAll) {
    // 添加到所有图片
    for (const pair of filePairs.value) {
      if (!pair.tags.some(t => t.tag === tag)) {
        pair.tags.push({
          id: `tag-${Date.now()}-${Math.random()}`,
          tag,
          translation: translation || ''
        })
      }
    }
  } else {
    // 只添加到当前图片
    addTagToCurrent(tag, translation)
  }
  
  // 重置表单
  newTag.value = { tag: '', translation: '', applyToAll: false }
}

function focusTranslation() {
  nextTick(() => {
    transInput.value?.focus()
  })
}

// 导出文件
async function exportFiles() {
  if (filePairs.value.length === 0) return
  
  const zip = new JSZip()
  
  for (const pair of filePairs.value) {
    // 添加图片
    const imageExt = pair.imageFile.name.split('.').pop()
    zip.file(`${pair.name}.${imageExt}`, pair.imageFile)
    
    // 添加 TXT 文件（使用编辑后的 TAG）
    const tagsText = pair.tags.map(t => t.tag).join(', ')
    zip.file(`${pair.name}.txt`, tagsText)
  }
  
  // 生成并下载
  const content = await zip.generateAsync({ type: 'blob' })
  const url = URL.createObjectURL(content)
  const a = document.createElement('a')
  a.href = url
  a.download = `lora_tags_${new Date().toISOString().slice(0, 10)}.zip`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// 处理图片加载错误
function handleImageError(pair) {
  pair.imageError = true
}

// 清理所有文件对（释放 URL 对象）
function cleanupFilePairs() {
  for (const pair of filePairs.value) {
    if (pair.preview && pair.preview.startsWith('blob:')) {
      URL.revokeObjectURL(pair.preview)
    }
  }
  filePairs.value = []
}

// 删除单个文件对
function removeFilePair(index) {
  const pair = filePairs.value[index]
  if (pair.preview && pair.preview.startsWith('blob:')) {
    URL.revokeObjectURL(pair.preview)
  }
  filePairs.value.splice(index, 1)
  
  // 调整当前选中索引
  if (currentIndex.value >= filePairs.value.length) {
    currentIndex.value = filePairs.value.length - 1
  }
}

// 组件卸载时清理
onUnmounted(() => {
  cleanupFilePairs()
})
</script>

<style scoped>
.lora-tag-cleaner {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #0f172a;
}

/* 工具栏 */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid #334155;
  background-color: #0f172a;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.subtitle {
  font-size: 13px;
  color: #64748b;
}

.toolbar-actions {
  display: flex;
  gap: 10px;
}

/* 按钮样式 */
.btn-primary, .btn-secondary {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background-color: #0284c7;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #0ea5e9;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #1e293b;
  color: #e2e8f0;
  border: 1px solid #334155;
}

.btn-secondary:hover {
  background-color: #334155;
}

.btn-danger {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  background-color: transparent;
  color: #ef4444;
  border: 1px solid #ef4444;
}

.btn-danger:hover {
  background-color: rgba(239, 68, 68, 0.1);
}

.btn-icon {
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

.btn-icon:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.btn-icon.danger:hover {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.btn-link {
  background: none;
  border: none;
  color: #0ea5e9;
  cursor: pointer;
  font-size: 13px;
  text-decoration: underline;
}

.btn-link:hover {
  color: #38bdf8;
}

/* 工作区 */
.workspace {
  flex: 1;
  display: grid;
  grid-template-columns: 240px 1fr 300px;
  overflow: hidden;
}

/* 面板通用样式 */
.left-panel, .center-panel, .right-panel {
  display: flex;
  flex-direction: column;
  border-right: 1px solid #334155;
  overflow: hidden;
}

.right-panel {
  border-right: none;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #334155;
  background-color: #1e293b;
}

.panel-title {
  font-size: 13px;
  font-weight: 600;
  color: #e2e8f0;
}

.header-actions {
  display: flex;
  gap: 4px;
}

/* 左侧面板 - 文件列表 */
.file-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 4px;
}

.file-item:hover {
  background-color: #1e293b;
}

.file-item.active {
  background-color: #0ea5e9;
}

.file-item.no-tags {
  opacity: 0.7;
}

.file-item.active.no-tags {
  opacity: 1;
}

.file-thumb-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.file-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  background-color: #334155;
}

.file-thumb-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #334155;
  color: #64748b;
}

.remove-file {
  opacity: 0;
  transition: opacity 0.2s;
}

.file-item:hover .remove-file {
  opacity: 1;
}

.file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-name {
  font-size: 12px;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-status {
  display: flex;
  align-items: center;
  color: #22c55e;
}

.no-tag-indicator {
  font-size: 11px;
  color: #64748b;
}

/* 中间面板 - TAG 编辑 */
.center-panel {
  background-color: #0f172a;
}

.tag-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.current-image {
  padding: 16px;
  display: flex;
  justify-content: center;
  background-color: #1e293b;
  border-bottom: 1px solid #334155;
}

.current-image img {
  max-height: 200px;
  max-width: 100%;
  object-fit: contain;
  border-radius: 8px;
}

.current-image-error {
  height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #64748b;
}

.tag-list-header {
  display: grid;
  grid-template-columns: 30px 1fr 1fr 36px;
  gap: 8px;
  padding: 10px 16px;
  background-color: #1e293b;
  border-bottom: 1px solid #334155;
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
}

.col-tag, .col-trans {
  padding-left: 8px;
}

.tag-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.tag-item {
  display: grid;
  grid-template-columns: 30px 1fr 1fr 36px;
  gap: 8px;
  align-items: center;
  padding: 6px 8px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  margin-bottom: 6px;
}

.drag-handle {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #475569;
  cursor: grab;
}

.drag-handle:active {
  cursor: grabbing;
}

.tag-input, .trans-input {
  padding: 6px 10px;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 4px;
  color: #e2e8f0;
  font-size: 13px;
  outline: none;
}

.tag-input:focus {
  border-color: #0ea5e9;
}

.trans-input {
  color: #94a3b8;
}

.empty-tags {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #64748b;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #64748b;
}

/* 右侧面板 */
.right-panel {
  display: flex;
  flex-direction: column;
}

.top-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-bottom: 1px solid #334155;
}

.all-tags-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-content: flex-start;
}

.tag-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.tag-chip:hover {
  border-color: #0ea5e9;
  background-color: #0f172a;
}

.tag-chip.active {
  background-color: #0ea5e9;
  border-color: #0ea5e9;
}

.tag-text {
  font-size: 12px;
  color: #e2e8f0;
}

.tag-chip.active .tag-text {
  color: white;
}

.tag-count {
  font-size: 10px;
  padding: 2px 6px;
  background-color: #334155;
  border-radius: 10px;
  color: #94a3b8;
}

.tag-chip.active .tag-count {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
}

.bottom-section {
  height: 280px;
  display: flex;
  flex-direction: column;
  background-color: #1e293b;
}

.new-tag-form {
  flex: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  font-size: 12px;
  color: #94a3b8;
}

.form-group input {
  padding: 8px 12px;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 13px;
  outline: none;
}

.form-group input:focus {
  border-color: #0ea5e9;
}

.form-options {
  display: flex;
  align-items: center;
}

.checkbox {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 12px;
  color: #94a3b8;
}

.checkbox input {
  width: 16px;
  height: 16px;
  accent-color: #0ea5e9;
}

.w-full {
  width: 100%;
}

.empty-hint {
  width: 100%;
  text-align: center;
  color: #64748b;
  font-size: 12px;
  padding: 20px;
}

/* 空工作区 */
.empty-workspace {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #64748b;
}

.empty-workspace h3 {
  font-size: 18px;
  color: #e2e8f0;
  margin: 0;
}

.empty-workspace p {
  margin: 0;
}

.upload-hint {
  text-align: center;
  font-size: 12px;
  color: #475569;
  margin-top: 16px;
}

.upload-hint p {
  margin: 4px 0;
}

/* 模态框 */
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
  max-width: 600px;
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
  margin: 0;
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
  justify-content: flex-end;
}

/* 上传区域 */
.upload-area {
  border: 2px dashed #334155;
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-area:hover, .upload-area.dragging {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.05);
}

.upload-area p {
  margin: 12px 0 4px;
  color: #e2e8f0;
  font-weight: 500;
}

.upload-area .hint {
  font-size: 12px;
  color: #64748b;
}

.hidden {
  display: none;
}

.upload-preview {
  margin-top: 20px;
}

.preview-title {
  font-size: 13px;
  font-weight: 500;
  color: #e2e8f0;
  margin-bottom: 10px;
}

.preview-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 8px;
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  font-size: 12px;
}

.preview-item:not(:last-child) {
  border-bottom: 1px solid #1e293b;
}

.preview-name {
  flex: 1;
  color: #e2e8f0;
}

.preview-size {
  color: #64748b;
}
</style>
