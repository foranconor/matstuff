import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { get, post } from '../api'
import Modal from '../components/Modal'

export default function Profiles() {
  const [profiles, setProfiles] = useState([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [newName, setNewName] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    get('/profiles')
      .then(p => setProfiles(p ?? []))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  async function handleCreate(e) {
    e.preventDefault()
    const p = await post('/profiles', { Name: newName })
    setShowModal(false)
    setNewName('')
    navigate(`/profiles/${p.ID}`)
  }

  const filtered = profiles.filter(p =>
    p.Name.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <h1>Profiles</h1>
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
          </tr>
        </thead>
        <tbody>
          {filtered.map(p => (
            <tr key={p.ID} onClick={() => navigate(`/profiles/${p.ID}`)}>
              <td>{p.Name}</td>
            </tr>
          ))}
          {filtered.length === 0 && (
            <tr><td>No profiles found</td></tr>
          )}
        </tbody>
      </table>
      {showModal && (
        <Modal title="New Profile" onClose={() => setShowModal(false)} onSubmit={handleCreate}>
          <label className="modal-label">
            Name
            <input autoFocus required value={newName} onChange={e => setNewName(e.target.value)} />
          </label>
        </Modal>
      )}
    </div>
  )
}
