import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { get, put, post, del } from '../api'
import InlineField from '../components/InlineField'

export default function BracketDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [bracket, setBracket] = useState(null)
  const [notes, setNotes] = useState([])
  const [tab, setTab] = useState('specs')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [newNote, setNewNote] = useState('')

  useEffect(() => {
    Promise.all([
      get(`/brackets/${id}`),
      get(`/brackets/${id}/notes`),
    ])
      .then(([b, nts]) => {
        setBracket(b)
        setNotes(nts ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  async function saveField(field, value) {
    const updated = { ...bracket, [field]: value }
    await put(`/brackets/${id}`, updated)
    setBracket(updated)
  }

  async function saveBracket(changes) {
    const updated = { ...bracket, ...changes }
    await put(`/brackets/${id}`, updated)
    setBracket(updated)
  }

  async function addNote() {
    if (!newNote.trim()) return
    const note = await post(`/brackets/${id}/notes`, { Content: newNote })
    setNotes(n => [note, ...n])
    setNewNote('')
  }

  async function deleteNote(noteID) {
    if (!confirm('Delete this note?')) return
    await del(`/brackets/${id}/notes/${noteID}`)
    setNotes(n => n.filter(x => x.ID !== noteID))
  }

  async function saveNote(noteID, content) {
    const note = await put(`/brackets/${id}/notes/${noteID}`, { Content: content })
    setNotes(n => n.map(x => x.ID === noteID ? note : x))
  }

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <button className="btn-back" onClick={() => navigate('/brackets')}>← Back</button>
        <h1>{bracket.Name}</h1>
        <div className="header-actions">
          <span className={`badge ${bracket.Archived ? 'archived' : ''}`}>
            {bracket.Archived ? 'Archived' : bracket.Published ? 'Published' : 'Draft'}
          </span>
        </div>
      </header>

      <div className="tabs">
        {['specs', 'notes'].map(t => (
          <button
            key={t}
            className={`tab ${tab === t ? 'active' : ''}`}
            onClick={() => setTab(t)}
          >
            {t.charAt(0).toUpperCase() + t.slice(1)}
          </button>
        ))}
      </div>

      {tab === 'specs' && (
        <div className="detail-grid">
          {[['Name', 'text'], ['Nickname', 'text'], ['Price', 'number']].map(([field, type]) => (
            <div key={field} className="field-row">
              <label>{field}</label>
              <InlineField value={bracket[field]} type={type} onSave={v => saveField(field, v)} />
            </div>
          ))}
          <div className="field-row">
            <label>Published</label>
            <input
              type="checkbox"
              checked={bracket.Published}
              onChange={e => saveBracket({ Published: e.target.checked })}
            />
          </div>
          <div className="field-row">
            <label>Archived</label>
            <input
              type="checkbox"
              checked={bracket.Archived}
              onChange={e => saveBracket({ Archived: e.target.checked })}
            />
          </div>
        </div>
      )}

      {tab === 'notes' && (
        <div className="notes">
          <div className="note-add">
            <textarea
              placeholder="Add a note..."
              value={newNote}
              onChange={e => setNewNote(e.target.value)}
            />
            <button className="btn-primary" onClick={addNote}>Add</button>
          </div>
          {notes.map(n => (
            <div key={n.ID} className="note">
              <InlineField value={n.Content} onSave={v => saveNote(n.ID, v)} />
              <button className="btn-danger-sm" onClick={() => deleteNote(n.ID)}>×</button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
