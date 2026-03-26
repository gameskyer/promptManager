"""
PromptMaster for ComfyUI
A plugin that integrates PromptMaster's atom search into ComfyUI
"""

from .prompt_search_node import PromptSearchNode, PromptSearchAdvanced

NODE_CLASS_MAPPINGS = {
    "PromptSearch": PromptSearchNode,
    "PromptSearchAdvanced": PromptSearchAdvanced,
}

NODE_DISPLAY_NAME_MAPPINGS = {
    "PromptSearch": "🔍 Prompt Search",
    "PromptSearchAdvanced": "🔍 Prompt Search (Advanced)",
}

__all__ = ['NODE_CLASS_MAPPINGS', 'NODE_DISPLAY_NAME_MAPPINGS']

# Web 扩展路径
WEB_DIRECTORY = "./web"
