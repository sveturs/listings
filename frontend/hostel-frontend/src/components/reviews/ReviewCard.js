// src/components/reviews/ReviewCard.js
import React from 'react';
import { Card, CardContent, Typography, Rating } from '@mui/material';

const ReviewCard = ({ review }) => {
    return (
        <Card>
            <CardContent>
                <Rating value={review.rating} readOnly />
                <Typography>{review.comment}</Typography>
            </CardContent>
        </Card>
    );
};

export default ReviewCard;