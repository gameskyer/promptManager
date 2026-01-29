<template>
  <div 
    class="atom-card"
    :class="{ 
      selected: isSelected,
      'in-workbench': inWorkbench 
    }"
  >
    <div class="card-header">
      <div class="atom-type" :class="atom.type?.toLowerCase()">
        {{ atom.type === 'Positive' ? '正' : '负' }}
      </div>
      <div class="header-actions">
        <button class="edit-btn" @click.stop="$emit('edit', atom)">
          <PencilIcon class="w-3 h-3" />
        </button>
      </div>
    </div>
    
    <div class="card-body" @click="handleClick">
      <div class="atom-value">{{ atom.value }}</div>
      <div class="atom-label">{{ atom.label }}</div>
    </div>
    
    <!-- 近义词区域 -->
    <div v-if="atom.synonyms?.length > 0" class="card-footer" @click="handleClick">
      <div class="synonyms" @click.stop="handleSynonymsClick">
        <span 
          v-for="(syn, index) in displayedSynonyms" 
          :key="index"
          class="synonym-tag"
        >
          {{ syn }}
        </span>
        <span v-if="hasMoreSynonyms" class="synonym-more">
          +{{ remainingCount }}
        </span>
      </div>
    </div>
    
    <!-- 添加到工作区按钮 -->
    <div class="card-actions">
      <button 
        v-if="!isSelected"
        class="action-btn add"
        @click.stop="handleAdd"
        title="添加到工作区"
      >
        <PlusIcon class="w-4 h-4" />
      </button>
      <button 
        v-else
        class="action-btn remove"
        @click.stop="handleRemove"
        title="从工作区移除"
      >
        <CheckIcon class="w-4 h-4" />
      </button>
    </div>
    
    <!-- 近义词完整列表弹窗 -->
    <div 
      v-if="showSynonymsModal" 
      class="synonyms-modal"
      @click.stop="showSynonymsModal = false"
    >
      <div class="synonyms-modal-content" @click.stop>
        <div class="synonyms-modal-header">
          <span>近义词列表</span>
          <button @click="showSynonymsModal = false">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
        <div class="synonyms-modal-body">
          <div class="synonym-item main">
            <span class="label">原词</span>
            <span class="value">{{ atom.value }}</span>
          </div>
          <div class="synonym-item main">
            <span class="label">中文</span>
            <span class="value">{{ atom.label }}</span>
          </div>
          <div class="synonyms-divider"></div>
          <div 
            v-for="(syn, index) in atom.synonyms" 
            :key="index"
            class="synonym-item"
          >
            <span class="label">{{ index + 1 }}</span>
            <span class="value">{{ syn }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { PlusIcon, CheckIcon, PencilIcon, XMarkIcon } from '@heroicons/vue/24/solid'

const props = defineProps({
  atom: {
    type: Object,
    required: true,
  },
  isSelected: {
    type: Boolean,
    default: false,
  },
  inWorkbench: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['toggle', 'add', 'remove', 'edit'])

const showSynonymsModal = ref(false)
const MAX_SYNONYMS_DISPLAY = 3

// 显示的近义词（最多3条）
const displayedSynonyms = computed(() => {
  if (!props.atom.synonyms || props.atom.synonyms.length === 0) return []
  return props.atom.synonyms.slice(0, MAX_SYNONYMS_DISPLAY)
})

// 是否有更多近义词
const hasMoreSynonyms = computed(() => {
  return props.atom.synonyms && props.atom.synonyms.length > MAX_SYNONYMS_DISPLAY
})

// 剩余近义词数量
const remainingCount = computed(() => {
  if (!props.atom.synonyms) return 0
  return props.atom.synonyms.length - MAX_SYNONYMS_DISPLAY
})

function handleClick() {
  emit('toggle', props.atom)
}

function handleAdd() {
  emit('add', props.atom)
}

function handleRemove() {
  emit('remove', props.atom.id)
}

function handleSynonymsClick() {
  if (hasMoreSynonyms.value) {
    showSynonymsModal.value = true
  }
}
</script>

<style scoped>
.atom-card {
  position: relative;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  overflow: hidden;
  /* 高度自适应 */
  height: auto;
  display: flex;
  flex-direction: column;
}

.atom-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  border-color: #475569;
}

.atom-card.selected {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.1);
}

.atom-card.in-workbench {
  background-color: #0f172a;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.atom-type {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.atom-type.positive {
  background-color: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.atom-type.negative {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.edit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  opacity: 0;
  transition: all 0.2s;
}

.atom-card:hover .edit-btn {
  opacity: 1;
}

.edit-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.card-body {
  flex: 1;
  min-height: 0;
  margin-bottom: 6px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.atom-value {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
  line-height: 1.3;
  word-break: break-all;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.atom-label {
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.3;
  word-break: break-all;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  margin-top: 4px;
  flex-shrink: 0;
}

.synonyms {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
}

.synonym-tag {
  font-size: 10px;
  color: #64748b;
  background-color: rgba(255, 255, 255, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.synonym-more {
  font-size: 10px;
  color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.15);
  padding: 2px 6px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.synonym-more:hover {
  background-color: rgba(14, 165, 233, 0.25);
  color: #38bdf8;
}

.card-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.atom-card:hover .card-actions {
  opacity: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.add {
  background-color: #0284c7;
  color: white;
}

.action-btn.add:hover {
  background-color: #0ea5e9;
}

.action-btn.remove {
  background-color: #22c55e;
  color: white;
}

.action-btn.remove:hover {
  background-color: #16a34a;
}

.atom-card.selected .action-btn.remove {
  opacity: 1;
}

/* 近义词弹窗 */
.synonyms-modal {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.synonyms-modal-content {
  background-color: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  min-width: 240px;
  max-width: 320px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.synonyms-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid #334155;
}

.synonyms-modal-header span {
  font-size: 13px;
  font-weight: 600;
  color: #e2e8f0;
}

.synonyms-modal-header button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.synonyms-modal-header button:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.synonyms-modal-body {
  padding: 8px;
  max-height: 240px;
  overflow-y: auto;
}

.synonym-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 4px;
}

.synonym-item:hover {
  background-color: rgba(255, 255, 255, 0.03);
}

.synonym-item.main {
  background-color: rgba(14, 165, 233, 0.1);
  margin-bottom: 4px;
}

.synonym-item .label {
  font-size: 11px;
  color: #64748b;
  min-width: 32px;
}

.synonym-item.main .label {
  color: #0ea5e9;
  font-weight: 500;
}

.synonym-item .value {
  font-size: 13px;
  color: #e2e8f0;
}

.synonyms-divider {
  height: 1px;
  background-color: #334155;
  margin: 6px 0;
}
</style>
