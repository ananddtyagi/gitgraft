"""Configuration management for GitGraft"""

import os
import json
from pathlib import Path


class Config:
    """Manages GitGraft configuration"""
    
    def __init__(self):
        self.config_dir = Path.home() / ".config" / "gitgraft"
        self.config_file = self.config_dir / "config.json"
        self.config_dir.mkdir(parents=True, exist_ok=True)
        self._load_config()
    
    def _load_config(self):
        """Load configuration from file"""
        if self.config_file.exists():
            with open(self.config_file, 'r') as f:
                self.data = json.load(f)
        else:
            self.data = {}
    
    def save(self):
        """Save configuration to file"""
        with open(self.config_file, 'w') as f:
            json.dump(self.data, f, indent=2)
    
    def is_configured(self):
        """Check if onboarding has been completed"""
        return self.data.get("onboarding_complete", False)
    
    def set_onboarding_complete(self):
        """Mark onboarding as complete"""
        self.data["onboarding_complete"] = True
        self.save()
    
    def set_alias_preference(self, create_alias):
        """Save user's alias preference"""
        self.data["create_alias"] = create_alias
        self.save()
    
    def get_alias_preference(self):
        """Get user's alias preference"""
        return self.data.get("create_alias", False)
