import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        background: 'var(--background)',
        foreground: 'var(--foreground)',
      },
    },
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        light: {
          primary: '#570df8',
          secondary: '#f000b8',
          accent: '#37cdbe',
          neutral: '#3d4451',
          'base-100': '#ffffff',
          info: '#3abff8',
          success: '#36d399',
          warning: '#fbbd23',
          error: '#f87272',
        },
      },
      {
        dark: {
          primary: '#793ef9',
          secondary: '#f000b8',
          accent: '#37cdbe',
          neutral: '#2a2e37',
          'base-100': '#1f2937',
          'base-200': '#111827',
          'base-300': '#0f172a',
          info: '#3abff8',
          success: '#36d399',
          warning: '#fbbd23',
          error: '#f87272',
        },
      },
    ],
  },
} as any;

export default config;
