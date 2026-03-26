/**
 * PromptMaster Search Extension for ComfyUI
 */

import { app } from "../../scripts/app.js";
import { api } from "../../scripts/api.js";

console.log("[PromptMaster] Extension script loading...");

// Global state
let searchCache = new Map();
let debounceTimer = null;

/**
 * Search atoms via API
 */
async function searchAtoms(query, limit = 10) {
    if (!query || query.length < 1) {
        return [];
    }
    
    const cacheKey = `${query}_${limit}`;
    if (searchCache.has(cacheKey)) {
        return searchCache.get(cacheKey);
    }
    
    try {
        console.log("[PromptMaster] Searching for:", query);
        
        const response = await api.fetchApi("/promptmaster/search", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, limit })
        });
        
        if (!response.ok) {
            console.error("[PromptMaster] Search API error:", response.status);
            return [];
        }
        
        const data = await response.json();
        console.log("[PromptMaster] Got results:", data.count);
        
        if (data.success) {
            searchCache.set(cacheKey, data.results);
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
 * Create search popup
 */
function createSearchPopup(x, y, results, onSelect) {
    // Remove existing
    const existing = document.getElementById("pm-search-popup");
    if (existing) existing.remove();
    
    const popup = document.createElement("div");
    popup.id = "pm-search-popup";
    popup.style.cssText = `
        position: fixed;
        left: ${x}px;
        top: ${y + 5}px;
        background: #2a2a2a;
        border: 2px solid #4a9eff;
        border-radius: 6px;
        min-width: 300px;
        max-height: 300px;
        overflow-y: auto;
        z-index: 99999;
        box-shadow: 0 4px 12px rgba(0,0,0,0.5);
    `;
    
    if (results.length === 0) {
        popup.innerHTML = `<div style="padding: 12px; color: #888; text-align: center;">无搜索结果</div>`;
    } else {
        const header = document.createElement("div");
        header.textContent = `找到 ${results.length} 个结果`;
        header.style.cssText = "padding: 8px 12px; color: #4a9eff; font-size: 11px; border-bottom: 1px solid #444;";
        popup.appendChild(header);
        
        results.forEach(item => {
            const div = document.createElement("div");
            const typeColor = item.type === "Negative" ? "#ff6b6b" : "#51cf66";
            const typeLabel = item.type === "Negative" ? "负" : "正";
            
            div.innerHTML = `
                <div style="display: flex; align-items: center; gap: 8px; padding: 8px 12px; cursor: pointer;">
                    <span style="background: ${typeColor}; color: white; font-size: 10px; padding: 2px 6px; border-radius: 4px;">${typeLabel}</span>
                    <div>
                        <div style="color: #fff; font-size: 13px;">${item.label || item.value}</div>
                        <div style="color: #888; font-size: 11px;">${item.value}</div>
                    </div>
                </div>
            `;
            
            div.addEventListener("mouseenter", () => {
                div.style.background = "#3a3a3a";
            });
            div.addEventListener("mouseleave", () => {
                div.style.background = "transparent";
            });
            div.addEventListener("click", () => {
                onSelect(item);
                popup.remove();
            });
            
            popup.appendChild(div);
        });
    }
    
    document.body.appendChild(popup);
    
    // Close on click outside
    setTimeout(() => {
        const closeHandler = (e) => {
            if (!popup.contains(e.target)) {
                popup.remove();
                document.removeEventListener("click", closeHandler);
            }
        };
        document.addEventListener("click", closeHandler);
    }, 10);
}

/**
 * Enhance a widget with search
 */
function enhanceWidget(node, widget) {
    const inputEl = widget.inputEl || widget.element;
    if (!inputEl || inputEl._pmEnhanced) return;
    
    // Only enhance STRING widgets in PromptMaster or CLIPTextEncode nodes
    if (widget.type !== "STRING") return;
    
    const nodeType = node.comfyClass || node.type;
    const isTargetNode = nodeType === "PromptSearch" || 
                         nodeType === "PromptSearchAdvanced" ||
                         nodeType === "CLIPTextEncode";
    
    if (!isTargetNode) return;
    
    // Skip certain widgets
    if (widget.name === "selected_prompt") return;
    
    console.log("[PromptMaster] Enhancing widget:", widget.name, "in node:", nodeType);
    
    inputEl._pmEnhanced = true;
    
    // Add search icon
    const wrapper = document.createElement("div");
    wrapper.style.cssText = "position: relative; display: flex; align-items: center;";
    
    const searchBtn = document.createElement("span");
    searchBtn.innerHTML = "🔍";
    searchBtn.style.cssText = `
        position: absolute;
        right: 6px;
        cursor: pointer;
        font-size: 14px;
        opacity: 0.5;
        user-select: none;
    `;
    searchBtn.title = "搜索提示词 (Ctrl+Space)";
    
    searchBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        const rect = inputEl.getBoundingClientRect();
        const query = inputEl.value?.trim() || "";
        doSearch(query, rect.left, rect.bottom, inputEl);
    });
    
    if (inputEl.parentNode) {
        inputEl.parentNode.insertBefore(wrapper, inputEl);
        wrapper.appendChild(inputEl);
        wrapper.appendChild(searchBtn);
    }
    
    // Input handler
    inputEl.addEventListener("input", (e) => {
        const value = e.target.value || "";
        const cursorPos = e.target.selectionStart || 0;
        const textBefore = value.substring(0, cursorPos);
        const words = textBefore.split(/[,，\s]+/);
        const currentWord = words[words.length - 1];
        
        if (currentWord.length >= 1) {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                const rect = inputEl.getBoundingClientRect();
                doSearch(currentWord, rect.left, rect.bottom, inputEl, (item) => {
                    // Replace word
                    words[words.length - 1] = item.value;
                    const newText = words.join(", ") + value.substring(cursorPos);
                    inputEl.value = newText;
                    inputEl.dispatchEvent(new Event("input", { bubbles: true }));
                });
            }, 300);
        }
    });
    
    // Keyboard shortcut
    inputEl.addEventListener("keydown", (e) => {
        if (e.ctrlKey && e.code === "Space") {
            e.preventDefault();
            const rect = inputEl.getBoundingClientRect();
            doSearch(inputEl.value?.trim() || "", rect.left, rect.bottom, inputEl);
        }
        if (e.key === "Escape") {
            const popup = document.getElementById("pm-search-popup");
            if (popup) popup.remove();
        }
    });
}

/**
 * Perform search
 */
async function doSearch(query, x, y, inputEl, onSelect) {
    const results = await searchAtoms(query, 10);
    const defaultSelect = (item) => {
        const current = inputEl.value?.trim() || "";
        inputEl.value = current ? current + ", " + item.value : item.value;
        inputEl.dispatchEvent(new Event("input", { bubbles: true }));
    };
    createSearchPopup(x, y, results, onSelect || defaultSelect);
}

// Register extension
app.registerExtension({
    name: "PromptMaster.Search",
    
    async setup() {
        console.log("[PromptMaster] Extension setup called");
        
        // Test API
        try {
            const resp = await api.fetchApi("/promptmaster/health");
            if (resp.ok) {
                const data = await resp.json();
                console.log("[PromptMaster] API ready:", data);
            } else {
                console.log("[PromptMaster] API health check failed:", resp.status);
            }
        } catch (e) {
            console.log("[PromptMaster] API not ready:", e.message);
        }
    },
    
    async nodeCreated(node) {
        // Enhance widgets after a short delay
        setTimeout(() => {
            if (node.widgets) {
                node.widgets.forEach(w => enhanceWidget(node, w));
            }
        }, 100);
    }
});

console.log("[PromptMaster] Extension registered");
