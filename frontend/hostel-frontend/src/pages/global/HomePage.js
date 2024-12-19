import React from "react";
import { Container } from "@mui/material";
import RoomList from "../../components/accommodation/RoomList";

const HomePage = () => (
  <Container sx={{ marginTop: 4 }}>
    <RoomList />
  </Container>
);

export default HomePage;
