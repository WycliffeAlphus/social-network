export const showFieldError = (fieldId, message) => {
    const errorElement = document.getElementById(`${fieldId}-error`);
    const inputElement = document.getElementById(fieldId);

    if (message) {
        errorElement.textContent = message;
        errorElement.style.display = 'block';
        if (fieldId != "title" && fieldId != "content") {
            inputElement.classList.add('input-error');
        } else {
            inputElement.classList.add('input-error-bottom');
        }
    } else {
        errorElement.textContent = '';
        errorElement.style.display = 'none';
        if (fieldId != "title" && fieldId != "content") {
            inputElement.classList.remove('input-error');
        } else {
            inputElement.classList.remove('input-error-bottom');
        }
    }
}
