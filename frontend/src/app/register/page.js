"use client"

import { useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import { showFieldError } from '@/lib/auth'
import Link from 'next/link'

export default function Register() {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        firstName: '',
        lastName: '',
        dateOfBirth: '',
        avatarImage: null,
        nickname: '',
        aboutMe: '',
        profileVisibility: 'public',
    })
    const [error, setError] = useState('')
    const [preview, setPreview] = useState('') // for image preview
    const router = useRouter()

    // initial password requirements state
    const [passwordRequirements, setPasswordRequirements] = useState({
        length: false,
        digit: false,
        uppercase: false,
        lowercase: false,
        special: false
    })

    // memoize required fields to avoid recreation on every render
    const requiredFields = useMemo(() => [
        'email', 'password', 'firstName', 'lastName', 'dateOfBirth', 'profileVisibility'
    ], [])

    const handleChange = (e) => {
        const { name, value } = e.target
        showFieldError(name, '');

        if (requiredFields.includes(name) && !value) {
            showFieldError(name, "This field is needed")
        }

        setFormData(prev => ({ ...prev, [name]: value }))

        // validate password if it's the password field
        if (name === 'password') {
            validatePassword(value)
        }
    }

    const validatePassword = (password) => {
        const requirements = {
            length: password.length >= 8 && password.length <= 16,
            digit: /\d/.test(password),
            uppercase: /[A-Z]/.test(password),
            lowercase: /[a-z]/.test(password),
            special: /[!@#$%^&*(),.?":{}|<>]/.test(password)
        }
        setPasswordRequirements(requirements)
    }

    const handleFileChange = (e) => {
        showFieldError('image', '');
        const file = e.target.files[0]
        if (file) {
            // validate file type
            const validTypes = ['image/jpeg', 'image/png', 'image/gif'];
            if (!validTypes.includes(file.type)) {
                showFieldError('image', 'Only JPEG, PNG, and GIF images are allowed');
                e.target.value = ''; // clear the file input
                return;
            }

            // validate file size (20MB max)
            if (file.size > 20 * 1000 * 1000) {
                showFieldError('image', 'Image must be 20MB or smaller');
                e.target.value = '';
                return;
            }

            setFormData(prev => ({ ...prev, avatarImage: file }))

            // create selected image preview
            const reader = new FileReader()
            reader.onloadend = () => {
                setPreview(reader.result)
            }
            reader.readAsDataURL(file)
        }
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')

        // check if all password requirements are met
        const allRequirementsMet = Object.values(passwordRequirements).every(req => req);
        if (!allRequirementsMet) {
            showFieldError('password', 'Please meet all password requirements');
            return;
        }

        try {
            const formDataToSend = new FormData()

            // append all form data to FormData object
            Object.entries(formData).forEach(([key, value]) => {
                if (value !== null && value !== undefined) {
                    if (key === 'avatarImage' && value instanceof File) {
                        formDataToSend.append(key, value)
                    } else {
                        formDataToSend.append(key, value)
                    }
                }
            })

            const response = await fetch('http://localhost:8080/api/register', {
                method: 'POST',
                body: formDataToSend,
            })

            if (!response.ok) {
                const data = await response.json();

                if (data.Email) {
                    showFieldError('email', data.Email);
                }
                if (data.Nickname) {
                    showFieldError('nickname', data.Nickname);
                }
                if (data.DateOfBirth) {
                    showFieldError('dateOfBirth', data.DateOfBirth);
                }
                if (data.AboutMe) {
                    showFieldError('aboutMe', data.AboutMe);
                }
            }

            if (response.ok) {
                router.push('/login') // redirect to login page after registration
            }
        } catch (err) {
            setError('Registration failed. Please try again.')
            console.error('Registration error:', err)
        }
    }

    return (
        <div className="container mx-auto p-4 max-w-md">
            <h1 className="text-2xl font-bold mb-4">Register</h1>
            {error && <p className="text-red-500 mb-4">{error}</p>}
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label className="block mb-1">Email <span className="text-red-500">*</span></label>
                    <input
                        type="email"
                        name="email"
                        value={formData.email}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="email"
                        required
                    />
                    <div id="email-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">First Name <span className="text-red-500">*</span></label>
                    <input
                        type="text"
                        name="firstName"
                        value={formData.firstName}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="firstName"
                        required
                    />
                    <div id="firstName-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">Last Name <span className="text-red-500">*</span></label>
                    <input
                        type="text"
                        name="lastName"
                        value={formData.lastName}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="lastName"
                        required
                    />
                    <div id="lastName-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">Date of Birth  <span className="text-red-500">*</span></label>
                    <input
                        type="date"
                        name="dateOfBirth"
                        value={formData.dateOfBirth}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="dateOfBirth"
                        required
                    />
                    <div id="dateOfBirth-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">Password <span className="text-red-500">*</span></label>
                    <input
                        type="password"
                        name="password"
                        value={formData.password}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="password"
                        required
                    />
                    <div id="password-error" className="text-red-500"></div>

                    {/* Password Requirements */}
                    <div className="password-requirements mt-2 text-sm">
                        <p className={`requirement ${passwordRequirements.length ? 'text-green-500' : 'text-gray-500'}`}>
                            • Your password length must be 8-16 characters
                        </p>
                        <p className={`requirement ${passwordRequirements.digit ? 'text-green-500' : 'text-gray-500'}`}>
                            • It must have at least one digit
                        </p>
                        <p className={`requirement ${passwordRequirements.uppercase ? 'text-green-500' : 'text-gray-500'}`}>
                            • Must contain at least one uppercase letter
                        </p>
                        <p className={`requirement ${passwordRequirements.lowercase ? 'text-green-500' : 'text-gray-500'}`}>
                            • Must contain at least one lowercase letter
                        </p>
                        <p className={`requirement ${passwordRequirements.special ? 'text-green-500' : 'text-gray-500'}`}>
                            • Must have at least one of the special characters (!@#$%^&*(),.?":&#123;&#125;|&lt;&gt;)
                        </p>
                    </div>
                </div>
                 <div>
                    <label className="block mb-1">Password <span className="text-red-500">*</span></label>
                    <input
                        type="password"
                        name="password"
                        value={formData.password}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="password"
                        required
                    />
                    <div id="password-error" className="text-red-500"></div>
                    </div>
                <div>
                    <label className="block mb-1">Avatar Image</label>
                    {preview && (
                        <div className="mb-2">
                            <img
                                src={preview}
                                alt="Preview"
                                className="w-25 object-cover rounded-full"
                            />
                        </div>
                    )}
                    <input
                        type="file"
                        name="avatarImage"
                        onChange={handleFileChange}
                        className="w-full p-2 border rounded"
                        id="image"
                        accept="image/*"
                    />
                    <div id="image-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">Nickname</label>
                    <input
                        type="text"
                        name="nickname"
                        value={formData.nickname}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="nickname"
                    />
                    <div id="nickname-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">About Me</label>
                    <textarea
                        name="aboutMe"
                        value={formData.aboutMe}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="aboutMe"
                        rows="3"
                    />
                    <div id="aboutMe-error" className="text-red-500"></div>
                </div>
                {/* profile visibility */}
                <div>
                    <label className="block mb-1">Profile Visibility*</label>
                    <select
                        name="profileVisibility"
                        value={formData.profileVisibility}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        id="profileVisibility"
                        required
                    >
                        <option value="public" className="bg-gray-900">Public (Anyone can view your profile)</option>
                        <option value="private" className="bg-gray-900">Private (Only you can view your profile)</option>
                    </select>
                    <div id="profileVisibility-error" className="text-red-500"></div>
                </div>
                <button
                    type="submit"
                    className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600"
                >
                    Register
                </button>
            </form>
            <p className="mt-4">
                Already have an account?{' '}
                <Link href="/login" className="text-blue-500 hover:underline">
                    Login here
                </Link>
            </p>
        </div>
    )
}