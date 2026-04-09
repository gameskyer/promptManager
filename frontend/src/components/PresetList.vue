<template>
  <div class="preset-list-container">
    <!-- 预设列表视图 -->
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
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import {
  SwatchIcon,
  PlusIcon,
  MagnifyingGlassIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, usePresetStore, useCategoryStore, useVersionStore, useImageStore } from '../stores'
import PresetCard from './PresetCard.vue'
import PresetDialog from './PresetDialog.vue'
import PresetDetailModal from './PresetDetailModal.vue'

const props = defineProps({
  selectedCategory: {
    type: Object,
    default: () => ({ categoryId: 0, subCategoryId: null, childIds: [] })
  }
})

const appStore = useAppStore()
const presetStore = usePresetStore()
const categoryStore = useCategoryStore()
const versionStore = useVersionStore()
const imageStore = useImageStore()

const { activePresets: presets } = storeToRefs(presetStore)
const { categories } = storeToRefs(categoryStore)

// 当前选中的分类
const currentCategoryId = ref(0) // 0=全部, -1=未分类, >0=具体分类

// 监听 props 变化
watch(() => props.selectedCategory, (newVal) => {
  if (newVal) {
    if (newVal.subCategoryId) {
      // 选中二级分类
      currentCategoryId.value = newVal.subCategoryId
    } else if (newVal.categoryId !== undefined) {
      // 选中一级分类或全部
      currentCategoryId.value = newVal.categoryId
    }
  }
}, { immediate: true, deep: true })

// 根据ID获取分类名称
function getCategoryNameById(categoryId) {
  if (categoryId === 0) return '全部预设'
  if (categoryId === -1) return '未分类'
  const cat = categories.value.find(c => c.id === categoryId)
  return cat?.name || '预设'
}

const currentCategoryName = computed(() => {
  return getCategoryNameById(currentCategoryId.value)
})

// 获取分类的子分类
function getPresetChildren(parentId) {
  return categories.value.filter(c => c.type === 'PRESET' && c.parent_id === parentId)
}

// 搜索和筛选
const searchQuery = ref('')

const filteredPresets = computed(() => {
  let result = presets.value
  
  // 分类筛选
  if (currentCategoryId.value === -1) {
    result = result.filter(p => !p.category_id || p.category_id === 0)
  } else if (currentCategoryId.value > 0) {
    // 检查是否是一级分类
    const isRootCategory = getPresetChildren(0).some(c => c.id === currentCategoryId.value)
    
    if (isRootCategory) {
      // 一级分类：包含该分类及其子分类的预设
      const children = getPresetChildren(currentCategoryId.value)
      const categoryIds = [currentCategoryId.value, ...children.map(c => c.id)]
      result = result.filter(p => categoryIds.includes(p.category_id))
    } else {
      // 二级分类：只包含该分类的预设
      result = result.filter(p => p.category_id === currentCategoryId.value)
    }
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
const editingPreset = ref(null)
const viewingPreset = ref(null)

// CURD 操作
function createPreset() {
  editingPreset.value = null
  showPresetDialog.value = true
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
    existingPreviews: preset.previews?.map(p => p?.substring(0, 50)),
  })
  
  try {
    // 如果没有传入 newPreviews，只是切换当前显示的封面图
    // 保持原有的 previews 列表不变
    if (!newPreviews) {
      // 只更新封面图（thumbnail），保持 previews 不变
      preset.thumbnail = thumbnailUrl
      
      // 保存到后端（只更新封面图）
      await versionStore.updateVersionPreview(presetId, thumbnailUrl, preset.previews)
      
      console.log('[PresetList] updateThumbnail (switch only):', {
        thumbnailUrl,
        previews: preset.previews,
      })
      return
    }
    
    // 处理新上传的图片（data URL 格式）
    const newImages = []
    const existingPaths = []
    
    for (const p of newPreviews) {
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

/* 预设列表视图 */
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
