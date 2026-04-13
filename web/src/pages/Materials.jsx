import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { get, post } from '../api'
import Modal from '../components/Modal'

export default function Materials() {
  const [materials, setMaterials] = useState([])
  const [search, setSearch] = useState('')
  const [showArchived, setShowArchived] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [newName, setNewName] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    const url = '/materials' + (showArchived ? '?all=true' : '')
    get(url)
      .then(data => setMaterials(data ?? []))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [showArchived])

  async function handleCreate(e) {
    e.preventDefault()
    const mwu = await post('/materials', { Name: newName })
    setShowModal(false)
    setNewName('')
    navigate(`/materials/${mwu.Material.ID}`)
  }

  const filtered = materials.filter(m =>
    m.Name.toLowerCase().includes(search.toLowerCase()) ||
    m.Nickname.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <h1>Materials</h1>
        <input
          type="search"
          placeholder="Search..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <label>
          <input
            type="checkbox"
            checked={showArchived}
            onChange={e => setShowArchived(e.target.checked)}
          />
          {' '}Show archived
        </label>
        <button className="btn-primary" onClick={() => setShowModal(true)}>+ New</button>
      </header>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Nickname</th>
            <th>Treatment</th>
            <th>Thickness</th>
            <th>Units</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map(m => (
            <tr key={m.ID} className={m.Archived ? 'row-archived' : ''} onClick={() => navigate(`/materials/${m.ID}`)}>
              <td>{m.Name}</td>
              <td>{m.Nickname}</td>
              <td>{m.Treatment}</td>
              <td>{m.Thickness}</td>
              <td>{m.Units}</td>
            </tr>
          ))}
          {filtered.length === 0 && (
            <tr><td colSpan={5}>No materials found</td></tr>
          )}
        </tbody>
      </table>
      {showModal && (
        <Modal title="New Material" onClose={() => setShowModal(false)} onSubmit={handleCreate}>
          <label className="modal-label">
            Name
            <input autoFocus required value={newName} onChange={e => setNewName(e.target.value)} />
          </label>
        </Modal>
      )}
    </div>
  )
}
