import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { get, put, post, del } from '../api'
import InlineField from '../components/InlineField'
import InlineSelect from '../components/InlineSelect'

const TREATMENTS = ['UT', 'FDA', 'MCA', 'H1.2', 'H3.2', 'H4', 'H5']
const UNITS = ['mm', 'sheet']

export default function MaterialDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [material, setMaterial] = useState(null)
  const [uses, setUses] = useState(null)
  const [suppliers, setSuppliers] = useState([])
  const [notes, setNotes] = useState([])
  const [tab, setTab] = useState('specs')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [newNote, setNewNote] = useState('')

  useEffect(() => {
    Promise.all([
      get(`/materials/${id}`),
      get(`/materials/${id}/suppliers`),
      get(`/materials/${id}/notes`),
    ])
      .then(([mwu, sups, nts]) => {
        setMaterial(mwu.Material)
        setUses(mwu.Uses)
        setSuppliers(sups ?? [])
        setNotes(nts ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  async function saveField(field, value) {
    const updated = { ...material, [field]: value }
    await put(`/materials/${id}`, updated)
    setMaterial(updated)
  }

  async function saveUses(field, value) {
    const updated = { ...uses, [field]: value }
    await put(`/materials/${id}/uses`, updated)
    setUses(updated)
  }

  async function addNote() {
    if (!newNote.trim()) return
    const note = await post(`/materials/${id}/notes`, { Content: newNote })
    setNotes(n => [note, ...n])
    setNewNote('')
  }

  async function deleteNote(noteID) {
    if (!confirm('Delete this note?')) return
    await del(`/materials/${id}/notes/${noteID}`)
    setNotes(n => n.filter(x => x.ID !== noteID))
  }

  async function saveNote(noteID, content) {
    const note = await put(`/materials/${id}/notes/${noteID}`, { Content: content })
    setNotes(n => n.map(x => x.ID === noteID ? note : x))
  }

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  const usesFlags = [
    'Exterior', 'Interior', 'EarlyAccess', 'Stringers',
    'Risers', 'Treads', 'Timber', 'Panel', 'Handrail',
  ]

  return (
    <div className="page">
      <header>
        <button className="btn-back" onClick={() => navigate('/materials')}>←Back</button>
        <h1>{material.Name}</h1>
        <div className="header-actions">
          <span className={`badge ${uses.Archived ? 'archived' : ''}`}>
            {uses.Archived ? 'Archived' : uses.Published ? 'Published' : 'Draft'}
          </span>
        </div>
      </header>

      <div className="tabs">
        {['specs', 'uses', 'suppliers', 'notes'].map(t => (
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
          {[
            ['Name', 'text'], ['Nickname', 'text'],
            ['Length', 'number'], ['Width', 'number'],
            ['Thickness', 'number'], ['MaxSpan', 'number'], ['Density', 'number'],
            ['MaxOverhang', 'number'], ['Radius', 'number'],
          ].map(([field, type]) => (
            <div key={field} className="field-row">
              <label>{field}</label>
              <InlineField value={material[field]} type={type} onSave={v => saveField(field, v)} />
            </div>
          ))}
          <div className="field-row">
            <label>Treatment</label>
            <InlineSelect value={material.Treatment} options={TREATMENTS} onSave={v => saveField('Treatment', v)} />
          </div>
          <div className="field-row">
            <label>Units</label>
            <InlineSelect value={material.Units} options={UNITS} onSave={v => saveField('Units', v)} />
          </div>
          <div className="field-row field-row--full">
            <label>Blurb</label>
            <InlineField value={material.Blurb} onSave={v => saveField('Blurb', v)} />
          </div>
        </div>
      )}

      {tab === 'uses' && (
        <div className="uses-grid">
          {usesFlags.map(flag => (
            <label key={flag} className="toggle-row">
              <input
                type="checkbox"
                checked={uses[flag]}
                onChange={e => saveUses(flag, e.target.checked)}
              />
              {flag}
            </label>
          ))}
          <div className="uses-divider" />
          <label className="toggle-row">
            <input
              type="checkbox"
              checked={uses.Published}
              onChange={e => saveUses('Published', e.target.checked)}
            />
            Published
          </label>
          <label className="toggle-row">
            <input
              type="checkbox"
              checked={uses.Archived}
              onChange={e => saveUses('Archived', e.target.checked)}
            />
            Archived
          </label>
        </div>
      )}

      {tab === 'suppliers' && (
        <table>
          <thead>
            <tr>
              <th>Supplier</th>
              <th>Priority</th>
              <th>Price</th>
              <th>Lead Time</th>
            </tr>
          </thead>
          <tbody>
            {suppliers.map(s => (
              <tr key={s.ID}>
                <td>{s.SupplierName}</td>
                <td>{s.Priority}</td>
                <td>{s.Price}</td>
                <td>{s.LeadTime}</td>
              </tr>
            ))}
            {suppliers.length === 0 && (
              <tr><td colSpan={4}>No suppliers linked</td></tr>
            )}
          </tbody>
        </table>
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
