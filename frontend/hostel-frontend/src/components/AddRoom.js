// src/components/AddRoom.js
import React, { useState } from "react";
import { 
    TextField, 
    Button, 
    Container, 
    Typography,
    Box, 
    Alert,
    Grid, 
    IconButton,
    FormControl, 
    InputLabel, 
    Select, 
    MenuItem,
    FormControlLabel, 
    Switch 
} from "@mui/material";
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import axios from "../api/axios";

const AddRoom = () => {
    const [room, setRoom] = useState({
        name: "",
        accommodation_type: "room",
        capacity: 0,
        price_per_night: 0,
        address_street: "",
        address_city: "",
        address_state: "",
        address_country: "",
        address_postal_code: "",
        is_shared: false,
        total_beds: null,
        available_beds: null,
        has_private_bathroom: true
    });
    
    const [beds, setBeds] = useState([
        { bed_number: "1", price_per_night: 0 }
    ]);
    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [errorMessage, setErrorMessage] = useState("");
    const [isSuccess, setIsSuccess] = useState(false);

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files || []);
        if (files.length === 0) {
            return;
        }

        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                setErrorMessage("Можно загружать только изображения");
                return false;
            }
            if (file.size > 5 * 1024 * 1024) {
                setErrorMessage("Размер файла не должен превышать 5MB");
                return false;
            }
            return true;
        });

        if (validFiles.length === 0) {
            return;
        }

        setImages(prev => [...prev, ...validFiles]);

        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrls(prev => [...prev, reader.result]);
            };
            reader.onerror = () => {
                setErrorMessage("Ошибка при чтении файла: " + file.name);
            };
            reader.readAsDataURL(file);
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setErrorMessage("");
        setIsSuccess(false);

        try {
            if (room.accommodation_type === 'bed' && (!room.total_beds || beds.length === 0)) {
                setErrorMessage("Добавьте информацию о кроватях");
                return;
            }

            // ... rest of the submit logic ...

            setIsSuccess(true);
            // Reset form
            setRoom({
                name: "",
                accommodation_type: "room",
                // ... rest of the initial state
            });
            setBeds([{ bed_number: "1", price_per_night: 0 }]);
            setImages([]);
            setPreviewUrls([]);
        } catch (error) {
            console.error('Ошибка при добавлении:', error);
            setErrorMessage(error.response?.data || "Ошибка при добавлении объекта размещения");
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>
                Добавить объект размещения
            </Typography>
            {errorMessage && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {errorMessage}
                </Alert>
            )}
            {isSuccess && (
                <Alert severity="success" sx={{ mb: 2 }}>
                    Объект размещения успешно добавлен
                </Alert>
            )}
            {/* Rest of the JSX */}
        </Container>
    );
};

export default AddRoom;