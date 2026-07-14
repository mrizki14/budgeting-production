export default function Field({ label, error, children }) {
  return (
    <label className="budget-field">
      <span>{label}</span>
      {children}
      {error && <span className="mt-2 block text-sm text-rose-600">{error}</span>}
    </label>
  )
}
