import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { get, put, post, del } from '../api'
import InlineField from '../components/InlineField'

export default function SupplierDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [supplier, setSupplier] = useState(null)
  const [contact, setContact] = useState(null)
  const [materials, setMaterials] = useState([])
  const [allMaterials, setAllMaterials] = useState([])
  const [allContacts, setAllContacts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [newMaterialID, setNewMaterialID] = useState('')

  useEffect(() => {
    Promise.all([
      get(`/suppliers/${id}`),
      get(`/suppliers/${id}/materials`),
      get('/materials?all=true'),
      get('/contacts'),
    ])
      .then(([data, mats, allMats, contacts]) => {
        setSupplier(data.Supplier)
        setContact(data.Contact)
        setMaterials(mats ?? [])
        setAllMaterials(allMats ?? [])
        setAllContacts(contacts ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  async function saveSupplier(field, value) {
    const updated = { ...supplier, [field]: value }
    await put(`/suppliers/${id}`, updated)
    setSupplier(updated)
  }

  async function swapContact(contactID) {
    const newContact = allContacts.find(c => c.ID === contactID)
    const updatedSupplier = { ...supplier, ContactID: contactID }
    await put(`/suppliers/${id}`, updatedSupplier)
    setSupplier(updatedSupplier)
    setContact(newContact)
  }

  async function saveContact(field, value) {
    const updated = { ...contact, [field]: value }
    await put(`/contacts/${contact.ID}`, updated)
    setContact(updated)
  }

  async function addMaterial() {
    if (!newMaterialID) { return }
    const ms = await post(`/suppliers/${id}/materials`, { MaterialID: Number(newMaterialID) })
    setMaterials(m => [...m, ms])
    setNewMaterialID('')
  }

  async function saveMaterialLink(stmID, field, value) {
    const current = materials.find(m => m.ID === stmID)
    const updated = { ...current, [field]: value }
    await put(`/supplier-materials/${stmID}`, updated)
    setMaterials(m => m.map(x => x.ID === stmID ? updated : x))
  }

  async function removeMaterialLink(stmID) {
    if (!confirm('Remove this material link?')) { return }
    await del(`/supplier-materials/${stmID}`)
    setMaterials(m => m.filter(x => x.ID !== stmID))
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
        <div className="link-add">
          <select value={contact.ID} onChange={e => swapContact(Number(e.target.value))}>
            {allContacts.map(c => (
              <option key={c.ID} value={c.ID}>{c.Name}</option>
            ))}
          </select>
        </div>
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
        <div className="link-add">
          <select value={newMaterialID} onChange={e => setNewMaterialID(e.target.value)}>
            <option value="">— select material —</option>
            {allMaterials.map(m => {
              let label = m.Name
              if (m.Nickname) { label = m.Nickname }
              return <option key={m.ID} value={m.ID}>{label}</option>
            })}
          </select>
          <button className="btn-primary" onClick={addMaterial} disabled={!newMaterialID}>
            Link Material
          </button>
        </div>
        <table>
          <thead>
            <tr>
              <th>Material</th>
              <th>Priority</th>
              <th>Price</th>
              <th>Lead Time (working days)</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {materials.map(m => (
              <tr key={m.ID}>
                <td className="cell-link" onClick={() => navigate(`/materials/${m.MaterialID}`)}>{m.MaterialName}</td>
                <td><InlineField value={m.Priority} type="number" onSave={v => saveMaterialLink(m.ID, 'Priority', v)} /></td>
                <td><InlineField value={m.Price} type="number" onSave={v => saveMaterialLink(m.ID, 'Price', v)} /></td>
                <td><InlineField value={m.LeadTime} type="number" onSave={v => saveMaterialLink(m.ID, 'LeadTime', v)} /></td>
                <td><button className="btn-danger-sm" onClick={() => removeMaterialLink(m.ID)}>×</button></td>
              </tr>
            ))}
            {materials.length === 0 && (
              <tr><td colSpan={5}>No materials linked</td></tr>
            )}
          </tbody>
        </table>
      </section>
    </div>
  )
}
