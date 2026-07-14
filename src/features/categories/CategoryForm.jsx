import Field from '../../components/ui/Field'

export default function CategoryForm({ value, setValue, errors }) {
  return <><Field label="Category Name" error={errors.name?.[0]}><input type="text" value={value.name} onChange={(e) => setValue({ ...value, name: e.target.value })} placeholder="Contoh: Groceries" required /></Field><Field label="Type" error={errors.type?.[0]}><select value={value.type} onChange={(e) => setValue({ ...value, type: e.target.value })} required><option value="" disabled>Pilih tipe kategori</option><option value="income">Income</option><option value="expense">Expense</option></select></Field></>
}
