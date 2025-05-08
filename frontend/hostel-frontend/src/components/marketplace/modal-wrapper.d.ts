// modal-wrapper.d.ts
import React from 'react';
import { ModalProps } from '@mui/material/Modal';

// This declaration augments Material UI's Modal to make 'children' optional
declare module '@mui/material/Modal' {
  interface ModalProps {
    children?: React.ReactNode;
  }
}