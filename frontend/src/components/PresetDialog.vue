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
          
          <!-- 导入 ComfyUI 信息 -->
          <div class="form-group import-section">
            <label>快速导入</label>
            <button type="button" class="btn-import" @click="importComfyUIFile">
              <DocumentArrowUpIcon class="w-4 h-4" />
              导入 ComfyUI 图片信息
            </button>
            <span class="import-hint">选择 ComfyUI 生成的 txt 文件，自动填充提示词和参数</span>
          </div>
          
          <div class="form-row">
            <div class="form-group" style="flex: 1;">
              <label>预设标题 <span class="required">*</span></label>
              <input 
                v-model="form.title" 
                type="text" 
                placeholder="如：动漫女孩基础预设"
                required
              />
            </div>
            
            <div class="form-group" style="width: 160px;">
              <label>分类</label>
              <select v-model.number="form.category_id">
                <option :value="0">未分类</option>
                <option 
                  v-for="cat in presetCategories" 
                  :key="cat.id" 
                  :value="cat.id"
                >
                  {{ cat.name }}
                </option>
              </select>
            </div>
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
import { ref, computed, watch, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import {
  XMarkIcon,
  PlusIcon,
  TrashIcon,
  CheckIcon,
  StarIcon,
  DocumentArrowUpIcon,
} from '@heroicons/vue/24/outline'
import ImageViewer from './ImageViewer.vue'
import { useCategoryStore } from '../stores'
import { parseComfyUIFile, readFileContent } from '../utils/comfyuiParser'


const props = defineProps({
  preset: {
    type: Object,
    default: null,
  },
  defaultCategoryId: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['close', 'save', 'delete'])

const categoryStore = useCategoryStore()
const { categories } = storeToRefs(categoryStore)

const isEdit = computed(() => !!props.preset)

const presetCategories = computed(() =>
  categories.value.filter(c => c.type === 'PRESET')
)

const form = ref({
  title: '',
  category_id: 0,
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
  console.log('[PresetDialog] preset loaded:', {
    id: newPreset?.id,
    title: newPreset?.title,
    thumbnail: newPreset?.thumbnail,
    previews: newPreset?.previews,
    previewsLength: newPreset?.previews?.length,
  })
  if (newPreset) {
    form.value = {
      title: newPreset.title || '',
      category_id: newPreset.category_id || 0,
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
    console.log('[PresetDialog] form initialized:', {
      thumbnail: form.value.thumbnail,
      previews: form.value.previews,
    })
  } else {
    form.value = {
      title: '',
      category_id: props.defaultCategoryId,
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

onMounted(async () => {
  if (categoryStore.categories.length === 0) {
    await categoryStore.fetchCategories()
  }
})

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

// 导入 ComfyUI 图片信息文件
async function importComfyUIFile() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.txt'
  
  input.onchange = async (e) => {
    const file = e.target.files[0]
    if (!file) return
    
    try {
      console.log('[PresetDialog] Importing ComfyUI file:', file.name)
      const content = await readFileContent(file)
      const parsed = parseComfyUIFile(content)
      
      console.log('[PresetDialog] Parsed ComfyUI data:', {
        posText: parsed.posText?.substring(0, 50) + '...',
        negText: parsed.negText?.substring(0, 50) + '...',
        model: parsed.model,
        params: parsed.params,
      })
      
      // 填充表单数据
      if (parsed.posText) {
        form.value.pos_text = parsed.posText
      }
      if (parsed.negText) {
        form.value.neg_text = parsed.negText
      }
      if (parsed.model) {
        form.value.params.model = parsed.model
      }
      
      // 更新参数
      if (parsed.params) {
        form.value.params.steps = parsed.params.steps || 30
        form.value.params.cfg = parsed.params.cfg || 7
        form.value.params.sampler = parsed.params.sampler || 'DPM++ 2M Karras'
      }
      
      // 如果没有标题，尝试从文件名生成
      if (!form.value.title.trim()) {
        const fileName = file.name.replace(/\.txt$/i, '').replace(/^ComfyUI_/, '')
        form.value.title = fileName
      }
      
      alert('导入成功！已自动填充提示词和参数')
    } catch (err) {
      console.error('[PresetDialog] Import failed:', err)
      alert('导入失败：' + err.message)
    }
  }
  
  input.click()
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
  
  // 处理图片：保持原始格式发送，让后端判断
  // - 新上传的: data:image/xxx;base64,xxxxx
  // - 已有的: /images/xxx.png
  const previewData = form.value.previews
    .filter(p => p && p.trim() !== '')
    .map(p => p) // 保持原样，不做转换
  
  // 处理封面：保持原始格式
  let thumbnailData = form.value.thumbnail || ''
  
  const submitData = {
    ...form.value,
    loras: validLoras,
    previews: previewData,
    thumbnail: thumbnailData,
    id: props.preset?.id,
  }
  
  console.log('[PresetDialog] handleSubmit:', {
    id: submitData.id,
    thumbnail: submitData.thumbnail?.substring(0, 50),
    previews: submitData.previews?.map(p => p?.substring(0, 50)),
    previewsLength: submitData.previews?.length,
  })
  
  emit('save', submitData)
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

/* 导入区域样式 */
.import-section {
  background-color: rgba(14, 165, 233, 0.05);
  border: 1px dashed #0ea5e9;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 16px;
}

.import-section label {
  color: #0ea5e9;
  font-weight: 500;
}

.btn-import {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 10px 16px;
  background-color: rgba(14, 165, 233, 0.1);
  border: 1px solid #0ea5e9;
  border-radius: 8px;
  color: #0ea5e9;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 8px;
}

.btn-import:hover {
  background-color: #0ea5e9;
  color: white;
}

.import-hint {
  display: block;
  font-size: 12px;
  color: #64748b;
  margin-top: 8px;
  text-align: center;
}
</style>
