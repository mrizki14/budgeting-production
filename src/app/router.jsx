import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom'
import AppShell from '../components/layout/AppShell'
import AuthPage from '../features/auth/AuthPage'
import GuestRoute from '../features/auth/GuestRoute'
import ProtectedRoute from '../features/auth/ProtectedRoute'
import LandingPage from '../features/landing/LandingPage'
import DashboardPage from '../features/dashboard/DashboardPage'
import CategoriesPage from '../features/categories/CategoriesPage'
import CategoryEditorPage from '../features/categories/CategoryEditorPage'
import BudgetsPage from '../features/budgets/BudgetsPage'
import BudgetEditorPage from '../features/budgets/BudgetEditorPage'
import TransactionsPage from '../features/transactions/TransactionsPage'
import TransactionEditorPage from '../features/transactions/TransactionEditorPage'
import ReportsPage from '../features/reports/ReportsPage'
import SettingsPage from '../features/settings/SettingsPage'

const router = createBrowserRouter([
  { path: '/', element: <LandingPage /> },
  {
    element: <GuestRoute />,
    children: [{ path: '/login', element: <AuthPage /> }],
  },
  {
    element: <ProtectedRoute />,
    children: [{
      element: <AppShell />,
      children: [
        { path: '/dashboard', element: <DashboardPage /> },
        { path: '/categories', element: <CategoriesPage /> },
        { path: '/categories/create', element: <CategoryEditorPage /> },
        { path: '/categories/:id/edit', element: <CategoryEditorPage edit /> },
        { path: '/budgets', element: <BudgetsPage /> },
        { path: '/budgets/create', element: <BudgetEditorPage /> },
        { path: '/budgets/:id/edit', element: <BudgetEditorPage edit /> },
        { path: '/transactions', element: <TransactionsPage /> },
        { path: '/transactions/create', element: <TransactionEditorPage /> },
        { path: '/transactions/:id/edit', element: <TransactionEditorPage edit /> },
        { path: '/reports', element: <ReportsPage /> },
        { path: '/settings', element: <SettingsPage /> },
      ],
    }],
  },
  { path: '*', element: <Navigate to="/" replace /> },
])

export default function AppRouter() {
  return <RouterProvider router={router} />
}
