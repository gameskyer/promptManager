<template>
  <main class="category-management">
    <!-- Header -->
    <div class="management-header">
      <div class="header-left">
        <h2 class="page-title">分类管理</h2>
        <span class="count-badge">共 {{ categoryStore.categories.length }} 个</span>
      </div>
      <div class="header-actions">
        <button class="btn-secondary" @click="expandAll">
          <ChevronDownIcon class="w-4 h-4" />
          展开全部
        </button>
        <button class="btn-secondary" @click="collapseAll">
          <ChevronRightIcon class="w-4 h-4" />
          折叠全部
        </button>
        <button class="btn-primary" @click="createCategory">
          <PlusIcon class="w-4 h-4" />
          新建分类
        </button>
      </div>
    </div>

    <!-- Category Tree -->
    <div class="tree-container">
      <div v-if="categoryTree.length === 0" class="empty-state">
        <FolderIcon class="w-16 h-16 text-slate-600" />
        <p>暂无分类</p>
        <button class="btn-primary" @click="createCategory">
          <PlusIcon class="w-4 h-4" />
          创建第一个分类
        </button>
      </div>

      <div v-else class="category-tree">
        <CategoryTreeNode
          v-for="node in categoryTree"
          :key="node.id"
          :node="node"
          :expanded="expandedIds"
          :level="0"
          @toggle="toggleExpand"
          @edit="editCategory"
          @delete="deleteCategory"
          @add-child="addChildCategory"
          @move="moveCategory"
        />
      </div>
    </div>

    <!-- Category Dialog -->
    <CategoryDialog
      v-if="showCategoryDialog"
      :category="editingCategory"
      @close="closeCategoryDialog"
      @save="saveCategory"
      @delete="confirmDeleteCategory"
    />

    <!-- Delete Confirm Dialog -->
    <div v-if="showDeleteConfirm" class="modal-overlay" @click="showDeleteConfirm = false">
      <div class="confirm-dialog" @click.stop>
        <div class="confirm-header">
          <ExclamationTriangleIcon class="w-6 h-6 text-amber-500" />
          <h3>确认删除</h3>
        </div>
        <div class="confirm-body">
          <p>确定要删除分类 "<strong>{{ deletingCategory?.name }}</strong>" 吗？</p>
          <p v-if="deletingCategory?.children?.length > 0" class="warning">
            <strong>警告：</strong>此分类下有 {{ deletingCategory.children.length }} 个子分类，删除后将一并删除！
          </p>
          <p class="hint">此操作不可撤销。</p>
        </div>
        <div class="confirm-footer">
          <button class="btn-secondary" @click="showDeleteConfirm = false">取消</button>
          <button class="btn-danger" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import {
  PlusIcon,
  FolderIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  ExclamationTriangleIcon,
} from '@heroicons/vue/24/outline'
import { useCategoryStore } from '../stores'
import CategoryDialog from './CategoryDialog.vue'
import CategoryTreeNode from './CategoryTreeNode.vue'

const categoryStore = useCategoryStore()
const { categoryTree, categories } = storeToRefs(categoryStore)

// State
const expandedIds = ref([])
const showCategoryDialog = ref(false)
const showDeleteConfirm = ref(false)
const editingCategory = ref(null)
const deletingCategory = ref(null)

// Initialize expanded state with root categories
onMounted(() => {
  // Expand all root categories by default
  const rootIds = categoryTree.value
    .filter(node => node.parent_id === 0 || node.parent_id === null)
    .map(node => node.id)
  expandedIds.value = rootIds
})

function expandAll() {
  expandedIds.value = categories.value.map(c => c.id)
}

function collapseAll() {
  expandedIds.value = []
}

function toggleExpand(nodeId) {
  const index = expandedIds.value.indexOf(nodeId)
  if (index === -1) {
    expandedIds.value.push(nodeId)
  } else {
    expandedIds.value.splice(index, 1)
  }
}

function createCategory() {
  editingCategory.value = {
    type: 'ATOM',
    parent_id: 0,
  }
  showCategoryDialog.value = true
}

function addChildCategory(parentId) {
  editingCategory.value = {
    type: 'ATOM',
    parent_id: parentId,
  }
  showCategoryDialog.value = true
  // Auto expand parent
  if (!expandedIds.value.includes(parentId)) {
    expandedIds.value.push(parentId)
  }
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
  try {
    if (data.id) {
      await categoryStore.updateCategory(data.id, data)
    } else {
      const newCategory = await categoryStore.createCategory(data.name, data.type, data.parent_id)
      // Auto expand parent if creating child
      if (newCategory.parent_id > 0 && !expandedIds.value.includes(newCategory.parent_id)) {
        expandedIds.value.push(newCategory.parent_id)
      }
    }
    closeCategoryDialog()
  } catch (error) {
    console.error('Failed to save category:', error)
    alert('保存失败: ' + error.message)
  }
}

function deleteCategory(category) {
  deletingCategory.value = category
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  if (!deletingCategory.value) return
  try {
    await categoryStore.deleteCategory(deletingCategory.value.id)
    showDeleteConfirm.value = false
    deletingCategory.value = null
  } catch (error) {
    console.error('Failed to delete category:', error)
    alert('删除失败: ' + error.message)
  }
}

function confirmDeleteCategory(id) {
  const category = categories.value.find(c => c.id === id)
  if (category) deleteCategory(category)
}

async function moveCategory(data) {
  try {
    const { categoryId, newParentId } = data
    await categoryStore.moveCategory(categoryId, newParentId)
  } catch (error) {
    console.error('Failed to move category:', error)
    alert('移动失败: ' + error.message)
  }
}
</script>

<style scoped>
.category-management {
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

.btn-secondary {
  background-color: #1e293b;
  color: #94a3b8;
  border: 1px solid #334155;
}

.btn-secondary:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.btn-danger {
  background-color: #dc2626;
  color: white;
}

.btn-danger:hover {
  background-color: #ef4444;
}

.tree-container {
  flex: 1;
  overflow: auto;
  padding: 20px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 80px 20px;
  color: #64748b;
}

.empty-state p {
  font-size: 14px;
}

.category-tree {
  max-width: 800px;
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

.confirm-dialog {
  width: 100%;
  max-width: 400px;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  overflow: hidden;
}

.confirm-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-bottom: 1px solid #334155;
}

.confirm-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
}

.confirm-body {
  padding: 20px;
}

.confirm-body p {
  font-size: 14px;
  color: #e2e8f0;
  margin-bottom: 8px;
}

.confirm-body .warning {
  padding: 10px 12px;
  background-color: rgba(245, 158, 11, 0.15);
  border-radius: 6px;
  color: #f59e0b;
  font-size: 13px;
}

.confirm-body .hint {
  font-size: 12px;
  color: #64748b;
  margin-top: 12px;
}

.confirm-footer {
  display: flex;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid #334155;
}

.confirm-footer button {
  flex: 1;
}
</style>
