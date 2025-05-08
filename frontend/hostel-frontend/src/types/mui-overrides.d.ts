// src/types/mui-overrides.d.ts
import React from 'react';

declare module '@mui/material/Tooltip' {
  interface TooltipProps {
    children?: React.ReactElement;
  }
}

declare module '@mui/material/Modal' {
  interface ModalProps {
    children?: React.ReactNode;
  }
}

declare module '@mui/material/InputLabel' {
  interface InputLabelProps {
    children?: React.ReactNode;
  }
}

// Extend all MUI components to allow an 'as' prop for polymorphic components
declare module '@mui/material/styles' {
  interface ComponentsPropsList {
    MuiTooltip: TooltipProps;
    MuiModal: ModalProps;
    MuiInputLabel: InputLabelProps;
  }
}