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

declare module '@mui/material/List' {
  interface ListProps {
    children?: React.ReactNode;
  }
}

declare module '@mui/material/ListItem' {
  interface ListItemProps {
    children?: React.ReactNode;
  }
}

// Extend MUI styled components
declare module '@mui/material/styles' {
  interface ComponentsPropsList {
    MuiTooltip: TooltipProps;
    MuiModal: ModalProps;
    MuiInputLabel: InputLabelProps;
  }

  export interface Theme {
    breakpoints: {
      down: (key: string) => string;
      up: (key: string) => string;
    };
    spacing: (factor: number) => number | string;
    palette: {
      primary: {
        light: string;
        main: string;
        dark: string;
        contrastText: string;
      };
      secondary: {
        light: string;
        main: string;
        dark: string;
        contrastText: string;
      };
      error: {
        light: string;
        main: string;
        dark: string;
        contrastText: string;
      };
      warning: {
        light: string;
        main: string;
        dark: string;
        contrastText: string;
      };
      info: {
        light: string;
        main: string;
        dark: string;
        contrastText: string;
      };
      success: {
        light: string;
        main: string;
        dark: string;
        contrastText: string;
      };
      action: {
        hover: string;
        disabled: string;
      };
      divider: string;
      text: {
        primary: string;
        secondary: string;
        disabled: string;
      };
      grey: {
        50: string;
        100: string;
        200: string;
        300: string;
        400: string;
        500: string;
        600: string;
        700: string;
        800: string;
        900: string;
        A100: string;
        A200: string;
        A400: string;
        A700: string;
      };
      background: {
        default: string;
        paper: string;
      };
    };
    shape: {
      borderRadius: number;
    };
  }

  // Styled component definitions
  export interface StyledComponent<P = {}> extends React.ComponentType<P> {
    (props: P): JSX.Element;
    propTypes?: any;
    displayName?: string;
  }

  export function styled<C extends React.ComponentType<React.ComponentProps<C>>, P = {}>(
    Component: C,
    options?: any
  ): StyledComponent<P & Omit<React.ComponentProps<C>, keyof P>>;
}