// src/components/reviews/ReviewForm.js
import React, { useState } from 'react';
import { Box, Rating, TextField, Button } from '@mui/material';

const ReviewForm = ({ onSubmit }) => {
    const [review, setReview] = useState({
        rating: 0,
        comment: ''
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        onSubmit(review);
    };

    return (
        <Box component="form" onSubmit={handleSubmit}>
            <Rating
                value={review.rating}
                onChange={(e, value) => setReview({ ...review, rating: value })}
            />
            <TextField
                fullWidth
                multiline
                rows={4}
                value={review.comment}
                onChange={(e) => setReview({ ...review, comment: e.target.value })}
                margin="normal"
            />
            <Button type="submit" variant="contained">
                Отправить
            </Button>
        </Box>
    );
};

export default ReviewForm;