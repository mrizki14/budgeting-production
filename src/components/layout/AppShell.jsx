import { Outlet, useLocation } from 'react-router-dom'
import Flash from '../ui/Flash'
import MobileHeader from './MobileHeader'
import Sidebar from './Sidebar'

export default function AppShell() {
  const location = useLocation()
  return (
    <div className="min-h-screen bg-slate-100 text-slate-900 antialiased lg:grid lg:grid-cols-[17rem_minmax(0,1fr)]">
      <Sidebar />
      <div className="flex min-h-screen flex-col">
        <MobileHeader />
        <main className="flex-1">
          {location.state?.status && (
            <div className="mx-auto max-w-7xl px-4 pt-6 md:px-8">
              <Flash message={location.state.status} />
            </div>
          )}
          <Outlet />
        </main>
      </div>
    </div>
  )
}
