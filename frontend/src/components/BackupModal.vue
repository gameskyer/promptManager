<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>
          <ArrowDownTrayIcon class="w-5 h-5" />
          数据导入/导出
        </h3>
        <button class="close-btn" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      
      <!-- 标签页 -->
      <div class="tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'export' }"
          @click="activeTab = 'export'"
        >
          <ArrowUpTrayIcon class="w-4 h-4" />
          导出备份
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'import' }"
          @click="activeTab = 'import'"
        >
          <ArrowDownTrayIcon class="w-4 h-4" />
          导入数据
        </button>
      </div>
      
      <div class="modal-body">
        <!-- 导出页面 -->
        <template v-if="activeTab === 'export'">
          <div class="section">
            <div class="section-title">导出数据</div>
            <p class="section-desc">将所有数据（分类、原子词、预设）导出为 JSON 文件</p>
            
            <div class="export-options">
              <button class="export-btn" @click="exportToJSON" :disabled="loading">
                <DocumentTextIcon class="w-5 h-5" />
                <div class="btn-content">
                  <span class="btn-title">导出 JSON</span>
                  <span class="btn-desc">纯文本格式，便于查看和编辑</span>
                </div>
              </button>
              
              <button class="export-btn" @click="exportToZip" :disabled="loading">
                <FolderIcon class="w-5 h-5" />
                <div class="btn-content">
                  <span class="btn-title">导出 ZIP</span>
                  <span class="btn-desc">包含数据 + 图片的完整备份</span>
                </div>
              </button>
            </div>
          </div>
          
          <!-- 备份历史 -->
          <div class="section">
            <div class="section-header">
              <div class="section-title">历史备份</div>
              <button class="refresh-btn" @click="loadBackupList" :disabled="loading">
                <ArrowPathIcon class="w-4 h-4" :class="{ spinning: loading }" />
              </button>
            </div>
            
            <div v-if="backupList.length === 0" class="empty-state">
              <FolderOpenIcon class="w-8 h-8 text-slate-600" />
              <span>暂无备份文件</span>
            </div>
            
            <div v-else class="backup-list">
              <div 
                v-for="backup in backupList" 
                :key="backup"
                class="backup-item"
              >
                <DocumentIcon class="w-4 h-4 text-slate-500" />
                <span class="backup-name">{{ backup }}</span>
                <button class="action-btn" @click="downloadBackup(backup)">
                  <ArrowDownTrayIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </template>
        
        <!-- 导入页面 -->
        <template v-if="activeTab === 'import'">
          <div class="section">
            <div class="section-title">导入数据</div>
            <p class="section-desc">从 JSON 文件导入数据到系统中</p>
            
            <div class="import-options">
              <label class="import-mode">
                <input 
                  type="radio" 
                  v-model="importMode" 
                  value="merge"
                />
                <span class="radio-text">
                  <span class="radio-title">合并导入</span>
                  <span class="radio-desc">保留现有数据，导入新数据或更新已有数据</span>
                </span>
              </label>
              
              <label class="import-mode">
                <input 
                  type="radio" 
                  v-model="importMode" 
                  value="replace"
                />
                <span class="radio-text">
                  <span class="radio-title">覆盖导入</span>
                  <span class="radio-desc">删除所有现有数据，仅保留导入的数据</span>
                </span>
              </label>
            </div>
            
            <div class="file-input-area">
              <input
                ref="fileInput"
                type="file"
                accept=".json"
                @change="handleFileSelect"
                class="hidden"
              />
              <div 
                class="drop-zone"
                :class="{ dragging: isDragging }"
                @click="$refs.fileInput.click()"
                @dragover.prevent="isDragging = true"
                @dragleave="isDragging = false"
                @drop.prevent="handleDrop"
              >
                <CloudArrowUpIcon class="w-8 h-8" />
                <span>点击选择文件或拖拽到此处</span>
                <span class="hint">支持 .json 格式</span>
              </div>
            </div>
            
            <div v-if="importPreview" class="import-preview">
              <div class="preview-title">数据预览</div>
              <div class="preview-stats">
                <div class="stat-item">
                  <span class="stat-label">版本</span>
                  <span class="stat-value">{{ importPreview.version }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">分类</span>
                  <span class="stat-value">{{ importPreview.categories?.length || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">原子词</span>
                  <span class="stat-value">{{ importPreview.atoms?.length || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">预设</span>
                  <span class="stat-value">{{ importPreview.presets?.length || 0 }}</span>
                </div>
              </div>
              <button 
                class="import-btn" 
                @click="confirmImport"
                :disabled="loading"
              >
                <ArrowDownTrayIcon class="w-4 h-4" />
                {{ loading ? '导入中...' : '确认导入' }}
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  XMarkIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  DocumentTextIcon,
  FolderIcon,
  ArrowPathIcon,
  FolderOpenIcon,
  DocumentIcon,
  CloudArrowUpIcon,
} from '@heroicons/vue/24/outline'
import { useBackupStore } from '../stores'

const emit = defineEmits(['close', 'imported'])

const backupStore = useBackupStore()
const loading = ref(false)
const activeTab = ref('export')
const backupList = ref([])
const importMode = ref('merge')
const isDragging = ref(false)
const fileInput = ref(null)
const importPreview = ref(null)
const importData = ref('')

onMounted(() => {
  loadBackupList()
})

async function loadBackupList() {
  backupList.value = await backupStore.fetchBackupList()
}

async function exportToJSON() {
  const result = await backupStore.exportData()
  if (result.success) {
    const filename = `promptmaster_backup_${formatDate(new Date())}.json`
    backupStore.downloadAsFile(result.data, filename)
    await loadBackupList()
  } else {
    alert('导出失败: ' + result.error)
  }
}

async function exportToZip() {
  const result = await backupStore.exportToZip()
  if (result.success) {
    alert(`备份已保存到: ${result.path}`)
    await loadBackupList()
  } else {
    alert('导出失败: ' + result.error)
  }
}

async function downloadBackup(filename) {
  const data = await backupStore.readBackupFile(filename)
  if (data) {
    backupStore.downloadAsFile(data, filename)
  } else {
    alert('读取备份文件失败')
  }
}

function handleFileSelect(e) {
  const file = e.target.files[0]
  if (file) {
    readFile(file)
  }
}

function handleDrop(e) {
  isDragging.value = false
  const file = e.dataTransfer.files[0]
  if (file && file.name.endsWith('.json')) {
    readFile(file)
  } else {
    alert('请选择 JSON 文件')
  }
}

function readFile(file) {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const data = JSON.parse(e.target.result)
      importPreview.value = data
      importData.value = e.target.result
    } catch (err) {
      alert('无效的 JSON 文件: ' + err.message)
    }
  }
  reader.readAsText(file)
}

async function confirmImport() {
  if (!importData.value) return
  
  const confirmed = confirm(
    importMode.value === 'merge' 
      ? '确定要导入数据吗？现有数据将与导入数据合并。'
      : '确定要覆盖导入吗？这将删除所有现有数据！'
  )
  
  if (!confirmed) return
  
  loading.value = true
  const result = await backupStore.importData(importData.value, importMode.value === 'merge')
  loading.value = false
  
  if (result.success) {
    alert('导入成功！')
    emit('imported')
    emit('close')
  } else {
    alert('导入失败: ' + result.error)
  }
}

function formatDate(date) {
  return date.toISOString().slice(0, 19).replace(/:/g, '-')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  width: 100%;
  max-width: 500px;
  max-height: 85vh;
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #334155;
}

.modal-header h3 {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.close-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.tabs {
  display: flex;
  gap: 4px;
  padding: 0 20px;
  border-bottom: 1px solid #334155;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  color: #94a3b8;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: #e2e8f0;
}

.tab-btn.active {
  color: #0ea5e9;
  border-bottom-color: #0ea5e9;
}

.modal-body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.section {
  margin-bottom: 24px;
}

.section:last-child {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-desc {
  font-size: 13px;
  color: #94a3b8;
  margin-bottom: 16px;
}

.refresh-btn {
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

.refresh-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  color: #64748b;
  font-size: 13px;
}

/* Export Options */
.export-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.export-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 10px;
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.export-btn:hover:not(:disabled) {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.1);
}

.export-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
  text-align: left;
}

.btn-title {
  font-size: 14px;
  font-weight: 600;
}

.btn-desc {
  font-size: 12px;
  color: #64748b;
}

/* Backup List */
.backup-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.backup-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background-color: #1e293b;
  border-radius: 6px;
  transition: all 0.2s;
}

.backup-item:hover {
  background-color: #334155;
}

.backup-name {
  flex: 1;
  font-size: 13px;
  color: #e2e8f0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-btn {
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

.action-btn:hover {
  background-color: #0ea5e9;
  color: white;
}

/* Import Options */
.import-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.import-mode {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px;
  background-color: #1e293b;
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.import-mode:hover {
  border-color: #475569;
}

.import-mode input[type="radio"] {
  margin-top: 2px;
}

.radio-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.radio-title {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
}

.radio-desc {
  font-size: 12px;
  color: #64748b;
}

/* File Input */
.file-input-area {
  margin-bottom: 20px;
}

.hidden {
  display: none;
}

.drop-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 32px;
  background-color: #1e293b;
  border: 2px dashed #475569;
  border-radius: 10px;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s;
}

.drop-zone:hover,
.drop-zone.dragging {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.05);
  color: #0ea5e9;
}

.drop-zone .hint {
  font-size: 12px;
  color: #64748b;
}

/* Import Preview */
.import-preview {
  padding: 16px;
  background-color: #1e293b;
  border-radius: 10px;
}

.preview-title {
  font-size: 13px;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 12px;
}

.preview-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 10px;
  background-color: #0f172a;
  border-radius: 6px;
}

.stat-label {
  font-size: 11px;
  color: #64748b;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
}

.import-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  background-color: #22c55e;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.import-btn:hover:not(:disabled) {
  background-color: #16a34a;
}

.import-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
