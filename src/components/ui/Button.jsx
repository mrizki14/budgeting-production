export default function Button({ tone = 'primary', className = '', type = 'button', ...props }) {
  return <button type={type} className={`budget-button budget-button-${tone} ${className}`} {...props} />
}
