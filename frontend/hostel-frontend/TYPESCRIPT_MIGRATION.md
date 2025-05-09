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
- `src/components/reviews/PhotoViewer.tsx` - Photo viewer with navigation
- `src/components/reviews/ReviewComponents.tsx` - Review form, card, and stats components
- `src/components/reviews/ReviewsSection.tsx` - Review section management component

## Latest Component Updates

### Review Components Converted (May 2025)

#### PhotoViewer.tsx
- Converted from PhotoViewer.js to TypeScript with proper interface
- Added strong typing for component props and state
- Added event type annotations for click handlers
- Implemented type-safe image rendering and navigation

#### ReviewComponents.tsx
- Converted from ReviewComponents.jsx to TypeScript with comprehensive interfaces
- Created detailed interfaces for review data, votes, and responses
- Added proper typing for form handling and event callbacks
- Implemented type-safe internationalization with i18next
- Added type safety for conditional rendering based on mobile detection
- Created proper interfaces for review form data, stats, and card props

#### ReviewsSection.tsx
- Converted from ReviewsSection.js to TypeScript with comprehensive interfaces
- Created detailed interfaces for Review, ReviewStat, ReviewResponse, and other data structures
- Added proper typing for API interactions and state management
- Implemented type-safe event handling for various user actions
- Created strongly typed SnackbarState interface for notification handling
- Added proper return type annotations for all async functions
- Implemented type safety for different vote types with union type 'helpful' | 'not_helpful'

### Chat Components Migration Completed

Chat system components have been fully migrated to TypeScript:

1. **Chat UI Components**:
   - ChatButton.tsx - Button for starting new chat conversations
   - ChatComponents.tsx - Chat interface components (ChatWindow, ChatList, ChatHeader)

2. **Chat Service and Context**:
   - ChatService.ts - WebSocket and HTTP chat communication service
   - ChatContext.tsx - Global chat state and service management

### Global Components Migration in Progress

We've started migrating global components that are used across the application:

1. **Core Layout Component**:
   - Layout.tsx - Main application layout with header and navigation
   - LocationPicker.tsx - Map-based location selector with geocoding

2. **User Interface Components**:
   - CitySelector.tsx - City selector with geolocation and search
   - NewMessageIndicator.tsx - Notification badge for new messages
   
3. **Authentication Components**:
   - PrivateRoute.tsx - Protected route wrapper for authentication
   - AdminRoute.tsx - Admin route protection wrapper

4. **Context Providers**:
   - NotificationContext.tsx - Global notification state management and API integrations

### Balance Components Migration in Progress

Financial components have also been migrated:

1. **Balance Management**:
   - BalanceWidget.tsx - User balance display widget
   - DepositDialog.tsx - Payment deposit dialog with method selection

### Notification Components Migration Completed

Notification system components have been fully migrated to TypeScript:

1. **Notification UI Components**:
   - NotificationBadge.tsx - Notification count indicator with badge
   - NotificationDrawer.tsx - Slide-out drawer to display notifications
   - NotificationSettings.tsx - Settings panel for configuring notification preferences

2. **Notification Hook**:
   - useNotifications.tsx - Custom hook for notification management with TypeScript interfaces

### Map Components Migration Completed

Map components for location display have been migrated to TypeScript:

1. **Map Display Components**:
   - MiniMap.tsx - Compact map component with expandable fullscreen view
   - FullscreenMap.tsx - Large map view with interactive markers and listing previews
   
2. **Map Utilities**:
   - map-constants.ts - Constants for tile providers and attribution
   - leaflet-icons.ts - Custom Leaflet icon configurations and utilities

### Fixed TypeScript Errors in CitySelector.tsx

- Fixed issue with Tooltip component requiring children by adding placement prop
- Fixed auth check by deriving isAuthenticated from user presence

### Fixed TypeScript Errors in MiniMap.tsx

- Completely removed Material UI Modal component and replaced it with Box-based custom modal
- Added proper event handler for modal backdrop click with stopPropagation on content
- Created a lightweight, TypeScript-compatible custom modal implementation that works with strict typing

### Fixed TypeScript Errors in Layout.tsx

- Fixed type comparison between string and number in ChatMessage handling
- Replaced Material UI Modal with conditional rendering to avoid TypeScript errors
- Replaced Slide component with conditional rendering for better TypeScript compatibility
- Improved type safety for event handlers using React.MouseEvent

## Component Migration Progress

### Marketplace Components Migration Completed

All marketplace components have now been successfully migrated to TypeScript:

1. **Core Interface Components**:
   - ListingCard.tsx
   - MapView.tsx
   - MarketplaceListingsList.tsx
   - MarketplaceFilters.tsx
   - MobileComponents.tsx (including MobileHeader, MobileListingCard, MobileListingGrid, and MobileFilters)
   - CategoryMenu.tsx
   - AttributeFilters.tsx
   - CentralAttributeFilters.tsx
   - VirtualizedCategoryTree.tsx

2. **Auxiliary Components**:
   - PhonePopup.tsx
   - ImageEnhancementOffer.tsx
   - ImageUploader.tsx
   - ShareButton.tsx
   - Breadcrumbs.tsx
   - CategorySelect.tsx
   - InfiniteScroll.tsx
   - SimilarListings.tsx
   - PriceHistoryChart.tsx
   - CategoryTree.tsx
   - AttributeFields.tsx
   - ItemDetails.tsx
   - AutoDetails.tsx
   - CallButton.tsx

3. **Common Issues Resolved**:
   - Properly typed API responses with interfaces
   - Fixed array type safety issues with union and intersection types
   - Added proper type predicates for runtime type checking
   - Created comprehensive interfaces for complex nested data structures
   - Addressed Material UI component typing challenges
   - Fixed event handling with proper React event types
   - Added proper typing for internationalization with i18next
   - Implemented type-safe state management with useState<T>
   - Used type assertions to handle dynamic content from backend APIs

The migration has significantly improved code quality and developer experience by catching type-related errors at compile time rather than runtime.

### Review Components Migration Completed

All review components have been successfully migrated to TypeScript:

1. **Core Review Components**:
   - PhotoViewer.tsx - Full-screen image viewer with navigation
   - ReviewComponents.tsx - Container for review form, card, and stats components
   - ReviewsSection.tsx - Review section management with API integration

2. **Review Sub-Components**:
   - ReviewForm - Dynamic form for creating and editing reviews
   - ReviewCard - Individual review display with rating, content, and interaction options
   - RatingStats - Statistical visualization of ratings distribution

3. **Common Patterns**:
   - Created comprehensive interfaces for review data structures
   - Added proper typing for form submissions and API interactions
   - Implemented type-safe event handling for various user actions
   - Added proper typing for conditional rendering based on device type
   - Created reusable interfaces for review stats, votes, and responses
   - Implemented proper error handling with type-safe API error responses

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

### i18next Library Type Definition Issues

When working with i18next and TypeScript, you may encounter type definition errors in the library itself. These are common issues when the library's type definitions are not fully compatible with the TypeScript version used.

To work around these issues:

1. **Use type assertions** for i18next functions:
```typescript
const { t } = useTranslation() as any;
```

2. **Skip library checking** when running TypeScript compiler:
```bash
npx tsc --noEmit --skipLibCheck src/components/*.tsx
```

3. **Create custom type definitions** for i18next in your project:
```typescript
// src/types/i18next.d.ts
import 'i18next';

declare module 'i18next' {
  // Add your custom type definitions here
}
```

4. **Add i18next typescript definitions to the exclude list** in tsconfig.json:
```json
{
  "exclude": [
    "node_modules/i18next/typescript/**/*"
  ]
}
```

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

## Known Issues and Solutions

### lucide-react Named Export Issue

When building the project with TypeScript, you may encounter the following error:

```
Failed to compile.

./node_modules/lucide-react/dist/esm/createLucideIcon.mjs
Can't import the named export 'createElement' from non EcmaScript module (only default export is available)
```

This is a common issue when mixing ESM and CommonJS modules. There are several solutions:

1. **Use a specific version of lucide-react**:
```bash
npm install lucide-react@0.263.1
```

2. **Add a transpile module in webpack config**:
```js
// In react-scripts eject or craco config
module.exports = {
  webpack: {
    configure: (webpackConfig) => {
      webpackConfig.module.rules.push({
        test: /\.mjs$/,
        include: /node_modules/,
        type: 'javascript/auto',
      });
      return webpackConfig;
    },
  },
};
```

3. **Use icon components from a different library** like @mui/icons-material if available

### Material-UI Type Errors

When using Material-UI components like Modal and Slide, you may encounter TypeScript errors related to the children prop. Solutions include:

1. **Use conditional rendering instead of Modal/Slide components**:
```tsx
{isModalOpen && (
  <Box sx={{ /* modal styles */ }}>
    {children}
  </Box>
)}
```

2. **Use a wrapper component with proper type declarations**:
```tsx
const ModalWrapper: React.FC<{ open: boolean; onClose: () => void; children: React.ReactNode }> = ({
  open,
  onClose,
  children
}) => {
  return (
    <Modal open={open} onClose={onClose}>
      <Box>{children}</Box>
    </Modal>
  );
};
```

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [TypeScript with React](https://www.typescriptlang.org/docs/handbook/react.html)
- [Material UI with TypeScript](https://mui.com/material-ui/guides/typescript/)