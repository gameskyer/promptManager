<template>
  <div class="app-container">
    <!-- Top Bar -->
    <TopBar />
    
    <!-- Main Layout -->
    <div class="main-layout">
      <!-- Side Menu -->
      <SideMenu ref="sideMenu" @view-change="handleViewChange" @select-preset-category="handlePresetCategoryChange" />
      
      <!-- Main Content - 根据当前视图显示不同内容 -->
      <MainContent v-if="currentView === 'atoms'" />
      <PresetList v-else-if="currentView === 'presets'" :selected-category="selectedPresetCategory" />
      <AtomManagement v-else-if="currentView === 'atom-management'" />
      <CategoryManagement v-else-if="currentView === 'category-management'" />
      <LoraTagCleaner v-else-if="currentView === 'lora-tag-cleaner'" />
      <TransSession v-else-if="currentView === 'trans-session'" />
      
      <!-- Right Workbench -->
      <Workbench v-if="currentView === 'atoms'" />
    </div>
    
    <!-- Bottom Timeline - 在原子词或预设视图中都可以显示 -->
    <Timeline v-if="showTimeline" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useAppStore, useCategoryStore, useAtomStore, usePresetStore, useAIStore } from './stores'

import TopBar from './components/TopBar.vue'
import SideMenu from './components/SideMenu.vue'
import MainContent from './components/MainContent.vue'
import PresetList from './components/PresetList.vue'
import AtomManagement from './components/AtomManagement.vue'
import CategoryManagement from './components/CategoryManagement.vue'
import LoraTagCleaner from './components/LoraTagCleaner.vue'
import TransSession from './components/TransSession.vue'
import Workbench from './components/Workbench.vue'
import Timeline from './components/Timeline.vue'

const appStore = useAppStore()
const categoryStore = useCategoryStore()
const atomStore = useAtomStore()
const presetStore = usePresetStore()
const aiStore = useAIStore()

const { showTimeline } = storeToRefs(appStore)

const currentView = ref('atoms') // 'atoms' | 'presets' | 'atom-management' | 'category-management'
const sideMenu = ref(null)
const selectedPresetCategory = ref({ categoryId: 0, subCategoryId: null, childIds: [] })

function handleViewChange(view) {
  currentView.value = view
}

function handlePresetCategoryChange(categoryInfo) {
  selectedPresetCategory.value = categoryInfo
}

onMounted(async () => {
  // Initialize theme
  appStore.applyTheme(appStore.theme)
  
  // Initialize AI config (async)
  await aiStore.init()
  
  // Initialize data
  await categoryStore.fetchCategories()
  await atomStore.fetchAtoms()
  await presetStore.fetchPresets()
})
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #020617;
  color: #e2e8f0;
}

.main-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
}
</style>
