# TypeScript Migration Guide

This document outlines the plan and guidelines for migrating the Hostel Booking System frontend from JavaScript to TypeScript.

## Current Progress

The following components have been migrated to TypeScript:

- `src/components/icons/SveTuLogo.tsx`
- `src/components/shared/LanguageSwitcher.tsx`
- `src/components/shared/AutocompleteInput.tsx`

## Migration Strategy

We're following a gradual migration approach, converting components one by one while maintaining compatibility with the existing JavaScript codebase. This approach allows us to:

1. Migrate components incrementally without disrupting the application
2. Learn from each migration to improve subsequent migrations
3. Validate that each converted component works correctly before moving on

## Migration Steps for Each Component

1. **Create TypeScript interfaces** for component props and state
2. **Create a new TypeScript file** with the same name but `.tsx` extension
3. **Convert JavaScript to TypeScript** and add appropriate type annotations
4. **Update import statements** in files that use the converted component
5. **Test the component** to ensure it works correctly

## TypeScript Configuration

The project is configured with the following key TypeScript settings:

```json
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": false,
    "forceConsistentCasingInFileNames": true,
    "noFallthroughCasesInSwitch": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve"
  },
  "include": ["src"]
}
```

Key points:
- `allowJs: true` - Allows JavaScript files to be imported
- `strict: false` - For now, we're using a more relaxed type checking to ease migration
- `jsx: "preserve"` - Preserves JSX for Babel to process

## Common Type Definitions

We've created common type definitions in `src/types/index.d.ts` for:

- Image imports (PNG, JPG, etc.)
- SVG files as React components
- CSS modules
- JSON imports
- Custom window properties

Additionally, for i18next, we've added supplementary type definitions in `src/types/i18next.d.ts`.

## Build and Development Setup

When running the development server or building the project, use the following scripts to handle Node.js OpenSSL issues:

- For development: `./dev.sh` - Sets `NODE_OPTIONS=--openssl-legacy-provider` and runs `npm start`
- For building: `./build.sh` - Sets `NODE_OPTIONS=--openssl-legacy-provider` and runs `npm run build`

## Migration Priority

We recommend migrating components in the following order:

1. Simple, self-contained components (like icons, buttons)
2. Shared components that are used in multiple places
3. Page-specific components
4. Context providers and hooks
5. Utility functions

## Best Practices for TypeScript Migration

1. **Start with interfaces**: Define the component props and any other types/interfaces first
2. **Use React's built-in types**: Leverage `React.FC<Props>`, `React.ChangeEvent<HTMLInputElement>`, etc.
3. **Be explicit with event types**: Use specific event types like `React.MouseEvent` instead of `any`
4. **Add utility types for API responses**: Create interfaces for API responses and data structures
5. **Comment your types**: Add JSDoc comments to complex types for better documentation
6. **Use union types** for props that can accept multiple value types
7. **Be careful with null/undefined**: Use optional chaining and nullish coalescing where appropriate

## Examples

### Component Props Interface
```typescript
interface ButtonProps {
  text: string;
  onClick: () => void;
  variant?: 'primary' | 'secondary' | 'tertiary';
  disabled?: boolean;
  size?: 'small' | 'medium' | 'large';
}
```

### Function Component
```typescript
const Button: React.FC<ButtonProps> = ({ 
  text, 
  onClick, 
  variant = 'primary', 
  disabled = false,
  size = 'medium' 
}) => {
  // Component implementation
};
```

### Event Handling
```typescript
const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
  const newValue = e.target.value;
  // Handler implementation
};
```

## Troubleshooting

### Type Errors in Dependencies
If you encounter type errors in dependencies (like i18next), you can:
1. Create custom type definitions in `src/types/`
2. Use the `skipLibCheck: true` option in tsconfig.json
3. Use type assertions where necessary (`as any` or more specific types)

### JSX Compilation Issues
Ensure that the `jsx` option in tsconfig.json is set to `"preserve"` for use with Babel.

### Build Issues
Use the `NODE_OPTIONS=--openssl-legacy-provider` environment variable when building with Node.js 18 or higher to solve OpenSSL-related issues.

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [TypeScript with React](https://www.typescriptlang.org/docs/handbook/react.html)