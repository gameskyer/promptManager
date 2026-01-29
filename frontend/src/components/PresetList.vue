<template>
  <div class="preset-list-view">
    <!-- 头部工具栏 -->
    <div class="toolbar">
      <h2 class="page-title">
        <SwatchIcon class="w-5 h-5 text-amber-400" />
        我的预设
        <span class="count">({{ presets.length }})</span>
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
        <span>点击上方按钮创建你的第一个预设</span>
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
</template>

<script setup>
import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import {
  SwatchIcon,
  PlusIcon,
  MagnifyingGlassIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, usePresetStore, useVersionStore } from '../stores'
import PresetCard from './PresetCard.vue'
import PresetDialog from './PresetDialog.vue'
import PresetDetailModal from './PresetDetailModal.vue'

const appStore = useAppStore()
const presetStore = usePresetStore()
const versionStore = useVersionStore()

const { activePresets: presets } = storeToRefs(presetStore)

const searchQuery = ref('')
const showPresetDialog = ref(false)
const editingPreset = ref(null)
const viewingPreset = ref(null)

const filteredPresets = computed(() => {
  if (!searchQuery.value) return presets.value
  const query = searchQuery.value.toLowerCase()
  return presets.value.filter(p => 
    p.title.toLowerCase().includes(query) ||
    p.pos_text?.toLowerCase().includes(query)
  )
})

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
  if (data.id) {
    await presetStore.updatePreset(data.id, data.title)
    // 更新其他字段
    const preset = presets.value.find(p => p.id === data.id)
    if (preset) {
      preset.pos_text = data.pos_text
      preset.neg_text = data.neg_text
      preset.params = data.params
      preset.loras = data.loras
    }
  } else {
    await presetStore.createPreset(
      data.title,
      data.pos_text,
      data.neg_text,
      [],
      data.params,
      data.loras,
      data.previews || []
    )
  }
  closePresetDialog()
  await presetStore.fetchPresets()
}

async function deletePreset(id) {
  await presetStore.softDeletePreset(id)
  await presetStore.fetchPresets()
}

function viewPreset(preset) {
  viewingPreset.value = preset
}

function usePreset(preset) {
  // 加载预设到工作区
  appStore.setCurrentPreset(preset)
  // 如果提示词不为空，解析并添加到工作区
  if (preset.pos_text) {
    // 解析提示词为原子词（简化处理）
    const atomTexts = preset.pos_text.split(',').map(s => s.trim()).filter(Boolean)
    console.log('使用预设:', preset.title, '原子词:', atomTexts)
  }
}

// 查看预设版本历史
async function viewPresetHistory(preset) {
  // 关闭详情弹窗
  viewingPreset.value = null
  // 设置当前预设
  appStore.setCurrentPreset(preset)
  // 加载版本历史
  await versionStore.fetchVersions(preset.id)
  // 显示时间线
  appStore.showTimeline = true
}

// 更新缩略图（不触发版本变动）
async function updateThumbnail(presetId, thumbnailUrl, newPreviews = null) {
  const preset = presets.value.find(p => p.id === presetId)
  if (preset) {
    preset.thumbnail = thumbnailUrl
    if (newPreviews) {
      preset.previews = newPreviews
    }
    // 这里调用一个不创建版本的服务方法
    await presetStore.updateThumbnailOnly(presetId, thumbnailUrl, newPreviews)
  }
}
</script>

<style scoped>
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
