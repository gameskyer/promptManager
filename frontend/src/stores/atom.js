import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  CreateAtom,
  GetAtomByID,
  GetAtomsByCategory,
  UpdateAtom,
  DeleteAtom,
  RecordUsage,
  FindAtomsBySynonym,
  GetPopularAtoms,
  BatchImportAtoms,
  GetAllAtomsPaginated,
  ExportAtoms,
} from '../lib/wailsjs/go/handlers/AtomHandler'
import {
  BatchMoveCategory,
  BatchUpdateType,
  BatchDelete,
  BatchAddSynonyms,
  BatchClearCategory,
} from '../lib/wailsjs/go/handlers/BatchHandler'

export const useAtomStore = defineStore('atom', () => {
  // State
  const atoms = ref([])
  const loading = ref(false)
  const error = ref(null)
  const currentPage = ref(1)
  const totalCount = ref(0)
  const pageSize = ref(50)
  
  // Batch operation state
  const selectedAtoms = ref([])
  const isBatchMode = ref(false)

  // Getters
  const getAtomsByCategory = computed(() => (categoryId) =>
    atoms.value.filter(a => a.category_id === categoryId)
  )

  const getAtomById = computed(() => (id) =>
    atoms.value.find(a => a.id === id)
  )

  const popularAtoms = computed(() =>
    [...atoms.value].sort((a, b) => (b.usage_count || 0) - (a.usage_count || 0)).slice(0, 20)
  )

  // Actions
  async function fetchAtoms(categoryId = 0, page = 1, size = 50) {
    loading.value = true
    error.value = null
    try {
      let response
      if (categoryId > 0) {
        response = await GetAtomsByCategory({
          category_id: categoryId,
          page: page,
          page_size: size,
        })
      } else {
        response = await GetAllAtomsPaginated(page, size)
      }
      
      if (response.success) {
        if (categoryId > 0) {
          atoms.value = response.data.atoms || []
          totalCount.value = response.data.total || 0
        } else {
          atoms.value = response.data || []
          totalCount.value = response.data?.length || 0
        }
        currentPage.value = page
        pageSize.value = size
      } else {
        error.value = response.error || '获取原子词失败'
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch atoms:', e)
    } finally {
      loading.value = false
    }
  }

  async function searchAtoms(searchTerm, type = '', categoryId = 0, limit = 50) {
    loading.value = true
    error.value = null
    try {
      const response = await FindAtomsBySynonym(searchTerm)
      
      if (response.success) {
        let results = response.data || []
        
        // 客户端过滤
        if (type) {
          results = results.filter(a => a.type === type)
        }
        if (categoryId > 0) {
          results = results.filter(a => a.category_id === categoryId)
        }
        
        atoms.value = results.slice(0, limit)
        totalCount.value = results.length
      } else {
        error.value = response.error || '搜索失败'
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to search atoms:', e)
    } finally {
      loading.value = false
    }
  }

  async function createAtom(atomData) {
    loading.value = true
    error.value = null
    try {
      const response = await CreateAtom({
        value: atomData.value,
        label: atomData.label,
        type: atomData.type || 'Positive',
        category_id: atomData.category_id,
        synonyms: atomData.synonyms || [],
      })
      
      if (response.success) {
        atoms.value.push(response.data)
        return response.data
      } else {
        throw new Error(response.error || '创建失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateAtom(id, updates) {
    loading.value = true
    error.value = null
    try {
      const response = await UpdateAtom({
        id: id,
        updates: updates,
      })
      
      if (response.success) {
        const index = atoms.value.findIndex(a => a.id === id)
        if (index !== -1) {
          atoms.value[index] = { ...atoms.value[index], ...response.data }
        }
        return response.data
      } else {
        throw new Error(response.error || '更新失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteAtom(id) {
    loading.value = true
    error.value = null
    try {
      const response = await DeleteAtom(id)
      
      if (response.success) {
        atoms.value = atoms.value.filter(a => a.id !== id)
      } else {
        throw new Error(response.error || '删除失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function recordUsage(atomId) {
    try {
      const response = await RecordUsage(atomId)
      
      if (response.success) {
        const atom = atoms.value.find(a => a.id === atomId)
        if (atom) {
          atom.usage_count = (atom.usage_count || 0) + 1
          atom.last_used_at = new Date().toISOString()
        }
      }
    } catch (e) {
      console.error('Failed to record usage:', e)
    }
  }

  async function fetchPopularAtoms(limit = 20) {
    loading.value = true
    error.value = null
    try {
      const response = await GetPopularAtoms(limit)
      
      if (response.success) {
        return response.data || []
      } else {
        error.value = response.error || '获取热门原子词失败'
        return []
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch popular atoms:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  async function batchImport(jsonData) {
    loading.value = true
    error.value = null
    try {
      const response = await BatchImportAtoms({ json_data: jsonData })
      
      if (response.success) {
        // 重新加载数据
        await fetchAtoms()
        return response.data?.imported || 0
      } else {
        throw new Error(response.error || '导入失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function exportAtoms() {
    try {
      const response = await ExportAtoms()
      
      if (response.success) {
        return response.data
      } else {
        throw new Error(response.error || '导出失败')
      }
    } catch (e) {
      console.error('Failed to export atoms:', e)
      throw e
    }
  }

  // Batch operations
  function toggleAtomSelection(atomId) {
    const index = selectedAtoms.value.indexOf(atomId)
    if (index === -1) {
      selectedAtoms.value.push(atomId)
    } else {
      selectedAtoms.value.splice(index, 1)
    }
  }

  function selectAll(atomIds) {
    selectedAtoms.value = [...atomIds]
  }

  function clearSelection() {
    selectedAtoms.value = []
  }

  function setBatchMode(mode) {
    isBatchMode.value = mode
    if (!mode) {
      clearSelection()
    }
  }

  async function batchMoveCategory(categoryId) {
    loading.value = true
    error.value = null
    try {
      const response = await BatchMoveCategory({
        atom_ids: selectedAtoms.value,
        category_id: categoryId,
      })
      
      if (response.success) {
        clearSelection()
        await fetchAtoms()
        return response.data?.moved_count || 0
      } else {
        throw new Error(response.error || '批量移动失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function batchUpdateType(type) {
    loading.value = true
    error.value = null
    try {
      const response = await BatchUpdateType({
        atom_ids: selectedAtoms.value,
        type: type,
      })
      
      if (response.success) {
        clearSelection()
        await fetchAtoms()
        return response.data?.updated_count || 0
      } else {
        throw new Error(response.error || '批量更新类型失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function batchDelete() {
    loading.value = true
    error.value = null
    try {
      const response = await BatchDelete({
        atom_ids: selectedAtoms.value,
      })
      
      if (response.success) {
        const deletedCount = response.data?.deleted_count || 0
        atoms.value = atoms.value.filter(a => !selectedAtoms.value.includes(a.id))
        clearSelection()
        return deletedCount
      } else {
        throw new Error(response.error || '批量删除失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function batchAddSynonyms(synonyms) {
    loading.value = true
    error.value = null
    try {
      const response = await BatchAddSynonyms({
        atom_ids: selectedAtoms.value,
        synonyms: synonyms,
      })
      
      if (response.success) {
        clearSelection()
        await fetchAtoms()
        return response.data?.updated_count || 0
      } else {
        throw new Error(response.error || '批量添加近义词失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function batchClearCategory() {
    loading.value = true
    error.value = null
    try {
      const response = await BatchClearCategory({
        atom_ids: selectedAtoms.value,
      })
      
      if (response.success) {
        clearSelection()
        await fetchAtoms()
        return response.data?.cleared_count || 0
      } else {
        throw new Error(response.error || '批量清除分类失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    atoms,
    loading,
    error,
    currentPage,
    totalCount,
    pageSize,
    getAtomsByCategory,
    getAtomById,
    popularAtoms,
    fetchAtoms,
    searchAtoms,
    createAtom,
    updateAtom,
    deleteAtom,
    recordUsage,
    fetchPopularAtoms,
    batchImport,
    exportAtoms,
    // Batch operations
    selectedAtoms,
    isBatchMode,
    toggleAtomSelection,
    selectAll,
    clearSelection,
    setBatchMode,
    batchMoveCategory,
    batchUpdateType,
    batchDelete,
    batchAddSynonyms,
    batchClearCategory,
  }
})
