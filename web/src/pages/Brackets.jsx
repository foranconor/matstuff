import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { get, post } from '../api'
import Modal from '../components/Modal'

export default function Brackets() {
  const [brackets, setBrackets] = useState([])
  const [suppliers, setSuppliers] = useState([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showArchived, setShowArchived] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [newName, setNewName] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    const url = '/brackets' + (showArchived ? '?all=true' : '')
    Promise.all([get(url), get('/suppliers')])
      .then(([brks, sups]) => {
        setBrackets(brks ?? [])
        setSuppliers(sups ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [showArchived])

  async function handleCreate(e) {
    e.preventDefault()
    const b = await post('/brackets', { Name: newName })
    setShowModal(false)
    setNewName('')
    navigate(`/brackets/${b.ID}`)
  }

  const filtered = brackets.filter(b =>
    b.Name.toLowerCase().includes(search.toLowerCase()) ||
    b.Nickname.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <h1>Brackets</h1>
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
            <th>Supplier</th>
            <th>Price</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map(b => (
            <tr key={b.ID} className={b.Archived ? 'row-archived' : ''} onClick={() => navigate(`/brackets/${b.ID}`)}>
              <td>{b.Name}</td>
              <td>{b.Nickname}</td>
              <td>{b.SupplierName}</td>
              <td>{b.Price}</td>
            </tr>
          ))}
          {filtered.length === 0 && (
            <tr><td colSpan={4}>No brackets found</td></tr>
          )}
        </tbody>
      </table>
      {showModal && (
        <Modal title="New Bracket" onClose={() => setShowModal(false)} onSubmit={handleCreate}>
          <label className="modal-label">
            Name
            <input autoFocus required value={newName} onChange={e => setNewName(e.target.value)} />
          </label>
        </Modal>
      )}
    </div>
  )
}
