<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ isEdit ? '编辑原子词' : '新建原子词' }}</h3>
        <button class="close-btn" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      
      <div class="modal-body">
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label>英文原词 <span class="required">*</span></label>
            <input 
              v-model="form.value" 
              type="text" 
              placeholder="如：long hair"
              required
            />
          </div>
          
          <div class="form-group">
            <label>中文标签 <span class="required">*</span></label>
            <input 
              v-model="form.label" 
              type="text" 
              placeholder="如：长发"
              required
            />
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>类型</label>
              <select v-model="form.type">
                <option value="Positive">正向词</option>
                <option value="Negative">负向词</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>分类</label>
              <select v-model="form.category_id">
                <option :value="0">请选择分类</option>
                <optgroup v-for="parent in atomCategories" :key="parent.id" :label="parent.name">
                  <option v-for="child in getChildren(parent.id)" :key="child.id" :value="child.id">
                    {{ child.name }}
                  </option>
                </optgroup>
              </select>
            </div>
          </div>
          
          <div class="form-group">
            <label>近义词</label>
            <div class="synonyms-input">
              <input 
                v-model="synonymInput" 
                type="text" 
                placeholder="输入近义词按回车添加"
                @keyup.enter="addSynonym"
              />
              <button type="button" class="btn-add" @click="addSynonym">
                <PlusIcon class="w-4 h-4" />
              </button>
            </div>
            <div class="synonyms-list">
              <span 
                v-for="(syn, index) in form.synonyms" 
                :key="index"
                class="synonym-tag"
              >
                {{ syn }}
                <button type="button" @click="removeSynonym(index)">
                  <XMarkIcon class="w-3 h-3" />
                </button>
              </span>
            </div>
          </div>
        </form>
      </div>
      
      <div class="modal-footer">
        <button v-if="isEdit" type="button" class="btn-danger" @click="handleDelete">
          <TrashIcon class="w-4 h-4" />
          删除
        </button>
        <div class="spacer"></div>
        <button type="button" class="btn-secondary" @click="$emit('close')">取消</button>
        <button type="button" class="btn-primary" @click="handleSubmit">
          <CheckIcon class="w-4 h-4" />
          {{ isEdit ? '保存' : '创建' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import {
  XMarkIcon,
  PlusIcon,
  TrashIcon,
  CheckIcon,
} from '@heroicons/vue/24/outline'
import { useCategoryStore, useAtomStore } from '../stores'

const props = defineProps({
  atom: {
    type: Object,
    default: null,
  },
  defaultCategoryId: {
    type: Number,
    default: 0,
  },
})

// 调试：输出接收到的 props
console.log('[AtomDialog] props:', { 
  atom: props.atom, 
  defaultCategoryId: props.defaultCategoryId 
})

const emit = defineEmits(['close', 'save', 'delete'])

const categoryStore = useCategoryStore()
const atomStore = useAtomStore()

const isEdit = computed(() => !!props.atom)

const form = ref({
  value: '',
  label: '',
  type: 'Positive',
  category_id: 0,
  synonyms: [],
})

const synonymInput = ref('')

const atomCategories = computed(() => 
  categoryStore.rootCategories.filter(c => c.type === 'ATOM')
)

const getChildren = (parentId) => 
  categoryStore.getChildren(parentId)

// 监听 atom 变化（编辑模式）
watch(() => props.atom, (newAtom) => {
  console.log('[AtomDialog] atom watch triggered:', newAtom)
  if (newAtom) {
    form.value = {
      value: newAtom.value || '',
      label: newAtom.label || '',
      type: newAtom.type || 'Positive',
      category_id: newAtom.category_id || 0,
      synonyms: [...(newAtom.synonyms || [])],
    }
  }
}, { immediate: true })

// 监听 defaultCategoryId 变化（新建模式）
watch(() => props.defaultCategoryId, (newId) => {
  console.log('[AtomDialog] defaultCategoryId watch triggered:', newId, 'atom:', props.atom)
  if (!props.atom && newId > 0) {
    console.log('[AtomDialog] setting category_id to:', newId)
    form.value.category_id = newId
  }
}, { immediate: true })

// 调试：监听 form.category_id 变化
watch(() => form.value.category_id, (newVal) => {
  console.log('[AtomDialog] form.category_id changed to:', newVal)
})

function addSynonym() {
  const value = synonymInput.value.trim()
  if (value && !form.value.synonyms.includes(value)) {
    form.value.synonyms.push(value)
    synonymInput.value = ''
  }
}

function removeSynonym(index) {
  form.value.synonyms.splice(index, 1)
}

async function handleSubmit() {
  if (!form.value.value.trim() || !form.value.label.trim()) {
    alert('请填写必填项')
    return
  }
  
  emit('save', {
    ...form.value,
    id: props.atom?.id,
  })
}

async function handleDelete() {
  if (confirm('确定要删除这个原子词吗？')) {
    emit('delete', props.atom.id)
  }
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
  max-width: 480px;
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

.modal-body {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #94a3b8;
  margin-bottom: 6px;
}

.form-group .required {
  color: #ef4444;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group select:focus {
  border-color: #0ea5e9;
}

.form-group input::placeholder {
  color: #64748b;
}

.synonyms-input {
  display: flex;
  gap: 8px;
}

.synonyms-input input {
  flex: 1;
}

.btn-add {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background-color: #0284c7;
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-add:hover {
  background-color: #0ea5e9;
}

.synonyms-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.synonym-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  font-size: 12px;
  color: #e2e8f0;
}

.synonym-tag button {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: color 0.2s;
}

.synonym-tag button:hover {
  color: #ef4444;
}

.modal-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid #334155;
}

.spacer {
  flex: 1;
}

.btn-secondary,
.btn-primary,
.btn-danger {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-secondary {
  background-color: #1e293b;
  color: #e2e8f0;
  border: 1px solid #334155;
}

.btn-secondary:hover {
  background-color: #334155;
}

.btn-primary {
  background-color: #0284c7;
  color: white;
}

.btn-primary:hover {
  background-color: #0ea5e9;
}

.btn-danger {
  background-color: transparent;
  color: #ef4444;
  border: 1px solid #ef4444;
}

.btn-danger:hover {
  background-color: rgba(239, 68, 68, 0.1);
}
</style>
