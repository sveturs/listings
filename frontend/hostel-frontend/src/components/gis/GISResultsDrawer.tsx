import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Drawer,
  Box,
  Typography,
  IconButton,
  Divider,
  List,
  CircularProgress,
  Button,
  Fab,
  useMediaQuery,
  Zoom,
  Slide,
  Paper,
  ToggleButtonGroup,
  ToggleButton
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import {
  ChevronLeft as ChevronLeftIcon,
  ChevronRight as ChevronRightIcon,
  List as ListIcon,
  Map as MapIcon,
  SortRounded as SortIcon,
  FilterList as FilterIcon,
  Refresh as RefreshIcon,
  GridView as GridViewIcon,
  ViewList as ViewListIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import GISListingCard, { GISListing } from './GISListingCard';

// Width of the drawer - full width on the page
const drawerWidth = 400;

interface ResultsDrawerProps {
  children?: React.ReactNode;
}

const ResultsDrawer = styled(Box)<ResultsDrawerProps>(({ theme }) => ({
  width: drawerWidth,
  height: '100%',
  flexShrink: 0,
  position: 'relative',
  backgroundColor: theme.palette.background?.paper,
  boxShadow: '2px 0 5px rgba(0,0,0,0.1)',
  borderRight: `1px solid ${theme.palette.divider}`,
  zIndex: 10,
  [theme.breakpoints.down('md')]: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    height: '50%',
    width: '100%',
    zIndex: 1000
  }
})) as React.ComponentType<ResultsDrawerProps>;

interface DrawerContentProps {
  children?: React.ReactNode;
}

const DrawerContent = styled(Box)<DrawerContentProps>(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  height: '100%'
})) as React.ComponentType<DrawerContentProps>;

interface DrawerHeaderProps {
  children?: React.ReactNode;
}

const DrawerHeader = styled(Box)<DrawerHeaderProps>(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: theme.spacing(1, 2),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar
})) as React.ComponentType<DrawerHeaderProps>;

interface ListingsContainerProps {
  children?: React.ReactNode;
}

const ListingsContainer = styled(Box)<ListingsContainerProps>(({ theme }) => ({
  flexGrow: 1,
  padding: theme.spacing(2),
  overflowY: 'auto'
})) as React.ComponentType<ListingsContainerProps>;

interface ToggleDrawerButtonProps {
  drawerOpen?: boolean;
  color?: string;
  'aria-label'?: string;
  onClick?: () => void;
  size?: 'small' | 'medium' | 'large';
  sx?: object;
  children?: React.ReactNode;
}

const ToggleDrawerButton = styled(Fab, {
  shouldForwardProp: (prop) => prop !== 'drawerOpen'
})<ToggleDrawerButtonProps>(({ theme, drawerOpen }) => ({
  position: 'fixed', // Use fixed positioning to keep it visible regardless of drawer state
  top: theme.spacing(10), // Position below the top app bar
  left: drawerOpen ? drawerWidth : theme.spacing(2), // If drawer is open, place at edge of drawer, otherwise near left edge
  transform: drawerOpen ? 'translateX(-50%)' : 'none', // Center on edge only when drawer is open
  zIndex: 1300, // Higher z-index to stay above everything
  transition: 'left 0.3s', // Smooth transition when drawer opens/closes
  boxShadow: '0 2px 10px rgba(0,0,0,0.2)', // More visible shadow
  [theme.breakpoints.down('md')]: {
    display: 'none'
  }
})) as React.ComponentType<ToggleDrawerButtonProps>;

interface MobileDrawerToggleProps {
  children?: React.ReactNode;
}

const MobileDrawerToggle = styled(Paper)<MobileDrawerToggleProps>(({ theme }) => ({
  position: 'absolute',
  bottom: theme.spacing(12),
  left: '50%',
  transform: 'translateX(-50%)',
  zIndex: 1200,
  display: 'flex',
  padding: theme.spacing(0.5),
  borderRadius: theme.spacing(3),
  boxShadow: theme.shadows[3]
})) as React.ComponentType<MobileDrawerToggleProps>;

interface ViewToggleButtonsProps {
  value?: string;
  exclusive?: boolean;
  onChange?: (event: React.MouseEvent<HTMLElement>, value: string | null) => void;
  'aria-label'?: string;
  size?: 'small' | 'medium' | 'large';
  children?: React.ReactNode;
}

const ViewToggleButtons = styled(ToggleButtonGroup)<ViewToggleButtonsProps>(({ theme }) => ({
  marginLeft: theme.spacing(2)
})) as React.ComponentType<ViewToggleButtonsProps>;

interface GISResultsDrawerProps {
  open: boolean;
  onToggleDrawer: () => void;
  listings?: GISListing[];
  loading?: boolean;
  onShowOnMap?: (listing: GISListing) => void;
  onFilterClick?: () => void;
  onSortClick?: () => void;
  onRefresh?: () => void;
  favoriteListings?: (number | string)[];
  onFavoriteToggle?: (id: number | string, isFavorite: boolean) => void;
  onContactClick?: (listing: GISListing) => void;
  totalCount?: number;
  expandToEdge?: boolean;
}

const GISResultsDrawer: React.FC<GISResultsDrawerProps> = ({ 
  open, 
  onToggleDrawer, 
  listings = [], 
  loading = false, 
  onShowOnMap,
  onFilterClick,
  onSortClick,
  onRefresh,
  favoriteListings = [],
  onFavoriteToggle,
  onContactClick,
  totalCount = 0,
  expandToEdge = false // Add new prop
}) => {
  const { t } = useTranslation('gis');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [viewType, setViewType] = useState<string>('list');

  const handleViewTypeChange = (_event: React.MouseEvent<HTMLElement>, newViewType: string | null): void => {
    if (newViewType !== null) {
      setViewType(newViewType);
    }
  };

  // Determine drawer position based on screen size
  const drawerProps = isMobile ? { 
    anchor: 'bottom',
    variant: 'persistent',
    open: open,
    // Ensure proper width on mobile
    sx: { width: open ? '100%' : 0 }
  } : { 
    anchor: 'left',
    variant: 'persistent',
    open: open,
    // This ensures the drawer takes up width in the flex layout only when open
    sx: { width: open ? drawerWidth : 0, visibility: open ? 'visible' : 'hidden' }
  };

  // Toggle button is already positioned by the styled component

  return (
    <>
      {/* Only render the drawer if it's open */}
      {open && (
        <ResultsDrawer>
          <DrawerContent>
            <DrawerHeader>
              <Box display="flex" alignItems="center">
                <Typography variant="h6">
                  {t('drawer.results')}
                  {totalCount > 0 && (
                    <Typography component="span" variant="body2" sx={{ ml: 1 }}>
                      ({totalCount})
                    </Typography>
                  )}
                </Typography>
                <ViewToggleButtons
                  value={viewType}
                  exclusive
                  onChange={handleViewTypeChange}
                  aria-label="view type"
                  size="small"
                >
                  <ToggleButton value="list" aria-label="list view">
                    <ViewListIcon fontSize="small" />
                  </ToggleButton>
                  <ToggleButton value="grid" aria-label="grid view">
                    <GridViewIcon fontSize="small" />
                  </ToggleButton>
                </ViewToggleButtons>
              </Box>
              <Box>
                <IconButton onClick={onSortClick} size="small" sx={{ mr: 1 }}>
                  <SortIcon />
                </IconButton>
                <IconButton onClick={onFilterClick} size="small" sx={{ mr: 1 }}>
                  <FilterIcon />
                </IconButton>
                <IconButton onClick={onRefresh} size="small" sx={{ mr: 1 }}>
                  <RefreshIcon />
                </IconButton>
                <IconButton onClick={onToggleDrawer} edge="end">
                  <ChevronLeftIcon />
                </IconButton>
              </Box>
            </DrawerHeader>
            
            <Divider />
            
            <ListingsContainer>
              {loading ? (
                <Box display="flex" justifyContent="center" alignItems="center" height="100%">
                  <CircularProgress />
                </Box>
              ) : listings.length > 0 ? (
                <List disablePadding>
                  {listings.map(listing => (
                    <GISListingCard
                      key={listing.id}
                      listing={listing}
                      isFavorite={favoriteListings.includes(listing.id)}
                      onFavoriteToggle={onFavoriteToggle}
                      onShowOnMap={onShowOnMap}
                      onContactClick={onContactClick}
                      compact={viewType === 'grid'}
                    />
                  ))}
                </List>
              ) : (
                <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height="70%">
                  <Typography variant="body1" color="textSecondary" gutterBottom>
                    {t('drawer.noListingsFound')}
                  </Typography>
                  <Button
                    variant="outlined"
                    startIcon={<RefreshIcon />}
                    onClick={onRefresh}
                    sx={{ mt: 2 }}
                  >
                    {t('search.clearSearch')}
                  </Button>
                </Box>
              )}
            </ListingsContainer>
          </DrawerContent>
        </ResultsDrawer>
      )}

      {/* Always show the toggle button */}
      {isMobile ?
        <Zoom in={!open} children={
          <MobileDrawerToggle>
            <Button
              startIcon={<ListIcon />}
              onClick={onToggleDrawer}
              color="primary"
            >
              {t('drawer.results')}
              {totalCount > 0 && ` (${totalCount})`}
            </Button>
          </MobileDrawerToggle>
        } />
       :
        <Fab
          color="primary"
          aria-label={open ? t('drawer.close') : t('drawer.open')}
          onClick={onToggleDrawer}
          size="medium"
          sx={{
            position: 'fixed',
            top: '85px', // Increased by 5px to match main container
            left: open ? `${drawerWidth - 20}px` : '20px',
            zIndex: 1000,
            transition: 'left 0.3s ease'
          }}
        >
          {open ? <ChevronLeftIcon /> : <ChevronRightIcon />}
        </Fab>
      }
    </>
  );
};

export default GISResultsDrawer;