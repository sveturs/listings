// src/components/marketplace/ImageUploader.js
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import imageCompression from 'browser-image-compression';
import { Box, Button, Typography, CircularProgress, Alert } from '@mui/material';
import { CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import axios from '../../api/axios';

const ImageUploader = ({ onImagesSelected, maxImages = 10, maxSizeMB = 1 }) => {
  const { t } = useTranslation('marketplace');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');
  const [progress, setProgress] = useState(0);

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

      console.log("Начинаем обработку", validFiles.length, "изображений");

      const processPromises = validFiles.map(async (file, index) => {
        try {
          // Сжатие изображения перед загрузкой
          const compressionOptions = {
            maxSizeMB: maxSizeMB,
            maxWidthOrHeight: 1920,
            useWebWorker: true,
            fileType: 'image/jpeg',
            onProgress: (p) => setProgress(Math.round(p * 100))
          };

          console.log(`Сжимаем изображение ${index + 1}/${validFiles.length}`);
          const compressedFile = await imageCompression(file, compressionOptions);
          console.log(`Изображение ${index + 1} сжато: ${compressedFile.size} байт`);

          // Создаем объект Blob с корректным типом
          const blob = new Blob([compressedFile], { type: 'image/jpeg' });
          const newFile = new File([blob], `image-${index}.jpg`, { type: 'image/jpeg' });

          return {
            file: newFile,
            preview: URL.createObjectURL(newFile),
            isMain: index === 0
          };
        } catch (error) {
          console.error('Error processing image:', error);
          return null;
        }
      });

      const results = await Promise.all(processPromises);
      const filteredResults = results.filter(Boolean);
      console.log("Обработано изображений:", filteredResults.length);

      onImagesSelected(filteredResults);

    } catch (error) {
      console.error('Error processing images:', error);
      setError(t('listings.create.photos.processingError'));
    } finally {
      setUploading(false);
      setProgress(0);
      event.target.value = '';
    }
  };

  // Создаем ссылку на скрытый input
  const fileInputRef = React.useRef(null);

  // Функция для открытия окна выбора файлов
  const handleButtonClick = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
        {/* Скрытый input для выбора файлов */}
        <input
          type="file"
          ref={fileInputRef}
          style={{ display: 'none' }}
          multiple
          accept="image/*"
          onChange={handleImageChange}
        />

        {/* Кнопка для открытия окна выбора файлов */}
        <Button
          id="loadPhotoButton"
          variant="contained"
          onClick={handleButtonClick}
          startIcon={uploading ? <CircularProgress size={20} /> : <CloudUploadIcon />}
          disabled={uploading}
        >
          {uploading
            ? t('listings.create.photos.processing', { progress })
            : t('listings.create.photos.upload')
          }
        </Button>
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