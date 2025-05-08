# TypeScript Migration Guide

This document outlines the plan and guidelines for migrating the Hostel Booking System frontend from JavaScript to TypeScript.

## Current Progress

The following components have been migrated to TypeScript:

- `src/components/icons/SveTuLogo.tsx` - SVG logo component with animation
- `src/components/shared/LanguageSwitcher.tsx` - Language selection component
- `src/components/shared/AutocompleteInput.tsx` - Search input with autocompletion
- `src/components/user/UserProfile.tsx` - User profile component with form handling
- `src/components/global/CitySelector.tsx` - City selector with geolocation and search
- `src/components/global/NewMessageIndicator.tsx` - Notification badge for new messages
- `src/components/global/PrivateRoute.tsx` - Protected route wrapper for authentication
- `src/components/global/AdminRoute.tsx` - Admin route protection wrapper
- `src/contexts/AuthContext.tsx` - Authentication context provider and hook
- `src/contexts/LanguageContext.tsx` - Language selection context provider
- `src/contexts/LocationContext.tsx` - Geolocation context provider
- `src/contexts/ChatContext.tsx` - Chat messaging context provider
- `src/components/marketplace/chat/ChatService.ts` - WebSocket chat service implementation
- `src/components/marketplace/ShareButton.tsx` - Social media sharing component
- `src/components/marketplace/CallButton.tsx` - Phone call button with phone display
- `src/components/marketplace/Breadcrumbs.tsx` - Navigation breadcrumbs with category paths
- `src/components/marketplace/CategorySelect.tsx` - Category selection dropdown component
- `src/components/marketplace/InfiniteScroll.tsx` - Infinite scrolling loader component
- `src/components/marketplace/ListingCard.tsx` - Marketplace listing card component
- `src/components/marketplace/SimilarListings.tsx` - Similar listings component with lazy loading
- `src/components/marketplace/PriceHistoryChart.tsx` - Price history visualization chart component
- `src/components/marketplace/CategoryTree.tsx` - Hierarchical category tree navigation component
- `src/components/marketplace/AttributeFields.tsx` - Dynamic form fields for category attributes
- `src/components/marketplace/ItemDetails.tsx` - Detailed listing view with image gallery
- `src/components/marketplace/AutoDetails.tsx` - Vehicle details display component
- `src/components/marketplace/PhonePopup.tsx` - Phone number popup display component
- `src/components/marketplace/ImageEnhancementOffer.tsx` - Image enhancement service offer component
- `src/components/marketplace/ImageUploader.tsx` - Image upload and compression component
- `src/components/marketplace/VirtualizedCategoryTree.tsx` - Virtual scroll category tree component
- `src/components/marketplace/CentralAttributeFilters.tsx` - Centralized attribute filters component
- `src/components/marketplace/AttributeFilters.tsx` - Attribute-based filter controls component
- `src/components/marketplace/MarketplaceListingsList.tsx` - Sortable and filterable marketplace listings table component
- `src/components/marketplace/MobileComponents.tsx` - Mobile-optimized UI components for marketplace
- `src/components/marketplace/MarketplaceFilters.tsx` - Marketplace filtering and search component
- `src/hooks/useTranslatedContent.tsx` - Hook for automatic content translation

<!-- Новые компоненты (Май 2025) -->
- `src/components/marketplace/Breadcrumbs.tsx` - Navigation breadcrumbs with category paths
- `src/components/marketplace/CategorySelect.tsx` - Category selection dropdown component  
- `src/components/marketplace/InfiniteScroll.tsx` - Infinite scrolling loader component
- `src/components/marketplace/SimilarListings.tsx` - Similar listings component with lazy loading
- `src/components/marketplace/PriceHistoryChart.tsx` - Price history visualization chart component
- `src/components/marketplace/CategoryTree.tsx` - Hierarchical category tree navigation component
- `src/components/marketplace/AttributeFields.tsx` - Dynamic form fields for category attributes
- `src/components/marketplace/ItemDetails.tsx` - Detailed listing view with image gallery
- `src/components/marketplace/AutoDetails.tsx` - Vehicle details display component
- `src/components/marketplace/PhonePopup.tsx` - Phone number popup display component
- `src/components/marketplace/ImageEnhancementOffer.tsx` - Image enhancement service offer component
- `src/components/marketplace/ImageUploader.tsx` - Image upload and compression component
- `src/components/marketplace/VirtualizedCategoryTree.tsx` - Virtual scroll category tree component
- `src/components/marketplace/CentralAttributeFilters.tsx` - Centralized attribute filters component
- `src/components/marketplace/AttributeFilters.tsx` - Attribute-based filter controls component

## Latest Component Updates

### New TypeScript Components Converted (May 2025)

#### CategorySelect.tsx
- Converted to TypeScript with proper interfaces and type definitions
- Added proper typing for category objects including translations
- Added types for event handlers

#### InfiniteScroll.tsx
- Implemented TypeScript interfaces for component props
- Added proper typing for Intersection Observer
- Added types for refs and state

#### Breadcrumbs.tsx
- Created interfaces for category paths and props
- Added type safety to navigation and event handling
- Added types for translated content

#### SimilarListings.tsx
- Converted component to TypeScript with comprehensive interfaces
- Added proper typing for listings and image handling
- Implemented typed event handlers for load more functionality
- Reused existing Listing interface from ListingCard component
- Added proper typing for async data fetching

#### PriceHistoryChart.tsx
- Converted chart component to TypeScript with proper interfaces for price history data
- Added typed props for Recharts components including custom tooltip
- Implemented proper typing for async API calls and data transformation
- Added return type annotations for functions like formatPrice
- Created proper interfaces for component props and internal data structures

#### CategoryTree.tsx
- Converted tree component to TypeScript with proper interfaces for category data
- Extended base Category interface with tree-specific properties
- Added proper typing for recursive tree rendering functions
- Implemented type-safe event handlers
- Added proper return type annotations for utility functions

#### AttributeFields.tsx
- Converted complex dynamic form fields component to TypeScript
- Created comprehensive interfaces for attribute definitions and values
- Implemented type safety for various input field types
- Added proper typing for attribute transformations and validations
- Created type-safe event handlers for all form controls

#### ItemDetails.tsx
- Converted product details component with image gallery to TypeScript
- Reused Listing and ImageObject interfaces from ListingCard
- Added proper typing for component state management
- Implemented type-safe image URL handling logic
- Added null-safety with optional chaining

#### AutoDetails.tsx
- Converted vehicle details component to TypeScript
- Created comprehensive interface for auto properties
- Added proper typing for translation functions
- Implemented type-safe rendering with nullable properties
- Added return type annotations for all utility functions

#### PhonePopup.tsx
- Converted phone number popup component to TypeScript
- Added interface for component props
- Added proper typing for canvas ref using HTMLCanvasElement
- Implemented null safety for context initialization
- Added type-safe animation with Material UI keyframes

#### ImageEnhancementOffer.tsx
- Converted image enhancement service component to TypeScript
- Created interfaces for image objects and API responses
- Added proper typing for async functions with Promise<void> return types
- Implemented type-safe handling of API errors
- Added proper typing for form data and file uploads

#### ImageUploader.tsx
- Converted image upload component to TypeScript
- Created reusable ProcessedImage interface exported for other components
- Added proper typing for File objects and image compression
- Added event handler types for input events
- Integrated with ImageEnhancementOffer component
- Added type-safe refs with useRef<HTMLInputElement>

#### VirtualizedCategoryTree.tsx
- Converted virtualized tree component to TypeScript with proper interfaces
- Reused existing Category interface from HierarchicalCategorySelect
- Added proper typing for react-window List and ListChildComponentProps
- Implemented type-safe callback functions with useCallback
- Added type safety for Sets and Maps used for tracking expanded items
- Improved translation function with proper type checking

#### CentralAttributeFilters.tsx
- Converted attribute filters component to TypeScript with proper interfaces
- Created a DebugData interface for API response data structure
- Implemented type-safe event handlers for button clicks
- Added proper typing for useCallback and useEffect hooks
- Added type annotations for all state variables
- Added null safety checks for filter handling

#### AttributeFilters.tsx
- Converted complex filter component to TypeScript with comprehensive interfaces
- Created detailed interfaces for attribute structures, options, and translations
- Implemented strict typing for attribute ranges and option handling
- Added proper event typing for form controls and sliders
- Improved type safety for async data fetching and processing
- Added type-safe callbacks with proper parameter types

#### MarketplaceListingsList.tsx
- Converted complex marketplace listings table component to TypeScript
- Created extensive interfaces for listings, images, filters, and discount information
- Added proper typing for sorting and filtering functionality
- Implemented type-safe event handlers for row and checkbox clicks
- Added comprehensive typing for ratings and reviews data fetching
- Created proper typing for localization and column display logic
- Implemented typed utility functions for formatting prices and dates
- Added type-safe image URL processing with proper null handling

#### MobileComponents.tsx
- Converted mobile UI components to TypeScript with comprehensive interfaces
- Created multiple type-safe components in one file: MobileHeader, MobileListingCard, MobileListingGrid, and MobileFilters
- Added proper typing for event handlers with MouseEvent and TouchEvent types
- Implemented type-safe state management with useState<T> generics
- Created reusable interfaces for listings, categories, and filter options
- Added proper typing for UI-specific features like view modes ('grid' | 'list')
- Implemented proper type safety for browser detection and conditional rendering
- Added internationalization support with typed translation functions
- Created proper typing for touch events and click handling
- Improved type safety for dynamic data structures with optional chaining

#### MarketplaceFilters.tsx
- Converted marketplace filtering component to TypeScript with typed state and callbacks
- Created interfaces for filter options and attribute filters
- Added type safety for SelectChangeEvent and form events
- Implemented proper typing for async geolocation functions with Promise<void>
- Used typed generic React.FC<MarketplaceFiltersProps> for component definition
- Added proper type safety for Location context usage
- Implemented proper typing for callback functions with useCallback
- Added type-safe boolean checks for conditional rendering
- Added proper typing for i18n translations with string literals
- Created type-safe memoization with useMemo

### Fixed TypeScript Errors in ListingCard.tsx

- Added interfaces for listings, attributes, and other data structures
- Fixed Modal component to properly handle children with correct JSX structure
- Added proper event typing for all click handlers
- Added typings for image objects and discount data

### Fixed TypeScript Errors in CitySelector.tsx

- Fixed issue with Tooltip component requiring children by adding placement prop
- Fixed auth check by deriving isAuthenticated from user presence

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

### Window Properties
When accessing custom window properties (like `window.ENV`), use type assertions to avoid TypeScript errors:
```typescript
const backendUrl = (window as any).ENV?.REACT_APP_BACKEND_URL || '';
```

Or define a global interface extension in a `.d.ts` file:
```typescript
declare global {
  interface Window {
    ENV?: {
      REACT_APP_BACKEND_URL?: string;
      [key: string]: string | undefined;
    };
  }
}
```

### i18next and react-i18next
When working with i18next translations in TypeScript, there are several approaches:

1. **Use type assertions for simple components**:
```typescript
const { t } = useTranslation('namespace') as any;
```

2. **Define custom type declarations**:
```typescript
// src/types/i18next.d.ts
import 'i18next';

declare module 'i18next' {
  interface i18n {
    changeLanguage(lng?: string): Promise<TFunction>;
    language: string;
  }

  interface TFunction {
    (key: string | string[], options?: any): string;
    (key: string | string[], defaultValue: string, options?: any): string;
  }
}
```

3. **Use skipLibCheck flag for troublesome type checking**:
When type checking i18next in isolated components, use:
```bash
npx tsc --skipLibCheck your-component.tsx
```

### Import Path Extensions
When importing TypeScript files in other TypeScript files, don't include the `.tsx` extension:
```typescript
// Incorrect
import { useAuth } from '../../contexts/AuthContext.tsx';

// Correct
import { useAuth } from '../../contexts/AuthContext';
```

### Build Issues
Use the `NODE_OPTIONS=--openssl-legacy-provider` environment variable when building with Node.js 18 or higher to solve OpenSSL-related issues.

### Material UI Component Issues

#### Modal Component
When using Modal component from Material UI, the TypeScript typings can be problematic. There are several workarounds:

1. **Using a wrapper component**
```typescript
// Create a wrapper component that handles the Modal and its children
const ModalWrapper: React.FC<Props> = ({ open, onClose, children }) => {
  return (
    // @ts-ignore - Ignoring TypeScript errors for Modal
    <Modal open={open} onClose={onClose}>
      <Box sx={{ /* your styles */ }}>
        {children}
      </Box>
    </Modal>
  );
};
```

2. **Using @ts-ignore**
```typescript
// Using @ts-ignore directive to bypass TypeScript checking
// @ts-ignore
<Modal
  open={isOpen}
  onClose={handleClose}
>
  <Box>Content</Box>
</Modal>
```

3. **Creating custom type declarations**
```typescript
// In a .d.ts file:
declare module '@mui/material/Modal' {
  interface ModalProps {
    children?: React.ReactNode;
  }
}
```

#### Tooltip Component
When using Tooltip component, there are multiple approaches to handle type issues:

1. **Using @ts-ignore**
```typescript
// @ts-ignore
<Tooltip title="Tooltip text" arrow placement="bottom">
  <Button>Click me</Button>
</Tooltip>
```

2. **Wrapping children in a span or div**
```typescript
<Tooltip title="Tooltip text" arrow placement="bottom">
  <span>
    <Button>Click me</Button>
  </span>
</Tooltip>
```

3. **Using Box with component="span"**
```typescript
<Tooltip title="Tooltip text" arrow placement="bottom">
  <Box component="span">
    <Button>Click me</Button>
  </Box>
</Tooltip>
```

4. **Creating custom type declarations**
```typescript
// In a .d.ts file:
declare module '@mui/material/Tooltip' {
  interface TooltipProps {
    children?: React.ReactElement;
  }
}
```

### Running TypeScript Check on Components

We've created a custom script `run-tsc.sh` that handles common TypeScript errors:

```bash
# Check specific files
./run-tsc.sh src/components/path/to/file.tsx

# Check all TypeScript files
./run-tsc.sh
```

This script includes flags to:
- Skip problematic library type checks
- Handle JSX properly
- Resolve module resolution issues
- Allow synthetic default imports

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [TypeScript with React](https://www.typescriptlang.org/docs/handbook/react.html)
- [Material UI with TypeScript](https://mui.com/material-ui/guides/typescript/)