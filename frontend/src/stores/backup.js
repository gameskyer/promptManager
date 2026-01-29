import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  ExportData,
  ImportData,
  ExportToFile,
  GetBackupList,
  ExportToZip,
  ReadBackupFile,
} from '../lib/wailsjs/go/handlers/BackupHandler'

export const useBackupStore = defineStore('backup', () => {
  // State
  const loading = ref(false)
  const error = ref(null)
  const backupList = ref([])

  // Actions
  async function exportData() {
    loading.value = true
    error.value = null

    try {
      const resp = await ExportData()
      if (resp.success) {
        return { success: true, data: resp.data }
      } else {
        throw new Error(resp.error || '导出失败')
      }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function importData(jsonData, merge = true) {
    loading.value = true
    error.value = null

    try {
      const resp = await ImportData({
        data: jsonData,
        merge: merge,
      })
      if (resp.success) {
        return { success: true }
      } else {
        throw new Error(resp.error || '导入失败')
      }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function exportToFile() {
    loading.value = true
    error.value = null

    try {
      const resp = await ExportToFile()
      if (resp.success) {
        return { success: true, path: resp.data }
      } else {
        throw new Error(resp.error || '导出失败')
      }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function exportToZip() {
    loading.value = true
    error.value = null

    try {
      const resp = await ExportToZip()
      if (resp.success) {
        return { success: true, path: resp.data }
      } else {
        throw new Error(resp.error || '导出失败')
      }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function fetchBackupList() {
    try {
      const resp = await GetBackupList()
      if (resp.success) {
        backupList.value = resp.data || []
        return backupList.value
      }
    } catch (err) {
      console.error('Failed to fetch backup list:', err)
    }
    return []
  }

  async function readBackupFile(filename) {
    try {
      const resp = await ReadBackupFile(filename)
      if (resp.success) {
        return resp.data
      }
    } catch (err) {
      console.error('Failed to read backup file:', err)
    }
    return null
  }

  // Download data as file
  function downloadAsFile(content, filename, type = 'application/json') {
    const blob = new Blob([content], { type })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  return {
    loading,
    error,
    backupList,
    exportData,
    importData,
    exportToFile,
    exportToZip,
    fetchBackupList,
    readBackupFile,
    downloadAsFile,
  }
})
