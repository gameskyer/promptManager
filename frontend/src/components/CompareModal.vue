<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>版本对比</h3>
        <button class="close-btn" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        正在对比...
      </div>
      
      <div v-else-if="diff" class="diff-content">
        <div class="diff-header">
          <div class="version-header">
            <span class="version-badge old">V{{ version1 }}</span>
            <span class="version-date">{{ formatDate(diff.version1.created_at) }}</span>
          </div>
          <div class="diff-summary">
            <div class="summary-item added">
              <PlusIcon class="w-4 h-4" />
              {{ diff.added.length }}
            </div>
            <div class="summary-item removed">
              <MinusIcon class="w-4 h-4" />
              {{ diff.removed.length }}
            </div>
          </div>
          <div class="version-header">
            <span class="version-badge new">V{{ version2 }}</span>
            <span class="version-date">{{ formatDate(diff.version2.created_at) }}</span>
          </div>
        </div>
        
        <div class="diff-body">
          <!-- Added Atoms -->
          <div v-if="diff.added.length > 0" class="diff-section">
            <div class="section-title added">
              <PlusIcon class="w-4 h-4" />
              新增原子词 ({{ diff.added.length }})
            </div>
            <div class="atom-list">
              <span v-for="atom in diff.added" :key="atom" class="atom-tag added">
                {{ atom }}
              </span>
            </div>
          </div>
          
          <!-- Removed Atoms -->
          <div v-if="diff.removed.length > 0" class="diff-section">
            <div class="section-title removed">
              <MinusIcon class="w-4 h-4" />
              删除原子词 ({{ diff.removed.length }})
            </div>
            <div class="atom-list">
              <span v-for="atom in diff.removed" :key="atom" class="atom-tag removed">
                {{ atom }}
              </span>
            </div>
          </div>
          
          <!-- Parameter Changes -->
          <div v-if="diff.params_changed" class="diff-section">
            <div class="section-title changed">
              <AdjustmentsHorizontalIcon class="w-4 h-4" />
              参数变更
            </div>
            <div class="param-changes">
              <div class="param-row">
                <span class="param-name">Steps</span>
                <span class="param-old">{{ diff.version1.snapshot.params?.steps || 20 }}</span>
                <ArrowRightIcon class="w-4 h-4" />
                <span class="param-new">{{ diff.version2.snapshot.params?.steps || 20 }}</span>
              </div>
              <div class="param-row">
                <span class="param-name">CFG Scale</span>
                <span class="param-old">{{ diff.version1.snapshot.params?.cfg || 7 }}</span>
                <ArrowRightIcon class="w-4 h-4" />
                <span class="param-new">{{ diff.version2.snapshot.params?.cfg || 7 }}</span>
              </div>
            </div>
          </div>
          
          <!-- Prompt Comparison -->
          <div class="diff-section">
            <div class="section-title">
              <DocumentTextIcon class="w-4 h-4" />
              完整提示词
            </div>
            <div class="prompt-comparison">
              <div class="prompt-box old">
                <div class="prompt-label">V{{ version1 }}</div>
                <div class="prompt-text">{{ diff.version1.snapshot.pos_text }}</div>
              </div>
              <div class="prompt-box new">
                <div class="prompt-label">V{{ version2 }}</div>
                <div class="prompt-text">{{ diff.version2.snapshot.pos_text }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  XMarkIcon,
  PlusIcon,
  MinusIcon,
  AdjustmentsHorizontalIcon,
  ArrowRightIcon,
  DocumentTextIcon,
} from '@heroicons/vue/24/outline'
import { useVersionStore } from '../stores'

const props = defineProps({
  presetId: Number,
  version1: Number,
  version2: Number,
})

defineEmits(['close'])

const versionStore = useVersionStore()
const loading = ref(true)
const diff = ref(null)

onMounted(async () => {
  try {
    diff.value = await versionStore.compareVersions(
      props.presetId,
      props.version1,
      props.version2
    )
  } finally {
    loading.value = false
  }
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
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

.modal-content {
  width: 100%;
  max-width: 800px;
  max-height: 80vh;
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
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
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

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 60px 20px;
  color: #64748b;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #334155;
  border-top-color: #0ea5e9;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.diff-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.diff-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background-color: #1e293b;
  border-radius: 8px;
  margin-bottom: 20px;
}

.version-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.version-badge {
  font-size: 14px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
}

.version-badge.old {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.version-badge.new {
  background-color: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.version-date {
  font-size: 11px;
  color: #64748b;
}

.diff-summary {
  display: flex;
  gap: 16px;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  font-weight: 600;
}

.summary-item.added {
  color: #22c55e;
}

.summary-item.removed {
  color: #ef4444;
}

.diff-body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.diff-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #94a3b8;
}

.section-title.added {
  color: #22c55e;
}

.section-title.removed {
  color: #ef4444;
}

.section-title.changed {
  color: #f59e0b;
}

.atom-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.atom-tag {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.atom-tag.added {
  background-color: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.atom-tag.removed {
  background-color: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  text-decoration: line-through;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.param-changes {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.param-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background-color: #1e293b;
  border-radius: 6px;
}

.param-name {
  width: 80px;
  font-size: 12px;
  color: #64748b;
}

.param-old {
  flex: 1;
  font-size: 13px;
  color: #ef4444;
  text-decoration: line-through;
}

.param-new {
  flex: 1;
  font-size: 13px;
  color: #22c55e;
}

.param-row .w-4 {
  color: #475569;
}

.prompt-comparison {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.prompt-box {
  padding: 12px;
  border-radius: 8px;
  background-color: #1e293b;
  border: 1px solid #334155;
}

.prompt-box.old {
  border-color: rgba(239, 68, 68, 0.3);
}

.prompt-box.new {
  border-color: rgba(34, 197, 94, 0.3);
}

.prompt-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  margin-bottom: 8px;
  text-transform: uppercase;
}

.prompt-text {
  font-size: 12px;
  line-height: 1.6;
  color: #cbd5e1;
  max-height: 120px;
  overflow-y: auto;
}
</style>
