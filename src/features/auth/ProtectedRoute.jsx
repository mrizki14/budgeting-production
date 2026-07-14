import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from './AuthContext'

export default function ProtectedRoute() {
  const { ready, token } = useAuth()

  if (!ready) return <div className="app-loading">Loading...</div>
  if (!token) return <Navigate to="/login" replace />
  return <Outlet />
}
