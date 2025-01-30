// frontend/hostel-frontend/src/components/marketplace/ImageUploader.js
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import imageCompression from 'browser-image-compression';
import { Box, Button, IconButton, Typography, CircularProgress } from '@mui/material';
import { CloudUpload as CloudUploadIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { addWatermark } from '../../utils/imageUtils';

const ImageUploader = ({ onImagesSelected, maxImages = 10, maxSizeMB = 1 }) => {
    const { t } = useTranslation('marketplace');
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState('');
    const [progress, setProgress] = useState(0);

    const processImage = async (file) => {
        const compressionOptions = {
            maxSizeMB: maxSizeMB,
            maxWidthOrHeight: 1920,
            useWebWorker: true,
            fileType: 'image/jpeg',
            onProgress: (p) => setProgress(Math.round(p * 100))
        };

        try {
            const compressedFile = await imageCompression(file, compressionOptions);
            const watermarkedBlob = await addWatermark(compressedFile);
            
            return new File([watermarkedBlob], file.name, {
                type: 'image/jpeg',
                lastModified: new Date().getTime()
            });
        } catch (error) {
            console.error('Error processing image:', error);
            throw error;
        }
    };

    const handleImageChange = async (event) => {
        const files = Array.from(event.target.files || []);
        setError('');
        setUploading(true);
        setProgress(0);

        try {
            if (files.length > maxImages) {
                setError(t('listings.create.photos.maxCount', { count: maxImages }));
                return;
            }

            const validFiles = files.filter(file => {
                if (!file.type.startsWith('image/')) {
                    setError(t('listings.create.photos.onlyImages'));
                    return false;
                }
                return true;
            });

            if (validFiles.length === 0) {
                setUploading(false);
                return;
            }

            const processPromises = validFiles.map(async (file) => {
                const processedFile = await processImage(file);
                return {
                    file: processedFile,
                    preview: URL.createObjectURL(processedFile)
                };
            });

            const processedImages = await Promise.all(processPromises);
            onImagesSelected(processedImages);

        } catch (error) {
            console.error('Error processing images:', error);
            setError(t('listings.create.photos.processingError'));
        } finally {
            setUploading(false);
            setProgress(0);
            event.target.value = '';
        }
    };

    return (
        <Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <Button
                    id="loadPhotoButton"
                    variant="contained"
                    component="label"
                    startIcon={uploading ? <CircularProgress size={20} /> : <CloudUploadIcon />}
                    disabled={uploading}
                >
                    {uploading 
                        ? t('listings.create.photos.processing', { progress })
                        : t('listings.create.photos.upload')
                    }
                    <input
                        type="file"
                        hidden
                        multiple
                        accept="image/*"
                        onChange={handleImageChange}
                    />
                </Button>
                {uploading && (
                    <Typography variant="body2" color="text.secondary">
                        {t('listings.create.photos.addWatermark')}
                    </Typography>
                )}
            </Box>
            {error && (
                <Typography color="error" variant="body2" sx={{ mt: 1 }}>
                    {error}
                </Typography>
            )}
        </Box>
    );
};

export default ImageUploader;