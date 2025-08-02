export const validateAndPreviewFile = (file, setPreview, maxSizeMB = 20) => {
  const validTypes = ['image/jpeg', 'image/png', 'image/gif'];
  
  // validate file type
  if (!validTypes.includes(file.type)) {
    return { valid: false, error: 'Only JPEG, PNG, and GIF images are allowed' };
  }

  // validate file size
  if (file.size > maxSizeMB * 1000 * 1000) {
    return { valid: false, error: `Image must be ${maxSizeMB}MB or smaller` };
  }

  // create preview if callback provided
  if (setPreview) {
    const reader = new FileReader();
    reader.onloadend = () => setPreview(reader.result);
    reader.readAsDataURL(file);
  }

  return { valid: true, file };
};