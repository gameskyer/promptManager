import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  UploadImage,
  GetImageByID,
  GetImagesByPreset,
  DeleteImage,
  GetImageData,
} from '../lib/wailsjs/go/handlers/ImageHandler'

export const useImageStore = defineStore('image', () => {
  // State
  const loading = ref(false)
  const error = ref(null)

  // Actions
  async function uploadImage(file, presetID = 0, versionID = 0) {
    loading.value = true
    error.value = null

    try {
      // 读取文件为 base64
      const base64 = await fileToBase64(file)
      
      const req = {
        preset_id: presetID,
        version_id: versionID,
        data: base64,
      }

      const resp = await UploadImage(req)
      return resp
    } catch (err) {
      error.value = err.message || '上传失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function getImageByID(id) {
    loading.value = true
    error.value = null

    try {
      const resp = await GetImageByID(id)
      return resp
    } catch (err) {
      error.value = err.message || '获取图片失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function getImagesByPreset(presetID) {
    loading.value = true
    error.value = null

    try {
      const resp = await GetImagesByPreset(presetID)
      return resp
    } catch (err) {
      error.value = err.message || '获取预设图片列表失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function deleteImage(id) {
    loading.value = true
    error.value = null

    try {
      const resp = await DeleteImage(id)
      return resp
    } catch (err) {
      error.value = err.message || '删除失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function getImageData(id) {
    loading.value = true
    error.value = null

    try {
      const resp = await GetImageData(id)
      return resp
    } catch (err) {
      error.value = err.message || '获取图片数据失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  // Helper function
  function fileToBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => {
        // 移除 data:image/xxx;base64, 前缀
        const base64 = reader.result.split(',')[1]
        resolve(base64)
      }
      reader.onerror = reject
      reader.readAsDataURL(file)
    })
  }

  return {
    loading,
    error,
    uploadImage,
    getImageByID,
    getImagesByPreset,
    deleteImage,
    getImageData,
  }
})
