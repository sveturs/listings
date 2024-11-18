import React from "react";
import AddRoom from "../components/AddRoom";
import AddUser from "../components/AddUser";
import AddBooking from "../components/AddBooking";

const AdminPanelPage = () => (
  <div>
    <h1>Админ-панель</h1>
    <AddRoom />
    <AddUser />
    <AddBooking />
  </div>
);

export default AdminPanelPage;
