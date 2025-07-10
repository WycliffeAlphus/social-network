export const showFieldError = (fieldId, message) => {
    const errorElement = document.getElementById(`${fieldId}-error`);
    const inputElement = document.getElementById(fieldId);

    if (message) {
        errorElement.textContent = message;
        errorElement.style.display = 'block';
        inputElement.classList.add('input-error');
    } else {
        errorElement.textContent = '';
        errorElement.style.display = 'none';
        inputElement.classList.remove('input-error');
    }
}
