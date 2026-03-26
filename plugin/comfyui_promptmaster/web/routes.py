"""
ComfyUI API Routes for PromptMaster
"""

import json
import os
import sys

# Add parent directory to path to import database
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from database import get_db


class PromptMasterRoutes:
    """API routes for PromptMaster search functionality"""
    
    def __init__(self):
        self.db = get_db()
    
    def search(self, request):
        """
        Handle search request
        
        POST /promptmaster/search
        Body: {"query": "裙", "limit": 10}
        """
        try:
            if hasattr(request, 'get_json'):
                # Flask-style request
                data = request.get_json()
            else:
                # Handle raw data
                data = json.loads(request)
            
            query = data.get('query', '')
            limit = data.get('limit', 10)
            
            if not query:
                return json.dumps({
                    "success": True,
                    "query": query,
                    "results": [],
                    "count": 0
                }, ensure_ascii=False)
            
            results = self.db.search_atoms(query, limit)
            
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
                "query": data.get('query', '') if 'data' in locals() else '',
                "results": [],
                "count": 0
            }, ensure_ascii=False)
    
    def categories(self, request=None):
        """
        Get all categories
        
        GET /promptmaster/categories
        """
        try:
            categories = self.db.get_all_categories()
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
    
    def recent(self, request=None):
        """
        Get recent atoms
        
        GET /promptmaster/recent
        """
        try:
            atoms = self.db.get_recent_atoms(20)
            return json.dumps({
                "success": True,
                "atoms": atoms,
                "count": len(atoms)
            }, ensure_ascii=False)
        except Exception as e:
            return json.dumps({
                "success": False,
                "error": str(e),
                "atoms": [],
                "count": 0
            }, ensure_ascii=False)


# Global routes instance
routes = PromptMasterRoutes()


def setup_routes(server):
    """
    Setup routes with ComfyUI server
    
    This should be called from ComfyUI's server initialization
    """
    try:
        # Try to add routes using ComfyUI's route system
        if hasattr(server, 'add_route'):
            server.add_route("POST", "/promptmaster/search", routes.search)
            server.add_route("GET", "/promptmaster/categories", routes.categories)
            server.add_route("GET", "/promptmaster/recent", routes.recent)
        elif hasattr(server, 'app') and hasattr(server.app, 'route'):
            # Flask-style routing
            from flask import request as flask_request
            
            @server.app.route('/promptmaster/search', methods=['POST'])
            def search():
                return routes.search(flask_request)
            
            @server.app.route('/promptmaster/categories', methods=['GET'])
            def get_categories():
                return routes.categories()
            
            @server.app.route('/promptmaster/recent', methods=['GET'])
            def get_recent():
                return routes.recent()
                
        print("[PromptMaster] API routes registered successfully")
    except Exception as e:
        print(f"[PromptMaster] Failed to register routes: {e}")


# For ComfyUI's route extension
def get_custom_routes():
    """Return custom routes for ComfyUI"""
    return [
        ("POST", "/promptmaster/search", routes.search),
        ("GET", "/promptmaster/categories", routes.categories),
        ("GET", "/promptmaster/recent", routes.recent),
    ]
