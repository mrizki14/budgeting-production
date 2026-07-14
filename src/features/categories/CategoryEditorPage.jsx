import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import Flash from '../../components/ui/Flash'
import { apiPost, apiPut, apiRequest } from '../../lib/api'
import CategoryForm from './CategoryForm'

export default function CategoryEditorPage({ edit = false }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const [value, setValue] = useState({ name: '', type: '' })
  const [errors, setErrors] = useState({})
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(edit)
  const [submitting, setSubmitting] = useState(false)
  useEffect(() => { if (edit) apiRequest(`/categories/${id}`).then((item) => setValue({ name: item.name, type: item.type })).catch((e) => setMessage(e.message)).finally(() => setLoading(false)) }, [edit, id])
  async function submit(event) { event.preventDefault(); setSubmitting(true); setErrors({}); try { if (edit) await apiPut(`/categories/${id}`, value); else await apiPost('/categories', value); navigate('/categories', { state: { status: edit ? 'Category updated successfully' : 'Category created successfully' } }) } catch (e) { setErrors(e.errors ?? {}); setMessage(e.message) } finally { setSubmitting(false) } }
  const title = edit ? 'Edit Category' : 'Create Category'
  const copy = edit ? 'Perbarui kategori sesuai kebutuhan pencatatan Anda.' : 'Tambahkan kategori baru untuk transaksi atau budget Anda.'
  if (loading) return <div className="app-loading">Loading...</div>
  return <section className="mx-auto max-w-3xl space-y-8 px-4 py-6 md:px-8 md:py-8"><div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between"><div><h1 className="text-3xl font-bold text-slate-900">{title}</h1><p className="mt-1 text-sm text-slate-500">{copy}</p></div><Link to="/categories" className="budget-button budget-button-secondary">Back to Categories</Link></div>{message && <Flash message={message} tone="error" />}<section className="budget-panel"><form onSubmit={submit} className="space-y-5"><CategoryForm value={value} setValue={setValue} errors={errors} /><div className="flex flex-col gap-3 pt-2 sm:flex-row"><button disabled={submitting} type="submit" className="budget-button budget-button-primary">{edit ? 'Update Category' : 'Save Category'}</button><Link to="/categories" className="budget-button budget-button-secondary">Cancel</Link></div></form></section></section>
}
