// frontend/hostel-frontend/src/pages/global/HomePage.js
import React from "react";
import { Container } from "@mui/material";
import ListingCard from "../../components/marketplace/ListingCard";

const HomePage = () => (
  <Container sx={{ marginTop: 4 }}>
    <ListingCard />
  </Container>
);

export default HomePage;
