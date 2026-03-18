"""
cache-redis-config

A Python package for configuring Redis cache.

Usage
-----

To use this package, you'll need to install it first:

    pip install cache-redis-config

Then, you can use it in your code like this:

    from cache_redis_config import Config

    config = Config()
    config.set('cache_ttl', 3600)
    config.set('cache_max_connections', 100)
    config.save()

"""

class Config:
    def __init__(self, filename='config.json'):
        self.filename = filename
        self.config = {}

        try:
            with open(self.filename, 'r') as f:
                self.config = json.load(f)
        except FileNotFoundError:
            pass

    def set(self, key, value):
        self.config[key] = value

    def save(self):
        with open(self.filename, 'w') as f:
            json.dump(self.config, f, indent=4)

import json