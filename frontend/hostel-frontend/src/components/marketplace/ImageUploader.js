import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import imageCompression from 'browser-image-compression';
import { Box, Button, Typography, CircularProgress, Alert } from '@mui/material';
import { CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import { addWatermark } from '../../utils/imageUtils';
import ImageEnhancementOffer from './ImageEnhancementOffer';
import axios from '../../api/axios';

const ImageUploader = ({ onImagesSelected, maxImages = 10, maxSizeMB = 1 }) => {
  const { t } = useTranslation('marketplace');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');
  const [progress, setProgress] = useState(0);
  const [processedImages, setProcessedImages] = useState([]);
  const [moderationWarnings, setModerationWarnings] = useState([]);

  const processImage = async (file) => {
    const formData = new FormData();
    formData.append('image', file);

    try {
      const moderationResponse = await axios.post('/api/v1/marketplace/moderate-image', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });

      if (!moderationResponse.data.data.safe) {
        setModerationWarnings(prev => [
          ...prev,
          {
            filename: file.name,
            reason: moderationResponse.data.data.reason || 'Запрещённый контент'
          }
        ]);
        return null;
      }

      const compressionOptions = {
        maxSizeMB: maxSizeMB,
        maxWidthOrHeight: 1920,
        useWebWorker: true,
        fileType: 'image/jpeg',
        onProgress: (p) => setProgress(Math.round(p * 100))
      };

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
    setModerationWarnings([]);

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
        if (!processedFile) return null;

        return {
          file: processedFile,
          preview: URL.createObjectURL(processedFile)
        };
      });

      const results = await Promise.all(processPromises);
      const filteredResults = results.filter(Boolean);

      setProcessedImages(filteredResults);
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

  const handleEnhancedImages = (enhancedImages) => {
    setProcessedImages(enhancedImages);
    onImagesSelected(enhancedImages);
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
      </Box>

      {error && (
        <Typography color="error" variant="body2" sx={{ mt: 1 }}>
          {error}
        </Typography>
      )}

      {moderationWarnings.length > 0 && (
        <Alert severity="error" sx={{ mt: 2 }}>
          <Typography variant="subtitle2">
            Некоторые фотографии были отклонены:
          </Typography>
          <ul>
            {moderationWarnings.map((warning, idx) => (
              <li key={idx}>
                {warning.filename}: {warning.reason}
              </li>
            ))}
          </ul>
        </Alert>
      )}

      {processedImages.length > 0 && (
        <ImageEnhancementOffer
          images={processedImages}
          onEnhanced={handleEnhancedImages}
        />
      )}
    </Box>
  );
};

export default ImageUploader;
