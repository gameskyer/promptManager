// Wails runtime binding helpers
// This file provides a clean interface for calling Go backend methods

// Check if Wails runtime is available
const isWails = () => typeof window !== 'undefined' && window.go

// Safe call helper that works with or without Wails
async function safeCall(service, method, ...args) {
  if (!isWails()) {
    console.warn(`Wails not available: ${service}.${method}`)
    return null
  }
  
  try {
    const handler = window.go?.backend?.handlers?.[service]
    if (!handler || typeof handler[method] !== 'function') {
      console.error(`Method not found: ${service}.${method}`)
      return null
    }
    return await handler[method](...args)
  } catch (error) {
    console.error(`Error calling ${service}.${method}:`, error)
    throw error
  }
}

// Atom API
export const AtomAPI = {
  create: (data) => safeCall('AtomHandler', 'CreateAtom', data),
  getByID: (id) => safeCall('AtomHandler', 'GetAtomByID', id),
  getByCategory: (data) => safeCall('AtomHandler', 'GetAtomsByCategory', data),
  update: (data) => safeCall('AtomHandler', 'UpdateAtom', data),
  delete: (id) => safeCall('AtomHandler', 'DeleteAtom', id),
  recordUsage: (id) => safeCall('AtomHandler', 'RecordUsage', id),
  findBySynonym: (term) => safeCall('AtomHandler', 'FindAtomsBySynonym', term),
  getPopular: (limit) => safeCall('AtomHandler', 'GetPopularAtoms', limit),
  batchImport: (data) => safeCall('AtomHandler', 'BatchImportAtoms', data),
  export: () => safeCall('AtomHandler', 'ExportAtoms'),
}

// Category API
export const CategoryAPI = {
  create: (data) => safeCall('CategoryHandler', 'CreateCategory', data),
  getByID: (id) => safeCall('CategoryHandler', 'GetCategoryByID', id),
  getByParent: (parentId, type) => safeCall('CategoryHandler', 'GetCategoriesByParent', parentId, type),
  getTree: (type) => safeCall('CategoryHandler', 'GetCategoryTree', type),
  update: (data) => safeCall('CategoryHandler', 'UpdateCategory', data),
  delete: (id) => safeCall('CategoryHandler', 'DeleteCategory', id),
  move: (data) => safeCall('CategoryHandler', 'MoveCategory', data),
  reorder: (data) => safeCall('CategoryHandler', 'ReorderCategories', data),
}

// Preset API
export const PresetAPI = {
  create: (data) => safeCall('PresetHandler', 'CreatePreset', data),
  getByID: (id) => safeCall('PresetHandler', 'GetPresetByID', id),
  getList: (data) => safeCall('PresetHandler', 'GetPresets', data),
  update: (data) => safeCall('PresetHandler', 'UpdatePreset', data),
  softDelete: (id) => safeCall('PresetHandler', 'SoftDeletePreset', id),
  restore: (id) => safeCall('PresetHandler', 'RestorePreset', id),
  buildPrompt: (data) => safeCall('PresetHandler', 'BuildPrompt', data),
  getWorkState: (presetId) => safeCall('PresetHandler', 'GetCurrentWorkState', presetId),
  fork: (data) => safeCall('PresetHandler', 'ForkPreset', data),
  cleanup: (data) => safeCall('PresetHandler', 'CleanupOldVersions', data),
}

// Version API
export const VersionAPI = {
  create: (data) => safeCall('VersionHandler', 'CreateVersion', data),
  get: (presetId, versionNum) => safeCall('VersionHandler', 'GetVersion', presetId, versionNum),
  getHistory: (presetId, limit) => safeCall('VersionHandler', 'GetVersionHistory', presetId, limit),
  getLatest: (presetId, count) => safeCall('VersionHandler', 'GetLatestVersions', presetId, count),
  star: (data) => safeCall('VersionHandler', 'StarVersion', data),
  rollback: (presetId, versionNum) => safeCall('VersionHandler', 'RollbackToVersion', presetId, versionNum),
  compare: (presetId, v1, v2) => safeCall('VersionHandler', 'CompareVersions', presetId, v1, v2),
  delete: (versionId) => safeCall('VersionHandler', 'DeleteVersion', versionId),
  getStarred: (presetId) => safeCall('VersionHandler', 'GetStarredVersions', presetId),
}

// Search API
export const SearchAPI = {
  search: (data) => safeCall('SearchHandler', 'Search', data),
  searchAtoms: (data) => safeCall('SearchHandler', 'SearchAtoms', data),
  quickSearch: (term) => safeCall('SearchHandler', 'QuickSearch', term),
  searchPresets: (data) => safeCall('SearchHandler', 'SearchPresets', data),
  reindex: () => safeCall('SearchHandler', 'ReindexAll'),
}

// AI API
export const AIAPI = {
  explode: (data) => safeCall('AIHandler', 'ExplodePrompt', data),
  importExtracted: (data) => safeCall('AIHandler', 'ImportExtractedAtoms', data),
  optimize: (data) => safeCall('AIHandler', 'OptimizePrompt', data),
  reverseImage: (data) => safeCall('AIHandler', 'ReverseImagePrompt', data),
  saveConfig: (data) => safeCall('AIHandler', 'SaveAIConfig', data),
  getConfig: () => safeCall('AIHandler', 'GetAIConfig'),
}

// Initialize Wails event listeners
export function initWailsListeners(callbacks) {
  if (!isWails() || !window.runtime) return
  
  // Listen for save shortcut
  window.runtime.EventsOn('save-preset', () => {
    callbacks.onSave?.()
  })
  
  // Listen for other events
  window.runtime.EventsOn('export-data', () => {
    callbacks.onExport?.()
  })
}
