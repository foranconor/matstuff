import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { get, post } from '../api'
import Modal from '../components/Modal'

export default function Contacts() {
  const [contacts, setContacts] = useState([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ Name: '', Phone: '', Email: '' })
  const navigate = useNavigate()

  useEffect(() => {
    get('/contacts')
      .then(data => setContacts(data ?? []))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  function setField(field) {
    return e => setForm(f => ({ ...f, [field]: e.target.value }))
  }

  async function handleCreate(e) {
    e.preventDefault()
    const c = await post('/contacts', form)
    setShowModal(false)
    setForm({ Name: '', Phone: '', Email: '' })
    navigate(`/contacts/${c.ID}`)
  }

  const filtered = contacts.filter(c =>
    c.Name.toLowerCase().includes(search.toLowerCase()) ||
    c.Email.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <h1>Contacts</h1>
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
            <th>Phone</th>
            <th>Email</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map(c => (
            <tr key={c.ID} onClick={() => navigate(`/contacts/${c.ID}`)}>
              <td>{c.Name}</td>
              <td>{c.Phone}</td>
              <td>{c.Email}</td>
            </tr>
          ))}
          {filtered.length === 0 && (
            <tr><td colSpan={3}>No contacts found</td></tr>
          )}
        </tbody>
      </table>
      {showModal && (
        <Modal title="New Contact" onClose={() => setShowModal(false)} onSubmit={handleCreate}>
          {['Name', 'Phone', 'Email'].map(field => (
            <label key={field} className="modal-label">
              {field}
              <input autoFocus={field === 'Name'} required value={form[field]} onChange={setField(field)} />
            </label>
          ))}
        </Modal>
      )}
    </div>
  )
}
