// src/components/marketplace/ImageUploader.tsx
import React, { useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import imageCompression from 'browser-image-compression';
import { Box, Button, Typography, CircularProgress } from '@mui/material';
import { CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import ImageEnhancementOffer from './ImageEnhancementOffer';

// Type definitions for browser-image-compression
interface CompressionOptions {
  maxSizeMB: number;
  maxWidthOrHeight: number;
  useWebWorker: boolean;
  fileType: string;
  onProgress: (progress: number) => void;
}

export interface ProcessedImage {
  file: File | string;
  preview: string;
  isMain?: boolean;
  cloudinary_id?: string;
}

interface ImageUploaderProps {
  onImagesSelected: (images: ProcessedImage[]) => void;
  maxImages?: number;
  maxSizeMB?: number;
  showEnhancementOffer?: boolean;
}

const ImageUploader: React.FC<ImageUploaderProps> = ({ 
  onImagesSelected, 
  maxImages = 10, 
  maxSizeMB = 1,
  showEnhancementOffer = false
}) => {
  const { t } = useTranslation('marketplace') as any;
  const [uploading, setUploading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [progress, setProgress] = useState<number>(0);
  const [processedImages, setProcessedImages] = useState<ProcessedImage[]>([]);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const handleImageChange = async (event: React.ChangeEvent<HTMLInputElement>): Promise<void> => {
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
        if (!(file as File).type.startsWith('image/')) {
          setError(t('listings.create.photos.onlyImages'));
          return false;
        }
        return true;
      }) as File[];

      if (validFiles.length === 0) {
        setUploading(false);
        return;
      }

      console.log("Начинаем обработку", validFiles.length, "изображений");

      const processPromises = validFiles.map(async (file, index) => {
        try {
          // Сжатие изображения перед загрузкой
          const compressionOptions: CompressionOptions = {
            maxSizeMB: maxSizeMB,
            maxWidthOrHeight: 1920,
            useWebWorker: true,
            fileType: 'image/jpeg',
            onProgress: (p: number) => setProgress(Math.round(p * 100))
          };

          console.log(`Сжимаем изображение ${index + 1}/${validFiles.length}`);
          const compressedFile = await imageCompression(file as File, compressionOptions);
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
      const filteredResults = results.filter(Boolean) as ProcessedImage[];
      console.log("Обработано изображений:", filteredResults.length);

      setProcessedImages(filteredResults);
      onImagesSelected(filteredResults);

    } catch (error) {
      console.error('Error processing images:', error);
      setError(t('listings.create.photos.processingError'));
    } finally {
      setUploading(false);
      setProgress(0);
      if (event.target) {
        event.target.value = '';
      }
    }
  };

  // Функция для открытия окна выбора файлов
  const handleButtonClick = (): void => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  // Обработчик для улучшенных изображений
  const handleEnhancedImages = (enhancedImages: ProcessedImage[]): void => {
    setProcessedImages(enhancedImages);
    onImagesSelected(enhancedImages);
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

      {showEnhancementOffer && processedImages.length > 0 && (
        <ImageEnhancementOffer 
          images={processedImages} 
          onEnhanced={handleEnhancedImages} 
        />
      )}
    </Box>
  );
};

export default ImageUploader;