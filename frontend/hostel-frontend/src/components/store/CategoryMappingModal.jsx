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
          width: { xs: '95%', sm: '90%', md: '80%', lg: '70%' },
          maxWidth: '1000px',
          bgcolor: 'background.paper',
          boxShadow: 24,
          p: { xs: 2, sm: 3, md: 4 },
          maxHeight: '90vh',
          overflow: 'auto',
          borderRadius: 2
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