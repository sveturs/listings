import React from "react";
import AddRoom from "../components/accommodation/AddRoom";
import AddUser from "../components/user/AddUser";
import AddBooking from "../components/accommodation/AddBooking";

const AdminPanelPage = () => (
  <div>
    <h1>Админ-панель</h1>
    <AddRoom />
    <AddUser />
    <AddBooking />
  </div>
);

export default AdminPanelPage;
