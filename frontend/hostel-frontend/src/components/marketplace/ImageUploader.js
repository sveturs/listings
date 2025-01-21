import React, { useState } from 'react';
import imageCompression from 'browser-image-compression';
import { Box, Button, IconButton, Typography } from '@mui/material';
import { CloudUpload as CloudUploadIcon, Delete as DeleteIcon } from '@mui/icons-material';

const ImageUploader = ({ onImagesSelected, maxImages = 10, maxSizeMB = 1 }) => {
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState('');

    const compressImage = async (file) => {
        const options = {
            maxSizeMB: maxSizeMB,
            maxWidthOrHeight: 1920,
            useWebWorker: true,
            fileType: file.type,
        };

        try {
            return await imageCompression(file, options);
        } catch (error) {
            console.error('Error compressing image:', error);
            throw error;
        }
    };

    const handleImageChange = async (event) => {
        const files = Array.from(event.target.files || []);
        setError('');
        setUploading(true);

        try {
            if (files.length > maxImages) {
                setError(`Максимальное количество фотографий: ${maxImages}`);
                return;
            }

            const validFiles = files.filter(file => {
                if (!file.type.startsWith('image/')) {
                    setError('Можно загружать только изображения');
                    return false;
                }
                return true;
            });

            if (validFiles.length === 0) {
                setUploading(false);
                return;
            }

            const compressPromises = validFiles.map(async (file) => {
                const compressedFile = await compressImage(file);
                return {
                    file: compressedFile,
                    preview: URL.createObjectURL(compressedFile)
                };
            });

            const compressedImages = await Promise.all(compressPromises);
            onImagesSelected(compressedImages);

        } catch (error) {
            console.error('Error processing images:', error);
            setError('Ошибка при обработке изображений');
        } finally {
            setUploading(false);
            event.target.value = ''; // Сброс input для возможности повторной загрузки
        }
    };

    return (
        <Box>
            <Button
                variant="contained"
                component="label"
                startIcon={<CloudUploadIcon />}
                disabled={uploading}
            >
                {uploading ? 'Обработка...' : 'Загрузить фото'}
                <input
                    type="file"
                    hidden
                    multiple
                    accept="image/*"
                    onChange={handleImageChange}
                />
            </Button>
            {error && (
                <Typography color="error" variant="body2" sx={{ mt: 1 }}>
                    {error}
                </Typography>
            )}
        </Box>
    );
};

export default ImageUploader;