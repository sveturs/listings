import React from 'react';
import { useTranslation } from 'react-i18next';
import {
    Select,
    MenuItem,
    FormControl,
    Box,
    Typography
} from '@mui/material';

const LANGUAGES = [
    { code: 'en', name: 'English', flag: '🇬🇧' },
    { code: 'ru', name: 'Русский', flag: '🇷🇺' },
    { code: 'sr', name: 'Српски', flag: '🇷🇸' }
];

const LanguageSwitcher = () => {
    const { i18n } = useTranslation();

    const handleLanguageChange = (event) => {
        const newLang = event.target.value;
        i18n.changeLanguage(newLang);
        localStorage.setItem('preferredLanguage', newLang);
        document.documentElement.lang = newLang;
    };

    return (
        <FormControl size="small">
            <Select
                value={i18n.language}
                onChange={handleLanguageChange}
                sx={{
                    minWidth: 120,
                    '& .MuiSelect-select': {
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1
                    }
                }}
            >
                {LANGUAGES.map((lang) => (
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