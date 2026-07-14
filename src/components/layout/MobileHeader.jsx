import { useAuth } from '../../features/auth/AuthContext'

export default function MobileHeader() {
  const { user, logout } = useAuth()
  return (
    <header className="sticky top-0 z-10 border-b border-slate-200 bg-white/90 backdrop-blur lg:hidden">
      <div className="flex h-16 items-center justify-between px-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-blue-600 text-sm font-bold text-white shadow-sm">B</div>
          <div>
            <p className="text-sm font-semibold text-slate-900">Budget</p>
            <p className="text-xs text-slate-500">{user?.name}</p>
          </div>
        </div>
        <button type="button" onClick={logout} className="rounded-full bg-rose-50 px-3 py-1 text-xs font-medium text-rose-600">Logout</button>
      </div>
    </header>
  )
}
