import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // State
  const currentCategory = ref(null)
  const currentSubCategory = ref(null)
  const selectedAtoms = ref([])
  const currentPreset = ref(null)
  const currentVersion = ref(null)
  const showTimeline = ref(false)
  const searchQuery = ref('')
  const isLoading = ref(false)
  const activeTab = ref('atoms') // atoms, presets, settings
  const theme = ref('dark') // dark, light, auto
  


  // Getters
  const selectedAtomIDs = computed(() => selectedAtoms.value.map(a => a.id))
  const selectedCount = computed(() => selectedAtoms.value.length)
  const hasSelection = computed(() => selectedAtoms.value.length > 0)

  // Actions
  function setCategory(category) {
    currentCategory.value = category
    currentSubCategory.value = null
  }

  function setSubCategory(subCategory) {
    currentSubCategory.value = subCategory
  }

  function toggleAtom(atom) {
    const index = selectedAtoms.value.findIndex(a => a.id === atom.id)
    if (index === -1) {
      selectedAtoms.value.push(atom)
    } else {
      selectedAtoms.value.splice(index, 1)
    }
  }

  function addAtom(atom) {
    if (!selectedAtoms.value.find(a => a.id === atom.id)) {
      selectedAtoms.value.push(atom)
    }
  }

  function removeAtom(atomId) {
    const index = selectedAtoms.value.findIndex(a => a.id === atomId)
    if (index !== -1) {
      selectedAtoms.value.splice(index, 1)
    }
  }

  function moveAtom(atomId, direction) {
    const index = selectedAtoms.value.findIndex(a => a.id === atomId)
    if (index === -1) return

    if (direction === 'up' && index > 0) {
      const temp = selectedAtoms.value[index]
      selectedAtoms.value[index] = selectedAtoms.value[index - 1]
      selectedAtoms.value[index - 1] = temp
    } else if (direction === 'down' && index < selectedAtoms.value.length - 1) {
      const temp = selectedAtoms.value[index]
      selectedAtoms.value[index] = selectedAtoms.value[index + 1]
      selectedAtoms.value[index + 1] = temp
    }
  }

  function clearSelection() {
    selectedAtoms.value = []
  }

  function setCurrentPreset(preset) {
    currentPreset.value = preset
    currentVersion.value = preset?.current_version || null
  }

  function setCurrentVersion(version) {
    currentVersion.value = version
  }

  function toggleTimeline() {
    showTimeline.value = !showTimeline.value
  }

  function setSearchQuery(query) {
    searchQuery.value = query
  }

  function setLoading(loading) {
    isLoading.value = loading
  }

  function setActiveTab(tab) {
    activeTab.value = tab
  }

  function setTheme(newTheme) {
    theme.value = newTheme
    applyTheme(newTheme)
  }

  function applyTheme(themeValue) {
    const root = document.documentElement
    if (themeValue === 'dark') {
      root.classList.remove('light')
      root.classList.add('dark')
    } else if (themeValue === 'light') {
      root.classList.remove('dark')
      root.classList.add('light')
    } else {
      // auto - 跟随系统
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      if (prefersDark) {
        root.classList.remove('light')
        root.classList.add('dark')
      } else {
        root.classList.remove('dark')
        root.classList.add('light')
      }
    }
  }

  function buildPromptText() {
    return selectedAtoms.value.map(a => a.value).join(', ')
  }

  return {
    // State
    currentCategory,
    currentSubCategory,
    selectedAtoms,
    currentPreset,
    currentVersion,
    showTimeline,
    searchQuery,
    isLoading,
    activeTab,
    theme,

    
    // Getters
    selectedAtomIDs,
    selectedCount,
    hasSelection,
    
    // Actions
    setCategory,
    setSubCategory,
    toggleAtom,
    addAtom,
    removeAtom,
    moveAtom,
    clearSelection,
    setCurrentPreset,
    setCurrentVersion,
    toggleTimeline,
    setSearchQuery,
    setLoading,
    setActiveTab,
    setTheme,
    applyTheme,

    buildPromptText,
  }
})
