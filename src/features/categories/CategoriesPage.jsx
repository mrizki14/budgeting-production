import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import Flash from '../../components/ui/Flash'
import { apiDelete, apiRequest } from '../../lib/api'

function createdAt(value) {
  const date = new Date(value)
  const day = String(date.getUTCDate()).padStart(2, '0')
  const month = date.toLocaleString('en-US', { month: 'short', timeZone: 'UTC' })
  return `${day} ${month} ${date.getUTCFullYear()} ${String(date.getUTCHours()).padStart(2, '0')}:${String(date.getUTCMinutes()).padStart(2, '0')}`
}

export default function CategoriesPage() {
  const [categories, setCategories] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    apiRequest('/categories').then(setCategories).catch((e) => setError(e.message))
  }, [])

  async function destroy(category) {
    if (!window.confirm('Hapus kategori ini?')) return
    try {
      await apiDelete(`/categories/${category.id}`)
      setCategories((items) => items.filter((item) => item.id !== category.id))
    } catch (e) {
      setError(e.message)
    }
  }

  return (
    <section className="mx-auto max-w-6xl space-y-8 px-4 py-6 md:px-8 md:py-8">
      <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div><h1 className="text-3xl font-bold text-slate-900">Categories</h1><p className="mt-1 text-sm text-slate-500">Kelola kategori pemasukan dan pengeluaran Anda.</p></div>
        <Link to="/categories/create" className="budget-button budget-button-primary">Add Category</Link>
      </div>
      {error && <Flash message={error} tone="error" />}
      {!categories ? <div className="budget-panel text-sm text-slate-500">Loading...</div> : (
        <section className="budget-panel overflow-hidden p-0"><div className="overflow-x-auto"><table className="min-w-full divide-y divide-slate-200 text-sm">
          <thead className="bg-slate-50 text-left text-xs font-semibold uppercase tracking-[0.18em] text-slate-500"><tr><th className="px-6 py-4">Name</th><th className="px-6 py-4">Type</th><th className="px-6 py-4">Created</th><th className="px-6 py-4 text-right">Action</th></tr></thead>
          <tbody className="divide-y divide-slate-200 bg-white">{categories.length ? categories.map((category) => (
            <tr key={category.id} className="hover:bg-slate-50/80"><td className="px-6 py-4 font-medium text-slate-800">{category.name}</td><td className="px-6 py-4"><span className={`budget-badge ${category.type === 'income' ? 'budget-badge-success' : 'budget-badge-danger'}`}>{category.type === 'income' ? 'Income' : 'Expense'}</span></td><td className="px-6 py-4 text-slate-500">{createdAt(category.created_at)}</td><td className="px-6 py-4"><div className="flex justify-end gap-3"><Link to={`/categories/${category.id}/edit`} className="budget-button budget-button-secondary px-4 py-2">Edit</Link><button onClick={() => destroy(category)} className="budget-button budget-button-danger px-4 py-2">Delete</button></div></td></tr>
          )) : <tr><td colSpan="4" className="px-6 py-10 text-center text-sm text-slate-500">Belum ada kategori. Tambahkan kategori pertama Anda.</td></tr>}</tbody>
        </table></div></section>
      )}
    </section>
  )
}
