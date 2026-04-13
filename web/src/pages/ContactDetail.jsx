import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { get, put } from '../api'
import InlineField from '../components/InlineField'

export default function ContactDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [contact, setContact] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    get(`/contacts/${id}`)
      .then(data => setContact(data))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  async function saveField(field, value) {
    const updated = { ...contact, [field]: value }
    await put(`/contacts/${id}`, updated)
    setContact(updated)
  }

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <button className="btn-back" onClick={() => navigate('/contacts')}>← Back</button>
        <h1>{contact.Name}</h1>
      </header>
      <div className="detail-grid">
        {['Name', 'Phone', 'Email'].map(field => (
          <div key={field} className="field-row">
            <label>{field}</label>
            <InlineField value={contact[field]} onSave={v => saveField(field, v)} />
          </div>
        ))}
      </div>
    </div>
  )
}
