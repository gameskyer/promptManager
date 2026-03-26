"""
PromptMaster for ComfyUI
A plugin that integrates PromptMaster's atom search into ComfyUI
"""

import os
import sys
import json

# Add current directory to path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from .prompt_search_node import PromptSearchNode, PromptSearchAdvanced
from .database import get_db

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


# =============================================================================
# API Routes Setup - ComfyUI aiohttp server integration
# =============================================================================

def setup_comfyui_routes(server):
    """Setup routes with ComfyUI's aiohttp server"""
    try:
        import aiohttp
        from aiohttp import web
        
        db = get_db()
        
        @server.routes.get("/promptmaster/health")
        async def health_check(request):
            """Health check endpoint"""
            return web.json_response({
                "status": "ok",
                "version": "1.0.0"
            })
        
        @server.routes.post("/promptmaster/search")
        async def search_prompts(request):
            """Search atoms endpoint"""
            try:
                data = await request.json()
                query = data.get('query', '')
                limit = data.get('limit', 10)
                
                if not query:
                    return web.json_response({
                        "success": True,
                        "query": query,
                        "results": [],
                        "count": 0
                    })
                
                results = db.search_atoms(query, limit)
                
                return web.json_response({
                    "success": True,
                    "query": query,
                    "results": results,
                    "count": len(results)
                })
            except Exception as e:
                print(f"[PromptMaster] Search error: {e}")
                return web.json_response({
                    "success": False,
                    "error": str(e),
                    "results": [],
                    "count": 0
                }, status=500)
        
        @server.routes.get("/promptmaster/categories")
        async def get_categories(request):
            """Get all categories"""
            try:
                categories = db.get_all_categories()
                return web.json_response({
                    "success": True,
                    "categories": categories
                })
            except Exception as e:
                return web.json_response({
                    "success": False,
                    "error": str(e),
                    "categories": []
                }, status=500)
        
        @server.routes.get("/promptmaster/recent")
        async def get_recent(request):
            """Get recent atoms"""
            try:
                limit = int(request.query.get('limit', 20))
                atoms = db.get_recent_atoms(limit)
                return web.json_response({
                    "success": True,
                    "atoms": atoms,
                    "count": len(atoms)
                })
            except Exception as e:
                return web.json_response({
                    "success": False,
                    "error": str(e),
                    "atoms": [],
                    "count": 0
                }, status=500)
        
        print("[PromptMaster] API routes registered successfully")
        
    except Exception as e:
        print(f"[PromptMaster] Failed to register routes: {e}")
        import traceback
        traceback.print_exc()


# Try to get the server instance and setup routes
# This will be called when ComfyUI loads the extension
try:
    # Method 1: Try to get server from comfyui module
    import server
    if hasattr(server, 'PromptServer') and hasattr(server.PromptServer, 'instance'):
        setup_comfyui_routes(server.PromptServer.instance)
except Exception as e:
    print(f"[PromptMaster] Note: Routes will be set up later: {e}")


# Also export a setup function that can be called manually if needed
def setup(server_instance):
    """Manual setup function"""
    setup_comfyui_routes(server_instance)
