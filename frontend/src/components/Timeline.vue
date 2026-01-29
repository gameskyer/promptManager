<template>
  <div class="timeline-drawer" :style="drawerStyle">
    <!-- 高度调整手柄（顶部） -->
    <div class="resize-handle-top" @mousedown="startResize" title="拖拽调整高度">
      <div class="resize-indicator"></div>
    </div>
    
    <div class="timeline-header">
      <h3 class="timeline-title">
        <ClockIcon class="w-5 h-5" />
        版本历史
      </h3>
      <div class="timeline-actions">
        <button v-if="selectedVersions.length === 2" class="compare-btn" @click="showCompare = true">
          <ArrowsRightLeftIcon class="w-4 h-4" />
          对比
        </button>
        <button class="icon-btn" @click="toggleTimeline">
          <ChevronDownIcon class="w-5 h-5" />
        </button>
      </div>
    </div>
    
    <div class="timeline-content">
      <!-- Current Editing State -->
      <div class="timeline-section">
        <div class="section-label">现在（当前编辑）</div>
        <div class="current-state">
          <div class="state-dot current"></div>
          <div class="state-info">
            <span class="state-label">未保存的更改</span>
            <span class="state-hint">按 Ctrl+S 保存为新版本</span>
          </div>
        </div>
      </div>
      
      <!-- Version List -->
      <div class="timeline-list">
        <div v-for="version in versionHistory" :key="version.id" class="timeline-item" :class="{ 'is-starred': version.is_starred, 'is-selected': selectedVersions.includes(version.version_num), 'is-current': currentVersionNum === version.version_num, 'is-expanded': expandedVersions.includes(version.id) }">
          <!-- 主要信息行 -->
          <div class="item-main" @click="toggleExpand(version)">
            <div class="item-left">
              <!-- Selection checkbox for compare -->
              <div class="select-box" :class="{ selected: selectedVersions.includes(version.version_num) }" @click.stop="toggleSelect(version.version_num)">
                <CheckIcon v-if="selectedVersions.includes(version.version_num)" class="w-3 h-3" />
              </div>
              
              <!-- Star status -->
              <button class="star-btn" :class="{ starred: version.is_starred }" @click.stop="toggleStar(version)">
                <StarIcon class="w-4 h-4" />
              </button>
              
              <!-- Version number -->
              <div class="version-label">
                <span class="version-badge">V{{ version.version_num }}</span>
              </div>
            </div>
            
            <div class="item-center">
              <!-- Thumbnail - 支持显示版本的预览图 -->
              <div class="thumbnail" @click.stop="viewVersionImages(version)">
                <img v-if="version.thumbnail_path" :src="version.thumbnail_path" alt="" />
                <div v-else-if="version.previews?.length > 0" class="thumbnail-stack">
                  <img :src="version.previews[0]" alt="" />
                  <span class="preview-count">{{ version.previews.length }}</span>
                </div>
                <PhotoIcon v-else class="w-6 h-6 text-slate-600" />
              </div>
              
              <!-- Timestamp -->
              <div class="timestamp">
                {{ version.formattedTime }}
              </div>
            </div>
            
            <div class="item-expand">
              <!-- 展开/折叠图标 -->
              <ChevronDownIcon 
                class="expand-icon"
                :class="{ expanded: expandedVersions.includes(version.id) }"
              />
            </div>
            
            <div class="item-right">
              <!-- Diff stats -->
              <div class="diff-badge" :title="getDiffTooltip(version)">
                <span v-if="version.diff_stats" class="diff-stats">
                  <span class="added">+{{ getAdded(version.diff_stats) }}</span>
                  <span class="separator">/</span>
                  <span class="removed">-{{ getRemoved(version.diff_stats) }}</span>
                </span>
              </div>
              
              <!-- Actions -->
              <div class="item-actions">
                <button class="action-btn view-detail" @click.stop="viewVersionDetail(version)" title="查看版本详情">
                  <EyeIcon class="w-4 h-4" />
                </button>
                <button class="action-btn view-images" @click.stop="viewVersionImages(version)" title="查看预览图">
                  <PhotoIcon class="w-4 h-4" />
                </button>
                <button class="action-btn rollback" @click.stop="rollbackTo(version)" title="回滚到此版本">
                  <ArrowUturnLeftIcon class="w-4 h-4" />
              </button>
              <button class="action-btn fork" @click.stop="forkFrom(version)" title="基于此创建新预设">
                <DocumentDuplicateIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
          
          <!-- 展开的详情内容 -->
          <div v-if="expandedVersions.includes(version.id)" class="item-details">
            <div class="details-content">
              <!-- 正向提示词 -->
              <div class="detail-row">
                <span class="detail-label">正向词</span>
                <span class="detail-value positive">{{ version.snapshot?.pos_text || '无' }}</span>
              </div>
              <!-- 负向提示词 -->
              <div class="detail-row">
                <span class="detail-label">负向词</span>
                <span class="detail-value negative">{{ version.snapshot?.neg_text || '无' }}</span>
              </div>
              <!-- 参数 -->
              <div v-if="version.snapshot?.params" class="detail-params">
                <span 
                  v-for="(value, key) in version.snapshot.params" 
                  :key="key"
                  class="param-tag"
                >
                  {{ key }}: {{ value }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Compare Modal -->
    <CompareModal
      v-if="showCompare && selectedVersions.length === 2"
      :preset-id="currentPreset?.id"
      :version1="selectedVersions[0]"
      :version2="selectedVersions[1]"
      @close="showCompare = false"
    />
    
    <!-- Version Images Viewer -->
    <ImageViewer
      v-if="viewingVersion && viewingVersion.previews?.length > 0"
      :images="viewingVersion.previews"
      :initial-index="0"
      @close="viewingVersion = null"
    />
    
    <!-- Version Detail Modal -->
    <VersionDetailModal
      v-if="viewingVersionDetail"
      :version="viewingVersionDetail"
      :preset="currentPreset"
      @close="viewingVersionDetail = null"
      @rollback="handleRollback"
      @use="handleUseVersion"
      @toggle-star="handleToggleStar"
    />
  </div>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { storeToRefs } from 'pinia'
import {
  ClockIcon,
  ArrowsRightLeftIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  StarIcon,
  CheckIcon,
  PhotoIcon,
  ArrowUturnLeftIcon,
  DocumentDuplicateIcon,
  EyeIcon,
} from '@heroicons/vue/24/outline'
import { useAppStore, useVersionStore, usePresetStore } from '../stores'
import CompareModal from './CompareModal.vue'
import ImageViewer from './ImageViewer.vue'
import VersionDetailModal from './VersionDetailModal.vue'

const appStore = useAppStore()
const versionStore = useVersionStore()
const presetStore = usePresetStore()

const { currentPreset, currentVersion } = storeToRefs(appStore)
const { versionHistory } = storeToRefs(versionStore)

const selectedVersions = ref([])
const showCompare = ref(false)
const viewingVersion = ref(null)
const viewingVersionDetail = ref(null)
const expandedVersions = ref([])

// 高度调整相关
const drawerHeight = ref(400) // 默认高度
const isResizing = ref(false)
const startY = ref(0)
const startHeight = ref(0)

const drawerStyle = computed(() => {
  console.log('[Timeline] drawerStyle computed, height:', drawerHeight.value)
  return {
    height: `${drawerHeight.value}px`
  }
})

function startResize(e) {
  isResizing.value = true
  startY.value = e.clientY
  startHeight.value = drawerHeight.value
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
  e.preventDefault()
}

function onResize(e) {
  if (!isResizing.value) return
  const delta = startY.value - e.clientY
  const newHeight = startHeight.value + delta
  // 限制最小和最大高度
  drawerHeight.value = Math.max(150, Math.min(newHeight, window.innerHeight * 0.7))
}

function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
}

onUnmounted(() => {
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
})

const currentVersionNum = computed(() => {
  return currentVersion.value?.version_num || currentPreset.value?.current_version
})

function toggleTimeline() {
  appStore.toggleTimeline()
}

function handleVersionClick(version) {
  viewVersionDetail(version)
}

async function handleVersionDoubleClick(version) {
  if (confirm(`基于 V${version.version_num} 创建新版本？`)) {
    await versionStore.rollbackToVersion(currentPreset.value.id, version.version_num)
  }
}

// 展开/折叠版本详情
function toggleExpand(version) {
  console.log('[Timeline] toggleExpand called for version:', version.id, 'current expanded:', expandedVersions.value)
  const index = expandedVersions.value.indexOf(version.id)
  if (index === -1) {
    expandedVersions.value = [...expandedVersions.value, version.id]
  } else {
    expandedVersions.value = expandedVersions.value.filter(id => id !== version.id)
  }
  console.log('[Timeline] new expanded:', expandedVersions.value)
}

// 获取版本的提示词摘要
function getPromptSummary(version) {
  const snapshot = version.snapshot || {}
  const posText = snapshot.pos_text || ''
  if (!posText) return '无提示词'
  // 截取前 50 个字符
  if (posText.length <= 50) return posText
  return posText.substring(0, 50) + '...'
}

function toggleSelect(versionNum) {
  const index = selectedVersions.value.indexOf(versionNum)
  if (index === -1) {
    if (selectedVersions.value.length < 2) {
      selectedVersions.value.push(versionNum)
    } else {
      selectedVersions.value.shift()
      selectedVersions.value.push(versionNum)
    }
  } else {
    selectedVersions.value.splice(index, 1)
  }
}

async function toggleStar(version) {
  await versionStore.starVersion(version.id, !version.is_starred)
}

async function rollbackTo(version) {
  if (confirm(`基于 V${version.version_num} 创建新版本？这将保留当前历史。`)) {
    const newVersion = await versionStore.rollbackToVersion(
      currentPreset.value.id, 
      version.version_num
    )
    if (newVersion) {
      alert(`已创建新版本 V${newVersion.version_num}`)
    }
  }
}

async function forkFrom(version) {
  const newTitle = prompt('输入新预设名称：', `${currentPreset.value.title} - Fork`)
  if (newTitle) {
    await presetStore.forkPreset(currentPreset.value.id, version.version_num, newTitle)
  }
}

function viewVersionImages(version) {
  if (version.previews?.length > 0 || version.thumbnail_path) {
    // 构建预览图列表
    const previews = version.previews?.length > 0 
      ? version.previews 
      : [version.thumbnail_path]
    viewingVersion.value = { ...version, previews }
  }
}

// View version details
function viewVersionDetail(version) {
  viewingVersionDetail.value = version
}

// Handle rollback from detail modal
async function handleRollback(version) {
  if (confirm(`基于 V${version.version_num} 创建新版本？这将保留当前历史。`)) {
    const newVersion = await versionStore.rollbackToVersion(
      currentPreset.value.id, 
      version.version_num
    )
    if (newVersion) {
      alert(`已创建新版本 V${newVersion.version_num}`)
      viewingVersionDetail.value = null
    }
  }
}

// Handle use version from detail modal
async function handleUseVersion(version) {
  appStore.setCurrentVersion(version)
  if (version.snapshot) {
    // Restore workspace with version snapshot
    // TODO: Integrate with workspace store
  }
  viewingVersionDetail.value = null
  appStore.toggleTimeline() // Close timeline
}

// Handle toggle star from detail modal
async function handleToggleStar(version) {
  await versionStore.starVersion(version.id, !version.is_starred)
  // Update local version state
  version.is_starred = !version.is_starred
}

function getAdded(diffStats) {
  if (!diffStats) return 0
  const match = diffStats.match(/\+(\d+)/)
  return match ? parseInt(match[1]) : 0
}

function getRemoved(diffStats) {
  if (!diffStats) return 0
  const match = diffStats.match(/-(\d+)/)
  return match ? parseInt(match[1]) : 0
}

function getDiffTooltip(version) {
  const added = getAdded(version.diff_stats)
  const removed = getRemoved(version.diff_stats)
  if (added === 0 && removed === 0) return '无变化'
  return `添加: ${added}, 删除: ${removed}`
}
</script>

<style scoped>
.timeline-drawer {
  background-color: #0f172a;
  border-top: 1px solid #1e293b;
  display: flex;
  flex-direction: column;
  position: relative;
  min-height: 150px;
}

/* 高度调整手柄（顶部） */
.resize-handle-top {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 12px;
  cursor: ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  transition: background-color 0.2s;
  z-index: 10;
  transform: translateY(-50%);
}

.resize-handle-top:hover {
  background-color: rgba(14, 165, 233, 0.1);
}

.resize-handle-top .resize-indicator {
  width: 60px;
  height: 4px;
  background-color: #475569;
  border-radius: 2px;
  transition: background-color 0.2s;
}

.resize-handle-top:hover .resize-indicator {
  background-color: #0ea5e9;
}

.timeline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-bottom: 1px solid #1e293b;
}

.timeline-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.timeline-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.compare-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  background-color: #7c3aed;
  border: none;
  color: white;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.compare-btn:hover {
  background-color: #8b5cf6;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background-color: #1e293b;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.timeline-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
}

.timeline-section {
  margin-bottom: 16px;
}

.section-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.current-state {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background-color: rgba(14, 165, 233, 0.1);
  border: 1px dashed #0ea5e9;
  border-radius: 8px;
}

.state-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: #0ea5e9;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.state-info {
  display: flex;
  flex-direction: column;
}

.state-label {
  font-size: 13px;
  font-weight: 500;
  color: #e2e8f0;
}

.state-hint {
  font-size: 11px;
  color: #64748b;
}

.timeline-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.timeline-item {
  display: flex;
  flex-direction: column;
  padding: 0;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
}

.timeline-item:hover {
  border-color: #475569;
  background-color: #252f47;
}

.timeline-item.is-selected {
  border-color: #7c3aed;
  background-color: rgba(124, 58, 237, 0.1);
}

.timeline-item.is-starred {
  border-color: #f59e0b;
}

.timeline-item.is-current {
  border-left: 3px solid #0ea5e9;
}

.item-main {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
}

.item-expand {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  margin-left: 12px;
  background: linear-gradient(135deg, #10b981, #3b82f6);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.4);
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.item-expand:hover {
  background: linear-gradient(135deg, #34d399, #60a5fa);
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.6);
}

.expand-icon {
  width: 20px;
  height: 20px;
  color: #ffffff;
  transition: transform 0.2s;
  display: block;
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

.item-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.select-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 4px;
  border: 2px solid #475569;
  background-color: transparent;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
}

.select-box.selected {
  background-color: #7c3aed;
  border-color: #7c3aed;
}

.star-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.star-btn:hover,
.star-btn.starred {
  color: #f59e0b;
}

.version-label {
  min-width: 40px;
}

.version-badge {
  font-size: 12px;
  font-weight: 600;
  color: #e2e8f0;
  font-variant-numeric: tabular-nums;
}

.item-center {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.thumbnail {
  position: relative;
  width: 48px;
  height: 48px;
  background-color: #0f172a;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid #334155;
  transition: all 0.2s;
}

.thumbnail:hover {
  border-color: #0ea5e9;
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-stack {
  position: relative;
  width: 100%;
  height: 100%;
}

.thumbnail-stack img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-count {
  position: absolute;
  bottom: 2px;
  right: 2px;
  padding: 1px 4px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  font-size: 10px;
  border-radius: 3px;
}

.timestamp {
  font-size: 12px;
  color: #94a3b8;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.diff-badge {
  font-size: 12px;
  font-weight: 500;
}

.diff-stats {
  display: flex;
  align-items: center;
  gap: 2px;
}

.diff-stats .added {
  color: #22c55e;
}

.diff-stats .removed {
  color: #ef4444;
}

.diff-stats .separator {
  color: #475569;
}

.item-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.timeline-item:hover .item-actions {
  opacity: 1;
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
  background-color: #334155;
  color: #e2e8f0;
}

.action-btn.view-detail:hover {
  background-color: rgba(168, 85, 247, 0.2);
  color: #a855f7;
}

.action-btn.view-images:hover {
  background-color: rgba(14, 165, 233, 0.2);
  color: #0ea5e9;
}

.action-btn.rollback:hover {
  background-color: rgba(14, 165, 233, 0.2);
  color: #0ea5e9;
}

.action-btn.fork:hover {
  background-color: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

/* 展开详情 */
.item-details {
  border-top: 1px solid #334155;
  background-color: #0f172a;
  padding: 12px;
}

.details-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12px;
  line-height: 1.5;
}

.detail-label {
  flex-shrink: 0;
  padding: 2px 8px;
  background-color: #334155;
  border-radius: 4px;
  color: #94a3b8;
  font-weight: 500;
}

.detail-value {
  flex: 1;
  color: #cbd5e1;
  word-break: break-word;
}

.detail-value.positive {
  color: #86efac;
}

.detail-value.negative {
  color: #fca5a5;
}

.detail-params {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.param-tag {
  padding: 2px 8px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 4px;
  font-size: 11px;
  color: #64748b;
}
</style>
