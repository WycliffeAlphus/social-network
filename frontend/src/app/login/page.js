"use client"

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { showFieldError } from '@/lib/auth'

export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const router = useRouter()

    const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    const formData = new FormData(e.target)
    const actualEmail = formData.get('email')?.toString() || email
    const actualPassword = formData.get('password')?.toString() || password

    try {
        const response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: actualEmail,
                password: actualPassword,
            }),
            credentials: 'include',
        })

        if (!response.ok) {
            const data = await response.json()

            if (data.Email) {
                showFieldError('email', data.Email)
                return
            }
            if (data.Password) {
                showFieldError('login-password', data.Password)
            }
        }

        if (response.ok) {
            router.push('/')
        }
    } catch (err) {
        setError('Invalid credentials')
        console.error('Login error:', err)
    }
}


    return (
        <div className="container mx-auto p-4 max-w-md">
            <h1 className="text-2xl font-bold mb-4">Login</h1>
            {error && <p className="text-red-500 mb-4">{error}</p>}
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label className="block mb-1">Email</label>
                    <input
                        type="text"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className="w-full p-2 border rounded"
                        id="email"
                        name = "email"
                        required
                    />
                    <div id="email-error" className="text-red-500"></div>
                </div>
                <div>
                    <label className="block mb-1">Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="w-full p-2 border rounded"
                        id="login-password"
                        name = "password"
                        required
                    />
                    <div id="login-password-error" className="text-red-500"></div>
                </div>
                <button
                    type="submit"
                    className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600"
                >
                    Login
                </button>
            </form>
            <p className="mt-4">
                Don't have an account?{' '}
                <Link href="/register" className="text-blue-500 hover:underline">
                    Register here
                </Link>
            </p>
        </div>
    )
}