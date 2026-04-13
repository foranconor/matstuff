import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { get, put } from '../api'
import InlineField from '../components/InlineField'

export default function SupplierDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [supplier, setSupplier] = useState(null)
  const [contact, setContact] = useState(null)
  const [materials, setMaterials] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    Promise.all([
      get(`/suppliers/${id}`),
      get(`/suppliers/${id}/materials`),
    ])
      .then(([data, mats]) => {
        setSupplier(data.Supplier)
        setContact(data.Contact)
        setMaterials(mats ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  async function saveSupplier(field, value) {
    const updated = { ...supplier, [field]: value }
    await put(`/suppliers/${id}`, updated)
    setSupplier(updated)
  }

  async function saveContact(field, value) {
    const updated = { ...contact, [field]: value }
    await put(`/contacts/${contact.ID}`, updated)
    setContact(updated)
  }

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <button className="btn-back" onClick={() => navigate('/suppliers')}>← Back</button>
        <h1>{supplier.Name}</h1>
      </header>

      <section className="detail-section">
        <h2>Supplier</h2>
        <div className="detail-grid">
          {['Name', 'Website'].map(field => (
            <div key={field} className="field-row">
              <label>{field}</label>
              <InlineField value={supplier[field]} onSave={v => saveSupplier(field, v)} />
            </div>
          ))}
        </div>
      </section>

      <section className="detail-section">
        <h2>Contact</h2>
        <div className="detail-grid">
          {['Name', 'Phone', 'Email'].map(field => (
            <div key={field} className="field-row">
              <label>{field}</label>
              <InlineField value={contact[field]} onSave={v => saveContact(field, v)} />
            </div>
          ))}
        </div>
      </section>

      <section className="detail-section">
        <h2>Materials</h2>
        <table>
          <thead>
            <tr>
              <th>Material</th>
              <th>Priority</th>
              <th>Price</th>
              <th>Lead Time</th>
            </tr>
          </thead>
          <tbody>
            {materials.map(m => (
              <tr key={m.ID} onClick={() => navigate(`/materials/${m.MaterialID}`)}>
                <td>{m.MaterialName}</td>
                <td>{m.Priority}</td>
                <td>{m.Price}</td>
                <td>{m.LeadTime}</td>
              </tr>
            ))}
            {materials.length === 0 && (
              <tr><td colSpan={4}>No materials linked</td></tr>
            )}
          </tbody>
        </table>
      </section>
    </div>
  )
}
