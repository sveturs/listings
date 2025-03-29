// frontend/hostel-frontend/src/components/store/CategoryMappingModal.jsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Modal,
  Box,
  Paper
} from '@mui/material';
import CategoryMappingEditor from './CategoryMappingEditor';

const CategoryMappingModal = ({ open, onClose, sourceId, onSave }) => {
  const { t } = useTranslation(['common', 'marketplace']);

  return (
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby="category-mapping-modal-title"
    >
      <Paper
        sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: { xs: '90%', sm: '80%', md: '70%' },
          bgcolor: 'background.paper',
          boxShadow: 24,
          p: 4,
          maxHeight: '90vh',
          overflow: 'auto'
        }}
      >
        <CategoryMappingEditor 
          sourceId={sourceId} 
          onClose={onClose} 
          onSave={onSave} 
        />
      </Paper>
    </Modal>
  );
};

export default CategoryMappingModal;