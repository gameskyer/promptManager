<template>
  <aside class="side-menu">
    <div class="menu-section">
      <div class="section-header">
        <span class="section-title">提示词库</span>
        <button class="add-btn" @click="createCategory('ATOM')" title="新建分类">
          <PlusIcon class="w-3 h-3" />
        </button>
      </div>
      <div class="menu-tree">
        <div
          v-for="category in atomRootCategories"
          :key="category.id"
          class="menu-item"
          :class="{ active: currentCategory?.id === category.id }"
        >
          <div 
            class="menu-item-header"
            @click="selectCategory(category)"
          >
            <ChevronDownIcon 
              v-if="expandedCategories.includes(category.id)"
              class="w-4 h-4 chevron"
              @click.stop="toggleExpand(category.id)"
            />
            <ChevronRightIcon 
              v-else
              class="w-4 h-4 chevron"
              @click.stop="toggleExpand(category.id)"
            />
            <FolderIcon class="w-4 h-4 icon" :class="{ 'text-sky-400': currentCategory?.id === category.id }" />
            <span class="item-label">{{ category.name }}</span>
            <button 
              class="edit-btn"
              @click.stop="editCategory(category)"
            >
              <PencilIcon class="w-3 h-3" />
            </button>
          </div>
          
          <div v-if="expandedCategories.includes(category.id)" class="sub-menu">
            <div
              v-for="child in getCategoryChildren(category.id)"
              :key="child.id"
              class="sub-menu-item"
              :class="{ active: currentSubCategory?.id === child.id }"
              @click="selectSubCategory(child)"
            >
              <DocumentIcon class="w-3 h-3 icon" />
              <span class="item-label">{{ child.name }}</span>
              <button 
                class="edit-btn"
                @click.stop="editCategory(child)"
              >
                <PencilIcon class="w-3 h-3" />
              </button>
            </div>
            <button class="add-sub-btn" @click="createSubCategory(category.id)">
              <PlusIcon class="w-3 h-3" />
              添加子分类
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="menu-divider"></div>
    
    <div class="menu-section">
      <div class="section-header">
        <span class="section-title">预设库</span>
      </div>
      <div class="menu-tree">
        <div
          class="menu-item"
          :class="{ active: currentView === 'presets' }"
          @click="selectPresetsView"
        >
          <div class="menu-item-header">
            <SwatchIcon class="w-4 h-4 icon text-amber-400" />
            <span class="item-label">我的预设</span>
            <span v-if="presetCount > 0" class="count-badge">{{ presetCount }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="menu-divider"></div>
    
    <div class="menu-section">
      <div class="section-header">
        <span class="section-title">管理</span>
      </div>
      <div class="menu-tree">
        <div
          class="menu-item"
          :class="{ active: currentView === 'atom-management' }"
          @click="selectAtomManagementView"
        >
          <div class="menu-item-header">
            <DocumentTextIcon class="w-4 h-4 icon text-emerald-400" />
            <span class="item-label">原子词管理</span>
            <span v-if="atomCount > 0" class="count-badge">{{ atomCount }}</span>
          </div>
        </div>
        <div
          class="menu-item"
          :class="{ active: currentView === 'category-management' }"
          @click="selectCategoryManagementView"
        >
          <div class="menu-item-header">
            <FolderIcon class="w-4 h-4 icon text-violet-400" />
            <span class="item-label">分类管理</span>
            <span v-if="categoryCount > 0" class="count-badge">{{ categoryCount }}</span>
          </div>
        </div>
      </div>
    </div>
  </aside>
  
  <!-- CURD Dialogs -->
  <CategoryDialog
    v-if="showCategoryDialog"
    :category="editingCategory"
    @close="closeCategoryDialog"
    @save="saveCategory"
    @delete="deleteCategory"
  />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import {
  ChevronRightIcon,
  ChevronDownIcon,
  FolderIcon,
  DocumentIcon,
  PlusIcon,
  PencilIcon,
  SwatchIcon,
  DocumentTextIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, useCategoryStore, useAtomStore, usePresetStore } from '../stores'
import CategoryDialog from './CategoryDialog.vue'

const emit = defineEmits(['view-change'])

const appStore = useAppStore()
const categoryStore = useCategoryStore()
const atomStore = useAtomStore()
const presetStore = usePresetStore()

const { currentCategory, currentSubCategory } = storeToRefs(appStore)
const { rootCategories, categories } = storeToRefs(categoryStore)
const { activePresets } = storeToRefs(presetStore)
const { atoms } = storeToRefs(atomStore)

const expandedCategories = ref([1])
const currentView = ref('atoms') // 'atoms' | 'presets' | 'atom-management' | 'category-management'

// Dialog states
const showCategoryDialog = ref(false)
const editingCategory = ref(null)

const presetCount = computed(() => activePresets.value.length)

const atomRootCategories = computed(() =>
  rootCategories.value.filter(c => c.type === 'ATOM')
)

const atomCount = computed(() => atoms.value.length)
const categoryCount = computed(() => categories.value.length)

function toggleExpand(categoryId) {
  const index = expandedCategories.value.indexOf(categoryId)
  if (index === -1) {
    expandedCategories.value.push(categoryId)
  } else {
    expandedCategories.value.splice(index, 1)
  }
}

function selectCategory(category) {
  currentView.value = 'atoms'
  emit('view-change', 'atoms')
  appStore.setCategory(category)
  appStore.setSubCategory(null)
  
  if (!expandedCategories.value.includes(category.id)) {
    expandedCategories.value.push(category.id)
  }
}

async function selectSubCategory(subCategory) {
  currentView.value = 'atoms'
  emit('view-change', 'atoms')
  appStore.setSubCategory(subCategory)
  await atomStore.fetchAtoms(subCategory.id)
}

function selectPresetsView() {
  currentView.value = 'presets'
  emit('view-change', 'presets')
  appStore.setCurrentPreset(null)
  appStore.setCategory(null)
  appStore.setSubCategory(null)
  appStore.activeTab = 'presets'
}

function selectAtomManagementView() {
  currentView.value = 'atom-management'
  emit('view-change', 'atom-management')
  appStore.setCategory(null)
  appStore.setSubCategory(null)
}

function selectCategoryManagementView() {
  currentView.value = 'category-management'
  emit('view-change', 'category-management')
  appStore.setCategory(null)
  appStore.setSubCategory(null)
}

// Category CURD
function createCategory(type) {
  editingCategory.value = { type, parent_id: 0 }
  showCategoryDialog.value = true
}

function createSubCategory(parentId) {
  editingCategory.value = { type: 'ATOM', parent_id: parentId }
  showCategoryDialog.value = true
}

function editCategory(category) {
  editingCategory.value = category
  showCategoryDialog.value = true
}

function closeCategoryDialog() {
  showCategoryDialog.value = false
  editingCategory.value = null
}

async function saveCategory(data) {
  if (data.id) {
    await categoryStore.updateCategory(data.id, data)
  } else {
    const newCategory = await categoryStore.createCategory(data.name, data.type, data.parent_id)
    // 如果是子分类，自动展开父分类
    if (newCategory.parent_id > 0 && !expandedCategories.value.includes(newCategory.parent_id)) {
      expandedCategories.value.push(newCategory.parent_id)
    }
  }
  closeCategoryDialog()
}

async function deleteCategory(id) {
  await categoryStore.deleteCategory(id)
  closeCategoryDialog()
}

function getCategoryChildren(parentId) {
  return categoryStore.categories.filter(c => c.parent_id === parentId)
}

onMounted(async () => {
  // 确保数据已加载（App.vue 已调用 fetchCategories）
  if (categoryStore.categories.length === 0) {
    await categoryStore.fetchCategories()
  }
  
  if (atomRootCategories.value.length > 0 && expandedCategories.value.length === 0) {
    expandedCategories.value.push(atomRootCategories.value[0].id)
  }
})
</script>

<style scoped>
.side-menu {
  width: 200px;
  background-color: #0f172a;
  border-right: 1px solid #1e293b;
  padding: 16px 0;
  overflow-y: auto;
  flex-shrink: 0;
}

.menu-section {
  padding: 0 12px;
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8px;
  margin-bottom: 8px;
}

.section-title {
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

.menu-tree {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.menu-item {
  border-radius: 6px;
  overflow: hidden;
}

.menu-item-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
  user-select: none;
}

.menu-item-header:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.menu-item.active > .menu-item-header {
  background-color: #0284c7;
  color: white;
}

.menu-item.active > .menu-item-header .icon {
  color: white !important;
}

.chevron {
  flex-shrink: 0;
  color: #64748b;
}

.icon {
  flex-shrink: 0;
  color: #94a3b8;
}

.item-label {
  flex: 1;
  font-size: 13px;
  color: #e2e8f0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.menu-item.active .item-label {
  color: white;
}

.count-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  color: white;
}

.edit-btn {
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
  opacity: 0;
  transition: all 0.2s;
}

.menu-item-header:hover .edit-btn,
.sub-menu-item:hover .edit-btn {
  opacity: 1;
}

.edit-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.menu-item.active .edit-btn {
  color: rgba(255, 255, 255, 0.7);
}

.menu-item.active .edit-btn:hover {
  color: white;
  background-color: rgba(255, 255, 255, 0.2);
}

.sub-menu {
  padding-left: 20px;
  padding-top: 2px;
}

.sub-menu-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
}

.sub-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.sub-menu-item.active {
  background-color: rgba(2, 132, 199, 0.3);
}

.sub-menu-item .icon {
  color: #64748b;
}

.sub-menu-item .item-label {
  font-size: 12px;
  color: #94a3b8;
}

.sub-menu-item.active .item-label {
  color: #38bdf8;
}

.add-sub-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  margin-top: 4px;
  border-radius: 4px;
  background-color: transparent;
  border: 1px dashed #334155;
  color: #64748b;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-sub-btn:hover {
  border-color: #475569;
  color: #94a3b8;
}

.menu-divider {
  height: 1px;
  background-color: #1e293b;
  margin: 16px 12px;
}
</style>
