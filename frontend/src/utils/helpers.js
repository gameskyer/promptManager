import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// Format relative time
export function formatRelativeTime(date) {
  if (!date) return ''
  return dayjs(date).fromNow()
}

// Format date
export function formatDate(date, format = 'YYYY-MM-DD HH:mm') {
  if (!date) return ''
  return dayjs(date).format(format)
}

// Debounce function
export function debounce(fn, delay = 300) {
  let timer = null
  return function (...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

// Throttle function
export function throttle(fn, limit = 300) {
  let inThrottle = false
  return function (...args) {
    if (!inThrottle) {
      fn.apply(this, args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

// Copy to clipboard
export async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (err) {
    console.error('Failed to copy:', err)
    return false
  }
}

// Generate unique ID
export function generateId() {
  return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

// Parse diff stats string (+2/-1)
export function parseDiffStats(stats) {
  if (!stats) return { added: 0, removed: 0 }
  const added = stats.match(/\+(\d+)/)?.[1] || 0
  const removed = stats.match(/-(\d+)/)?.[1] || 0
  return { added: parseInt(added), removed: parseInt(removed) }
}

// Build prompt text from atoms
export function buildPrompt(atoms) {
  if (!Array.isArray(atoms)) return ''
  return atoms.map(a => a.value).join(', ')
}

// Parse prompt text to extract atoms
export function parsePrompt(prompt) {
  if (!prompt) return []
  return prompt.split(/[,，;；]/).map(s => s.trim()).filter(Boolean)
}

// Validate atom value
export function validateAtomValue(value) {
  if (!value || value.trim().length === 0) {
    return { valid: false, error: '原子词不能为空' }
  }
  if (value.length > 200) {
    return { valid: false, error: '原子词长度不能超过200字符' }
  }
  // Check for invalid characters
  if (/[<>\"'&]/.test(value)) {
    return { valid: false, error: '原子词包含非法字符' }
  }
  return { valid: true }
}

// Sanitize string for storage
export function sanitizeString(str) {
  if (!str) return ''
  return str.replace(/[<>\"'&]/g, '').trim()
}

// Get pinyin initials (simplified)
export function getPinyinInitials(str) {
  // This is a placeholder - in production use a proper pinyin library
  // For now, just return the lowercase string
  return str.toLowerCase()
}

// Group array by key
export function groupBy(array, key) {
  return array.reduce((result, item) => {
    const groupKey = item[key]
    if (!result[groupKey]) {
      result[groupKey] = []
    }
    result[groupKey].push(item)
    return result
  }, {})
}

// Sort array by multiple criteria
export function sortBy(array, ...criteria) {
  return [...array].sort((a, b) => {
    for (const criterion of criteria) {
      const { key, order = 'asc' } = criterion
      const aVal = a[key]
      const bVal = b[key]
      
      if (aVal < bVal) return order === 'asc' ? -1 : 1
      if (aVal > bVal) return order === 'asc' ? 1 : -1
    }
    return 0
  })
}

// Deep clone
export function deepClone(obj) {
  return JSON.parse(JSON.stringify(obj))
}

// Is object empty
export function isEmpty(obj) {
  if (!obj) return true
  if (Array.isArray(obj)) return obj.length === 0
  if (typeof obj === 'object') return Object.keys(obj).length === 0
  return false
}

// File size formatter
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Truncate text
export function truncate(text, length = 50, suffix = '...') {
  if (!text) return ''
  if (text.length <= length) return text
  return text.substring(0, length) + suffix
}
