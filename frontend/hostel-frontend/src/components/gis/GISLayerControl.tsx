import React, { ChangeEvent, MouseEvent, SyntheticEvent } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Card,
  CardContent,
  Typography,
  ToggleButtonGroup,
  ToggleButton,
  FormControlLabel,
  Switch,
  Divider,
  Box,
  IconButton,
  Tooltip
} from '@mui/material';
import {
  Map as MapIcon,
  Satellite as SatelliteIcon,
  Terrain as TerrainIcon,
  FilterCenterFocus as HeatmapIcon,
  TrafficOutlined as TrafficIcon,
  GridOn as ClusteringIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

export type LayerType = 'standard' | 'satellite' | 'terrain' | 'traffic' | 'heatmap' | 'carto';

interface GISLayerControlProps {
  layers: LayerType;
  onLayerChange?: (layer: LayerType) => void;
  clusterMarkers: boolean;
  onClusteringToggle?: (clustered: boolean) => void;
  onClose?: () => void;
  minimized?: boolean;
  onToggleMinimize?: () => void;
}

interface LayerCardProps {
  children?: React.ReactNode;
}

const LayerCard = styled(Card)<LayerCardProps>(({ theme }) => ({
  position: 'absolute',
  zIndex: 1000,
  bottom: theme.spacing(2),
  right: theme.spacing(2),
  width: 250,
  [theme.breakpoints.down('sm')]: {
    bottom: theme.spacing(8),
    right: theme.spacing(1),
    width: 'calc(100% - 16px)',
    maxWidth: 250
  }
})) as React.ComponentType<LayerCardProps>;

interface StyledToggleButtonGroupProps {
  value?: string;
  exclusive?: boolean;
  onChange?: (event: React.MouseEvent<HTMLElement>, value: any) => void;
  'aria-label'?: string;
  children?: React.ReactNode;
}

const StyledToggleButtonGroup = styled(ToggleButtonGroup)<StyledToggleButtonGroupProps>(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'center',
  width: '100%',
  '& .MuiToggleButtonGroup-grouped': {
    margin: theme.spacing(0.5),
    border: 0,
    '&:not(:first-of-type)': {
      borderRadius: theme.shape.borderRadius,
    },
    '&:first-of-type': {
      borderRadius: theme.shape.borderRadius,
    },
    minWidth: 50,
    padding: theme.spacing(0.5)
  },
})) as React.ComponentType<StyledToggleButtonGroupProps>;

interface LayerToggleButtonProps {
  value: string;
  'aria-label'?: string;
  children?: React.ReactNode;
}

const LayerToggleButton = styled(ToggleButton)<LayerToggleButtonProps>(({ theme }) => ({
  width: 70,
  height: 70,
  flexDirection: 'column',
  '& .MuiSvgIcon-root': {
    fontSize: 24,
    marginBottom: theme.spacing(0.5)
  },
  '& .MuiTypography-root': {
    fontSize: 12
  }
})) as React.ComponentType<LayerToggleButtonProps>;

const GISLayerControl: React.FC<GISLayerControlProps> = ({ 
  layers, 
  onLayerChange, 
  clusterMarkers, 
  onClusteringToggle,
  onClose,
  minimized = false,
  onToggleMinimize
}) => {
  const { t } = useTranslation('gis');

  const handleLayerChange = (_event: MouseEvent<HTMLElement>, newLayer: LayerType | null) => {
    if (newLayer && onLayerChange) {
      onLayerChange(newLayer);
    }
  };

  const handleClusteringToggle = (event: ChangeEvent<HTMLInputElement>) => {
    if (onClusteringToggle) {
      onClusteringToggle(event.target.checked);
    }
  };

  if (minimized) {
    return (
      <Tooltip title={t('layers.title')}>
        <IconButton
          sx={{
            position: 'absolute',
            zIndex: 1000,
            bottom: 16,
            right: 16,
            backgroundColor: 'white',
            boxShadow: 2,
            '&:hover': {
              backgroundColor: 'white',
            }
          }}
          onClick={onToggleMinimize}
        >
          <MapIcon />
        </IconButton>
      </Tooltip>
    );
  }

  return (
    <LayerCard>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
          <Typography variant="h6" component="h2">
            {t('layers.title')}
          </Typography>
          <IconButton size="small" onClick={onToggleMinimize}>
            <MapIcon fontSize="small" />
          </IconButton>
        </Box>
        
        <Divider sx={{ mb: 2 }} />
        
        <StyledToggleButtonGroup
          value={layers}
          exclusive
          onChange={handleLayerChange}
          aria-label="map layers"
        >
          <LayerToggleButton value="standard" aria-label="standard map">
            <MapIcon />
            <Typography>{t('layers.standard')}</Typography>
          </LayerToggleButton>

          <LayerToggleButton value="satellite" aria-label="satellite map">
            <SatelliteIcon />
            <Typography>{t('layers.satellite')}</Typography>
          </LayerToggleButton>

          <LayerToggleButton value="terrain" aria-label="terrain map">
            <TerrainIcon />
            <Typography>{t('layers.terrain')}</Typography>
          </LayerToggleButton>

          <LayerToggleButton value="traffic" aria-label="traffic map">
            <TrafficIcon />
            <Typography>{t('layers.traffic')}</Typography>
          </LayerToggleButton>

          <LayerToggleButton value="heatmap" aria-label="heatmap">
            <HeatmapIcon />
            <Typography>{t('layers.heatmap')}</Typography>
          </LayerToggleButton>

          <LayerToggleButton value="carto" aria-label="carto map">
            <MapIcon color="primary" />
            <Typography>{t('layers.carto') || 'CartoDB'}</Typography>
          </LayerToggleButton>
        </StyledToggleButtonGroup>
        
        <Box mt={2}>
          <FormControlLabel
            control={
              <Switch
                checked={clusterMarkers}
                onChange={handleClusteringToggle}
                name="clustering"
                color="primary"
              />
            }
            label={
              <Box display="flex" alignItems="center">
                <ClusteringIcon sx={{ mr: 1, fontSize: 20 }} />
                <Typography variant="body2">{t('layers.clustering')}</Typography>
              </Box>
            }
          />
        </Box>
      </CardContent>
    </LayerCard>
  );
};

export default GISLayerControl;