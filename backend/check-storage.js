// Script to check localStorage and sessionStorage
console.log('=== Checking localStorage ===');
console.log('localStorage keys:', Object.keys(localStorage));
for (let key of Object.keys(localStorage)) {
    console.log(`localStorage['${key}']:`, localStorage.getItem(key));
}

console.log('\n=== Checking sessionStorage ===');
console.log('sessionStorage keys:', Object.keys(sessionStorage));
for (let key of Object.keys(sessionStorage)) {
    console.log(`sessionStorage['${key}']:`, sessionStorage.getItem(key));
}

console.log('\n=== Checking form field values ===');
const emailField = document.querySelector('input[type="email"]');
const passwordField = document.querySelector('input[type="password"]');
console.log('Email field value:', emailField?.value);
console.log('Password field value:', passwordField?.value);

console.log('\n=== Checking for React/Next.js specific storage ===');
// Check for any keys that might be storing form data
const possibleKeys = ['loginEmail', 'loginPassword', 'authForm', 'formData', 'userEmail', 'rememberMe'];
possibleKeys.forEach(key => {
    if (localStorage.getItem(key)) {
        console.log(`Found in localStorage - ${key}:`, localStorage.getItem(key));
    }
    if (sessionStorage.getItem(key)) {
        console.log(`Found in sessionStorage - ${key}:`, sessionStorage.getItem(key));
    }
});