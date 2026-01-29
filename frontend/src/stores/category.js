import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  CreateCategory,
  GetCategoryByID,
  GetCategoriesByParent,
  GetCategoryTree,
  UpdateCategory,
  DeleteCategory,
  MoveCategory,
  ReorderCategories,
} from '../lib/wailsjs/go/handlers/CategoryHandler'

export const useCategoryStore = defineStore('category', () => {
  // State
  const categories = ref([])
  const loading = ref(false)
  const error = ref(null)
  const categoryTree = ref([])

  // Getters
  const rootCategories = computed(() => 
    categories.value.filter(c => c.parent_id === 0 || c.parent_id === null)
  )

  const atomCategories = computed(() =>
    categories.value.filter(c => c.type === 'ATOM')
  )

  const presetCategories = computed(() =>
    categories.value.filter(c => c.type === 'PRESET')
  )

  const getChildren = computed(() => (parentId) =>
    categories.value.filter(c => c.parent_id === parentId)
  )

  const getCategoryById = computed(() => (id) =>
    categories.value.find(c => c.id === id)
  )

  // Actions
  async function fetchCategories() {
    loading.value = true
    error.value = null
    try {
      const response = await GetCategoryTree('')
      
      if (response.success) {
        categoryTree.value = response.data || []
        // 将树形结构扁平化为列表
        categories.value = flattenTree(categoryTree.value)
      } else {
        error.value = response.error || '获取分类失败'
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch categories:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchCategoriesByParent(parentId, type = '') {
    loading.value = true
    error.value = null
    try {
      const response = await GetCategoriesByParent(parentId, type)
      
      if (response.success) {
        return response.data || []
      } else {
        error.value = response.error || '获取分类失败'
        return []
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch categories by parent:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  async function createCategory(name, type, parentId = 0) {
    loading.value = true
    error.value = null
    try {
      const response = await CreateCategory({
        name: name,
        type: type,
        parent_id: parentId,
      })
      
      if (response.success) {
        categories.value.push(response.data)
        // 重新获取树形结构
        await fetchCategories()
        return response.data
      } else {
        throw new Error(response.error || '创建分类失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateCategory(id, updates) {
    loading.value = true
    error.value = null
    try {
      const response = await UpdateCategory({
        id: id,
        updates: updates,
      })
      
      if (response.success) {
        const index = categories.value.findIndex(c => c.id === id)
        if (index !== -1) {
          categories.value[index] = { ...categories.value[index], ...response.data }
        }
        // 重新获取树形结构
        await fetchCategories()
        return response.data
      } else {
        throw new Error(response.error || '更新分类失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteCategory(id) {
    loading.value = true
    error.value = null
    try {
      const response = await DeleteCategory(id)
      
      if (response.success) {
        // 递归删除子分类
        const deleteChildren = (parentId) => {
          const children = categories.value.filter(c => c.parent_id === parentId)
          children.forEach(child => {
            deleteChildren(child.id)
          })
        }
        deleteChildren(id)
        
        categories.value = categories.value.filter(c => c.id !== id)
        // 重新获取树形结构
        await fetchCategories()
      } else {
        throw new Error(response.error || '删除分类失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function moveCategory(id, newParentId) {
    loading.value = true
    error.value = null
    try {
      const response = await MoveCategory({
        id: id,
        new_parent_id: newParentId,
      })
      
      if (response.success) {
        const category = categories.value.find(c => c.id === id)
        if (category) {
          category.parent_id = newParentId
        }
        // 重新获取树形结构
        await fetchCategories()
      } else {
        throw new Error(response.error || '移动分类失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function reorderCategories(ids) {
    loading.value = true
    error.value = null
    try {
      const response = await ReorderCategories({ ids: ids })
      
      if (response.success) {
        // 重新获取分类
        await fetchCategories()
      } else {
        throw new Error(response.error || '排序失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  // Helper function to flatten tree structure
  function flattenTree(tree, result = []) {
    for (const node of tree) {
      const { children, ...category } = node
      result.push(category)
      if (children && children.length > 0) {
        flattenTree(children, result)
      }
    }
    return result
  }

  return {
    categories,
    categoryTree,
    loading,
    error,
    rootCategories,
    atomCategories,
    presetCategories,
    getChildren,
    getCategoryById,
    fetchCategories,
    fetchCategoriesByParent,
    createCategory,
    updateCategory,
    deleteCategory,
    moveCategory,
    reorderCategories,
  }
})
