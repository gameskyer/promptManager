<template>
  <header class="top-bar">
    <div class="logo-section">
      <div class="logo">
        <SparklesIcon class="w-6 h-6 text-sky-400" />
        <span class="logo-text">PromptMaster</span>
      </div>
      <span class="version">v2.0</span>
    </div>
    
    <div class="search-section">
      <div class="search-box">
        <MagnifyingGlassIcon class="search-icon" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索原子词（支持拼音，如：fz 搜索服装）..."
          @keyup.enter="handleSearch"
        />
        <button v-if="searchQuery" class="clear-btn" @click="clearSearch">
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
    
    <div class="actions-section">
      <!-- 创建下拉菜单 -->
      <div class="dropdown">
        <button class="action-btn primary" @click="showCreateMenu = !showCreateMenu">
          <PlusIcon class="w-5 h-5" />
          <span>新建</span>
          <ChevronDownIcon class="w-4 h-4" />
        </button>
        <div v-if="showCreateMenu" class="dropdown-menu">
          <button @click="openAtomDialog">
            <DocumentPlusIcon class="w-4 h-4" />
            原子词
          </button>
          <button @click="openCategoryDialog">
            <FolderPlusIcon class="w-4 h-4" />
            分类
          </button>
          <button @click="openPresetDialog">
            <SwatchIcon class="w-4 h-4" />
            预设
          </button>
        </div>
      </div>
      
      <button class="action-btn" @click="showAIExplode = true" title="AI拆解">
        <BoltIcon class="w-5 h-5" />
      </button>
      
      <button class="action-btn" @click="showBackup = true" title="备份/恢复">
        <ArrowDownTrayIcon class="w-5 h-5" />
      </button>
      
      <button class="action-btn" @click="showSettings = true" title="设置">
        <Cog6ToothIcon class="w-5 h-5" />
      </button>
      

    </div>
  </header>
  
  <!-- CURD Dialogs -->
  <AtomDialog
    v-if="showAtomDialog"
    :atom="editingAtom"
    :default-category-id="currentSelectedCategoryId"
    @close="closeAtomDialog"
    @save="saveAtom"
    @delete="deleteAtom"
  />
  
  <CategoryDialog
    v-if="showCategoryDialog"
    :category="editingCategory"
    @close="closeCategoryDialog"
    @save="saveCategory"
    @delete="deleteCategory"
  />
  
  <PresetDialog
    v-if="showPresetDialog"
    :preset="editingPreset"
    @close="closePresetDialog"
    @save="savePreset"
    @delete="deletePreset"
  />
  
  <AIModal
    v-if="showAIExplode"
    @close="showAIExplode = false"
    @import="handleAIImport"
    @open-settings="showAISettings"
  />
  
  <!-- 设置弹窗 -->
  <SettingsModal
    v-if="showSettings"
    @close="showSettings = false"
  />
  
  <!-- 备份/恢复弹窗 -->
  <BackupModal
    v-if="showBackup"
    @close="showBackup = false"
    @imported="handleImported"
  />
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { 
  SparklesIcon, 
  MagnifyingGlassIcon, 
  XMarkIcon,
  PlusIcon,
  ChevronDownIcon,
  DocumentPlusIcon,
  FolderPlusIcon,
  SwatchIcon,
  BoltIcon,
  Cog6ToothIcon,
  ArrowDownTrayIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, useAtomStore, useCategoryStore, usePresetStore } from '../stores'
import AtomDialog from './AtomDialog.vue'
import CategoryDialog from './CategoryDialog.vue'
import PresetDialog from './PresetDialog.vue'
import AIModal from './AIModal.vue'
import SettingsModal from './SettingsModal.vue'
import BackupModal from './BackupModal.vue'

const appStore = useAppStore()
const atomStore = useAtomStore()
const categoryStore = useCategoryStore()
const presetStore = usePresetStore()

const { searchQuery: storeSearchQuery, currentCategory, currentSubCategory } = storeToRefs(appStore)
const searchQuery = ref('')

// 获取当前选中的分类ID（优先使用子分类）
const currentSelectedCategoryId = computed(() => {
  const id = currentSubCategory.value?.id || currentCategory.value?.id || 0
  console.log('[TopBar] currentSelectedCategoryId:', id, 
    'subCategory:', currentSubCategory.value?.name, 
    'category:', currentCategory.value?.name)
  return id
})

const showCreateMenu = ref(false)
const showAIExplode = ref(false)
const showSettings = ref(false)
const showBackup = ref(false)

// Dialog states
const showAtomDialog = ref(false)
const showCategoryDialog = ref(false)
const showPresetDialog = ref(false)
const editingAtom = ref(null)
const editingCategory = ref(null)
const editingPreset = ref(null)

watch(storeSearchQuery, (val) => {
  searchQuery.value = val
})

watch(searchQuery, (val) => {
  appStore.setSearchQuery(val)
  // 实时搜索
  if (val.trim()) {
    atomStore.fetchAtoms() // 先获取全部
  }
})

function handleSearch() {
  // 回车时触发，已由 watch 处理
}

function clearSearch() {
  searchQuery.value = ''
  appStore.setSearchQuery('')
  atomStore.fetchAtoms()
}

// Atom CURD
function openAtomDialog() {
  console.log('[TopBar] openAtomDialog called')
  console.log('[TopBar] currentCategory:', currentCategory.value)
  console.log('[TopBar] currentSubCategory:', currentSubCategory.value)
  console.log('[TopBar] currentSelectedCategoryId:', currentSelectedCategoryId.value)
  editingAtom.value = null
  showAtomDialog.value = true
  showCreateMenu.value = false
}

function closeAtomDialog() {
  showAtomDialog.value = false
  editingAtom.value = null
}

async function saveAtom(data) {
  try {
    if (data.id) {
      await atomStore.updateAtom(data.id, data)
    } else {
      await atomStore.createAtom(data)
    }
    closeAtomDialog()
    await atomStore.fetchAtoms()
  } catch (error) {
    console.error('Failed to save atom:', error)
    alert('保存失败: ' + error.message)
  }
}

async function deleteAtom(id) {
  await atomStore.deleteAtom(id)
  closeAtomDialog()
}

// Category CURD
function openCategoryDialog() {
  editingCategory.value = null
  showCategoryDialog.value = true
  showCreateMenu.value = false
}

function closeCategoryDialog() {
  showCategoryDialog.value = false
  editingCategory.value = null
}

async function saveCategory(data) {
  if (data.id) {
    await categoryStore.updateCategory(data.id, data)
  } else {
    await categoryStore.createCategory(data.name, data.type, data.parent_id)
  }
  closeCategoryDialog()
  await categoryStore.fetchCategories()
}

async function deleteCategory(id) {
  await categoryStore.deleteCategory(id)
  closeCategoryDialog()
  await categoryStore.fetchCategories()
}

// Preset CURD
function openPresetDialog() {
  editingPreset.value = null
  showPresetDialog.value = true
  showCreateMenu.value = false
}

function closePresetDialog() {
  showPresetDialog.value = false
  editingPreset.value = null
}

async function savePreset(data) {
  if (data.id) {
    await presetStore.updatePreset(data.id, data.title)
    // 更新其他字段
    const preset = presetStore.presets.find(p => p.id === data.id)
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
  closePresetDialog()
  await presetStore.fetchPresets()
}

function handleAIImport(result) {
  console.log('AI Import:', result)
}

async function handleImported() {
  // 刷新所有数据
  await atomStore.fetchAtoms()
  await categoryStore.fetchCategories()
  await presetStore.fetchPresets()
}

function showAISettings() {
  showAIExplode.value = false
  showSettings.value = true
}

if (typeof window !== 'undefined') {
  window.addEventListener('click', (e) => {
    if (!e.target.closest('.dropdown')) {
      showCreateMenu.value = false
    }
  })
}
</script>

<style scoped>
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 20px;
  background-color: #0f172a;
  border-bottom: 1px solid #1e293b;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #38bdf8, #818cf8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.version {
  font-size: 11px;
  color: #64748b;
  padding: 2px 6px;
  background-color: #1e293b;
  border-radius: 4px;
}

.search-section {
  flex: 1;
  max-width: 600px;
  margin: 0 20px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 8px 12px;
  transition: all 0.2s;
}

.search-box:focus-within {
  border-color: #0ea5e9;
  box-shadow: 0 0 0 2px rgba(14, 165, 233, 0.2);
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
  padding: 0;
}

.search-box input::placeholder {
  color: #64748b;
}

.clear-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px;
  border-radius: 4px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.actions-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dropdown {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 6px;
  min-width: 160px;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.dropdown-menu button {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  border-radius: 6px;
  background-color: transparent;
  border: none;
  color: #e2e8f0;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.dropdown-menu button:hover {
  background-color: #334155;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 6px;
  background-color: transparent;
  border: 1px solid #334155;
  color: #94a3b8;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #1e293b;
  color: #e2e8f0;
  border-color: #475569;
}

.action-btn.primary {
  background-color: #0284c7;
  border-color: #0284c7;
  color: white;
}

.action-btn.primary:hover {
  background-color: #0ea5e9;
  border-color: #0ea5e9;
}

.save-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 6px;
  background-color: #7c3aed;
  border: none;
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.save-btn:hover {
  background-color: #8b5cf6;
}
</style>
