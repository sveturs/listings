'use client';

import { Suspense, useState } from 'react';
import { useTranslations } from 'next-intl';
import { 
  Box, 
  Container, 
  Grid, 
  TextField, 
  InputAdornment,
  ToggleButtonGroup,
  ToggleButton,
  Button,
  Typography,
  Paper
} from '@mui/material';
import { Search, Grid as GridIcon, List, Map } from 'lucide-react';
import ListingGrid from '@/components/marketplace/ListingGrid';
import ListingFilters from '@/components/marketplace/ListingFilters';
import Loading from './loading';

export default function MarketplacePage() {
  const t = useTranslations('marketplace');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState({});
  
  const handleViewModeChange = (
    event: React.MouseEvent<HTMLElement>,
    newViewMode: 'grid' | 'list' | null,
  ) => {
    if (newViewMode !== null) {
      setViewMode(newViewMode);
    }
  };

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value);
  };

  const handleFiltersChange = (newFilters: any) => {
    setFilters(newFilters);
  };
  
  return (
    <Container maxWidth="xl" sx={{ py: 2 }}>
      {/* Search Bar and View Controls */}
      <Box sx={{ mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              placeholder={t('listings.searchPlaceholder')}
              value={searchQuery}
              onChange={handleSearchChange}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
              <Button 
                variant="outlined" 
                startIcon={<Map />}
                onClick={() => {/* TODO: Map view */}}
              >
                {t('filters.map')}
              </Button>
              <ToggleButtonGroup
                value={viewMode}
                exclusive
                onChange={handleViewModeChange}
                aria-label="view mode"
              >
                <ToggleButton value="grid" aria-label="grid view">
                  <GridIcon />
                </ToggleButton>
                <ToggleButton value="list" aria-label="list view">
                  <List />
                </ToggleButton>
              </ToggleButtonGroup>
            </Box>
          </Grid>
        </Grid>
      </Box>

      <Grid container spacing={3}>
        {/* Filters Sidebar */}
        <Grid item xs={12} md={3}>
          <ListingFilters onFiltersChange={handleFiltersChange} />
        </Grid>
        
        {/* Listings Grid */}
        <Grid item xs={12} md={9}>
          <Typography variant="h4" component="h1" gutterBottom>
            {t('title')}
          </Typography>
          <Suspense fallback={<Loading />}>
            <ListingGrid 
              viewMode={viewMode}
              searchQuery={searchQuery}
              filters={filters}
            />
          </Suspense>
        </Grid>
      </Grid>
    </Container>
  );
}