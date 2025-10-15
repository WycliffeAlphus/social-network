export const showFieldError = (fieldId, message) => {
    const errorElement = document.getElementById(`${fieldId}-error`);
    const inputElement = document.getElementById(fieldId);

    if (errorElement) {
        if (message) {
            errorElement.textContent = message;
            errorElement.style.display = 'block';
            if (fieldId != "title" && fieldId != "content") {
                if (inputElement) inputElement.classList.add('input-error');
            } else {
                if (inputElement) inputElement.classList.add('input-error-bottom');
            }
        } else {
            errorElement.textContent = '';
            errorElement.style.display = 'none';
            if (fieldId != "title" && fieldId != "content") {
                if (inputElement) inputElement.classList.remove('input-error');
            } else {
                if (inputElement) inputElement.classList.remove('input-error-bottom');
            }
        }
    }
}
