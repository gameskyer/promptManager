<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ isEdit ? '编辑预设' : '新建预设' }}</h3>
        <button class="close-btn" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      
      <div class="modal-body">
        <form @submit.prevent="handleSubmit">
          <!-- 预览图区域 -->
          <div class="form-group preview-section">
            <label>预览图</label>
            <div class="preview-gallery">
              <div 
                v-for="(preview, index) in form.previews" 
                :key="index"
                class="preview-item"
              >
                <img 
                  :src="preview" 
                  :alt="`预览 ${index + 1}`"
                  @click="openImageViewer(index)"
                />
                <div class="preview-actions">
                  <button 
                    type="button"
                    class="preview-set-cover"
                    :class="{ active: form.thumbnail === preview }"
                    @click="setAsCover(index)"
                    title="设为封面"
                  >
                    <StarIcon class="w-3 h-3" />
                  </button>
                  <button 
                    type="button"
                    class="preview-delete"
                    @click="removePreview(index)"
                    title="删除"
                  >
                    <XMarkIcon class="w-3 h-3" />
                  </button>
                </div>
                <div v-if="form.thumbnail === preview" class="cover-badge">封面</div>
              </div>
              
              <button type="button" class="preview-add" @click="addPreview">
                <PlusIcon class="w-8 h-8" />
                <span>添加图片</span>
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label>预设标题 <span class="required">*</span></label>
            <input 
              v-model="form.title" 
              type="text" 
              placeholder="如：动漫女孩基础预设"
              required
            />
          </div>
          
          <div class="form-group">
            <label>正向提示词</label>
            <textarea 
              v-model="form.pos_text" 
              rows="3"
              placeholder="正向提示词，多个词用逗号分隔"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>负向提示词</label>
            <textarea 
              v-model="form.neg_text" 
              rows="2"
              placeholder="负向提示词，多个词用逗号分隔"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>模型 (Model)</label>
            <input 
              v-model="form.params.model" 
              type="text" 
              placeholder="如：animeModel_v20.safetensors"
            />
          </div>
          
          <!-- LoRA 列表 -->
          <div class="form-group lora-section">
            <label>LoRA</label>
            <div class="lora-list">
              <div 
                v-for="(lora, index) in form.loras" 
                :key="index"
                class="lora-row"
              >
                <input 
                  v-model="lora.name" 
                  type="text" 
                  placeholder="LoRA名称"
                  class="lora-name"
                />
                <input 
                  v-model.number="lora.weight" 
                  type="number" 
                  step="0.1"
                  min="0"
                  max="2"
                  placeholder="强度"
                  class="lora-weight"
                />
                <button 
                  type="button" 
                  class="lora-remove-btn"
                  @click="removeLora(index)"
                >
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
            <button type="button" class="lora-add-btn" @click="addLora">
              <PlusIcon class="w-4 h-4" />
              添加 LoRA
            </button>
          </div>
          
          <div class="form-section-title">生成参数</div>
          
          <div class="form-row">
            <div class="form-group">
              <label>采样步数 (Steps)</label>
              <input 
                v-model.number="form.params.steps" 
                type="number" 
                min="1"
                max="150"
              />
            </div>
            
            <div class="form-group">
              <label>CFG Scale</label>
              <input 
                v-model.number="form.params.cfg" 
                type="number" 
                min="1"
                max="30"
                step="0.5"
              />
            </div>
          </div>
          
          <div class="form-group">
            <label>采样器 (Sampler)</label>
            <select v-model="form.params.sampler">
              <option value="Euler">Euler</option>
              <option value="Euler a">Euler a</option>
              <option value="DPM++ 2M">DPM++ 2M</option>
              <option value="DPM++ 2M Karras">DPM++ 2M Karras</option>
              <option value="DDIM">DDIM</option>
            </select>
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
    
    <!-- 图片查看器 -->
    <ImageViewer
      v-if="viewerVisible"
      :images="form.previews"
      :initial-index="viewerIndex"
      @close="viewerVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import {
  XMarkIcon,
  PlusIcon,
  TrashIcon,
  CheckIcon,
  StarIcon,
} from '@heroicons/vue/24/outline'
import ImageViewer from './ImageViewer.vue'


const props = defineProps({
  preset: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['close', 'save', 'delete'])

const isEdit = computed(() => !!props.preset)

const form = ref({
  title: '',
  pos_text: '',
  neg_text: '',
  params: {
    steps: 30,
    cfg: 7,
    sampler: 'DPM++ 2M Karras',
    model: '',
  },
  loras: [],
  previews: [],
  thumbnail: '',
})

const viewerVisible = ref(false)
const viewerIndex = ref(0)

watch(() => props.preset, (newPreset) => {
  if (newPreset) {
    form.value = {
      title: newPreset.title || '',
      pos_text: newPreset.pos_text || '',
      neg_text: newPreset.neg_text || '',
      params: {
        steps: newPreset.params?.steps ?? 30,
        cfg: newPreset.params?.cfg ?? 7,
        sampler: newPreset.params?.sampler || 'DPM++ 2M Karras',
        model: newPreset.params?.model || '',
      },
      loras: newPreset.loras?.length > 0 
        ? [...newPreset.loras] 
        : [],
      previews: newPreset.previews?.length > 0
        ? [...newPreset.previews]
        : [],
      thumbnail: newPreset.thumbnail || '',
    }
  } else {
    form.value = {
      title: '',
      pos_text: '',
      neg_text: '',
      params: {
        steps: 30,
        cfg: 7,
        sampler: 'DPM++ 2M Karras',
        model: '',
      },
      loras: [],
      previews: [],
      thumbnail: '',
    }
  }
}, { immediate: true })

// 预览图操作 - 使用 base64 存储，保存时提交到后端
function addPreview() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.multiple = true
  input.onchange = (e) => {
    const files = Array.from(e.target.files)
    files.forEach(file => {
      const reader = new FileReader()
      reader.onload = (event) => {
        form.value.previews.push(event.target.result)
        // 如果没有封面，自动设置第一张为封面
        if (!form.value.thumbnail) {
          form.value.thumbnail = event.target.result
        }
      }
      reader.readAsDataURL(file)
    })
  }
  input.click()
}

function removePreview(index) {
  const removed = form.value.previews[index]
  form.value.previews.splice(index, 1)
  // 如果删除的是封面，重新设置封面
  if (form.value.thumbnail === removed && form.value.previews.length > 0) {
    form.value.thumbnail = form.value.previews[0]
  } else if (form.value.previews.length === 0) {
    form.value.thumbnail = ''
  }
}

function setAsCover(index) {
  form.value.thumbnail = form.value.previews[index]
}

function openImageViewer(index) {
  viewerIndex.value = index
  viewerVisible.value = true
}

// LoRA 操作
function addLora() {
  form.value.loras.push({
    name: '',
    weight: 1.0,
  })
}

function removeLora(index) {
  form.value.loras.splice(index, 1)
}

async function handleSubmit() {
  if (!form.value.title.trim()) {
    alert('请填写预设标题')
    return
  }
  
  // 过滤掉空的 LoRA
  const validLoras = form.value.loras.filter(l => l.name.trim() !== '')
  
  // 提取 base64 图片数据（去掉 data:image/xxx;base64, 前缀）
  const previewBase64s = form.value.previews
    .filter(p => p && p.startsWith('data:'))
    .map(p => p.split(',')[1])
  
  emit('save', {
    ...form.value,
    loras: validLoras,
    previews: previewBase64s,
    thumbnail: '', // 后端会设置第一张为封面
    id: props.preset?.id,
  })
}

async function handleDelete() {
  if (confirm('确定要删除这个预设吗？所有版本历史将被一并删除。')) {
    emit('delete', props.preset.id)
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
  max-width: 600px;
  max-height: 90vh;
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

.form-section-title {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 20px 0 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #334155;
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
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
  resize: vertical;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: #0ea5e9;
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: #64748b;
}

/* 预览图区域 */
.preview-section {
  background-color: rgba(30, 41, 59, 0.5);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #334155;
}

.preview-gallery {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.preview-item {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid transparent;
  cursor: pointer;
}

.preview-item:hover {
  border-color: #0ea5e9;
}

.preview-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.preview-item:hover .preview-actions {
  opacity: 1;
}

.preview-set-cover,
.preview-delete {
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

.preview-set-cover {
  background-color: rgba(0, 0, 0, 0.6);
  color: #94a3b8;
}

.preview-set-cover:hover,
.preview-set-cover.active {
  background-color: #f59e0b;
  color: white;
}

.preview-delete {
  background-color: rgba(239, 68, 68, 0.8);
  color: white;
}

.preview-delete:hover {
  background-color: #ef4444;
}

.cover-badge {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 2px 4px;
  background-color: rgba(245, 158, 11, 0.9);
  color: white;
  font-size: 10px;
  font-weight: 600;
  text-align: center;
}

.preview-add {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 100px;
  background-color: transparent;
  border: 2px dashed #475569;
  border-radius: 8px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
  gap: 4px;
}

.preview-add:hover {
  border-color: #0ea5e9;
  color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.05);
}

.preview-add span {
  font-size: 12px;
}

/* LoRA 样式 */
.lora-section {
  background-color: rgba(30, 41, 59, 0.5);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #334155;
}

.lora-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 10px;
}

.lora-row {
  display: grid;
  grid-template-columns: 1fr 80px 32px;
  gap: 8px;
  align-items: center;
}

.lora-row input {
  margin-bottom: 0;
}

.lora-name {
  min-width: 0;
}

.lora-weight {
  text-align: center;
}

.lora-remove-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 36px;
  background-color: transparent;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.lora-remove-btn:hover {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
  color: #ef4444;
}

.lora-add-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 8px;
  background-color: transparent;
  border: 1px dashed #475569;
  border-radius: 6px;
  color: #94a3b8;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.lora-add-btn:hover {
  border-color: #0ea5e9;
  color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.05);
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
