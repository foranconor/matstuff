import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { get, post } from '../api'
import Modal from '../components/Modal'

export default function Suppliers() {
  const [suppliers, setSuppliers] = useState([])
  const [contacts, setContacts] = useState([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ Name: '', Website: '', ContactID: '' })
  const navigate = useNavigate()

  useEffect(() => {
    Promise.all([get('/suppliers'), get('/contacts')])
      .then(([sups, conts]) => {
        setSuppliers(sups ?? [])
        setContacts(conts ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  function setField(field) {
    return e => setForm(f => ({ ...f, [field]: e.target.value }))
  }

  async function handleCreate(e) {
    e.preventDefault()
    const s = await post('/suppliers', { ...form, ContactID: Number(form.ContactID) })
    setShowModal(false)
    setForm({ Name: '', Website: '', ContactID: '' })
    navigate(`/suppliers/${s.ID}`)
  }

  const filtered = suppliers.filter(s =>
    s.Name.toLowerCase().includes(search.toLowerCase()) ||
    s.ContactName.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <h1>Suppliers</h1>
        <input
          type="search"
          placeholder="Search..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <button className="btn-primary" onClick={() => setShowModal(true)}>+ New</button>
      </header>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Website</th>
            <th>Contact</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map(s => (
            <tr key={s.ID} onClick={() => navigate(`/suppliers/${s.ID}`)}>
              <td>{s.Name}</td>
              <td>{s.Website}</td>
              <td>{s.ContactName}</td>
            </tr>
          ))}
          {filtered.length === 0 && (
            <tr><td colSpan={3}>No suppliers found</td></tr>
          )}
        </tbody>
      </table>
      {showModal && (
        <Modal title="New Supplier" onClose={() => setShowModal(false)} onSubmit={handleCreate}>
          <label className="modal-label">
            Name
            <input autoFocus required value={form.Name} onChange={setField('Name')} />
          </label>
          <label className="modal-label">
            Website
            <input required value={form.Website} onChange={setField('Website')} />
          </label>
          <label className="modal-label">
            Contact
            <select required value={form.ContactID} onChange={setField('ContactID')}>
              <option value="">Select a contact...</option>
              {contacts.map(c => (
                <option key={c.ID} value={c.ID}>{c.Name}</option>
              ))}
            </select>
          </label>
        </Modal>
      )}
    </div>
  )
}
