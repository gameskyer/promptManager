<template>
  <main class="atom-management">
    <!-- Header -->
    <div class="management-header">
      <div class="header-left">
        <h2 class="page-title">原子词管理</h2>
        <span class="count-badge">共 {{ atomStore.totalCount }} 个</span>
      </div>
      <div class="header-actions">
        <button 
          class="btn-secondary" 
          :class="{ active: isBatchMode }"
          @click="toggleBatchMode"
        >
          <Squares2X2Icon class="w-4 h-4" />
          {{ isBatchMode ? '退出批量' : '批量操作' }}
        </button>
        <button class="btn-secondary" @click="exportData">
          <ArrowDownTrayIcon class="w-4 h-4" />
          导出
        </button>
        <button class="btn-primary" @click="createAtom">
          <PlusIcon class="w-4 h-4" />
          新建原子词
        </button>
      </div>
    </div>

    <!-- Batch Operations Bar -->
    <div v-if="isBatchMode" class="batch-bar">
      <div class="batch-info">
        <label class="checkbox-wrapper">
          <input 
            type="checkbox" 
            :checked="isAllSelected"
            @change="toggleSelectAll"
          />
          <span>全选本页</span>
        </label>
        <span class="batch-count">已选择 {{ selectedCount }} 项</span>
      </div>
      <div class="batch-actions">
        <button 
          class="btn-batch" 
          :disabled="selectedCount === 0"
          @click="showBatchMove = true"
        >
          <FolderArrowDownIcon class="w-4 h-4" />
          移动分类
        </button>
        <button 
          class="btn-batch" 
          :disabled="selectedCount === 0"
          @click="showBatchType = true"
        >
          <TagIcon class="w-4 h-4" />
          修改类型
        </button>
        <button 
          class="btn-batch" 
          :disabled="selectedCount === 0"
          @click="showBatchSynonyms = true"
        >
          <LinkIcon class="w-4 h-4" />
          添加近义词
        </button>
        <button 
          class="btn-batch danger" 
          :disabled="selectedCount === 0"
          @click="confirmBatchDelete"
        >
          <TrashIcon class="w-4 h-4" />
          批量删除
        </button>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-group category-filter">
        <label>一级分类</label>
        <select v-model="filterRootCategory" @change="handleRootCategoryChange">
          <option value="">全部分类</option>
          <option v-for="cat in rootAtomCategories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>
      <div class="filter-group category-filter">
        <label>二级分类</label>
        <select 
          v-model="filterCategory" 
          @change="handleFilterChange"
          :disabled="!filterRootCategory"
        >
          <option value="">{{ filterRootCategory ? '全部子分类' : '请先选择一级分类' }}</option>
          <option v-for="cat in childCategories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>
      <div class="filter-group">
        <label>类型</label>
        <select v-model="filterType" @change="handleFilterChange">
          <option value="">全部类型</option>
          <option value="Positive">正向</option>
          <option value="Negative">负向</option>
        </select>
      </div>
      <div class="filter-group search">
        <label>搜索</label>
        <div class="search-input">
          <MagnifyingGlassIcon class="w-4 h-4 search-icon" />
          <input 
            v-model="searchKeyword" 
            type="text" 
            placeholder="搜索原子词、标签..."
            @keyup.enter="handleFilterChange"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="clearSearch">
            <XMarkIcon class="w-3 h-3" />
          </button>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th v-if="isBatchMode" class="col-checkbox">
              <input 
                type="checkbox" 
                :checked="isAllSelected"
                @change="toggleSelectAll"
              />
            </th>
            <th class="col-id">ID</th>
            <th class="col-type">类型</th>
            <th class="col-value">英文原词</th>
            <th class="col-label">中文标签</th>
            <th class="col-category">所属分类</th>
            <th class="col-synonyms">近义词</th>
            <th class="col-usage">使用次数</th>
            <th v-if="!isBatchMode" class="col-actions">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="atom in filteredAtoms" 
            :key="atom.id" 
            :class="{ 
              'is-deleted': atom.is_deleted,
              'is-selected': isSelected(atom.id)
            }"
          >
            <td v-if="isBatchMode" class="col-checkbox">
              <input 
                type="checkbox" 
                :checked="isSelected(atom.id)"
                @change="toggleSelection(atom.id)"
              />
            </td>
            <td class="col-id">{{ atom.id }}</td>
            <td class="col-type">
              <span class="type-badge" :class="atom.type?.toLowerCase()">
                {{ atom.type === 'Positive' ? '正' : '负' }}
              </span>
            </td>
            <td class="col-value">
              <span class="atom-value" :title="atom.value">{{ atom.value }}</span>
            </td>
            <td class="col-label">{{ atom.label }}</td>
            <td class="col-category">{{ getCategoryName(atom.category_id) }}</td>
            <td class="col-synonyms">
              <div class="synonyms-tags">
                <span v-for="(syn, idx) in atom.synonyms?.slice(0, 2)" :key="idx" class="syn-tag">
                  {{ syn }}
                </span>
                <span v-if="atom.synonyms?.length > 2" class="syn-more">+{{ atom.synonyms.length - 2 }}</span>
              </div>
            </td>
            <td class="col-usage">
              <span class="usage-count">{{ atom.usage_count || 0 }}</span>
            </td>
            <td v-if="!isBatchMode" class="col-actions">
              <button class="action-btn" @click="editAtom(atom)" title="编辑">
                <PencilIcon class="w-4 h-4" />
              </button>
              <button class="action-btn" @click="addToWorkbench(atom)" title="添加到工作区">
                <PlusIcon class="w-4 h-4" />
              </button>
              <button class="action-btn danger" @click="deleteAtom(atom)" title="删除">
                <TrashIcon class="w-4 h-4" />
              </button>
            </td>
          </tr>
          <tr v-if="filteredAtoms.length === 0">
            <td :colspan="isBatchMode ? 9 : 8" class="empty-cell">
              <div class="empty-state">
                <DocumentMagnifyingGlassIcon class="w-12 h-12 text-slate-600" />
                <p>暂无数据</p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="pagination-bar">
      <div class="page-size">
        <label>每页</label>
        <select v-model="pageSize" @change="handlePageSizeChange">
          <option :value="20">20</option>
          <option :value="50">50</option>
          <option :value="100">100</option>
        </select>
        <span>条</span>
      </div>
      <div class="page-info">
        第 {{ currentPage }} / {{ totalPages }} 页
      </div>
      <div class="page-actions">
        <button :disabled="currentPage === 1" @click="currentPage--">
          <ChevronLeftIcon class="w-4 h-4" />
        </button>
        <button :disabled="currentPage >= totalPages" @click="currentPage++">
          <ChevronRightIcon class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Atom Dialog -->
    <AtomDialog
      v-if="showAtomDialog"
      :atom="editingAtom"
      :default-category-id="editingAtom?.category_id || 0"
      @close="closeAtomDialog"
      @save="saveAtom"
      @delete="confirmDeleteAtom"
    />

    <!-- Batch Move Dialog -->
    <div v-if="showBatchMove" class="modal-overlay" @click="showBatchMove = false">
      <div class="batch-dialog" @click.stop>
        <div class="batch-dialog-header">
          <FolderArrowDownIcon class="w-5 h-5 text-blue-500" />
          <h3>批量移动分类</h3>
        </div>
        <div class="batch-dialog-body">
          <p class="batch-desc">将选中的 {{ selectedCount }} 个原子词移动到：</p>
          <div class="batch-cascade-select">
            <div class="cascade-level">
              <label>一级分类</label>
              <select v-model="batchRootCategory" @change="batchTargetCategory = ''" class="batch-select">
                <option value="">请选择</option>
                <option v-for="cat in rootAtomCategories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </option>
              </select>
            </div>
            <div class="cascade-level">
              <label>二级分类</label>
              <select v-model="batchTargetCategory" class="batch-select" :disabled="!batchRootCategory">
                <option value="">{{ batchRootCategory ? '请选择（留空则移到一级分类）' : '请先选择一级分类' }}</option>
                <option v-for="cat in batchChildCategories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </option>
              </select>
            </div>
          </div>
        </div>
        <div class="batch-dialog-footer">
          <button class="btn-secondary" @click="showBatchMove = false">取消</button>
          <button class="btn-primary" :disabled="!batchRootCategory && !batchTargetCategory" @click="executeBatchMove">
            确认移动
          </button>
        </div>
      </div>
    </div>

    <!-- Batch Type Dialog -->
    <div v-if="showBatchType" class="modal-overlay" @click="showBatchType = false">
      <div class="batch-dialog" @click.stop>
        <div class="batch-dialog-header">
          <TagIcon class="w-5 h-5 text-green-500" />
          <h3>批量修改类型</h3>
        </div>
        <div class="batch-dialog-body">
          <p class="batch-desc">将选中的 {{ selectedCount }} 个原子词修改为：</p>
          <div class="batch-type-options">
            <label class="type-option">
              <input type="radio" v-model="batchTargetType" value="Positive" />
              <span class="type-badge positive">正向</span>
            </label>
            <label class="type-option">
              <input type="radio" v-model="batchTargetType" value="Negative" />
              <span class="type-badge negative">负向</span>
            </label>
          </div>
        </div>
        <div class="batch-dialog-footer">
          <button class="btn-secondary" @click="showBatchType = false">取消</button>
          <button class="btn-primary" :disabled="!batchTargetType" @click="executeBatchType">
            确认修改
          </button>
        </div>
      </div>
    </div>

    <!-- Batch Synonyms Dialog -->
    <div v-if="showBatchSynonyms" class="modal-overlay" @click="showBatchSynonyms = false">
      <div class="batch-dialog" @click.stop>
        <div class="batch-dialog-header">
          <LinkIcon class="w-5 h-5 text-purple-500" />
          <h3>批量添加近义词</h3>
        </div>
        <div class="batch-dialog-body">
          <p class="batch-desc">为选中的 {{ selectedCount }} 个原子词添加近义词：</p>
          <input 
            v-model="batchSynonymsInput" 
            type="text" 
            class="batch-input"
            placeholder="输入近义词，用逗号分隔"
          />
          <p class="batch-hint">多个近义词用逗号分隔，将合并到现有近义词中</p>
        </div>
        <div class="batch-dialog-footer">
          <button class="btn-secondary" @click="showBatchSynonyms = false">取消</button>
          <button class="btn-primary" :disabled="!batchSynonymsInput.trim()" @click="executeBatchSynonyms">
            确认添加
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Dialog -->
    <div v-if="showDeleteConfirm" class="modal-overlay" @click="showDeleteConfirm = false">
      <div class="confirm-dialog" @click.stop>
        <div class="confirm-header">
          <ExclamationTriangleIcon class="w-6 h-6 text-amber-500" />
          <h3>确认删除</h3>
        </div>
        <div class="confirm-body">
          <p>确定要删除原子词 "<strong>{{ deletingAtom?.value }}</strong>" 吗？</p>
          <p class="hint">此操作将软删除该原子词，可以在回收站中恢复。</p>
        </div>
        <div class="confirm-footer">
          <button class="btn-secondary" @click="showDeleteConfirm = false">取消</button>
          <button class="btn-danger" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>

    <!-- Batch Delete Confirm Dialog -->
    <div v-if="showBatchDeleteConfirm" class="modal-overlay" @click="showBatchDeleteConfirm = false">
      <div class="confirm-dialog" @click.stop>
        <div class="confirm-header">
          <ExclamationTriangleIcon class="w-6 h-6 text-red-500" />
          <h3>确认批量删除</h3>
        </div>
        <div class="confirm-body">
          <p>确定要删除选中的 <strong>{{ selectedCount }}</strong> 个原子词吗？</p>
          <p class="hint">此操作将软删除这些原子词，可以在回收站中恢复。</p>
        </div>
        <div class="confirm-footer">
          <button class="btn-secondary" @click="showBatchDeleteConfirm = false">取消</button>
          <button class="btn-danger" @click="executeBatchDelete">确认删除</button>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import {
  PlusIcon,
  PencilIcon,
  TrashIcon,
  MagnifyingGlassIcon,
  XMarkIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  DocumentMagnifyingGlassIcon,
  ArrowDownTrayIcon,
  ExclamationTriangleIcon,
  Squares2X2Icon,
  FolderArrowDownIcon,
  TagIcon,
  LinkIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, useAtomStore, useCategoryStore } from '../stores'
import AtomDialog from './AtomDialog.vue'

const appStore = useAppStore()
const atomStore = useAtomStore()
const categoryStore = useCategoryStore()

const { atoms } = storeToRefs(atomStore)
const { categories, categoryTree } = storeToRefs(categoryStore)
const { selectedAtoms, isBatchMode } = storeToRefs(atomStore)

// Filter state
const filterRootCategory = ref('')
const filterCategory = ref('')
const filterType = ref('')
const searchKeyword = ref('')

// Pagination
const currentPage = ref(1)
const pageSize = ref(50)

// Dialog state
const showAtomDialog = ref(false)
const showDeleteConfirm = ref(false)
const editingAtom = ref(null)
const deletingAtom = ref(null)

// Batch operation state
const showBatchMove = ref(false)
const showBatchType = ref(false)
const showBatchSynonyms = ref(false)
const showBatchDeleteConfirm = ref(false)
const batchRootCategory = ref('')
const batchTargetCategory = ref('')
const batchTargetType = ref('')
const batchSynonymsInput = ref('')

// Computed for batch move dialog
const batchChildCategories = computed(() => {
  if (!batchRootCategory.value) return []
  const rootCat = rootAtomCategories.value.find(c => c.id === batchRootCategory.value)
  return rootCat?.children || []
})

// Computed
const rootAtomCategories = computed(() => 
  categoryTree.value.filter(c => c.type === 'ATOM' && (c.parent_id === 0 || c.parent_id === null))
)

const childCategories = computed(() => {
  if (!filterRootCategory.value) return []
  const rootCat = rootAtomCategories.value.find(c => c.id === filterRootCategory.value)
  return rootCat?.children || []
})

const filteredAtoms = computed(() => {
  let result = [...atoms.value]

  // Category filter
  if (filterCategory.value) {
    // 二级分类筛选
    result = result.filter(a => a.category_id === filterCategory.value)
  } else if (filterRootCategory.value) {
    // 一级分类筛选 - 包含该分类下所有子分类的原子词
    const childIds = childCategories.value.map(c => c.id)
    childIds.push(filterRootCategory.value) // 也包含直接属于一级分类的原子词
    result = result.filter(a => childIds.includes(a.category_id))
  }

  // Type filter
  if (filterType.value) {
    result = result.filter(a => a.type === filterType.value)
  }

  // Search filter
  if (searchKeyword.value.trim()) {
    const term = searchKeyword.value.toLowerCase()
    result = result.filter(a =>
      a.value.toLowerCase().includes(term) ||
      a.label.includes(term) ||
      a.synonyms?.some(s => s.toLowerCase().includes(term))
    )
  }

  // Pagination
  const start = (currentPage.value - 1) * pageSize.value
  return result.slice(start, start + pageSize.value)
})

const totalPages = computed(() => {
  let count = atoms.value.length

  if (filterCategory.value) {
    count = atoms.value.filter(a => a.category_id === filterCategory.value).length
  } else if (filterRootCategory.value) {
    const childIds = childCategories.value.map(c => c.id)
    childIds.push(filterRootCategory.value)
    count = atoms.value.filter(a => childIds.includes(a.category_id)).length
  }
  if (filterType.value) {
    count = atoms.value.filter(a => a.type === filterType.value).length
  }
  if (searchKeyword.value.trim()) {
    const term = searchKeyword.value.toLowerCase()
    count = atoms.value.filter(a =>
      a.value.toLowerCase().includes(term) ||
      a.label.includes(term) ||
      a.synonyms?.some(s => s.toLowerCase().includes(term))
    ).length
  }

  return Math.ceil(count / pageSize.value) || 1
})

const selectedCount = computed(() => selectedAtoms.value.length)

const isAllSelected = computed(() => {
  if (filteredAtoms.value.length === 0) return false
  return filteredAtoms.value.every(atom => selectedAtoms.value.includes(atom.id))
})

const isSelected = (id) => selectedAtoms.value.includes(id)

// Methods
function handleRootCategoryChange() {
  // 当一级分类变化时，清除二级分类选择
  filterCategory.value = ''
  currentPage.value = 1
}

function toggleBatchMode() {
  atomStore.setBatchMode(!isBatchMode.value)
}

function toggleSelection(id) {
  atomStore.toggleAtomSelection(id)
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    // Deselect all on current page
    const currentIds = filteredAtoms.value.map(a => a.id)
    selectedAtoms.value = selectedAtoms.value.filter(id => !currentIds.includes(id))
  } else {
    // Select all on current page
    const currentIds = filteredAtoms.value.map(a => a.id)
    const newSelection = [...new Set([...selectedAtoms.value, ...currentIds])]
    atomStore.selectAll(newSelection)
  }
}

function getCategoryName(categoryId) {
  // 查找分类，包括子分类
  for (const root of categoryTree.value) {
    if (root.id === categoryId) {
      return root.name
    }
    if (root.children) {
      for (const child of root.children) {
        if (child.id === categoryId) {
          return `${root.name} > ${child.name}`
        }
      }
    }
  }
  return '-'
}

function handleFilterChange() {
  currentPage.value = 1
}

function handlePageSizeChange() {
  currentPage.value = 1
}

function clearSearch() {
  searchKeyword.value = ''
  handleFilterChange()
}

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
  try {
    const updateData = {
      value: data.value,
      label: data.label,
      type: data.type,
      category_id: data.category_id,
      synonyms: data.synonyms || [],
    }
    
    if (data.id) {
      await atomStore.updateAtom(data.id, updateData)
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

function deleteAtom(atom) {
  deletingAtom.value = atom
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  if (!deletingAtom.value) return
  try {
    await atomStore.deleteAtom(deletingAtom.value.id)
    showDeleteConfirm.value = false
    deletingAtom.value = null
  } catch (error) {
    console.error('Failed to delete atom:', error)
    alert('删除失败: ' + error.message)
  }
}

function confirmDeleteAtom(id) {
  const atom = atoms.value.find(a => a.id === id)
  if (atom) deleteAtom(atom)
}

function addToWorkbench(atom) {
  appStore.addAtom(atom)
  const toast = document.createElement('div')
  toast.className = 'toast-notification'
  toast.textContent = `已添加 "${atom.label || atom.value}" 到工作区`
  document.body.appendChild(toast)
  setTimeout(() => toast.remove(), 2000)
}

async function exportData() {
  try {
    const data = await atomStore.exportAtoms()
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `atoms-backup-${new Date().toISOString().slice(0, 10)}.json`
    a.click()
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Export failed:', error)
    alert('导出失败: ' + error.message)
  }
}

// Batch operations
function confirmBatchDelete() {
  showBatchDeleteConfirm.value = true
}

async function executeBatchMove() {
  // 如果没有选择二级分类，但有选择一级分类，则移动到一级分类
  const targetCategoryId = batchTargetCategory.value || batchRootCategory.value
  if (!targetCategoryId) return
  
  try {
    const count = await atomStore.batchMoveCategory(targetCategoryId)
    showBatchMove.value = false
    batchRootCategory.value = ''
    batchTargetCategory.value = ''
    alert(`成功移动 ${count} 个原子词`)
  } catch (error) {
    console.error('Batch move failed:', error)
    alert('批量移动失败: ' + error.message)
  }
}

async function executeBatchType() {
  if (!batchTargetType.value) return
  try {
    const count = await atomStore.batchUpdateType(batchTargetType.value)
    showBatchType.value = false
    batchTargetType.value = ''
    alert(`成功修改 ${count} 个原子词的类型`)
  } catch (error) {
    console.error('Batch type update failed:', error)
    alert('批量修改类型失败: ' + error.message)
  }
}

async function executeBatchSynonyms() {
  if (!batchSynonymsInput.value.trim()) return
  const synonyms = batchSynonymsInput.value.split(',').map(s => s.trim()).filter(Boolean)
  try {
    const count = await atomStore.batchAddSynonyms(synonyms)
    showBatchSynonyms.value = false
    batchSynonymsInput.value = ''
    alert(`成功为 ${count} 个原子词添加近义词`)
  } catch (error) {
    console.error('Batch add synonyms failed:', error)
    alert('批量添加近义词失败: ' + error.message)
  }
}

async function executeBatchDelete() {
  try {
    const count = await atomStore.batchDelete()
    showBatchDeleteConfirm.value = false
    alert(`成功删除 ${count} 个原子词`)
  } catch (error) {
    console.error('Batch delete failed:', error)
    alert('批量删除失败: ' + error.message)
  }
}

onMounted(async () => {
  await atomStore.fetchAtoms()
})

// Reset page when filters change
watch([filterCategory, filterType, searchKeyword], () => {
  currentPage.value = 1
})

// Watch for root category changes to reset filter
watch(filterRootCategory, () => {
  currentPage.value = 1
})
</script>

<style scoped>
.atom-management {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #020617;
  overflow: hidden;
}

.management-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #1e293b;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #e2e8f0;
}

.count-badge {
  font-size: 12px;
  color: #64748b;
  padding: 4px 10px;
  background-color: #1e293b;
  border-radius: 12px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn-primary,
.btn-secondary,
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
}

.btn-primary {
  background-color: #0284c7;
  color: white;
}

.btn-primary:hover {
  background-color: #0ea5e9;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #1e293b;
  color: #94a3b8;
  border: 1px solid #334155;
}

.btn-secondary:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.btn-secondary.active {
  background-color: rgba(14, 165, 233, 0.2);
  border-color: #0ea5e9;
  color: #0ea5e9;
}

.btn-danger {
  background-color: #dc2626;
  color: white;
}

.btn-danger:hover {
  background-color: #ef4444;
}

/* Batch Bar */
.batch-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background-color: rgba(14, 165, 233, 0.1);
  border-bottom: 1px solid #1e293b;
  gap: 16px;
}

.batch-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.checkbox-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 13px;
  color: #e2e8f0;
}

.checkbox-wrapper input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: #0ea5e9;
}

.batch-count {
  font-size: 13px;
  color: #0ea5e9;
  font-weight: 500;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.btn-batch {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-batch:hover:not(:disabled) {
  background-color: #334155;
  border-color: #475569;
}

.btn-batch:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-batch.danger {
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
}

.btn-batch.danger:hover:not(:disabled) {
  background-color: rgba(239, 68, 68, 0.2);
}

/* Filter Bar */
.filter-bar {
  display: flex;
  gap: 16px;
  padding: 12px 20px;
  border-bottom: 1px solid #1e293b;
  background-color: #0f172a;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-group label {
  font-size: 13px;
  color: #64748b;
}

.filter-group select {
  padding: 6px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 13px;
  outline: none;
}

.filter-group select:focus {
  border-color: #0ea5e9;
}

.filter-group.search {
  flex: 1;
  max-width: 300px;
}

.filter-group.category-filter select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: #0f172a;
}

.search-input {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
}

.search-input input {
  flex: 1;
  background: transparent;
  border: none;
  color: #e2e8f0;
  font-size: 13px;
  outline: none;
}

.search-icon {
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

/* Table */
.table-container {
  flex: 1;
  overflow: auto;
  padding: 0;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.data-table th {
  position: sticky;
  top: 0;
  padding: 12px 16px;
  background-color: #0f172a;
  border-bottom: 1px solid #334155;
  text-align: left;
  font-weight: 600;
  color: #94a3b8;
  white-space: nowrap;
}

.data-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
}

.data-table tbody tr:hover {
  background-color: rgba(255, 255, 255, 0.02);
}

.data-table tbody tr.is-deleted {
  opacity: 0.5;
}

.data-table tbody tr.is-selected {
  background-color: rgba(14, 165, 233, 0.1);
}

.col-checkbox {
  width: 40px;
  text-align: center;
}

.col-checkbox input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: #0ea5e9;
}

.col-id {
  width: 60px;
  color: #64748b;
}

.col-type {
  width: 60px;
}

.col-value {
  min-width: 150px;
  max-width: 200px;
}

.atom-value {
  font-family: monospace;
  color: #38bdf8;
}

.col-label {
  min-width: 120px;
}

.col-category {
  width: 120px;
  color: #94a3b8;
}

.col-synonyms {
  min-width: 150px;
}

.synonyms-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.syn-tag {
  font-size: 11px;
  padding: 2px 8px;
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  color: #64748b;
}

.syn-more {
  font-size: 11px;
  padding: 2px 8px;
  background-color: rgba(14, 165, 233, 0.15);
  border-radius: 4px;
  color: #0ea5e9;
}

.col-usage {
  width: 80px;
  text-align: center;
}

.usage-count {
  font-size: 12px;
  padding: 4px 10px;
  background-color: rgba(34, 197, 94, 0.15);
  border-radius: 12px;
  color: #22c55e;
}

.col-actions {
  width: 120px;
}

.col-actions .action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  margin-right: 4px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.col-actions .action-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.col-actions .action-btn.danger:hover {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.type-badge.positive {
  background-color: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.type-badge.negative {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.empty-cell {
  padding: 60px 20px;
  text-align: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #64748b;
}

/* Pagination */
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-top: 1px solid #1e293b;
  background-color: #0f172a;
}

.page-size {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #64748b;
}

.page-size select {
  padding: 4px 8px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 4px;
  color: #e2e8f0;
  font-size: 13px;
}

.page-info {
  font-size: 13px;
  color: #94a3b8;
}

.page-actions {
  display: flex;
  gap: 8px;
}

.page-actions button {
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

.page-actions button:hover:not(:disabled) {
  background-color: #334155;
}

.page-actions button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Modal Styles */
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

.confirm-dialog,
.batch-dialog {
  width: 100%;
  max-width: 400px;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  overflow: hidden;
}

.batch-dialog {
  max-width: 420px;
}

.confirm-header,
.batch-dialog-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-bottom: 1px solid #334155;
}

.confirm-header h3,
.batch-dialog-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
}

.confirm-body,
.batch-dialog-body {
  padding: 20px;
}

.confirm-body p {
  font-size: 14px;
  color: #e2e8f0;
  margin-bottom: 8px;
}

.confirm-body .hint {
  font-size: 12px;
  color: #64748b;
}

.batch-desc {
  font-size: 14px;
  color: #94a3b8;
  margin-bottom: 16px;
}

.batch-select,
.batch-input {
  width: 100%;
  padding: 10px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 14px;
  outline: none;
}

.batch-select:focus,
.batch-input:focus {
  border-color: #0ea5e9;
}

.batch-hint {
  font-size: 12px;
  color: #64748b;
  margin-top: 8px;
}

.batch-cascade-select {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cascade-level {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cascade-level label {
  font-size: 12px;
  color: #94a3b8;
}

.batch-type-options {
  display: flex;
  gap: 20px;
}

.type-option {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.type-option input[type="radio"] {
  accent-color: #0ea5e9;
}

.confirm-footer,
.batch-dialog-footer {
  display: flex;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid #334155;
}

.confirm-footer button,
.batch-dialog-footer button {
  flex: 1;
}

/* Toast Notification */
:global(.toast-notification) {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  background-color: #10b981;
  color: white;
  border-radius: 8px;
  font-size: 14px;
  z-index: 9999;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translate(-50%, 20px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style>
