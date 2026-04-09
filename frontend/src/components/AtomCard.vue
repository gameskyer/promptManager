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
        <!-- 复制按钮 -->
        <button 
          class="header-btn copy-btn" 
          @click.stop="copyAtomContent"
          title="复制内容"
        >
          <ClipboardDocumentIcon class="w-3 h-3" />
        </button>
        <button 
          class="header-btn edit-btn" 
          @click.stop="$emit('edit', atom)"
          title="编辑"
        >
          <PencilIcon class="w-3 h-3" />
        </button>
      </div>
    </div>
    
    <div class="card-body" @click="handleClick">
      <div class="atom-value">{{ atom.value }}</div>
      <div class="atom-label">{{ atom.label }}</div>
    </div>
    
    <!-- 近义词区域 - 每行显示一个 -->
    <div v-if="atom.synonyms?.length > 0" class="card-footer">
      <div class="synonyms-list">
        <div 
          v-for="(syn, index) in displayedSynonyms" 
          :key="index"
          class="synonym-row"
        >
          <span 
            class="synonym-text"
            @click.stop="addSynonym(syn)"
            title="点击添加"
          >
            {{ syn }}
          </span>
          <button 
            class="synonym-copy-btn"
            @click.stop="copyText(syn)"
            title="复制"
          >
            <ClipboardDocumentIcon class="w-3 h-3" />
          </button>
        </div>
        <!-- 查看更多按钮 -->
        <div 
          v-if="hasMoreSynonyms" 
          class="synonym-row more-row"
          @click.stop="showSynonymsModal = true"
        >
          <span class="synonym-more-text">
            查看更多 (+{{ remainingCount }})
          </span>
          <ChevronDownIcon class="w-3 h-3 more-icon" />
        </div>
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
          <!-- 原词信息 -->
          <div class="synonym-item main">
            <span class="label">原词</span>
            <span class="value">{{ atom.value }}</span>
            <button 
              class="item-copy-btn"
              @click="copyText(atom.value)"
              title="复制"
            >
              <ClipboardDocumentIcon class="w-3 h-3" />
            </button>
            <button 
              class="item-add-btn"
              @click="addAtom(atom)"
              title="添加"
            >
              <PlusIcon class="w-3 h-3" />
            </button>
          </div>
          <div class="synonym-item main">
            <span class="label">中文</span>
            <span class="value">{{ atom.label }}</span>
            <button 
              class="item-copy-btn"
              @click="copyText(atom.label)"
              title="复制"
            >
              <ClipboardDocumentIcon class="w-3 h-3" />
            </button>
          </div>
          <div class="synonyms-divider"></div>
          <!-- 所有近义词列表 -->
          <div 
            v-for="(syn, index) in atom.synonyms" 
            :key="index"
            class="synonym-item"
          >
            <span class="label">{{ index + 1 }}</span>
            <span class="value">{{ syn }}</span>
            <button 
              class="item-copy-btn"
              @click="copyText(syn)"
              title="复制"
            >
              <ClipboardDocumentIcon class="w-3 h-3" />
            </button>
            <button 
              class="item-add-btn"
              @click="addSynonym(syn)"
              title="添加"
            >
              <PlusIcon class="w-3 h-3" />
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 复制成功提示 -->
    <div v-if="showCopyToast" class="copy-toast">
      <CheckIcon class="w-4 h-4" />
      <span>已复制</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { 
  PlusIcon, 
  CheckIcon, 
  PencilIcon, 
  XMarkIcon,
  ClipboardDocumentIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/solid'

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

const emit = defineEmits(['toggle', 'add', 'remove', 'edit', 'add-synonym'])

const showSynonymsModal = ref(false)
const showCopyToast = ref(false)
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

function addAtom(atom) {
  emit('add', atom)
}

function addSynonym(synonym) {
  // 创建一个临时的原子词对象，用于添加
  const synonymAtom = {
    ...props.atom,
    value: synonym,
    label: synonym,
  }
  emit('add-synonym', synonymAtom)
  emit('add', synonymAtom)
}

// 复制文本到剪贴板
async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text)
    showCopyToastMessage()
  } catch (err) {
    // 降级方案
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    showCopyToastMessage()
  }
}

// 复制原子词完整内容
function copyAtomContent() {
  const content = `${props.atom.value} (${props.atom.label})`
  copyText(content)
}

// 显示复制成功提示
function showCopyToastMessage() {
  showCopyToast.value = true
  setTimeout(() => {
    showCopyToast.value = false
  }, 1500)
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
  gap: 4px;
}

.header-btn {
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

.atom-card:hover .header-btn {
  opacity: 1;
}

.header-btn:hover {
  background-color: #334155;
  color: #e2e8f0;
}

.header-btn.copy-btn:hover {
  background-color: rgba(14, 165, 233, 0.2);
  color: #0ea5e9;
}

.card-body {
  flex: 1;
  min-height: 0;
  margin-bottom: 8px;
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

/* 近义词列表 - 每行一个 */
.synonyms-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.synonym-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  padding: 4px 6px;
  border-radius: 4px;
  background-color: rgba(255, 255, 255, 0.03);
  transition: all 0.2s;
}

.synonym-row:hover {
  background-color: rgba(255, 255, 255, 0.06);
}

.synonym-row.more-row {
  justify-content: center;
  background-color: rgba(14, 165, 233, 0.1);
  cursor: pointer;
}

.synonym-row.more-row:hover {
  background-color: rgba(14, 165, 233, 0.2);
}

.synonym-text {
  flex: 1;
  font-size: 11px;
  color: #94a3b8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
  transition: color 0.2s;
}

.synonym-text:hover {
  color: #38bdf8;
}

.synonym-copy-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 3px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  opacity: 0;
  transition: all 0.2s;
  flex-shrink: 0;
}

.synonym-row:hover .synonym-copy-btn {
  opacity: 1;
}

.synonym-copy-btn:hover {
  background-color: #334155;
  color: #0ea5e9;
}

.synonym-more-text {
  font-size: 11px;
  color: #0ea5e9;
  font-weight: 500;
}

.more-icon {
  color: #0ea5e9;
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
  min-width: 280px;
  max-width: 360px;
  max-height: 400px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.synonyms-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid #334155;
  flex-shrink: 0;
}

.synonyms-modal-header span {
  font-size: 14px;
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
  padding: 10px;
  overflow-y: auto;
  flex: 1;
}

.synonym-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.synonym-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.synonym-item.main {
  background-color: rgba(14, 165, 233, 0.08);
  margin-bottom: 4px;
}

.synonym-item.main:hover {
  background-color: rgba(14, 165, 233, 0.15);
}

.synonym-item .label {
  font-size: 11px;
  color: #64748b;
  min-width: 32px;
  flex-shrink: 0;
}

.synonym-item.main .label {
  color: #0ea5e9;
  font-weight: 500;
}

.synonym-item .value {
  flex: 1;
  font-size: 13px;
  color: #e2e8f0;
  word-break: break-all;
}

.item-copy-btn,
.item-add-btn {
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
  opacity: 0;
  transition: all 0.2s;
  flex-shrink: 0;
}

.synonym-item:hover .item-copy-btn,
.synonym-item:hover .item-add-btn {
  opacity: 1;
}

.item-copy-btn:hover {
  background-color: #334155;
  color: #0ea5e9;
}

.item-add-btn:hover {
  background-color: #334155;
  color: #22c55e;
}

.synonyms-divider {
  height: 1px;
  background-color: #334155;
  margin: 8px 0;
}

/* 复制成功提示 */
.copy-toast {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background-color: rgba(34, 197, 94, 0.9);
  border-radius: 6px;
  color: white;
  font-size: 13px;
  font-weight: 500;
  z-index: 50;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
}
</style>
