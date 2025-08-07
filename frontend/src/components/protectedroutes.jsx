"use client"

import { useRouter, usePathname } from 'next/navigation'
import { useEffect, useState } from 'react'

const publicRoutes = ['/login', '/register']

export default function ProtectedRoute({ children }) {
    const router = useRouter()
    const pathname = usePathname()
    const [isAuthenticated, setIsAuthenticated] = useState(false)
    const [isLoading, setIsLoading] = useState(true)

    const isPublic = publicRoutes.includes(pathname)

    useEffect(() => {
        const verifyAuth = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/profile/currentuser', {
                    method: 'GET',
                    credentials: 'include',
                })
                // return response.ok ? await response.json() : null
                if (response.ok) {
                    setIsAuthenticated(true)
                    if (isPublic) {
                        router.push('/') // redirect to home if already authenticated
                        return
                    }
                } else {
                    setIsAuthenticated(false)
                    if (!isPublic) {
                        router.push('/login')
                        return
                    }
                }
            } catch (error) {
                console.error('Auth check failed:', error);
                setIsLoading(false);
                return null;
            }

            setIsLoading(false) // wait for auth verification to complete
        }
        verifyAuth()
    }, [pathname, router])

    if (isLoading) return null // don't render children until auth verification and any redirect is completed

    return isPublic || isAuthenticated ? children : null
}
