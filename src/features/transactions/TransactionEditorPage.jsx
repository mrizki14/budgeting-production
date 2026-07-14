import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import Flash from '../../components/ui/Flash'
import { apiPost, apiPut, apiRequest } from '../../lib/api'
import { dateInputValue } from '../../lib/format'
import TransactionForm from './TransactionForm'

export default function TransactionEditorPage({ edit = false }) {
  const { id } = useParams(); const navigate = useNavigate(); const [categories, setCategories] = useState(null); const [value, setValue] = useState({ type: 'income', category_id: '', amount: '', date: new Date().toISOString().slice(0, 10), description: '' }); const [errors, setErrors] = useState({}); const [message, setMessage] = useState(''); const [submitting, setSubmitting] = useState(false)
  useEffect(() => { Promise.all([apiRequest('/categories'), edit ? apiRequest(`/transactions/${id}`) : Promise.resolve(null)]).then(([items, transaction]) => { setCategories(items); if (transaction) setValue({ type: transaction.type, category_id: String(transaction.category_id), amount: transaction.amount, date: dateInputValue(transaction.date), description: transaction.description ?? '' }) }).catch((e) => setMessage(e.message)) }, [edit, id])
  async function submit(event) { event.preventDefault(); setSubmitting(true); setErrors({}); const payload = { type: value.type, category_id: Number(value.category_id), amount: Number(value.amount), date: value.date, description: value.description }; try { if (edit) await apiPut(`/transactions/${id}`, payload); else await apiPost('/transactions', payload); navigate('/transactions', { state: { status: edit ? 'Transaction updated successfully' : 'Transaction created successfully' } }) } catch (e) { setErrors(e.errors ?? {}); setMessage(e.message) } finally { setSubmitting(false) } }
  if (!categories) return <div className="app-loading">Loading...</div>
  return <section className="mx-auto max-w-3xl px-4 py-6 md:px-8 md:py-8"><div className="budget-panel"><h1 className="text-3xl font-bold text-slate-900">{edit ? 'Edit Transaction' : 'Add Transaction'}</h1><p className="mt-2 text-sm text-slate-500">{edit ? 'Perbarui detail transaksi Anda.' : 'Tambahkan pemasukan atau pengeluaran baru.'}</p>{message && <div className="mt-5"><Flash message={message} tone="error" /></div>}<form onSubmit={submit} className="mt-8 space-y-5"><TransactionForm value={value} setValue={setValue} categories={categories} errors={errors} submitLabel={edit ? 'Update Transaction' : 'Save Transaction'} submitting={submitting} /></form></div></section>
}
