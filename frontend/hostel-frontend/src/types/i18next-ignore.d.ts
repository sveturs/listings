// This file tells TypeScript to completely ignore the type checking for i18next
declare module "i18next/typescript/t" {
  export const t: any;
}

// Override the problematic type definition
declare module "i18next" {
  export function t(key: string, options?: any): string;
}