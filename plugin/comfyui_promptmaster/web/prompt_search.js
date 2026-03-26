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
        top: ${y + 25}px;
        background: #1e1e1e;
        border: 2px solid #4a9eff;
        border-radius: 8px;
        padding: 8px 0;
        min-width: 280px;
        max-width: 400px;
        max-height: 300px;
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
                            ${highlightMatch(item.label, lastQuery)}
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
            
            div.addEventListener("click", () => {
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
    if (!query) return text;
    const regex = new RegExp(`(${escapeRegExp(query)})`, "gi");
    return text.replace(regex, '<span style="color: #ffd43b; font-weight: bold;">$1</span>');
}

function escapeRegExp(string) {
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
        // Call the node's search function via API
        const response = await api.fetchApi("/promptmaster/search", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, limit })
        });
        
        const data = await response.json();
        
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
 * Enhance text input widget with search functionality
 */
function enhanceTextWidget(node, widget) {
    if (widget.type !== "STRING" || widget.name === "selected_prompt") {
        return;
    }
    
    const originalInput = widget.inputEl || widget.element;
    if (!originalInput || originalInput._promptmasterEnhanced) {
        return;
    }
    
    originalInput._promptmasterEnhanced = true;
    
    // Create wrapper
    const wrapper = document.createElement("div");
    wrapper.style.cssText = "position: relative; display: flex; align-items: center;";
    
    // Add search icon/button
    const searchBtn = document.createElement("span");
    searchBtn.innerHTML = "🔍";
    searchBtn.style.cssText = `
        position: absolute;
        right: 8px;
        cursor: pointer;
        font-size: 14px;
        opacity: 0.6;
        transition: opacity 0.2s;
        z-index: 10;
    `;
    searchBtn.title = "点击搜索提示词";
    
    searchBtn.addEventListener("mouseenter", () => searchBtn.style.opacity = "1");
    searchBtn.addEventListener("mouseleave", () => searchBtn.style.opacity = "0.6");
    searchBtn.addEventListener("click", () => {
        const rect = originalInput.getBoundingClientRect();
        lastQuery = originalInput.value.trim();
        if (lastQuery) {
            performSearch(lastQuery, rect.left, rect.top, originalInput);
        }
    });
    
    // Insert wrapper
    if (originalInput.parentNode) {
        originalInput.parentNode.insertBefore(wrapper, originalInput);
        wrapper.appendChild(originalInput);
        wrapper.appendChild(searchBtn);
    }
    
    // Add input event listener for real-time search
    originalInput.addEventListener("input", (e) => {
        const value = e.target.value;
        const cursorPosition = e.target.selectionStart;
        
        // Get the word being typed (space or comma separated)
        const textBeforeCursor = value.substring(0, cursorPosition);
        const words = textBeforeCursor.split(/[,，\s]+/);
        const currentWord = words[words.length - 1];
        
        if (currentWord.length >= 1 && /[\u4e00-\u9fa5]/i.test(currentWord)) {
            // Chinese input detected
            lastQuery = currentWord;
            
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                const rect = originalInput.getBoundingClientRect();
                performSearch(currentWord, rect.left, rect.bottom, originalInput, (item) => {
                    // Replace current word with selected value
                    words[words.length - 1] = item.value;
                    const newTextBeforeCursor = words.join(", ");
                    const textAfterCursor = value.substring(cursorPosition);
                    originalInput.value = newTextBeforeCursor + textAfterCursor;
                    
                    // Set cursor position after inserted text
                    const newPosition = newTextBeforeCursor.length;
                    originalInput.setSelectionRange(newPosition, newPosition);
                    
                    // Trigger change event
                    originalInput.dispatchEvent(new Event("input", { bubbles: true }));
                    originalInput.dispatchEvent(new Event("change", { bubbles: true }));
                });
            }, 200);
        }
    });
    
    // Close popup on escape
    originalInput.addEventListener("keydown", (e) => {
        if (e.key === "Escape") {
            removeExistingPopup();
        }
    });
}

/**
 * Perform search and show popup
 */
async function performSearch(query, x, y, inputElement, onSelectCallback = null) {
    if (!query.trim()) {
        removeExistingPopup();
        return;
    }
    
    const suggestions = await searchAtoms(query, 10);
    
    const onSelect = onSelectCallback || ((item) => {
        // Default behavior: append to input
        const currentValue = inputElement.value.trim();
        const separator = currentValue ? ", " : "";
        inputElement.value = currentValue + separator + item.value;
        inputElement.dispatchEvent(new Event("input", { bubbles: true }));
        inputElement.dispatchEvent(new Event("change", { bubbles: true }));
    });
    
    createSuggestionPopup(x, y, suggestions, onSelect, inputElement);
}

/**
 * Register with ComfyUI
 */
app.registerExtension({
    name: "PromptMaster.Search",
    
    async setup() {
        console.log("[PromptMaster] Search extension loaded");
        
        // Register API endpoint
        api.addEventListener("promptmaster.search", (e) => {
            console.log("[PromptMaster] Search event:", e.detail);
        });
    },
    
    async nodeCreated(node) {
        // Enhance all string widgets in CLIPTextEncode nodes and PromptMaster nodes
        if (node.comfyClass === "CLIPTextEncode" || node.comfyClass === "PromptSearch" || node.comfyClass === "PromptSearchAdvanced") {
            for (const widget of node.widgets || []) {
                if (widget.type === "STRING" && widget.inputEl) {
                    // Wait for DOM to be ready
                    setTimeout(() => {
                        enhanceTextWidget(node, widget);
                    }, 100);
                }
            }
        }
    },
    
    async beforeRegisterNodeDef(nodeType, nodeData, app) {
        // Hook into string widget creation
        if (nodeData.input && nodeData.input.required) {
            for (const [name, config] of Object.entries(nodeData.input.required)) {
                if (Array.isArray(config) && config[0] === "STRING") {
                    const originalOnDraw = nodeType.prototype.onDrawForeground;
                    
                    nodeType.prototype.onDrawForeground = function(ctx) {
                        if (originalOnDraw) {
                            originalOnDraw.apply(this, arguments);
                        }
                        
                        // Enhance widgets on first draw
                        if (!this._promptmasterEnhanced) {
                            this._promptmasterEnhanced = true;
                            for (const widget of this.widgets || []) {
                                if (widget.type === "STRING" && widget.inputEl) {
                                    enhanceTextWidget(this, widget);
                                }
                            }
                        }
                    };
                }
            }
        }
    }
});

// Export for other extensions
window.PromptMaster = {
    search: searchAtoms,
    createPopup: createSuggestionPopup,
    removePopup: removeExistingPopup
};
