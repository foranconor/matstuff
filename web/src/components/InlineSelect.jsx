export default function InlineSelect({ value, options, onSave }) {
  return (
    <select value={value} onChange={e => onSave(e.target.value)}>
      {options.map(o => <option key={o} value={o}>{o}</option>)}
    </select>
  )
}
