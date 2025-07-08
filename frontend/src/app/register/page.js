"use client"

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { showFieldError } from '@/lib/auth'

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
    })
    const [error, setError] = useState('')
    const [preview, setPreview] = useState('') // for image preview
    const router = useRouter()

    const handleChange = (e) => {
        const { name, value } = e.target
        setFormData(prev => ({ ...prev, [name]: value }))
    }

    const handleFileChange = (e) => {
        const file = e.target.files[0]
        if (file) {
            // validate file type
            const validTypes = ['image/jpeg', 'image/png', 'image/gif'];
            if (!validTypes.includes(file.type)) {
                setError('Only JPEG, PNG, and GIF images are allowed');
                e.target.value = ''; // clear the file input
                return;
            }

            // validate file size (20MB max)
            if (file.size > 20 * 1000 * 1000) {
                setError('Image must be 20MB or smaller');
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
            }

            // const data = await response.json()
            // console.log('Registration successful:', data)
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
                    <label className="block mb-1">Password <span className="text-red-500">*</span></label>
                    <input
                        type="password"
                        name="password"
                        value={formData.password}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        required
                    />
                </div>
                <div>
                    <label className="block mb-1">First Name <span className="text-red-500">*</span></label>
                    <input
                        type="text"
                        name="firstName"
                        value={formData.firstName}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        required
                    />
                </div>
                <div>
                    <label className="block mb-1">Last Name <span className="text-red-500">*</span></label>
                    <input
                        type="text"
                        name="lastName"
                        value={formData.lastName}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        required
                    />
                </div>
                <div>
                    <label className="block mb-1">Date of Birth (YYYY-MM-DD)  <span className="text-red-500">*</span></label>
                    <input
                        type="date"
                        name="dateOfBirth"
                        value={formData.dateOfBirth}
                        onChange={handleChange}
                        className="w-full p-2 border rounded"
                        required
                    />
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
                        accept="image/*"
                    />
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
                        rows="3"
                    />
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
                <a href="/login" className="text-blue-500 hover:underline">
                    Login here
                </a>
            </p>
        </div>
    )
}