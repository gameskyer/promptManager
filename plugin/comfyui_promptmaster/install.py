#!/usr/bin/env python3
"""
Installation script for PromptMaster ComfyUI Plugin
"""

import os
import sys
import shutil


def find_comfyui_path():
    """Find ComfyUI installation path"""
    possible_paths = [
        "./ComfyUI",
        "../ComfyUI",
        "../../ComfyUI",
        "./",
    ]
    
    for path in possible_paths:
        if os.path.exists(os.path.join(path, "main.py")) or \
           os.path.exists(os.path.join(path, "comfy", "__init__.py")):
            return os.path.abspath(path)
    
    return None


def install_plugin():
    """Install plugin to ComfyUI"""
    print("=" * 60)
    print("PromptMaster ComfyUI Plugin Installer")
    print("=" * 60)
    
    # Find ComfyUI
    comfyui_path = find_comfyui_path()
    
    if comfyui_path is None:
        print("\n❌ Could not find ComfyUI installation.")
        print("Please run this script from within the ComfyUI directory.")
        print("\nOr manually install by copying this folder to:")
        print("  ComfyUI/custom_nodes/comfyui_promptmaster")
        return False
    
    print(f"\n✓ Found ComfyUI at: {comfyui_path}")
    
    # Check if we're already in custom_nodes
    current_dir = os.path.dirname(os.path.abspath(__file__))
    if "custom_nodes" in current_dir:
        print("\n✓ Already installed in custom_nodes")
        print("\n⚠ Make sure the database path is correct in database.py")
        return True
    
    # Install to custom_nodes
    target_dir = os.path.join(comfyui_path, "custom_nodes", "comfyui_promptmaster")
    
    if os.path.exists(target_dir):
        print(f"\n⚠ Plugin already exists at: {target_dir}")
        response = input("Overwrite? (y/N): ").strip().lower()
        if response != 'y':
            print("Installation cancelled.")
            return False
        shutil.rmtree(target_dir)
    
    # Copy files
    try:
        source_dir = os.path.dirname(os.path.abspath(__file__))
        shutil.copytree(source_dir, target_dir)
        print(f"\n✓ Plugin installed to: {target_dir}")
    except Exception as e:
        print(f"\n❌ Installation failed: {e}")
        return False
    
    # Check database
    db_paths = [
        os.path.join(comfyui_path, "..", "promptmaster.db"),
        os.path.join(comfyui_path, "promptmaster.db"),
        os.path.expanduser("~/.promptmaster/promptmaster.db"),
    ]
    
    db_found = False
    for db_path in db_paths:
        if os.path.exists(db_path):
            print(f"\n✓ Found database at: {os.path.abspath(db_path)}")
            db_found = True
            break
    
    if not db_found:
        print("\n⚠ Warning: Could not find PromptMaster database (promptmaster.db)")
        print("Please make sure the database file exists.")
        print("The plugin will search for it at runtime.")
    
    print("\n" + "=" * 60)
    print("Installation Complete!")
    print("=" * 60)
    print("\nNext steps:")
    print("1. Restart ComfyUI")
    print("2. Add the '🔍 Prompt Search' node to your workflow")
    print("3. Type in Chinese to search for prompt atoms")
    print("\n" + "=" * 60)
    
    return True


if __name__ == "__main__":
    success = install_plugin()
    sys.exit(0 if success else 1)
