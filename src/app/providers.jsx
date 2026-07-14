import { AuthProvider } from '../features/auth/AuthContext'

export default function AppProviders({ children }) {
  return <AuthProvider>{children}</AuthProvider>
}
