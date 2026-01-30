<template>
  <div class="tree-node" :style="{ paddingLeft: level * 20 + 'px' }">
    <div 
      class="node-content"
      :class="{ 
        'is-expanded': isExpanded,
        'is-dragging': isDragging,
        'is-drop-target': isDropTarget,
        'is-root': level === 0
      }"
      draggable="true"
      @dragstart="handleDragStart"
      @dragend="handleDragEnd"
      @dragover.prevent="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <!-- Expand/Collapse Button -->
      <button 
        v-if="hasChildren"
        class="expand-btn"
        @click.stop="toggleExpand"
      >
        <ChevronDownIcon v-if="isExpanded" class="w-4 h-4" />
        <ChevronRightIcon v-else class="w-4 h-4" />
      </button>
      <span v-else class="expand-placeholder"></span>

      <!-- Icon -->
      <FolderIcon class="node-icon" :class="{ 'is-open': isExpanded }" />

      <!-- Name -->
      <span class="node-name">{{ node.name }}</span>

      <!-- Type Badge -->
      <span class="type-badge" :class="node.type?.toLowerCase()">
        {{ node.type === 'ATOM' ? '原子词' : '预设' }}
      </span>

      <!-- Count -->
      <span v-if="node.children?.length > 0" class="child-count">
        {{ node.children.length }} 个子分类
      </span>

      <!-- Actions -->
      <div class="node-actions">
        <button 
          class="action-btn" 
          @click.stop="$emit('add-child', node.id)"
          title="添加子分类"
        >
          <PlusIcon class="w-4 h-4" />
        </button>
        <button 
          class="action-btn" 
          @click.stop="$emit('edit', node)"
          title="编辑"
        >
          <PencilIcon class="w-4 h-4" />
        </button>
        <button 
          class="action-btn danger" 
          @click.stop="$emit('delete', node)"
          title="删除"
        >
          <TrashIcon class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Children -->
    <div v-if="hasChildren && isExpanded" class="node-children">
      <CategoryTreeNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :expanded="expanded"
        :level="level + 1"
        @toggle="$emit('toggle', $event)"
        @edit="$emit('edit', $event)"
        @delete="$emit('delete', $event)"
        @add-child="$emit('add-child', $event)"
        @move="$emit('move', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import {
  ChevronDownIcon,
  ChevronRightIcon,
  FolderIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  node: {
    type: Object,
    required: true,
  },
  expanded: {
    type: Array,
    default: () => [],
  },
  level: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['toggle', 'edit', 'delete', 'add-child', 'move'])

const isDragging = ref(false)
const isDropTarget = ref(false)
const dragCounter = ref(0)

const hasChildren = computed(() => {
  return props.node.children && props.node.children.length > 0
})

const isExpanded = computed(() => {
  return props.expanded.includes(props.node.id)
})

function toggleExpand() {
  emit('toggle', props.node.id)
}

// Drag and Drop
function handleDragStart(e) {
  isDragging.value = true
  e.dataTransfer.setData('categoryId', props.node.id.toString())
  e.dataTransfer.setData('categoryName', props.node.name)
  e.dataTransfer.effectAllowed = 'move'
}

function handleDragEnd() {
  isDragging.value = false
  dragCounter.value = 0
  isDropTarget.value = false
}

function handleDragOver(e) {
  e.preventDefault()
  e.dataTransfer.dropEffect = 'move'
}

function handleDragEnter() {
  dragCounter.value++
  if (dragCounter.value > 0) {
    isDropTarget.value = true
  }
}

function handleDragLeave() {
  dragCounter.value--
  if (dragCounter.value <= 0) {
    isDropTarget.value = false
    dragCounter.value = 0
  }
}

function handleDrop(e) {
  e.preventDefault()
  e.stopPropagation()
  
  const draggedId = parseInt(e.dataTransfer.getData('categoryId'))
  const targetId = props.node.id
  
  isDropTarget.value = false
  dragCounter.value = 0
  
  // Prevent dropping on itself or its descendants
  if (draggedId === targetId) return
  if (isDescendant(draggedId, targetId)) return
  
  emit('move', { categoryId: draggedId, newParentId: targetId })
}

function isDescendant(parentId, childId) {
  const findNode = (nodes, id) => {
    for (const node of nodes) {
      if (node.id === id) return node
      if (node.children) {
        const found = findNode(node.children, id)
        if (found) return found
      }
    }
    return null
  }
  
  const parent = findNode([props.node], parentId)
  if (!parent || !parent.children) return false
  
  const checkDescendant = (nodes, targetId) => {
    for (const node of nodes) {
      if (node.id === targetId) return true
      if (node.children && checkDescendant(node.children, targetId)) return true
    }
    return false
  }
  
  return checkDescendant(parent.children, childId)
}
</script>

<style scoped>
.tree-node {
  user-select: none;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  margin: 2px 0;
  border-radius: 8px;
  background-color: #1e293b;
  border: 1px solid transparent;
  cursor: grab;
  transition: all 0.2s;
}

.node-content:hover {
  background-color: #334155;
  border-color: #475569;
}

.node-content.is-root {
  background-color: #0f172a;
  border-color: #334155;
}

.node-content.is-root:hover {
  background-color: #1e293b;
}

.node-content.is-dragging {
  opacity: 0.5;
}

.node-content.is-drop-target {
  border-color: #0ea5e9;
  background-color: rgba(14, 165, 233, 0.1);
}

.expand-btn {
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
  transition: all 0.2s;
}

.expand-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.expand-placeholder {
  width: 20px;
}

.node-icon {
  width: 18px;
  height: 18px;
  color: #f59e0b;
  flex-shrink: 0;
}

.node-icon.is-open {
  color: #fbbf24;
}

.node-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
}

.type-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.type-badge.atom {
  background-color: rgba(14, 165, 233, 0.15);
  color: #0ea5e9;
}

.type-badge.preset {
  background-color: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.child-count {
  font-size: 11px;
  color: #64748b;
  padding: 2px 8px;
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.node-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.node-content:hover .node-actions {
  opacity: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  background-color: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #475569;
  color: #e2e8f0;
}

.action-btn.danger:hover {
  background-color: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.node-children {
  margin-top: 2px;
}
</style>
