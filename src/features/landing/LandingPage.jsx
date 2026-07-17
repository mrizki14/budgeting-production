import { Link } from 'react-router-dom'

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 antialiased">
      <header className="border-b border-slate-200 bg-white/90">
        <nav className="mx-auto flex h-20 max-w-7xl items-center justify-between px-4 md:px-8">
          <Link to="/" className="flex items-center gap-3">
            <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-blue-600 text-sm font-bold text-white shadow-sm">B</span>
            <span><span className="block text-lg font-semibold text-slate-900">Budget</span><span className="block text-xs font-medium text-slate-500">Personal finance dashboard</span></span>
          </Link>
          <div className="flex items-center gap-3">
            <a href="#fitur" className="hidden text-sm font-semibold text-slate-600 hover:text-slate-900 sm:inline">Fitur</a>
            <Link to="/login" className="budget-button budget-button-primary">Mulai Sekarang</Link>
          </div>
        </nav>
      </header>
      <main>
        <section className="mx-auto grid max-w-7xl items-center gap-12 px-4 py-14 md:px-8 lg:grid-cols-[minmax(0,1fr)_minmax(24rem,32rem)] lg:py-20">
          <div>
            <p className="budget-pill bg-blue-50 text-blue-700">Budgeting App</p>
            <h1 className="mt-6 max-w-3xl text-4xl font-bold tracking-tight text-slate-900 sm:text-5xl lg:text-6xl">Kelola budget pribadi dengan lebih rapi.</h1>
            <p className="mt-6 max-w-2xl text-base leading-7 text-slate-600 sm:text-lg">Pantau pemasukan, pengeluaran, kategori, dan budget bulanan dalam satu dashboard yang bersih dan mudah dipahami.</p>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row"><Link to="/login" className="budget-button budget-button-primary">Mulai Sekarang</Link><a href="#fitur" className="budget-button budget-button-secondary">Lihat Fitur</a></div>
          </div>
          <div className="rounded-[2rem] border border-slate-200 bg-white p-5 shadow-xl shadow-slate-200/70">
            <div className="rounded-[1.5rem] bg-slate-950 p-5 text-white"><div className="flex items-center justify-between"><div><p className="text-sm text-slate-300">Total Balance</p><p className="mt-2 text-3xl font-bold">Rp 12.450.000</p></div><span className="rounded-full bg-emerald-400/15 px-3 py-1 text-xs font-semibold text-emerald-200">+12%</span></div></div>
            <div className="mt-5 grid gap-4 sm:grid-cols-2"><div className="rounded-3xl border border-emerald-100 bg-emerald-50 p-4"><p className="text-sm font-medium text-emerald-700">Income</p><p className="mt-2 text-xl font-bold text-emerald-950">Rp 8.200.000</p></div><div className="rounded-3xl border border-rose-100 bg-rose-50 p-4"><p className="text-sm font-medium text-rose-700">Expenses</p><p className="mt-2 text-xl font-bold text-rose-950">Rp 3.150.000</p></div></div>
            <div className="mt-5 space-y-4"><div><div className="flex items-center justify-between text-sm"><span className="font-semibold text-slate-700">Dining Budget</span><span className="text-slate-500">62%</span></div><div className="budget-progress-track mt-2"><div className="budget-progress-fill bg-blue-500" style={{ width: '62%' }} /></div></div><div className="rounded-3xl border border-slate-200 p-4"><div className="flex items-center justify-between gap-4"><div><p className="font-semibold text-slate-800">Groceries</p><p className="text-sm text-slate-500">Pengeluaran terbaru</p></div><p className="font-semibold text-rose-600">-Rp 250.000</p></div></div></div>
          </div>
        </section>
        <section id="fitur" className="border-y border-slate-200 bg-white py-16"><div className="mx-auto max-w-7xl px-4 md:px-8"><div className="max-w-2xl"><h2 className="text-3xl font-bold tracking-tight text-slate-900">Fitur utama untuk rutinitas finansial Anda.</h2><p className="mt-3 text-sm leading-6 text-slate-500">Semua yang dibutuhkan untuk melihat arus uang harian dan menjaga budget tetap terkendali.</p></div><div className="mt-8 grid gap-4 md:grid-cols-3">{[
          ['Transaksi', 'Catat pemasukan dan pengeluaran', 'Simpan setiap transaksi dengan kategori agar alur uang lebih mudah dilacak.'],
          ['Budget', 'Atur budget per kategori', 'Tetapkan batas bulanan dan pantau progres pemakaian dari dashboard.'],
          ['Laporan', 'Pantau laporan bulanan', 'Lihat ringkasan income, expenses, dan balance untuk mengambil keputusan lebih cepat.'],
        ].map(([label, title, copy]) => <article key={label} className="budget-panel"><p className="text-sm font-semibold text-blue-600">{label}</p><h3 className="mt-3 text-lg font-bold text-slate-900">{title}</h3><p className="mt-3 text-sm leading-6 text-slate-500">{copy}</p></article>)}</div></div></section>
        <section className="mx-auto max-w-7xl px-4 py-14 md:px-8"><div className="flex flex-col gap-5 rounded-[2rem] border border-slate-200 bg-blue-600 p-6 text-white shadow-lg shadow-blue-200/70 md:flex-row md:items-center md:justify-between md:p-8"><div><h2 className="text-2xl font-bold tracking-tight">Mulai rapikan budget hari ini.</h2><p className="mt-2 max-w-2xl text-sm leading-6 text-blue-100">Masuk atau buat akun dari halaman auth yang sama, lalu mulai catat transaksi pertama Anda.</p></div><Link to="/login" className="budget-button bg-white text-blue-700 hover:bg-blue-50">Mulai Sekarang</Link></div></section>
      </main>
    </div>
  )
}
