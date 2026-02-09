<template>
  <div class="preset-list-container">
    <!-- 左侧分类侧边栏 -->
    <aside class="preset-sidebar">
      <div class="sidebar-header">
        <span class="sidebar-title">预设分类</span>
        <button class="add-btn" @click="createCategory" title="新建分类">
          <PlusIcon class="w-3 h-3" />
        </button>
      </div>
      
      <div class="category-list">
        <!-- 全部预设 -->
        <div 
          class="category-item"
          :class="{ active: currentCategoryId === 0 }"
          @click="selectCategory(0)"
        >
          <SwatchIcon class="w-4 h-4 icon" />
          <span class="category-name">全部预设</span>
          <span class="category-count">{{ totalPresetCount }}</span>
        </div>
        
        <!-- 分类列表 -->
        <div
          v-for="category in presetCategories"
          :key="category.id"
          class="category-item"
          :class="{ active: currentCategoryId === category.id }"
          @click="selectCategory(category.id)"
        >
          <FolderIcon class="w-4 h-4 icon" />
          <span class="category-name">{{ category.name }}</span>
          <span class="category-count">{{ getPresetCountByCategory(category.id) }}</span>
        </div>
        
        <!-- 无分类 -->
        <div 
          class="category-item"
          :class="{ active: currentCategoryId === -1 }"
          @click="selectCategory(-1)"
        >
          <QuestionMarkCircleIcon class="w-4 h-4 icon" />
          <span class="category-name">未分类</span>
          <span class="category-count">{{ uncategorizedCount }}</span>
        </div>
      </div>
    </aside>
    
    <!-- 右侧预设列表 -->
    <div class="preset-list-view">
      <!-- 头部工具栏 -->
      <div class="toolbar">
        <h2 class="page-title">
          <SwatchIcon class="w-5 h-5 text-amber-400" />
          {{ currentCategoryName }}
          <span class="count">({{ filteredPresets.length }})</span>
        </h2>
        <div class="toolbar-actions">
          <div class="search-box">
            <MagnifyingGlassIcon class="search-icon" />
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="搜索预设..."
            />
          </div>
          <button class="btn-primary" @click="createPreset">
            <PlusIcon class="w-4 h-4" />
            新建预设
          </button>
        </div>
      </div>
      
      <!-- 预设卡片网格 -->
      <div class="preset-grid">
        <PresetCard
          v-for="preset in filteredPresets"
          :key="preset.id"
          :preset="preset"
          @view="viewPreset"
          @edit="editPreset"
          @use="usePreset"
          @delete="deletePreset"
          @update-thumbnail="updateThumbnail"
        />
        
        <!-- 空状态 -->
        <div v-if="filteredPresets.length === 0" class="empty-state">
          <SwatchIcon class="w-16 h-16 text-slate-600" />
          <p>暂无预设</p>
          <span>{{ emptyStateText }}</span>
          <button class="btn-primary" @click="createPreset">
            <PlusIcon class="w-4 h-4" />
            创建预设
          </button>
        </div>
      </div>
      
      <!-- CURD Dialogs -->
      <PresetDialog
        v-if="showPresetDialog"
        :preset="editingPreset"
        :default-category-id="currentCategoryId > 0 ? currentCategoryId : 0"
        @close="closePresetDialog"
        @save="savePreset"
        @delete="deletePreset"
      />
      
      <!-- 预设详情弹窗 -->
      <PresetDetailModal
        v-if="viewingPreset"
        :preset="viewingPreset"
        @close="viewingPreset = null"
        @edit="editPreset(viewingPreset)"
        @use="usePreset(viewingPreset)"
        @view-history="viewPresetHistory"
      />
      
      <!-- 分类对话框 -->
      <CategoryDialog
        v-if="showCategoryDialog"
        :category="{ type: 'PRESET', parent_id: 0 }"
        @close="showCategoryDialog = false"
        @save="saveCategory"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import {
  SwatchIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  FolderIcon,
  QuestionMarkCircleIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, usePresetStore, useCategoryStore, useVersionStore, useImageStore } from '../stores'
import PresetCard from './PresetCard.vue'
import PresetDialog from './PresetDialog.vue'
import PresetDetailModal from './PresetDetailModal.vue'
import CategoryDialog from './CategoryDialog.vue'

const appStore = useAppStore()
const presetStore = usePresetStore()
const categoryStore = useCategoryStore()
const versionStore = useVersionStore()
const imageStore = useImageStore()

const { activePresets: presets } = storeToRefs(presetStore)
const { categories } = storeToRefs(categoryStore)

// 分类筛选
const currentCategoryId = ref(0) // 0=全部, -1=未分类, >0=具体分类

const presetCategories = computed(() =>
  categories.value.filter(c => c.type === 'PRESET')
)

const currentCategoryName = computed(() => {
  if (currentCategoryId.value === 0) return '全部预设'
  if (currentCategoryId.value === -1) return '未分类'
  const cat = presetCategories.value.find(c => c.id === currentCategoryId.value)
  return cat?.name || '预设'
})

const totalPresetCount = computed(() => presets.value.length)

const uncategorizedCount = computed(() => 
  presets.value.filter(p => !p.category_id || p.category_id === 0).length
)

function getPresetCountByCategory(categoryId) {
  return presets.value.filter(p => p.category_id === categoryId).length
}

async function selectCategory(categoryId) {
  currentCategoryId.value = categoryId
  // 根据分类加载预设
  if (categoryId === 0) {
    await presetStore.fetchPresets(1, 20, 0)
  } else if (categoryId === -1) {
    // 未分类，需要在客户端筛选
    await presetStore.fetchPresets(1, 100, 0)
  } else {
    await presetStore.fetchPresets(1, 20, categoryId)
  }
}

// 搜索和筛选
const searchQuery = ref('')

const filteredPresets = computed(() => {
  let result = presets.value
  
  // 分类筛选
  if (currentCategoryId.value === -1) {
    result = result.filter(p => !p.category_id || p.category_id === 0)
  } else if (currentCategoryId.value > 0) {
    result = result.filter(p => p.category_id === currentCategoryId.value)
  }
  
  // 搜索筛选
  if (!searchQuery.value) return result
  const query = searchQuery.value.toLowerCase()
  return result.filter(p => 
    p.title.toLowerCase().includes(query) ||
    p.pos_text?.toLowerCase().includes(query)
  )
})

const emptyStateText = computed(() => {
  if (searchQuery.value) return '没有匹配的预设'
  if (currentCategoryId.value > 0) return '该分类下暂无预设'
  if (currentCategoryId.value === -1) return '没有未分类的预设'
  return '点击上方按钮创建你的第一个预设'
})

// 对话框状态
const showPresetDialog = ref(false)
const showCategoryDialog = ref(false)
const editingPreset = ref(null)
const viewingPreset = ref(null)

// CURD 操作
function createPreset() {
  editingPreset.value = null
  showPresetDialog.value = true
}

function createCategory() {
  showCategoryDialog.value = true
}

async function saveCategory(data) {
  await categoryStore.createCategory(data.name, 'PRESET', 0)
  showCategoryDialog.value = false
}

function editPreset(preset) {
  editingPreset.value = preset
  showPresetDialog.value = true
  viewingPreset.value = null
}

function closePresetDialog() {
  showPresetDialog.value = false
  editingPreset.value = null
}

async function savePreset(data) {
  console.log('[PresetList] savePreset received:', {
    id: data.id,
    thumbnail: data.thumbnail,
    previews: data.previews,
    previewsLength: data.previews?.length,
  })
  
  try {
    if (data.id) {
      // 获取原始 preset 数据进行对比
      const originalPreset = editingPreset.value
      
      // 判断核心数据是否有变化（排除预览图相关字段）
      const hasCoreDataChanged = 
        data.title !== originalPreset.title ||
        data.category_id !== originalPreset.category_id ||
        data.pos_text !== originalPreset.pos_text ||
        data.neg_text !== originalPreset.neg_text ||
        JSON.stringify(data.params) !== JSON.stringify(originalPreset.params) ||
        JSON.stringify(data.loras) !== JSON.stringify(originalPreset.loras)
      
      // 判断预览图是否有变化
      const hasPreviewChanged = 
        data.thumbnail !== originalPreset.thumbnail ||
        JSON.stringify(data.previews) !== JSON.stringify(originalPreset.previews)
      
      console.log('[PresetList] change detection:', {
        hasCoreDataChanged,
        hasPreviewChanged,
      })
      
      // 更新预设基本信息
      if (data.title !== originalPreset.title || data.category_id !== originalPreset.category_id) {
        await presetStore.updatePreset(data.id, data.title, data.category_id)
      }
      
      // 处理预览图上传（无论是否只改预览图，都需要上传新图片）
      // 分离新图片(base64)和已有图片路径
      const newImages = []  // data:image/xxx;base64,xxxxx
      const existingPaths = []  // /images/xxx.png
      for (const p of data.previews || []) {
        if (p && p.startsWith('data:')) {
          // 新上传的 base64 图片（完整 data URL）
          newImages.push(p)
        } else if (p && p.startsWith('/images/')) {
          // 已有的图片路径
          existingPaths.push(p)
        } else if (p) {
          // 其他未知格式，记录日志
          console.warn('[PresetList] Unknown preview format:', p?.substring(0, 50))
        }
      }
      
      console.log('[PresetList] image classification:', {
        newImagesCount: newImages.length,
        existingPathsCount: existingPaths.length,
        existingPaths: existingPaths,
      })
      
      // 上传新图片到后端
      let uploadedPaths = []
      for (const base64Data of newImages) {
        const resp = await imageStore.uploadImageBase64(base64Data, data.id)
        console.log('[PresetList] uploadImageBase64 response:', resp)
        if (resp.success && resp.data?.file_path) {
          uploadedPaths.push(resp.data.file_path)
        }
      }
      
      // 合并所有图片路径（已有 + 新上传）
      const allPreviewPaths = [...existingPaths, ...uploadedPaths]
      
      // 确定封面路径
      let thumbnailPath = data.thumbnail || ''
      if (thumbnailPath.startsWith('data:') || thumbnailPath.length > 100) {
        // 如果封面是新上传的图片，使用第一张新上传的图片作为封面
        thumbnailPath = uploadedPaths[0] || allPreviewPaths[0] || ''
      }
      
      if (hasCoreDataChanged) {
        // 核心数据有变化，创建新版本
        console.log('[PresetList] creating version with:', {
          thumbnailPath,
          allPreviewPaths,
          uploadedPaths,
        })
        
        await versionStore.createVersion(
          data.id,
          {
            pos_text: data.pos_text,
            neg_text: data.neg_text,
            atom_ids: [],
            params: {
              ...data.params,
              loras: data.loras || [],
            },
          },
          thumbnailPath,
          allPreviewPaths
        )
      } else if (hasPreviewChanged) {
        // 只有预览图变化，更新当前版本的预览图（不创建新版本）
        console.log('[PresetList] updating preview only:', {
          thumbnailPath,
          allPreviewPaths,
        })
        
        await versionStore.updateVersionPreview(
          data.id,
          thumbnailPath,
          allPreviewPaths
        )
      }
    } else {
      // 新建预设
      await presetStore.createPreset(
        data.title,
        data.category_id,
        data.pos_text,
        data.neg_text,
        [],
        data.params,
        data.loras,
        data.previews || []
      )
    }
    closePresetDialog()
    await presetStore.fetchPresets(1, 20, currentCategoryId.value > 0 ? currentCategoryId.value : 0)
  } catch (error) {
    console.error('Failed to save preset:', error)
    alert('保存预设失败: ' + error.message)
  }
}

async function deletePreset(id) {
  await presetStore.softDeletePreset(id)
  await presetStore.fetchPresets(1, 20, currentCategoryId.value > 0 ? currentCategoryId.value : 0)
}

function viewPreset(preset) {
  viewingPreset.value = preset
}

function usePreset(preset) {
  appStore.setCurrentPreset(preset)
  if (preset.pos_text) {
    const atomTexts = preset.pos_text.split(',').map(s => s.trim()).filter(Boolean)
    console.log('使用预设:', preset.title, '原子词:', atomTexts)
  }
}

async function viewPresetHistory(preset) {
  viewingPreset.value = null
  appStore.setCurrentPreset(preset)
  await versionStore.fetchVersions(preset.id)
  appStore.showTimeline = true
}

async function updateThumbnail(presetId, thumbnailUrl, newPreviews = null) {
  const preset = presets.value.find(p => p.id === presetId)
  if (!preset) return
  
  console.log('[PresetList] updateThumbnail:', {
    presetId,
    thumbnailUrl: thumbnailUrl?.substring(0, 50),
    newPreviews: newPreviews?.map(p => p?.substring(0, 50)),
  })
  
  try {
    // 处理新上传的图片（data URL 格式）
    const newImages = []
    const existingPaths = []
    
    for (const p of newPreviews || []) {
      if (p && p.startsWith('data:')) {
        newImages.push(p)
      } else if (p && p.startsWith('/images/')) {
        existingPaths.push(p)
      }
    }
    
    // 上传新图片
    let uploadedPaths = []
    for (const base64Data of newImages) {
      const resp = await imageStore.uploadImageBase64(base64Data, presetId)
      if (resp.success && resp.data?.file_path) {
        uploadedPaths.push(resp.data.file_path)
      }
    }
    
    // 合并所有图片路径
    const allPreviewPaths = [...existingPaths, ...uploadedPaths]
    
    // 确定封面路径
    let thumbnailPath = thumbnailUrl
    if (thumbnailPath.startsWith('data:')) {
      // 如果封面是新上传的图片，使用第一张新上传的图片
      thumbnailPath = uploadedPaths[0] || allPreviewPaths[0] || ''
    }
    
    // 更新本地状态
    preset.thumbnail = thumbnailPath
    preset.previews = allPreviewPaths
    
    // 保存到后端（只更新预览图，不创建新版本）
    await versionStore.updateVersionPreview(presetId, thumbnailPath, allPreviewPaths)
    
    console.log('[PresetList] updateThumbnail saved:', {
      thumbnailPath,
      allPreviewPaths,
    })
  } catch (error) {
    console.error('[PresetList] updateThumbnail failed:', error)
    alert('保存预览图失败: ' + error.message)
  }
}

onMounted(async () => {
  await categoryStore.fetchCategories()
})
</script>

<style scoped>
.preset-list-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 左侧分类侧边栏 */
.preset-sidebar {
  width: 200px;
  background-color: #0f172a;
  border-right: 1px solid #1e293b;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 12px;
  border-bottom: 1px solid #1e293b;
}

.sidebar-title {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.add-btn {
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

.add-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.category-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 2px;
}

.category-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.category-item.active {
  background-color: #0284c7;
}

.category-item .icon {
  flex-shrink: 0;
  color: #94a3b8;
}

.category-item.active .icon {
  color: white;
}

.category-name {
  flex: 1;
  font-size: 13px;
  color: #e2e8f0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-item.active .category-name {
  color: white;
}

.category-count {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: #94a3b8;
}

.category-item.active .category-count {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
}

/* 右侧预设列表 */
.preset-list-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #020617;
  overflow: hidden;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #1e293b;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.page-title .count {
  font-size: 14px;
  color: #64748b;
  font-weight: 400;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 8px 12px;
  width: 240px;
}

.search-icon {
  width: 18px;
  height: 18px;
  color: #64748b;
  flex-shrink: 0;
}

.search-box input {
  flex: 1;
  background: transparent;
  border: none;
  color: #e2e8f0;
  font-size: 14px;
  outline: none;
}

.search-box input::placeholder {
  color: #64748b;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 8px;
  background-color: #0284c7;
  border: none;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  background-color: #0ea5e9;
}

.preset-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  grid-auto-rows: min-content;
  gap: 20px;
  padding: 20px;
  overflow-y: auto;
  align-items: start;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 80px 20px;
  color: #64748b;
}

.empty-state p {
  font-size: 16px;
  font-weight: 500;
  color: #94a3b8;
}

.empty-state span {
  font-size: 14px;
}
</style>
