// frontend/hostel-frontend/src/components/shared/LanguageSwitcher.js
import React from 'react';
import {
    Select,
    MenuItem,
    FormControl,
    Box,
    Typography
} from '@mui/material';
import { useLanguage } from '../../contexts/LanguageContext';

const LanguageSwitcher = () => {
    const { language, setLanguage, supportedLanguages } = useLanguage();

    return (
        <FormControl size="small">
            <Select
                value={language}
                onChange={(e) => setLanguage(e.target.value)}
                sx={{
                    minWidth: 120,
                    '& .MuiSelect-select': {
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1
                    }
                }}
            >
                {supportedLanguages.map((lang) => (
                    <MenuItem key={lang.code} value={lang.code}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Typography variant="body2" component="span">
                                {lang.flag}
                            </Typography>
                            <Typography variant="body2">
                                {lang.name}
                            </Typography>
                        </Box>
                    </MenuItem>
                ))}
            </Select>
        </FormControl>
    );
};

export default LanguageSwitcher;