import { useState, useEffect } from 'react'

export default function InlineField({ value, onSave, type = 'text' }) {
  const [val, setVal] = useState(value)

  useEffect(() => { setVal(value) }, [value])

  function handleChange(e) {
    const next = type === 'number' ? Number(e.target.value) : e.target.value
    setVal(next)
  }

  function handleBlur() {
    if (val !== value) onSave(val)
  }

  return (
    <input
      type={type}
      value={val}
      onChange={handleChange}
      onBlur={handleBlur}
    />
  )
}
