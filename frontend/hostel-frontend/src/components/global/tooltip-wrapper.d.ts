// tooltip-wrapper.d.ts
import React from 'react';

// This workaround makes Material UI's Tooltip accept a span as children
declare module '@mui/material/Tooltip' {
  interface TooltipProps {
    children: React.ReactElement;
  }
}