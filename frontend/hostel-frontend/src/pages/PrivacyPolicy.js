import React from 'react';
import { Container, Typography, Box, Paper } from '@mui/material';

const PrivacyPolicy = () => {
  return (
    <Container maxWidth="md">
      <Box py={4}>
        <Paper elevation={1} sx={{ p: 4 }}>
          <Typography variant="h3" component="h1" gutterBottom>
            Privacy Policy
          </Typography>
          
          <Typography variant="subtitle1" color="text.secondary" paragraph>
            Last updated: December 03, 2024
          </Typography>

          <Typography variant="h5" component="h2" gutterBottom sx={{ mt: 4 }}>
            1. Introduction
          </Typography>
          <Typography paragraph>
            LandHub.rs ("we", "our", or "us") is committed to protecting your privacy. 
            This Privacy Policy explains how we collect, use, disclose, and safeguard your 
            information when you use our service.
          </Typography>

          <Typography variant="h5" component="h2" gutterBottom sx={{ mt: 4 }}>
            2. Information We Collect
          </Typography>
          
          <Typography variant="h6" component="h3" gutterBottom>
            2.1 Information You Provide
          </Typography>
          <Typography component="div" sx={{ pl: 2 }}>
            <ul>
              <li>Name and contact details</li>
              <li>Booking information</li>
              <li>Payment information</li>
              <li>Communication preferences</li>
              <li>Account credentials</li>
            </ul>
          </Typography>

          <Typography variant="h6" component="h3" gutterBottom>
            2.2 Information Automatically Collected
          </Typography>
          <Typography component="div" sx={{ pl: 2 }}>
            <ul>
              <li>Device information</li>
              <li>IP address</li>
              <li>Browser type</li>
              <li>Usage data</li>
              <li>Location data (if permitted)</li>
            </ul>
          </Typography>

          <Typography variant="h6" component="h3" gutterBottom>
            2.3 Information from Third Parties
          </Typography>
          <Typography paragraph>
            We may receive information about properties and availability from Booking.com 
            through their API integration.
          </Typography>

          {/* Additional sections following the same pattern... */}

          <Typography variant="h5" component="h2" gutterBottom sx={{ mt: 4 }}>
            11. Contact Us
          </Typography>
          <Typography paragraph>
            For any questions about this Privacy Policy, please contact us at:
          </Typography>
          <Typography component="div" sx={{ pl: 2 }}>
            <ul>
              <li>Email: voroshilovdo@gmail.com</li>
              <li>Address: Srbia Novi Sad Vase Staica 18</li>
             </ul>
          </Typography>
        </Paper>
      </Box>
    </Container>
  );
};

export default PrivacyPolicy;