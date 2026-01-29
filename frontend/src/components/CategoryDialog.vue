<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ isEdit ? '编辑分类' : '新建分类' }}</h3>
        <button class="close-btn" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      
      <div class="modal-body">
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label>分类名称 <span class="required">*</span></label>
            <input 
              v-model="form.name" 
              type="text" 
              placeholder="如：发型"
              required
            />
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>类型</label>
              <select v-model="form.type">
                <option value="ATOM">提示词分类</option>
                <option value="PRESET">预设分类</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>父分类</label>
              <select v-model="form.parent_id">
                <option :value="0">无（作为一级分类）</option>
                <option 
                  v-for="cat in availableParents" 
                  :key="cat.id" 
                  :value="cat.id"
                >
                  {{ cat.name }}
                </option>
              </select>
            </div>
          </div>
          
          <div class="form-group">
            <label>排序权重</label>
            <input 
              v-model.number="form.sort_order" 
              type="number" 
              placeholder="数字越小越靠前"
            />
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
  TrashIcon,
  CheckIcon,
} from '@heroicons/vue/24/outline'
import { useCategoryStore } from '../stores'

const props = defineProps({
  category: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['close', 'save', 'delete'])

const categoryStore = useCategoryStore()

const isEdit = computed(() => !!props.category)

const form = ref({
  name: '',
  type: 'ATOM',
  parent_id: 0,
  sort_order: 0,
})

// 获取可作为父分类的选项（排除当前分类及其子分类）
const availableParents = computed(() => {
  const all = categoryStore.categories
  if (!props.category) return all.filter(c => c.parent_id === 0)
  
  // 排除当前分类和它的子分类
  const excludeIds = new Set([props.category.id])
  const findChildren = (parentId) => {
    all.filter(c => c.parent_id === parentId).forEach(c => {
      excludeIds.add(c.id)
      findChildren(c.id)
    })
  }
  findChildren(props.category.id)
  
  return all.filter(c => !excludeIds.has(c.id))
})

watch(() => props.category, (newCat) => {
  if (newCat) {
    form.value = {
      name: newCat.name || '',
      type: newCat.type || 'ATOM',
      parent_id: newCat.parent_id || 0,
      sort_order: newCat.sort_order || 0,
    }
  } else {
    form.value = {
      name: '',
      type: 'ATOM',
      parent_id: 0,
      sort_order: 0,
    }
  }
}, { immediate: true })

async function handleSubmit() {
  if (!form.value.name.trim()) {
    alert('请填写分类名称')
    return
  }
  
  emit('save', {
    ...form.value,
    id: props.category?.id,
  })
}

async function handleDelete() {
  if (confirm('确定要删除这个分类吗？分类下的所有提示词将被一并删除。')) {
    emit('delete', props.category.id)
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
  max-width: 440px;
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
