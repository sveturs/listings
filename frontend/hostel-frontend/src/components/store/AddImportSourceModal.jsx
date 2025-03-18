// frontend/hostel-frontend/src/components/store/AddImportSourceModal.jsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Modal,
  Box,
  Paper
} from '@mui/material';
import ImportSourceForm from './ImportSourceForm';

const AddImportSourceModal = ({ open, onClose, onSuccess, storefrontId, initialData = null }) => {
  const { t } = useTranslation(['common', 'marketplace']);

  return (
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby="add-import-source-modal-title"
    >
      <Paper
        sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: { xs: '90%', sm: 500 },
          maxHeight: '90vh',
          overflow: 'auto'
        }}
      >
        <Box p={3}>
          <ImportSourceForm 
            onClose={onClose} 
            onSuccess={onSuccess} 
            storefrontId={storefrontId} 
            initialData={initialData}
          />
        </Box>
      </Paper>
    </Modal>
  );
};

export default AddImportSourceModal;