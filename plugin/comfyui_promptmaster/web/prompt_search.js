/**
 * PromptMaster Search Extension for ComfyUI
 * 
 * This extension adds intelligent search functionality to ComfyUI's string widgets.
 * When typing in Chinese or English, it shows suggestions from the PromptMaster database.
 */

import { app } from "../../scripts/app.js";
import { api } from "../../scripts/api.js";

// Global state
let searchCache = new Map();
let lastQuery = "";
let debounceTimer = null;

/**
 * Create search suggestion popup
 */
function createSuggestionPopup(x, y, suggestions, onSelect, inputWidget) {
    // Remove existing popup
    removeExistingPopup();
    
    const popup = document.createElement("div");
    popup.id = "promptmaster-search-popup";
    popup.style.cssText = `
        position: fixed;
        left: ${x}px;
        top: ${y + 5}px;
        background: #1e1e1e;
        border: 2px solid #4a9eff;
        border-radius: 8px;
        padding: 8px 0;
        min-width: 320px;
        max-width: 450px;
        max-height: 350px;
        overflow-y: auto;
        z-index: 10000;
        box-shadow: 0 4px 12px rgba(0,0,0,0.5);
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    `;
    
    if (suggestions.length === 0) {
        const noResult = document.createElement("div");
        noResult.textContent = "无搜索结果";
        noResult.style.cssText = `
            padding: 12px 16px;
            color: #888;
            text-align: center;
            font-size: 13px;
        `;
        popup.appendChild(noResult);
    } else {
        // Header with result count
        const header = document.createElement("div");
        header.textContent = `找到 ${suggestions.length} 个结果`;
        header.style.cssText = `
            padding: 6px 16px;
            color: #4a9eff;
            font-size: 11px;
            border-bottom: 1px solid #333;
            margin-bottom: 4px;
        `;
        popup.appendChild(header);
        
        // Suggestion items
        suggestions.forEach((item, index) => {
            const div = document.createElement("div");
            div.className = "suggestion-item";
            
            // Type badge color
            const typeColor = item.type === "Negative" ? "#ff6b6b" : "#51cf66";
            const typeLabel = item.type === "Negative" ? "负" : "正";
            
            div.innerHTML = `
                <div style="display: flex; align-items: center; gap: 8px;">
                    <span style="
                        background: ${typeColor};
                        color: white;
                        font-size: 10px;
                        padding: 2px 6px;
                        border-radius: 4px;
                        min-width: 18px;
                        text-align: center;
                    ">${typeLabel}</span>
                    <div style="flex: 1; min-width: 0;">
                        <div style="font-size: 14px; color: #fff; font-weight: 500;">
                            ${highlightMatch(item.label || item.value, lastQuery)}
                        </div>
                        <div style="font-size: 12px; color: #888;">
                            ${item.value}
                            ${item.category_name ? `<span style="color: #4a9eff;"> · ${item.category_name}</span>` : ''}
                        </div>
                    </div>
                </div>
            `;
            
            div.style.cssText = `
                padding: 10px 16px;
                cursor: pointer;
                transition: background 0.15s;
                border-left: 3px solid transparent;
            `;
            
            div.addEventListener("mouseenter", () => {
                div.style.background = "#2a2a2a";
                div.style.borderLeftColor = "#4a9eff";
            });
            
            div.addEventListener("mouseleave", () => {
                div.style.background = "transparent";
                div.style.borderLeftColor = "transparent";
            });
            
            div.addEventListener("click", (e) => {
                e.stopPropagation();
                onSelect(item);
                removeExistingPopup();
            });
            
            popup.appendChild(div);
        });
    }
    
    document.body.appendChild(popup);
    
    // Close popup when clicking outside
    setTimeout(() => {
        document.addEventListener("click", closePopupOnClickOutside);
    }, 10);
    
    return popup;
}

/**
 * Highlight matching text
 */
function highlightMatch(text, query) {
    if (!query || !text) return text || "";
    try {
        const regex = new RegExp(`(${escapeRegExp(query)})`, "gi");
        return text.replace(regex, '<span style="color: #ffd43b; font-weight: bold;">$1</span>');
    } catch (e) {
        return text;
    }
}

function escapeRegExp(string) {
    if (!string) return "";
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Remove existing popup
 */
function removeExistingPopup() {
    const existing = document.getElementById("promptmaster-search-popup");
    if (existing) {
        existing.remove();
    }
    document.removeEventListener("click", closePopupOnClickOutside);
}

/**
 * Close popup when clicking outside
 */
function closePopupOnClickOutside(e) {
    const popup = document.getElementById("promptmaster-search-popup");
    if (popup && !popup.contains(e.target)) {
        removeExistingPopup();
    }
}

/**
 * Search atoms via API
 */
async function searchAtoms(query, limit = 10) {
    if (!query || query.length < 1) {
        return [];
    }
    
    // Check cache
    const cacheKey = `${query}_${limit}`;
    if (searchCache.has(cacheKey)) {
        return searchCache.get(cacheKey);
    }
    
    try {
        console.log("[PromptMaster] Searching for:", query);
        
        // Call the API endpoint
        const response = await api.fetchApi("/promptmaster/search", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, limit })
        });
        
        if (!response.ok) {
            console.error("[PromptMaster] Search API error:", response.status, response.statusText);
            return [];
        }
        
        const data = await response.json();
        console.log("[PromptMaster] Search results:", data);
        
        if (data.success) {
            // Cache results
            searchCache.set(cacheKey, data.results);
            // Limit cache size
            if (searchCache.size > 50) {
                const firstKey = searchCache.keys().next().value;
                searchCache.delete(firstKey);
            }
            return data.results;
        }
    } catch (e) {
        console.error("[PromptMaster] Search failed:", e);
    }
    
    return [];
}

/**
 * Get widget's input element
 */
function getWidgetInput(widget) {
    // Try different properties that might contain the input element
    if (widget.inputEl) return widget.inputEl;
    if (widget.element) return widget.element;
    if (widget.canvasEl) return widget.canvasEl;
    return null;
}

/**
 * Check if this is a searchable widget
 */
function isSearchableWidget(widget, node) {
    // Skip non-string widgets
    if (widget.type !== "STRING") return false;
    
    // Skip certain widget names
    const skipNames = ["selected_prompt", "selected_prompts"];
    if (skipNames.includes(widget.name)) return false;
    
    // Check if node is a PromptMaster search node
    const isPromptMasterNode = node.comfyClass === "PromptSearch" || 
                               node.comfyClass === "PromptSearchAdvanced";
    
    // Check if node is CLIPTextEncode
    const isClipNode = node.comfyClass === "CLIPTextEncode";
    
    return isPromptMasterNode || isClipNode;
}

/**
 * Enhance text input widget with search functionality
 */
function enhanceTextWidget(node, widget) {
    const inputEl = getWidgetInput(widget);
    
    if (!inputEl) {
        console.log("[PromptMaster] No input element found for widget:", widget.name);
        return;
    }
    
    // Check if already enhanced
    if (inputEl._promptmasterEnhanced) {
        return;
    }
    
    // Check if this is a searchable widget
    if (!isSearchableWidget(widget, node)) {
        return;
    }
    
    console.log("[PromptMaster] Enhancing widget:", widget.name, "in node:", node.comfyClass);
    
    inputEl._promptmasterEnhanced = true;
    
    // Add search icon if input has a parent
    const parent = inputEl.parentElement;
    if (parent && !inputEl._searchIconAdded) {
        inputEl._searchIconAdded = true;
        
        // Check if parent is already a wrapper
        const wrapper = document.createElement("div");
        wrapper.style.cssText = "position: relative; display: inline-block; width: 100%;";
        
        // Create search button
        const searchBtn = document.createElement("span");
        searchBtn.innerHTML = "🔍";
        searchBtn.style.cssText = `
            position: absolute;
            right: 8px;
            top: 50%;
            transform: translateY(-50%);
            cursor: pointer;
            font-size: 14px;
            opacity: 0.5;
            transition: opacity 0.2s;
            z-index: 10;
            pointer-events: auto;
        `;
        searchBtn.title = "点击搜索提示词 (Ctrl+Space)";
        
        searchBtn.addEventListener("mouseenter", () => searchBtn.style.opacity = "1");
        searchBtn.addEventListener("mouseleave", () => searchBtn.style.opacity = "0.5");
        searchBtn.addEventListener("click", (e) => {
            e.stopPropagation();
            const rect = inputEl.getBoundingClientRect();
            const value = inputEl.value || "";
            lastQuery = value.trim();
            if (lastQuery) {
                performSearch(lastQuery, rect.left, rect.bottom, inputEl);
            } else {
                // Show empty search with recent items
                performSearch("", rect.left, rect.bottom, inputEl);
            }
        });
        
        // Try to insert wrapper
        try {
            if (parent.contains(inputEl)) {
                parent.insertBefore(wrapper, inputEl);
                wrapper.appendChild(inputEl);
                wrapper.appendChild(searchBtn);
            }
        } catch (e) {
            console.log("[PromptMaster] Could not add search icon:", e);
        }
    }
    
    // Add input event listener for real-time search
    const originalInputHandler = (e) => {
        const value = e.target.value || "";
        const cursorPosition = e.target.selectionStart || 0;
        
        // Get the word being typed (space or comma separated)
        const textBeforeCursor = value.substring(0, cursorPosition);
        const words = textBeforeCursor.split(/[,，\s]+/);
        const currentWord = words[words.length - 1];
        
        // Trigger search if current word is not empty
        if (currentWord.length >= 1) {
            lastQuery = currentWord;
            
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                const rect = inputEl.getBoundingClientRect();
                performSearch(currentWord, rect.left, rect.bottom, inputEl, (item) => {
                    // Replace current word with selected value
                    words[words.length - 1] = item.value;
                    const newTextBeforeCursor = words.join(", ");
                    const textAfterCursor = value.substring(cursorPosition);
                    inputEl.value = newTextBeforeCursor + textAfterCursor;
                    
                    // Set cursor position after inserted text
                    const newPosition = newTextBeforeCursor.length;
                    inputEl.setSelectionRange(newPosition, newPosition);
                    
                    // Trigger change event
                    inputEl.dispatchEvent(new Event("input", { bubbles: true }));
                    inputEl.dispatchEvent(new Event("change", { bubbles: true }));
                });
            }, 300); // 300ms debounce
        } else {
            removeExistingPopup();
        }
    };
    
    inputEl.addEventListener("input", originalInputHandler);
    
    // Add keyboard shortcut (Ctrl+Space) to trigger search
    inputEl.addEventListener("keydown", (e) => {
        if (e.ctrlKey && e.code === "Space") {
            e.preventDefault();
            const rect = inputEl.getBoundingClientRect();
            const value = inputEl.value || "";
            lastQuery = value.trim();
            performSearch(lastQuery || "", rect.left, rect.bottom, inputEl);
        }
        if (e.key === "Escape") {
            removeExistingPopup();
        }
    });
}

/**
 * Perform search and show popup
 */
async function performSearch(query, x, y, inputElement, onSelectCallback = null) {
    // Ensure coordinates are valid
    x = Math.max(10, x);
    y = Math.max(10, y);
    
    const suggestions = await searchAtoms(query, 10);
    
    const onSelect = onSelectCallback || ((item) => {
        // Default behavior: append to input
        const currentValue = inputElement.value || "";
        const separator = currentValue ? ", " : "";
        inputElement.value = currentValue + separator + item.value;
        inputElement.dispatchEvent(new Event("input", { bubbles: true }));
        inputElement.dispatchEvent(new Event("change", { bubbles: true }));
    });
    
    createSuggestionPopup(x, y, suggestions, onSelect, inputElement);
}

/**
 * Try to enhance widgets in a node
 */
function tryEnhanceNodeWidgets(node) {
    if (!node.widgets) return;
    
    const isTargetNode = node.comfyClass === "PromptSearch" || 
                         node.comfyClass === "PromptSearchAdvanced" ||
                         node.comfyClass === "CLIPTextEncode";
    
    if (!isTargetNode) return;
    
    console.log("[PromptMaster] Checking node:", node.comfyClass);
    
    for (const widget of node.widgets) {
        if (widget.type === "STRING") {
            // Wait a bit for DOM to be ready
            setTimeout(() => {
                enhanceTextWidget(node, widget);
            }, 100);
        }
    }
}

/**
 * Register with ComfyUI
 */
app.registerExtension({
    name: "PromptMaster.Search",
    
    async setup() {
        console.log("[PromptMaster] Search extension loaded");
        
        // Test API connectivity
        try {
            const response = await api.fetchApi("/promptmaster/health");
            if (response.ok) {
                const data = await response.json();
                console.log("[PromptMaster] API health check:", data);
            }
        } catch (e) {
            console.log("[PromptMaster] API not ready yet, will retry on node creation");
        }
    },
    
    async nodeCreated(node) {
        // Try to enhance widgets when node is created
        tryEnhanceNodeWidgets(node);
    },
    
    async afterConfigureGraph() {
        // Enhance all existing nodes after graph is loaded
        for (const node of app.graph._nodes || []) {
            tryEnhanceNodeWidgets(node);
        }
    }
});

// Export for other extensions
window.PromptMaster = {
    search: searchAtoms,
    createPopup: createSuggestionPopup,
    removePopup: removeExistingPopup
};
