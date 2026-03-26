"""
Database access module - Direct SQLite access to PromptMaster database
"""

import sqlite3
import json
import os
from typing import List, Dict, Optional


class PromptMasterDB:
    """Direct database access to PromptMaster SQLite database"""
    
    def __init__(self, db_path: str = None):
        """
        Initialize database connection
        
        Args:
            db_path: Path to promptmaster.db. If None, searches in common locations
        """
        if db_path is None:
            db_path = self._find_database()
        
        self.db_path = db_path
        self.conn = None
        
    def _find_database(self) -> str:
        """Find database in common locations"""
        # Possible database locations (in order of preference)
        possible_paths = [
            # Same directory as plugin
            os.path.join(os.path.dirname(__file__), "..", "..", "promptmaster.db"),
            # Current working directory
            "./promptmaster.db",
            # Parent directory
            "../promptmaster.db",
            # ComfyUI root
            "../../promptmaster.db",
            # User home
            os.path.expanduser("~/.promptmaster/promptmaster.db"),
        ]
        
        for path in possible_paths:
            abs_path = os.path.abspath(path)
            if os.path.exists(abs_path):
                print(f"[PromptMaster] Found database at: {abs_path}")
                return abs_path
        
        # Default fallback
        return "./promptmaster.db"
    
    def connect(self):
        """Establish database connection"""
        if self.conn is None:
            self.conn = sqlite3.connect(self.db_path)
            self.conn.row_factory = sqlite3.Row
        return self.conn
    
    def close(self):
        """Close database connection"""
        if self.conn:
            self.conn.close()
            self.conn = None
    
    def search_atoms(self, query: str, limit: int = 10) -> List[Dict]:
        """
        Search atoms by value or label (supports Chinese)
        
        Args:
            query: Search query (e.g., "裙", "dress")
            limit: Maximum number of results
            
        Returns:
            List of matching atoms
        """
        conn = self.connect()
        cursor = conn.cursor()
        
        # Search by value (English) or label (Chinese)
        # Also search in synonyms
        sql = """
        SELECT 
            a.id,
            a.value,
            a.label,
            a.type,
            a.synonyms,
            a.category_id,
            c.name as category_name,
            c.parent_id
        FROM atoms a
        LEFT JOIN categories c ON a.category_id = c.id
        WHERE a.value LIKE ? 
           OR a.label LIKE ? 
           OR a.synonyms LIKE ?
        ORDER BY 
            CASE WHEN a.label = ? THEN 0 ELSE 1 END,
            CASE WHEN a.label LIKE ? THEN 0 ELSE 1 END,
            a.usage_count DESC
        LIMIT ?
        """
        
        # Build search patterns
        exact_match = query
        prefix_match = f"{query}%"
        contains_match = f"%{query}%"
        
        params = (
            contains_match,  # value LIKE
            contains_match,  # label LIKE
            contains_match,  # synonyms LIKE
            exact_match,     # exact match priority
            prefix_match,    # prefix match priority
            limit
        )
        
        cursor.execute(sql, params)
        rows = cursor.fetchall()
        
        results = []
        for row in rows:
            atom = dict(row)
            # Parse synonyms JSON if exists
            if atom.get('synonyms'):
                try:
                    if isinstance(atom['synonyms'], str):
                        if atom['synonyms'].startswith('['):
                            atom['synonyms'] = json.loads(atom['synonyms'])
                        else:
                            # Legacy format: comma-separated or single value
                            atom['synonyms'] = [s.strip() for s in atom['synonyms'].split(',') if s.strip()]
                    else:
                        atom['synonyms'] = []
                except:
                    atom['synonyms'] = []
            else:
                atom['synonyms'] = []
            
            results.append(atom)
        
        return results
    
    def get_all_categories(self) -> List[Dict]:
        """Get all categories (for filtering)"""
        conn = self.connect()
        cursor = conn.cursor()
        
        cursor.execute("""
            SELECT id, name, parent_id, type 
            FROM categories 
            WHERE type = 'ATOM'
            ORDER BY parent_id, sort_order
        """)
        
        return [dict(row) for row in cursor.fetchall()]
    
    def get_atoms_by_category(self, category_id: int, limit: int = 50) -> List[Dict]:
        """Get atoms by category"""
        conn = self.connect()
        cursor = conn.cursor()
        
        cursor.execute("""
            SELECT 
                a.id, a.value, a.label, a.type, a.synonyms,
                c.name as category_name
            FROM atoms a
            LEFT JOIN categories c ON a.category_id = c.id
            WHERE a.category_id = ?
            ORDER BY a.usage_count DESC
            LIMIT ?
        """, (category_id, limit))
        
        results = []
        for row in cursor.fetchall():
            atom = dict(row)
            if atom.get('synonyms'):
                try:
                    if isinstance(atom['synonyms'], str):
                        if atom['synonyms'].startswith('['):
                            atom['synonyms'] = json.loads(atom['synonyms'])
                        else:
                            atom['synonyms'] = [s.strip() for s in atom['synonyms'].split(',') if s.strip()]
                    else:
                        atom['synonyms'] = []
                except:
                    atom['synonyms'] = []
            else:
                atom['synonyms'] = []
            results.append(atom)
        
        return results
    
    def get_recent_atoms(self, limit: int = 20) -> List[Dict]:
        """Get recently used atoms"""
        conn = self.connect()
        cursor = conn.cursor()
        
        cursor.execute("""
            SELECT 
                a.id, a.value, a.label, a.type, a.synonyms,
                c.name as category_name
            FROM atoms a
            LEFT JOIN categories c ON a.category_id = c.id
            WHERE a.last_used_at IS NOT NULL
            ORDER BY a.last_used_at DESC
            LIMIT ?
        """, (limit,))
        
        return [dict(row) for row in cursor.fetchall()]


# Global database instance
db_instance = None

def get_db() -> PromptMasterDB:
    """Get global database instance"""
    global db_instance
    if db_instance is None:
        db_instance = PromptMasterDB()
    return db_instance
