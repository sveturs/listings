// src/components/marketplace/ModalWrapper.tsx
import React from 'react';
import { Modal, Box, Typography } from '@mui/material';

interface ModalWrapperProps {
  open: boolean;
  onClose: () => void;
  'aria-labelledby'?: string;
  children: React.ReactNode;
  isMobile?: boolean;
  title?: React.ReactNode;
}

const ModalWrapper: React.FC<ModalWrapperProps> = ({
  open,
  onClose,
  'aria-labelledby': ariaLabelledBy,
  children,
  isMobile = false,
  title
}) => {
  return (
    // @ts-ignore - Ignoring TypeScript error for Modal component
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby={ariaLabelledBy}
    >
      <Box sx={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        width: isMobile ? '90%' : 600,
        bgcolor: 'background.paper',
        borderRadius: 2,
        boxShadow: 24,
        p: 4,
      }}>
        {title && (
          <Typography variant="h6" component="h2" gutterBottom>
            {title}
          </Typography>
        )}
        {children}
      </Box>
    </Modal>
  );
};

export default ModalWrapper;