#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import time
import random
import urllib.request
import os

# URLs for real stock images (using placeholder image services that don't require auth)
image_urls = {
    'electronics': [
        'https://picsum.photos/800/600?random=1',  # For iPhone
        'https://picsum.photos/800/600?random=2',  # For Samsung
        'https://picsum.photos/800/600?random=3',  # For TV
    ],
    'fashion': [
        'https://picsum.photos/800/600?random=4',  # For suit
        'https://picsum.photos/800/600?random=5',  # For shoes
        'https://picsum.photos/800/600?random=6',  # For bag
    ],
    'home': [
        'https://picsum.photos/800/600?random=7',  # For tools
        'https://picsum.photos/800/600?random=8',  # For carpet
    ],
    'agriculture': [
        'https://picsum.photos/800/600?random=9',   # For tractor
        'https://picsum.photos/800/600?random=10',  # For seeds
        'https://picsum.photos/800/600?random=11',  # For cattle
        'https://picsum.photos/800/600?random=12',  # For honey
    ],
    'industry': [
        'https://picsum.photos/800/600?random=13',  # For CNC
        'https://picsum.photos/800/600?random=14',  # For safety
    ],
    'food': [
        'https://picsum.photos/800/600?random=15',  # For rakija
    ]
}

# Download and save images locally in MinIO directory structure
minio_path = '/data/minio/listings'
os.makedirs(minio_path, exist_ok=True)

def download_images():
    """Download placeholder images for listings"""
    listing_id = 1
    
    for category, urls in image_urls.items():
        for url in urls:
            try:
                # Create directory for each listing
                listing_dir = os.path.join(minio_path, str(listing_id))
                os.makedirs(listing_dir, exist_ok=True)
                
                # Download image
                image_path = os.path.join(listing_dir, 'main.jpg')
                urllib.request.urlretrieve(url, image_path)
                print(f"Downloaded image for listing {listing_id}")
                
                # Also create a second image (copy of first)
                image2_path = os.path.join(listing_dir, 'image2.jpg')
                urllib.request.urlretrieve(url + '&v=2', image2_path)
                
                listing_id += 1
                time.sleep(0.5)  # Rate limiting
                
            except Exception as e:
                print(f"Error downloading image: {e}")
                listing_id += 1
                continue

if __name__ == "__main__":
    print("Downloading placeholder images for listings...")
    download_images()
    print("Done!")