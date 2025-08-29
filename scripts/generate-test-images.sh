#!/bin/bash

echo "🖼️ Generating test images for MinIO..."

# Create temp directory
TEMP_DIR="/tmp/test-images"
rm -rf $TEMP_DIR
mkdir -p $TEMP_DIR

# Get all image paths from database
docker-compose -f docker-compose.dev.yml exec -T db psql -U svetu_dev_user -d svetu_dev_db -t -c "SELECT DISTINCT file_path FROM marketplace_images ORDER BY file_path;" | while read file_path; do
    if [ -n "$file_path" ]; then
        # Clean up the path
        file_path=$(echo $file_path | tr -d ' ')
        
        # Create directory structure
        dir_path=$(dirname "$file_path")
        mkdir -p "$TEMP_DIR/$dir_path"
        
        # Get listing ID from path (first part before /)
        listing_id=$(echo $file_path | cut -d'/' -f1)
        filename=$(basename "$file_path")
        
        # Generate random color
        COLOR=$(printf '#%06X' $((RANDOM % 16777216)))
        
        # Create test image with ImageMagick
        convert -size 800x600 xc:"$COLOR" \
            -gravity center -pointsize 48 -fill white \
            -annotate +0+0 "Product $listing_id\n$filename" \
            "$TEMP_DIR/$file_path"
        
        echo "Created: $file_path"
    fi
done

echo "✅ Generated $(find $TEMP_DIR -type f | wc -l) test images"

# Now upload to MinIO using mc
echo "📤 Uploading to MinIO..."

# Configure mc for MinIO
mc alias set devminio http://localhost:9002 miniodevadmin h8Puk2qhcadCazC78J1

# Create bucket if not exists
mc mb devminio/listings 2>/dev/null || true

# Upload all images
mc cp --recursive $TEMP_DIR/* devminio/listings/

# Set bucket to public
mc anonymous set download devminio/listings

echo "🎉 Done! Images uploaded to MinIO"

# Verify
echo "📊 Verification:"
mc ls --recursive devminio/listings/ | head -10
echo "..."
echo "Total: $(mc ls --recursive devminio/listings/ | wc -l) files in MinIO"

# Clean up
rm -rf $TEMP_DIR