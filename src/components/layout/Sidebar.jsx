import { NavLink } from 'react-router-dom'
import { useAuth } from '../../features/auth/AuthContext'

const items = [
  ['dashboard', 'Dashboard', '/dashboard', 'DB'],
  ['categories', 'Categories', '/categories', 'CT'],
  ['transactions', 'Transactions', '/transactions', 'TR'],
  ['budgets', 'Budgets', '/budgets', 'BG'],
  ['reports', 'Reports', '/reports', 'RP'],
  ['settings', 'Settings', '/settings', 'ST'],
]

export default function Sidebar() {
  const { user, logout } = useAuth()
  return (
    <aside className="hidden border-r border-slate-200 bg-white lg:flex lg:flex-col">
      <div className="flex h-20 items-center gap-3 border-b border-slate-200 px-6">
        <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-blue-600 text-sm font-bold text-white shadow-sm">B</div>
        <div>
          <p className="text-lg font-semibold text-slate-900">Budget</p>
          <p className="text-sm text-slate-500">Manage your finances</p>
        </div>
      </div>
      <nav className="flex-1 space-y-2 px-4 py-6">
        {items.map(([key, label, to, short]) => (
          <NavLink
            key={key}
            to={to}
            className={({ isActive }) => `flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-medium transition ${isActive ? 'bg-blue-50 text-blue-700 shadow-sm ring-1 ring-blue-100' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}`}
          >
            <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-slate-100 text-xs font-semibold text-slate-500">{short}</span>
            <span>{label}</span>
          </NavLink>
        ))}
      </nav>
      <div className="space-y-4 border-t border-slate-200 px-6 py-5">
        <div className="space-y-1 text-sm text-slate-500">
          <p className="font-medium text-slate-700">{user?.name}</p>
          <p>{user?.email}</p>
        </div>
        <button type="button" onClick={logout} className="budget-button budget-button-danger w-full">Logout</button>
      </div>
    </aside>
  )
}
