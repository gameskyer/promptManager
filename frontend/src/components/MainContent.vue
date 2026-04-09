<template>
  <main class="main-content">
    <!-- Sub-category tabs -->
    <div v-if="currentCategory && subCategories.length > 0" class="sub-category-tabs">
      <button
        v-for="sub in subCategories"
        :key="sub.id"
        class="tab-btn"
        :class="{ active: currentSubCategory?.id === sub.id }"
        @click="selectSubCategory(sub)"
      >
        {{ sub.name }}
        <button 
          class="edit-sub-btn"
          @click.stop="editCategory(sub)"
        >
          <PencilIcon class="w-3 h-3" />
        </button>
      </button>
      <button class="tab-btn add" @click="createSubCategory">
        <PlusIcon class="w-4 h-4" />
      </button>
    </div>
    
    <!-- Content Header -->
    <div class="content-header">
      <div class="header-left">
        <h2 class="content-title">
          {{ currentSubCategory?.name || currentCategory?.name || '全部原子词' }}
        </h2>
        <button 
          v-if="currentCategory"
          class="edit-category-btn"
          @click="editCategory(currentCategory)"
        >
          <PencilIcon class="w-4 h-4" />
        </button>
      </div>
      <div class="content-actions">
        <button class="action-btn" @click="createAtom">
          <PlusIcon class="w-4 h-4" />
          <span>新建原子词</span>
        </button>

      </div>
    </div>
    
    <!-- Atom Grid -->
    <div class="atom-grid">
      <AtomCard
        v-for="atom in displayedAtoms"
        :key="atom.id"
        :atom="atom"
        :is-selected="isSelected(atom.id)"
        @toggle="toggleAtom(atom)"
        @add="addAtom(atom)"
        @add-synonym="addSynonym"
        @edit="editAtom(atom)"
      />
      
      <div v-if="displayedAtoms.length === 0" class="empty-state">
        <DocumentMagnifyingGlassIcon class="w-12 h-12 text-slate-600" />
        <p>暂无原子词</p>
        <button class="btn-primary" @click="createAtom">
          <PlusIcon class="w-4 h-4" />
          添加原子词
        </button>
      </div>
    </div>
    
    <!-- Pagination -->
    <div v-if="atomStore.totalCount > pageSize" class="pagination">
      <button 
        :disabled="currentPage === 1" 
        @click="currentPage--"
      >
        <ChevronLeftIcon class="w-4 h-4" />
      </button>
      <span class="page-info">{{ currentPage }} / {{ Math.ceil(atomStore.totalCount / pageSize) }}</span>
      <button 
        :disabled="currentPage * pageSize >= atomStore.totalCount" 
        @click="currentPage++"
      >
        <ChevronRightIcon class="w-4 h-4" />
      </button>
    </div>
  </main>
  
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
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import {
  PlusIcon,
  FunnelIcon,
  PencilIcon,
  DocumentMagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, useCategoryStore, useAtomStore } from '../stores'
import AtomCard from './AtomCard.vue'
import AtomDialog from './AtomDialog.vue'
import CategoryDialog from './CategoryDialog.vue'

const appStore = useAppStore()
const categoryStore = useCategoryStore()
const atomStore = useAtomStore()

const { currentCategory, currentSubCategory, selectedAtomIDs } = storeToRefs(appStore)
const { getChildren } = storeToRefs(categoryStore)

// 获取当前选中的分类ID（优先使用子分类）
const currentSelectedCategoryId = computed(() => {
  return currentSubCategory.value?.id || currentCategory.value?.id || 0
})

const currentPage = ref(1)
const pageSize = 50

// Dialog states
const showAtomDialog = ref(false)
const showCategoryDialog = ref(false)
const editingAtom = ref(null)
const editingCategory = ref(null)

const subCategories = computed(() => {
  if (!currentCategory.value) return []
  return getChildren.value(currentCategory.value.id)
})

const atoms = computed(() => atomStore.atoms)

const displayedAtoms = computed(() => {
  let result = atoms.value
  
  // 搜索过滤
  if (appStore.searchQuery) {
    const term = appStore.searchQuery.toLowerCase()
    result = result.filter(a => 
      a.value.toLowerCase().includes(term) ||
      a.label.includes(term) ||
      a.synonyms?.some(s => s.toLowerCase().includes(term))
    )
  }
  
  // 分页
  const start = (currentPage.value - 1) * pageSize
  return result.slice(start, start + pageSize)
})

function isSelected(atomId) {
  return selectedAtomIDs.value.includes(atomId)
}

function toggleAtom(atom) {
  appStore.toggleAtom(atom)
}

function addAtom(atom) {
  appStore.addAtom(atom)
  atomStore.recordUsage(atom.id)
}

function addSynonym(synonymAtom) {
  appStore.addAtom(synonymAtom)
  // 如果是已有原子词，记录使用
  if (synonymAtom.id) {
    atomStore.recordUsage(synonymAtom.id)
  }
}

function selectSubCategory(sub) {
  appStore.setSubCategory(sub)
  currentPage.value = 1
  atomStore.fetchAtoms(sub.id)
}

// Atom CURD
function createAtom() {
  editingAtom.value = null
  showAtomDialog.value = true
}

function editAtom(atom) {
  editingAtom.value = atom
  showAtomDialog.value = true
}

function closeAtomDialog() {
  showAtomDialog.value = false
  editingAtom.value = null
}

async function saveAtom(data) {
  if (data.id) {
    await atomStore.updateAtom(data.id, data)
  } else {
    await atomStore.createAtom(data)
  }
  closeAtomDialog()
  // 模拟数据模式下，不调用 fetchAtoms 以免重置数据
  // 新添加的原子词已经通过响应式更新到列表中
}

async function deleteAtom(id) {
  await atomStore.deleteAtom(id)
  closeAtomDialog()
}

// Category CURD
function createSubCategory() {
  editingCategory.value = {
    type: 'ATOM',
    parent_id: currentCategory.value?.id || 0,
  }
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

watch(() => currentSubCategory.value, (newVal) => {
  if (newVal) {
    atomStore.fetchAtoms(newVal.id)
  }
})
</script>

<style scoped>
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #020617;
  overflow: hidden;
}

.sub-category-tabs {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  border-bottom: 1px solid #1e293b;
  overflow-x: auto;
  scrollbar-width: none;
}

.sub-category-tabs::-webkit-scrollbar {
  display: none;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 6px;
  background-color: transparent;
  border: 1px solid transparent;
  color: #94a3b8;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tab-btn:hover {
  background-color: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

.tab-btn.active {
  background-color: rgba(2, 132, 199, 0.2);
  border-color: #0284c7;
  color: #38bdf8;
}

.tab-btn.add {
  padding: 6px 10px;
  border-color: #334155;
}

.edit-sub-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  opacity: 0;
  transition: all 0.2s;
}

.tab-btn:hover .edit-sub-btn {
  opacity: 1;
}

.edit-sub-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #1e293b;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.content-title {
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
}

.edit-category-btn {
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

.edit-category-btn:hover {
  background-color: #1e293b;
  color: #e2e8f0;
}

.content-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 6px;
  background-color: #1e293b;
  border: 1px solid #334155;
  color: #94a3b8;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.action-btn {
  background-color: #0284c7;
  border-color: #0284c7;
  color: white;
}

.action-btn:hover {
  background-color: #0ea5e9;
  border-color: #0ea5e9;
}

.atom-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  grid-auto-rows: min-content;
  gap: 8px;
  padding: 12px;
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
  padding: 60px 20px;
  color: #64748b;
}

.empty-state p {
  font-size: 14px;
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
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  background-color: #0ea5e9;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 12px 16px;
  border-top: 1px solid #1e293b;
}

.pagination button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background-color: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  background-color: #334155;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 13px;
  color: #94a3b8;
}
</style>
