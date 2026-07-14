import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Field from '../../components/ui/Field'
import Flash from '../../components/ui/Flash'
import { apiPost } from '../../lib/api'
import { useAuth } from './AuthContext'

const initialLogin = { email: '', password: '', remember: false }
const initialRegister = { name: '', email: '', password: '', password_confirmation: '' }

export default function AuthPage() {
  const [mode, setMode] = useState('login')
  const [loginForm, setLoginForm] = useState(initialLogin)
  const [registerForm, setRegisterForm] = useState(initialRegister)
  const [errors, setErrors] = useState({})
  const [message, setMessage] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  function switchMode(nextMode) {
    setMode(nextMode)
    setErrors({})
    setMessage('')
  }

  async function submitLogin(event) {
    event.preventDefault()
    setSubmitting(true)
    setErrors({})
    setMessage('')
    try {
      const result = await apiPost('/auth/login', {
        email: loginForm.email,
        password: loginForm.password,
      })
      login(result.token, result.user, loginForm.remember)
      navigate('/dashboard', { replace: true })
    } catch (error) {
      setErrors(error.errors ?? {})
      setMessage(error.message ?? 'Login gagal.')
    } finally {
      setSubmitting(false)
    }
  }

  async function submitRegister(event) {
    event.preventDefault()
    if (registerForm.password !== registerForm.password_confirmation) {
      setErrors({ password_confirmation: ['Konfirmasi password tidak cocok.'] })
      return
    }
    setSubmitting(true)
    setErrors({})
    setMessage('')
    try {
      const result = await apiPost('/auth/register', {
        name: registerForm.name,
        email: registerForm.email,
        password: registerForm.password,
      })
      login(result.token, result.user, true)
      navigate('/dashboard', { replace: true })
    } catch (error) {
      setErrors(error.errors ?? {})
      setMessage(error.message ?? 'Pendaftaran gagal.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top,_#dbeafe,_#f8fafc_55%)] text-slate-900 antialiased">
      <main className="mx-auto flex min-h-screen w-full max-w-6xl items-center justify-center px-4 py-10">
        <div className="grid w-full overflow-hidden rounded-[2rem] border border-slate-200 bg-white shadow-xl shadow-slate-200/70 lg:grid-cols-[1.1fr_minmax(24rem,30rem)]">
          <section className="hidden bg-blue-600 px-10 py-12 text-white lg:flex lg:flex-col lg:justify-between">
            <div>
              <div className="flex h-14 w-14 items-center justify-center rounded-3xl bg-white/15 text-lg font-bold">B</div>
              <h1 className="mt-8 text-4xl font-bold tracking-tight">Budget Dashboard</h1>
              <p className="mt-4 max-w-md text-sm leading-6 text-blue-100">
                Kelola anggaran kuliah, pengeluaran harian, dan target finansial dalam satu dashboard yang rapi.
              </p>
            </div>
            <div className="space-y-4">
              <div className="rounded-[1.5rem] border border-white/15 bg-white/10 p-6">
                <p className="text-sm font-semibold">Akun demo seeder</p>
                <p className="mt-3 text-sm text-blue-100">Email: demo@budgeting-app.test</p>
                <p className="mt-1 text-sm text-blue-100">Password: password123</p>
              </div>
            </div>
          </section>
          <section className="px-6 py-8 sm:px-10 sm:py-12">
            <div className="mx-auto w-full max-w-md">
              <div className="mb-8">
                <div className="flex h-12 w-12 items-center justify-center rounded-3xl bg-blue-600 text-base font-bold text-white">B</div>
                <h1 className="mt-6 text-3xl font-bold tracking-tight text-slate-900">Akses Budget App</h1>
                <p className="budget-auth-copy">Masuk ke akun Anda atau daftar akun baru untuk mulai mengelola keuangan.</p>
              </div>
              {message && <Flash message={message} tone="error" />}
              <div className="budget-auth-switch">
                <button type="button" onClick={() => switchMode('login')} className={`budget-auth-switch-button ${mode === 'login' ? 'budget-auth-switch-button-active' : ''}`} aria-pressed={mode === 'login'}>Masuk</button>
                <button type="button" onClick={() => switchMode('register')} className={`budget-auth-switch-button ${mode === 'register' ? 'budget-auth-switch-button-active' : ''}`} aria-pressed={mode === 'register'}>Daftar</button>
              </div>
              <div className="mt-8">
                {mode === 'login' ? (
                  <form className="space-y-5" onSubmit={submitLogin}>
                    <div>
                      <h2 className="text-2xl font-bold tracking-tight text-slate-900">Masuk</h2>
                      <p className="budget-auth-copy">Gunakan email dan password untuk melanjutkan ke dashboard.</p>
                    </div>
                    <Field label="Email" error={errors.email?.[0]}>
                      <input type="email" value={loginForm.email} onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })} placeholder="nama@email.com" required autoFocus />
                    </Field>
                    <Field label="Password" error={errors.password?.[0]}>
                      <input type="password" value={loginForm.password} onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })} placeholder="Masukkan password" required />
                    </Field>
                    <label className="flex items-center gap-3 text-sm text-slate-600">
                      <input type="checkbox" checked={loginForm.remember} onChange={(e) => setLoginForm({ ...loginForm, remember: e.target.checked })} className="h-4 w-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                      <span>Ingat saya</span>
                    </label>
                    <button disabled={submitting} type="submit" className="budget-button budget-button-primary w-full">Login ke Dashboard</button>
                    <p className="text-center text-sm text-slate-500">Belum punya akun? <button type="button" className="font-semibold text-blue-600" onClick={() => switchMode('register')}>Buat Akun</button></p>
                  </form>
                ) : (
                  <form className="space-y-5" onSubmit={submitRegister}>
                    <div>
                      <h2 className="text-2xl font-bold tracking-tight text-slate-900">Daftar</h2>
                      <p className="budget-auth-copy">Buat akun baru dan langsung masuk ke dashboard budgeting Anda.</p>
                    </div>
                    <Field label="Nama Lengkap" error={errors.name?.[0]}>
                      <input type="text" value={registerForm.name} onChange={(e) => setRegisterForm({ ...registerForm, name: e.target.value })} placeholder="Nama Anda" required />
                    </Field>
                    <Field label="Email" error={errors.email?.[0]}>
                      <input type="email" value={registerForm.email} onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })} placeholder="nama@email.com" required />
                    </Field>
                    <Field label="Password" error={errors.password?.[0]}>
                      <input type="password" value={registerForm.password} onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })} placeholder="Minimal 8 karakter" required minLength={8} />
                    </Field>
                    <Field label="Konfirmasi Password" error={errors.password_confirmation?.[0]}>
                      <input type="password" value={registerForm.password_confirmation} onChange={(e) => setRegisterForm({ ...registerForm, password_confirmation: e.target.value })} placeholder="Ulangi password" required />
                    </Field>
                    <button disabled={submitting} type="submit" className="budget-button budget-button-primary w-full">Buat Akun</button>
                    <p className="text-center text-sm text-slate-500">Sudah punya akun? <button type="button" className="font-semibold text-blue-600" onClick={() => switchMode('login')}>Masuk</button></p>
                  </form>
                )}
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>
  )
}
