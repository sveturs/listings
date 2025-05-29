'use client';

import { useState } from 'react';
import Grid from '@mui/material/Grid';
import { Box, TextField, InputAdornment, IconButton, ToggleButton, ToggleButtonGroup } from '@mui/material';
import { Search as SearchIcon, ViewList, ViewModule } from '@mui/icons-material';
import { useTranslations } from 'next-intl';
import ListingGrid from './ListingGrid';
import ListingFilters from './ListingFilters';

export default function MarketplaceContent() {
  const t = useTranslations('marketplace.listings');
  
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState<Record<string, unknown>>({});
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  const handleFiltersChange = (newFilters: Record<string, unknown>) => {
    setFilters(newFilters);
  };

  const handleSearch = () => {
    // Force re-render of ListingGrid with new search query
    setFilters(prev => ({ ...prev, _searchTrigger: Date.now() }));
  };

  const handleSearchKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <>
      {/* Search Bar and View Controls */}
      <Box sx={{ mb: 3, display: 'flex', gap: 2, alignItems: 'center' }}>
        <TextField
          fullWidth
          placeholder={t('searchPlaceholder')}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          onKeyPress={handleSearchKeyPress}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton onClick={handleSearch} edge="end">
                  <SearchIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
        
        <ToggleButtonGroup
          value={viewMode}
          exclusive
          onChange={(_, newMode) => newMode && setViewMode(newMode)}
          aria-label="view mode"
        >
          <ToggleButton value="grid" aria-label="grid view">
            <ViewModule />
          </ToggleButton>
          <ToggleButton value="list" aria-label="list view">
            <ViewList />
          </ToggleButton>
        </ToggleButtonGroup>
      </Box>

      <Grid container spacing={3}>
        {/* Filters Sidebar */}
        <Grid size={{ xs: 12, md: 3 }}>
          <ListingFilters onFiltersChange={handleFiltersChange} />
        </Grid>
        
        {/* Listings Grid */}
        <Grid size={{ xs: 12, md: 9 }}>
          <ListingGrid 
            viewMode={viewMode}
            searchQuery={searchQuery}
            filters={filters}
          />
        </Grid>
      </Grid>
    </>
  );
}