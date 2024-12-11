import React, { useState } from 'react';
import {
    Box,
    List,
    ListItem,
    ListItemText,
    ListItemIcon,
    Collapse,
    Typography,
    Chip,
    IconButton,
} from '@mui/material';

import {
    ChevronRight,
    ChevronDown,
    Folder,
    Package
} from 'lucide-react';

const CategoryTreeItem = ({ category, onSelectCategory, selectedId }) => {
    const [open, setOpen] = useState(false);

    const hasChildren = category.children && category.children.length > 0;
    const isSelected = selectedId === category.id;

    return (
        <>
            <ListItem
                button
                onClick={() => {
                    if (hasChildren) {
                        setOpen(!open);
                    }
                    onSelectCategory(category.id);
                }}
                sx={{
                    bgcolor: isSelected ? 'action.selected' : 'transparent',
                    '&:hover': {
                        bgcolor: 'action.hover',
                    },
                }}
            >
                <ListItemIcon>
                    {hasChildren ? (
                        <IconButton
                            size="small"
                            onClick={(e) => {
                                e.stopPropagation();
                                setOpen(!open);
                            }}
                        >
                            {open ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
                        </IconButton>
                    ) : (
                        <Package size={18} />
                    )}
                </ListItemIcon>
                <ListItemText
                    primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            {category.name}
                            <Chip
                                label={category.listing_count}
                                size="small"
                                color={isSelected ? "primary" : "default"}
                                sx={{ ml: 'auto' }}
                            />
                        </Box>
                    }
                />
            </ListItem>
            {hasChildren && (
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        {category.children.map((child) => (
                            <Box key={child.id} sx={{ pl: 3 }}>
                                <CategoryTreeItem
                                    category={child}
                                    onSelectCategory={onSelectCategory}
                                    selectedId={selectedId}
                                />
                            </Box>
                        ))}
                    </List>
                </Collapse>
            )}
        </>
    );
};

const CategoryTree = ({ categories, onSelectCategory, selectedId }) => {
    if (!categories || categories.length === 0) {
        return (
            <Box sx={{ p: 2 }}>
                <Typography color="text.secondary">
                    Категории не найдены
                </Typography>
            </Box>
        );
    }

    return (
        <List
            sx={{
                width: '100%',
                maxWidth: 360,
                bgcolor: 'background.paper',
                borderRadius: 1,
                '& .MuiListItem-root': {
                    borderRadius: 1,
                },
            }}
        >
            <ListItem>
                <ListItemIcon>
                    <Folder size={20} />
                </ListItemIcon>
                <ListItemText
                    primary={
                        <Typography variant="subtitle1" fontWeight="medium">
                            Категории
                        </Typography>
                    }
                />
            </ListItem>
            {categories.map((category) => (
                <CategoryTreeItem
                    key={category.id}
                    category={category}
                    onSelectCategory={onSelectCategory}
                    selectedId={selectedId}
                />
            ))}
        </List>
    );
};

export default CategoryTree;