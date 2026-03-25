import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  CreatePreset,
  GetPresetByID,
  GetPresets,
  UpdatePreset,
  SoftDeletePreset,
  RestorePreset,
  BuildPrompt,
  GetCurrentWorkState,
  ForkPreset,
  CleanupOldVersions,
} from '../lib/wailsjs/go/handlers/PresetHandler'

export const usePresetStore = defineStore('preset', () => {
  // State
  const presets = ref([])
  const currentPreset = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const totalCount = ref(0)
  const currentPage = ref(1)

  // Getters
  const activePresets = computed(() =>
    presets.value.filter(p => !p.is_deleted)
  )

  const deletedPresets = computed(() =>
    presets.value.filter(p => p.is_deleted)
  )

  const getPresetById = computed(() => (id) =>
    presets.value.find(p => p.id === id)
  )

  // Actions
  async function fetchPresets(page = 1, pageSize = 20, categoryId = 0, includeDeleted = false) {
    loading.value = true
    error.value = null
    try {
      const response = await GetPresets({
        page: page,
        page_size: pageSize,
        category_id: categoryId,
        include_deleted: includeDeleted,
      })
      
      if (response.success) {
        presets.value = response.data.presets || []
        totalCount.value = response.data.total || 0
        currentPage.value = response.data.page || page
        
        // Debug: 打印第一个预设的封面信息
        if (presets.value.length > 0) {
          const first = presets.value[0]
          console.log('[presetStore] fetchPresets first preset:', {
            id: first.id,
            title: first.title,
            thumbnail: first.thumbnail,
            previews: first.previews,
          })
        }
      } else {
        error.value = response.error || '获取预设失败'
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch presets:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchPresetById(id) {
    loading.value = true
    error.value = null
    try {
      const response = await GetPresetByID(id)
      
      if (response.success) {
        return response.data
      } else {
        error.value = response.error || '获取预设详情失败'
        return null
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch preset:', e)
      return null
    } finally {
      loading.value = false
    }
  }

  async function createPreset(title, categoryId, posText, negText, atomIDs, params, loras = [], previews = []) {
    loading.value = true
    error.value = null
    try {
      // 将LoRAs合并到params中
      const paramsWithLoras = {
        ...params,
        loras: loras,
      }
      
      const response = await CreatePreset({
        title: title,
        category_id: categoryId,
        pos_text: posText,
        neg_text: negText,
        atom_ids: atomIDs,
        params: paramsWithLoras,
        previews: previews,
      })
      
      if (response.success) {
        console.log('[presetStore] createPreset response:', response.data)
        presets.value.unshift(response.data)
        return response.data
      } else {
        throw new Error(response.error || '创建预设失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updatePreset(id, title, categoryId = 0) {
    loading.value = true
    error.value = null
    try {
      const response = await UpdatePreset({
        id: id,
        title: title,
        category_id: categoryId,
      })
      
      if (response.success) {
        const preset = presets.value.find(p => p.id === id)
        if (preset) {
          preset.title = title
          preset.category_id = categoryId
          preset.updated_at = new Date().toISOString()
        }
        return response.data
      } else {
        throw new Error(response.error || '更新预设失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function softDeletePreset(id) {
    loading.value = true
    error.value = null
    try {
      const response = await SoftDeletePreset(id)
      
      if (response.success) {
        const preset = presets.value.find(p => p.id === id)
        if (preset) {
          preset.is_deleted = true
          preset.updated_at = new Date().toISOString()
        }
      } else {
        throw new Error(response.error || '删除预设失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function restorePreset(id) {
    loading.value = true
    error.value = null
    try {
      const response = await RestorePreset(id)
      
      if (response.success) {
        const preset = presets.value.find(p => p.id === id)
        if (preset) {
          preset.is_deleted = false
          preset.updated_at = new Date().toISOString()
        }
      } else {
        throw new Error(response.error || '恢复预设失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function setCurrentPreset(preset) {
    currentPreset.value = preset
  }

  async function buildPrompt(atomIDs) {
    loading.value = true
    error.value = null
    try {
      const response = await BuildPrompt({ atom_ids: atomIDs })
      
      if (response.success) {
        return {
          text: response.data.text,
          atomIds: response.data.atom_ids,
        }
      } else {
        throw new Error(response.error || '构建提示词失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function getCurrentWorkState(presetID) {
    loading.value = true
    error.value = null
    try {
      const response = await GetCurrentWorkState(presetID)
      
      if (response.success) {
        return {
          state: response.data.state,
          versionNum: response.data.version_num,
        }
      } else {
        throw new Error(response.error || '获取工作状态失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function forkPreset(presetId, versionNum, newTitle) {
    loading.value = true
    error.value = null
    try {
      const response = await ForkPreset({
        preset_id: presetId,
        version_num: versionNum,
        new_title: newTitle,
      })
      
      if (response.success) {
        presets.value.unshift(response.data)
        return response.data
      } else {
        throw new Error(response.error || 'Fork预设失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function cleanupOldVersions(presetId, keepCount) {
    loading.value = true
    error.value = null
    try {
      const response = await CleanupOldVersions({
        preset_id: presetId,
        keep_count: keepCount,
      })
      
      if (response.success) {
        return true
      } else {
        throw new Error(response.error || '清理旧版本失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 本地更新缩略图（不创建新版本）
  async function updateThumbnailOnly(id, thumbnail, previews = null) {
    try {
      const preset = presets.value.find(p => p.id === id)
      if (preset) {
        // 这里仅更新本地状态，实际缩略图应通过版本管理
        preset.thumbnail_path = thumbnail
        if (previews !== null) {
          preset.previews = previews
        }
        preset.updated_at = new Date().toISOString()
      }
      return true
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  return {
    presets,
    currentPreset,
    loading,
    error,
    totalCount,
    currentPage,
    activePresets,
    deletedPresets,
    getPresetById,
    fetchPresets,
    fetchPresetById,
    createPreset,
    updatePreset,
    softDeletePreset,
    restorePreset,
    setCurrentPreset,
    buildPrompt,
    getCurrentWorkState,
    forkPreset,
    cleanupOldVersions,
    updateThumbnailOnly,
  }
})
