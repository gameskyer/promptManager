"""
ComfyUI Node for Prompt Search
"""

import os
import json
from .database import get_db


class PromptSearchNode:
    """
    A ComfyUI node that provides Chinese/English prompt search functionality
    """
    
    @classmethod
    def INPUT_TYPES(cls):
        return {
            "required": {
                "search_query": ("STRING", {
                    "multiline": False,
                    "default": "",
                    "placeholder": "输入中文或英文搜索... (例如: 裙, dress)",
                }),
                "selected_prompt": ("STRING", {
                    "multiline": True,
                    "default": "",
                    "placeholder": "选择的提示词将显示在这里...",
                }),
            },
            "optional": {
                "category_filter": (["全部"] + cls._get_categories(), {
                    "default": "全部",
                }),
                "max_results": ("INT", {
                    "default": 10,
                    "min": 1,
                    "max": 50,
                }),
            }
        }
    
    @classmethod
    def _get_categories(cls):
        """Get category names for dropdown"""
        try:
            db = get_db()
            categories = db.get_all_categories()
            return [c['name'] for c in categories if c.get('name')]
        except:
            return []
    
    RETURN_TYPES = ("STRING",)
    RETURN_NAMES = ("prompt",)
    FUNCTION = "search_and_select"
    CATEGORY = "PromptMaster"
    
    def search_and_select(self, search_query, selected_prompt, category_filter="全部", max_results=10):
        """
        Search for atoms and return selected prompt
        
        Note: The actual search UI is handled by the JavaScript extension.
        This node serves as a container for the prompt value.
        """
        # Return the selected prompt (CLIP text)
        return (selected_prompt,)
    
    @classmethod
    def IS_CHANGED(cls, search_query, selected_prompt, **kwargs):
        """Force refresh when search query changes"""
        return f"{search_query}_{selected_prompt}_{json.dumps(kwargs, default=str)}"


class PromptSearchAdvanced:
    """
    Advanced version with multiple outputs and weight support
    """
    
    @classmethod
    def INPUT_TYPES(cls):
        return {
            "required": {
                "search_query": ("STRING", {
                    "multiline": False,
                    "default": "",
                    "placeholder": "搜索提示词...",
                }),
            },
            "optional": {
                "weight": ("FLOAT", {
                    "default": 1.0,
                    "min": 0.0,
                    "max": 2.0,
                    "step": 0.1,
                }),
                "prefix": ("STRING", {
                    "multiline": False,
                    "default": "",
                    "placeholder": "前缀 (如: best quality)",
                }),
                "suffix": ("STRING", {
                    "multiline": False,
                    "default": "",
                    "placeholder": "后缀 (如: masterpiece)",
                }),
            }
        }
    
    RETURN_TYPES = ("STRING", "STRING")
    RETURN_NAMES = ("positive_prompt", "raw_value")
    FUNCTION = "search_advanced"
    CATEGORY = "PromptMaster"
    
    def search_advanced(self, search_query, weight=1.0, prefix="", suffix=""):
        """
        Advanced search with weight support
        
        The actual search and selection is handled by JavaScript.
        This node processes the selected value with weight.
        """
        # The search_query from JavaScript will contain the selected atom value
        raw_value = search_query.strip()
        
        if not raw_value:
            return ("", "")
        
        # Apply weight if not 1.0
        if weight != 1.0:
            if weight < 1.0:
                # Lower weight: (value:0.8)
                prompt = f"({raw_value}:{weight:.1f})"
            else:
                # Higher weight: ((value))
                prompt = f"({raw_value})"
                # Double parentheses for weight > 1
                for _ in range(int(weight) - 1):
                    prompt = f"({prompt})"
        else:
            prompt = raw_value
        
        # Add prefix and suffix
        parts = []
        if prefix.strip():
            parts.append(prefix.strip())
        parts.append(prompt)
        if suffix.strip():
            parts.append(suffix.strip())
        
        final_prompt = ", ".join(parts)
        
        return (final_prompt, raw_value)


# API endpoint for JavaScript to call
def search_atoms_api(query: str, limit: int = 10):
    """
    API function for JavaScript to search atoms
    
    Returns JSON string with search results
    """
    try:
        db = get_db()
        results = db.search_atoms(query, limit)
        return json.dumps({
            "success": True,
            "query": query,
            "results": results,
            "count": len(results)
        }, ensure_ascii=False)
    except Exception as e:
        return json.dumps({
            "success": False,
            "error": str(e),
            "query": query,
            "results": [],
            "count": 0
        }, ensure_ascii=False)


def get_categories_api():
    """Get all categories for dropdown"""
    try:
        db = get_db()
        categories = db.get_all_categories()
        return json.dumps({
            "success": True,
            "categories": categories
        }, ensure_ascii=False)
    except Exception as e:
        return json.dumps({
            "success": False,
            "error": str(e),
            "categories": []
        }, ensure_ascii=False)
