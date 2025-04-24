// frontend/hostel-frontend/src/utils/adminUtils.js
export const getAdminEmails = () => {
    return process.env.REACT_APP_ADMIN_EMAILS
        ? process.env.REACT_APP_ADMIN_EMAILS.split(',')
        : ['voroshilovdo@gmail.com']; // Значение по умолчанию
};

export const isAdmin = (email) => {
    if (!email) return false;
    return getAdminEmails().includes(email);
};