<template>
  <div class="trans-session">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <h2 class="title">
          <LanguageIcon class="w-5 h-5 text-emerald-400" />
          提示词翻译工作台
        </h2>
        <span class="subtitle">粘贴提示词，按逗号分割，按需翻译</span>
      </div>
      <div class="toolbar-actions">
        <!-- 翻译方向切换 -->
        <div class="direction-toggle">
          <button
            class="dir-btn"
            :class="{ active: direction === 'cn2en' }"
            @click="direction = 'cn2en'"
          >中→英</button>
          <button
            class="dir-btn"
            :class="{ active: direction === 'en2cn' }"
            @click="direction = 'en2cn'"
          >英→中</button>
        </div>

        <button class="btn-secondary" @click="clearAll" :disabled="words.length === 0">
          <TrashIcon class="w-4 h-4" />
          清空
        </button>
        <button class="btn-primary" @click="translateUntranslated" :disabled="!hasUntranslated || isTranslating">
          <SparklesIcon class="w-4 h-4" />
          {{ isTranslating ? '翻译中...' : '翻译待译词' }}
        </button>
        <button class="btn-accent" @click="copyOutput" :disabled="joinedOutput.length === 0">
          <ClipboardIcon class="w-4 h-4" />
          复制结果
        </button>
      </div>
    </div>

    <!-- 工作区 -->
    <div class="workspace">
      <!-- 左侧：导入区和词汇列表 -->
      <div class="left-panel">
        <!-- 粘贴导入区 -->
        <div class="paste-area">
          <textarea
            v-model="pasteText"
            class="paste-input"
            :placeholder="pastePlaceholder"
            rows="5"
          ></textarea>
          <button class="btn-primary paste-btn" @click="doSplit" :disabled="!pasteText.trim()">
            <ArrowPathIcon class="w-4 h-4" />
            分割导入
          </button>
        </div>

        <!-- 词汇列表 -->
        <div v-if="words.length > 0" class="words-section">
          <div class="section-header">
            <span class="section-label">提示词列表 ({{ words.length }})</span>
            <span class="stats">
              <span class="stat-badge pending">待译: {{ needCount }}</span>
              <span class="stat-badge done">已译: {{ translatedCount }}</span>
              <span class="stat-badge skip">跳过: {{ skipCount }}</span>
            </span>
          </div>

          <div class="word-list">
            <div
              v-for="(word, index) in words"
              :key="word.id"
              class="word-row"
              :class="{ 'needs-trans': word.needs_trans && !word.is_translated, 'is-done': word.is_translated && word.needs_trans }"
            >
              <span class="word-index">{{ index + 1 }}</span>

              <div class="word-field original">
                <input
                  v-model="word.original"
                  class="word-input"
                  @input="onWordEdited(word, index)"
                />
              </div>

              <div class="word-arrow">
                <ArrowRightIcon class="w-4 h-4" />
              </div>

              <div class="word-field translated">
                <input
                  v-model="word.translated"
                  class="word-input"
                  :class="{ 'highlight': word.needs_trans && word.is_translated }"
                  @input="onTranslationEdited(word)"
                />
              </div>

              <span v-if="word.needs_trans && !word.is_translated" class="word-tag pending">待翻译</span>
              <span v-else-if="word.needs_trans && word.is_translated" class="word-tag done">已翻译</span>
              <span v-else class="word-tag skip">无需译</span>

              <button class="word-action translate-single" v-if="word.needs_trans && !word.is_translated" @click="translateSingleWord(word, index)" :disabled="isTranslating">
                <SparklesIcon class="w-3 h-3" />
              </button>

              <button class="word-action delete" @click="removeWord(index)">
                <XMarkIcon class="w-3 h-3" />
              </button>
            </div>
          </div>

          <!-- 添加新词 -->
          <div class="add-word-row">
            <input
              v-model="newWordText"
              class="add-word-input"
              :placeholder="addPlaceholder"
              @keyup.enter="addNewWord"
            />
            <button class="btn-secondary btn-sm" @click="addNewWord" :disabled="!newWordText.trim()">
              <PlusIcon class="w-3 h-3" />
              添加
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：输出区 -->
      <div class="right-panel">
        <div class="panel-header">
          <span class="panel-title">翻译结果输出</span>
          <button
            class="btn-secondary btn-sm"
            @click="copyOutput"
            :disabled="joinedOutput.length === 0"
          >
            <ClipboardIcon class="w-3 h-3" />
            复制
          </button>
        </div>

        <div class="output-area">
          <textarea
            readonly
            :value="joinedOutput"
            class="output-text"
            placeholder="翻译后的结果将在这里显示..."
            rows="15"
            @click="$event.target.select()"
          ></textarea>
        </div>

        <div class="output-options">
          <label class="output-option">
            <input type="radio" v-model="separator" value="," />
            英文逗号
          </label>
          <label class="output-option">
            <input type="radio" v-model="separator" value="，" />
            中文逗号
          </label>
          <label class="output-option">
            <input type="checkbox" v-model="autoJoin" />
            自动拼接
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  ArrowPathIcon,
  ArrowRightIcon,
  ClipboardIcon,
  LanguageIcon,
  PlusIcon,
  SparklesIcon,
  TrashIcon,
  XMarkIcon,
} from '@heroicons/vue/24/outline'
import {
  SplitPrompt,
  TranslateWords,
  TranslateSingle,
} from '../lib/wailsjs/go/handlers/TransSessionHandler'

const pasteText = ref('')
const words = ref([])
const newWordText = ref('')
const isTranslating = ref(false)
const separator = ref(',')
const autoJoin = ref(true)
const direction = ref('cn2en') // 'cn2en' | 'en2cn'

const pastePlaceholder = computed(() => {
  if (direction.value === 'cn2en') {
    return '在此粘贴提示词文本...\n支持中文逗号（，）和英文逗号（,）分隔\n\n例如：\n精美画作，一个女孩，蓝色眼睛，masterpiece,best quality,1girl,blue_eyes'
  }
  return '在此粘贴提示词文本...\n支持中文逗号（，）和英文逗号（,）分隔\n\n例如：\nmasterpiece,1girl,blue eyes,pale skin,杰作，一个女孩'
})

const addPlaceholder = computed(() => {
  return direction.value === 'cn2en' ? '输入中文提示词（自动翻译）...' : '输入英文提示词（自动翻译）...'
})

const needCount = computed(() => words.value.filter(w => w.needs_trans && !w.is_translated).length)
const translatedCount = computed(() => words.value.filter(w => w.needs_trans && w.is_translated).length)
const skipCount = computed(() => words.value.filter(w => !w.needs_trans).length)
const hasUntranslated = computed(() => words.value.some(w => w.needs_trans && !w.is_translated))

const joinedOutput = computed(() => {
  if (!autoJoin.value) return ''
  return words.value.map(w => w.translated).join(separator.value)
})

function needsTranslation(text) {
  if (direction.value === 'cn2en') {
    return /[\u4e00-\u9fff]/.test(text)
  }
  return /[a-zA-Z]/.test(text) && !/[\u4e00-\u9fff]/.test(text)
}

async function doSplit() {
  if (!pasteText.value.trim()) return

  try {
    const response = await SplitPrompt({ text: pasteText.value, direction: direction.value })
    if (response.success) {
      const existingMap = new Map(words.value.map(w => [w.original, w]))
      const newWords = response.data.filter(nw => {
        if (existingMap.has(nw.original)) {
          return false
        }
        existingMap.set(nw.original, nw)
        return true
      })
      words.value = [...words.value, ...newWords]
      pasteText.value = ''
    } else {
      alert('分割失败: ' + response.error)
    }
  } catch (e) {
    alert('分割失败: ' + e.message)
  }
}

async function translateUntranslated() {
  if (!hasUntranslated.value || isTranslating.value) return

  isTranslating.value = true
  try {
    const response = await TranslateWords({ words: words.value, direction: direction.value })
    if (response.success) {
      words.value = response.data
    } else {
      alert('翻译失败: ' + response.error)
    }
  } catch (e) {
    alert('翻译失败: ' + e.message)
  } finally {
    isTranslating.value = false
  }
}

async function translateSingleWord(word, index) {
  if (isTranslating.value) return

  isTranslating.value = true
  try {
    const response = await TranslateSingle({ text: word.original, direction: direction.value })
    if (response.success) {
      words.value[index].translated = response.data
      words.value[index].is_translated = true
    } else {
      alert('翻译失败: ' + response.error)
    }
  } catch (e) {
    alert('翻译失败: ' + e.message)
  } finally {
    isTranslating.value = false
  }
}

function addNewWord() {
  const text = newWordText.value.trim()
  if (!text) return

  const needTrans = needsTranslation(text)
  words.value.push({
    id: crypto.randomUUID ? crypto.randomUUID() : Date.now().toString() + Math.random(),
    original: text,
    translated: text,
    needs_trans: needTrans,
    is_translated: !needTrans,
  })
  newWordText.value = ''

  if (needTrans) {
    translateSingleWord(words.value[words.value.length - 1], words.value.length - 1)
  }
}

function removeWord(index) {
  words.value.splice(index, 1)
}

function onWordEdited(word, index) {
  const needTrans = needsTranslation(word.original)
  word.needs_trans = needTrans
  if (!needTrans) {
    word.translated = word.original
    word.is_translated = true
  } else if (word.translated === word.original && word.is_translated) {
    word.is_translated = false
  }
}

function onTranslationEdited(word) {
  if (word.translated !== word.original && word.needs_trans) {
    word.is_translated = true
  }
}

function clearAll() {
  words.value = []
  pasteText.value = ''
  newWordText.value = ''
}

async function copyOutput() {
  const text = joinedOutput.value
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    const ta = document.createElement('textarea')
    ta.value = text
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
  }
}
</script>

<style scoped>
.trans-session {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #0f172a;
  overflow: hidden;
}

/* 工具栏 */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #1e293b;
  border-bottom: 1px solid #334155;
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.subtitle {
  font-size: 12px;
  color: #94a3b8;
  padding-left: 24px;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.direction-toggle {
  display: flex;
  border: 1px solid #334155;
  border-radius: 8px;
  overflow: hidden;
}

.dir-btn {
  padding: 5px 12px;
  font-size: 12px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  background: transparent;
  color: #94a3b8;
  transition: all 0.15s;
}

.dir-btn.active {
  background: #3b82f6;
  color: #fff;
}

.dir-btn:hover:not(.active) {
  background: #1e293b;
  color: #e2e8f0;
}

/* 主工作区 */
.workspace {
  display: flex;
  flex: 1;
  overflow: hidden;
  gap: 16px;
  padding: 16px;
}

/* 左侧面板 */
.left-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

/* 粘贴区 */
.paste-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}

.paste-input {
  width: 100%;
  padding: 10px 12px;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 13px;
  line-height: 1.6;
  resize: vertical;
  font-family: inherit;
}

.paste-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.15);
}

.paste-input::placeholder {
  color: #64748b;
}

.paste-btn {
  align-self: flex-end;
}

/* 词汇列表 */
.words-section {
  margin-top: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.section-label {
  font-size: 13px;
  font-weight: 600;
  color: #94a3b8;
}

.stats {
  display: flex;
  gap: 8px;
}

.stat-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
}

.stat-badge.pending {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.stat-badge.done {
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
}

.stat-badge.skip {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
}

.word-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-right: 4px;
}

.word-list::-webkit-scrollbar {
  width: 4px;
}

.word-list::-webkit-scrollbar-track {
  background: transparent;
}

.word-list::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 2px;
}

/* 词汇行 */
.word-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: #1e293b;
  border-radius: 8px;
  border: 1px solid #334155;
  transition: border-color 0.15s;
  flex-shrink: 0;
}

.word-row.needs-trans {
  border-color: rgba(251, 191, 36, 0.3);
  background: rgba(251, 191, 36, 0.03);
}

.word-row.is-done {
  border-color: rgba(52, 211, 153, 0.3);
  background: rgba(52, 211, 153, 0.03);
}

.word-index {
  font-size: 11px;
  color: #475569;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.word-field {
  flex: 1;
  min-width: 0;
}

.word-input {
  width: 100%;
  padding: 5px 8px;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 5px;
  color: #e2e8f0;
  font-size: 12px;
  font-family: inherit;
}

.word-input:focus {
  outline: none;
  border-color: #3b82f6;
}

.word-input.highlight {
  color: #34d399;
  border-color: rgba(52, 211, 153, 0.4);
}

.word-arrow {
  color: #475569;
  flex-shrink: 0;
}

.word-tag {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 8px;
  white-space: nowrap;
  flex-shrink: 0;
}

.word-tag.pending {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.word-tag.done {
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
}

.word-tag.skip {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
}

.word-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s;
}

.word-action.translate-single {
  background: rgba(139, 92, 246, 0.15);
  color: #a78bfa;
}

.word-action.translate-single:hover {
  background: rgba(139, 92, 246, 0.3);
}

.word-action.delete {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.word-action.delete:hover {
  background: rgba(239, 68, 68, 0.25);
}

.word-action:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* 添加新词行 */
.add-word-row {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  padding: 8px;
  background: #1e293b;
  border-radius: 8px;
  border: 1px dashed #334155;
  flex-shrink: 0;
}

.add-word-input {
  flex: 1;
  padding: 6px 10px;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #e2e8f0;
  font-size: 13px;
  font-family: inherit;
}

.add-word-input:focus {
  outline: none;
  border-color: #3b82f6;
}

.add-word-input::placeholder {
  color: #475569;
}

/* 右侧面板 */
.right-panel {
  width: 360px;
  display: flex;
  flex-direction: column;
  background: #1e293b;
  border-radius: 10px;
  border: 1px solid #334155;
  flex-shrink: 0;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #334155;
}

.panel-title {
  font-size: 13px;
  font-weight: 600;
  color: #94a3b8;
}

.output-area {
  flex: 1;
  padding: 12px;
  display: flex;
  flex-direction: column;
}

.output-text {
  flex: 1;
  width: 100%;
  padding: 12px;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #34d399;
  font-size: 13px;
  line-height: 1.7;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  resize: none;
  cursor: text;
}

.output-text:focus {
  outline: none;
  border-color: #34d399;
}

.output-text::placeholder {
  color: #475569;
  font-family: inherit;
}

.output-options {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 16px;
  border-top: 1px solid #334155;
}

.output-option {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #94a3b8;
  cursor: pointer;
}

.output-option input {
  accent-color: #3b82f6;
}

/* 通用按钮样式 */
.btn-primary,
.btn-secondary,
.btn-danger,
.btn-accent {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: all 0.15s;
  white-space: nowrap;
}

.btn-primary {
  background: #3b82f6;
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-primary:disabled {
  background: #1e3a5f;
  color: #64748b;
  cursor: not-allowed;
}

.btn-secondary {
  background: #334155;
  color: #cbd5e1;
}

.btn-secondary:hover:not(:disabled) {
  background: #475569;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-danger {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.btn-danger:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.25);
}

.btn-accent {
  background: #059669;
  color: #fff;
}

.btn-accent:hover:not(:disabled) {
  background: #047857;
}

.btn-accent:disabled {
  background: #064e3b;
  color: #64748b;
  cursor: not-allowed;
}

.btn-sm {
  padding: 5px 10px;
  font-size: 12px;
  border-radius: 6px;
}
</style>
