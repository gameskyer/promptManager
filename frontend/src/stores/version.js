import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import {
  CreateVersion,
  GetVersion,
  GetVersionHistory,
  GetLatestVersions,
  StarVersion,
  RollbackToVersion,
  CompareVersions,
  DeleteVersion,
  GetStarredVersions,
  GetVersionDiffStats,
} from '../lib/wailsjs/go/handlers/VersionHandler'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

export const useVersionStore = defineStore('version', () => {
  // State
  const versions = ref([])
  const currentVersion = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const sortedVersions = computed(() =>
    [...versions.value].sort((a, b) => b.version_num - a.version_num)
  )

  const starredVersions = computed(() =>
    versions.value.filter(v => v.is_starred)
  )

  const latestVersion = computed(() =>
    sortedVersions.value[0] || null
  )

  const getVersionByNum = computed(() => (versionNum) =>
    versions.value.find(v => v.version_num === versionNum)
  )

  const versionHistory = computed(() => {
    return sortedVersions.value.map(v => ({
      ...v,
      formattedTime: v.created_at ? dayjs(v.created_at).fromNow() : '',
      fullTime: v.created_at ? dayjs(v.created_at).format('YYYY-MM-DD HH:mm') : '',
    }))
  })

  // Actions
  async function fetchVersions(presetId, limit = 20) {
    loading.value = true
    error.value = null
    try {
      const response = await GetVersionHistory(presetId, limit)
      
      if (response.success) {
        versions.value = response.data || []
      } else {
        error.value = response.error || '获取版本历史失败'
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch versions:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchLatestVersions(presetId, count = 5) {
    loading.value = true
    error.value = null
    try {
      const response = await GetLatestVersions(presetId, count)
      
      if (response.success) {
        return response.data || []
      } else {
        error.value = response.error || '获取最新版本失败'
        return []
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch latest versions:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  async function fetchStarredVersions(presetId) {
    loading.value = true
    error.value = null
    try {
      const response = await GetStarredVersions(presetId)
      
      if (response.success) {
        return response.data || []
      } else {
        error.value = response.error || '获取星标版本失败'
        return []
      }
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch starred versions:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  async function createVersion(presetId, snapshot, thumbnailPath = '', previewPaths = []) {
    loading.value = true
    error.value = null
    try {
      const response = await CreateVersion({
        preset_id: presetId,
        pos_text: snapshot.pos_text,
        neg_text: snapshot.neg_text,
        atom_ids: snapshot.atom_ids || [],
        params: snapshot.params || {},
        preview_paths: previewPaths,
        thumbnail_path: thumbnailPath,
      })
      
      if (response.success) {
        versions.value.unshift(response.data)
        return response.data
      } else {
        throw new Error(response.error || '创建版本失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function starVersion(versionId, starred) {
    loading.value = true
    error.value = null
    try {
      const response = await StarVersion({
        version_id: versionId,
        starred: starred,
      })
      
      if (response.success) {
        const version = versions.value.find(v => v.id === versionId)
        if (version) {
          version.is_starred = starred
        }
      } else {
        throw new Error(response.error || '星标操作失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function rollbackToVersion(presetId, targetVersionNum) {
    loading.value = true
    error.value = null
    try {
      const response = await RollbackToVersion(presetId, targetVersionNum)
      
      if (response.success) {
        // 将新版本添加到列表
        versions.value.unshift(response.data)
        return response.data
      } else {
        throw new Error(response.error || '回滚失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function compareVersions(presetId, v1, v2) {
    loading.value = true
    error.value = null
    try {
      const response = await CompareVersions(presetId, v1, v2)
      
      if (response.success) {
        return {
          version1: response.data.version1,
          version2: response.data.version2,
          addedAtoms: response.data.added_atoms || [],
          removedAtoms: response.data.removed_atoms || [],
          paramChanges: response.data.param_changes || [],
          posTextDiff: response.data.pos_text_diff || {},
          negTextDiff: response.data.neg_text_diff || {},
        }
      } else {
        throw new Error(response.error || '版本对比失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteVersion(versionId) {
    loading.value = true
    error.value = null
    try {
      const response = await DeleteVersion(versionId)
      
      if (response.success) {
        versions.value = versions.value.filter(v => v.id !== versionId)
      } else {
        throw new Error(response.error || '删除版本失败')
      }
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function getVersionDiffStats(presetId, versionNum) {
    loading.value = true
    error.value = null
    try {
      const response = await GetVersionDiffStats(presetId, versionNum)
      
      if (response.success) {
        return response.data?.diff_stats || '+0/-0'
      } else {
        return '+0/-0'
      }
    } catch (e) {
      console.error('Failed to get diff stats:', e)
      return '+0/-0'
    } finally {
      loading.value = false
    }
  }

  function setCurrentVersion(version) {
    currentVersion.value = version
  }

  return {
    versions,
    currentVersion,
    loading,
    error,
    sortedVersions,
    starredVersions,
    latestVersion,
    getVersionByNum,
    versionHistory,
    fetchVersions,
    fetchLatestVersions,
    fetchStarredVersions,
    createVersion,
    starVersion,
    rollbackToVersion,
    compareVersions,
    deleteVersion,
    getVersionDiffStats,
    setCurrentVersion,
  }
})
