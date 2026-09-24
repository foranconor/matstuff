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
  const [allSuppliers, setAllSuppliers] = useState([])
  const [linkedProfiles, setLinkedProfiles] = useState([])
  const [allProfiles, setAllProfiles] = useState([])
  const [notes, setNotes] = useState([])
  const [tab, setTab] = useState('specs')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [newNote, setNewNote] = useState('')
  const [newSupplierID, setNewSupplierID] = useState('')
  const [newProfileID, setNewProfileID] = useState('')

  useEffect(() => {
    Promise.all([
      get(`/materials/${id}`),
      get(`/materials/${id}/suppliers`),
      get(`/materials/${id}/notes`),
      get('/suppliers'),
      get(`/materials/${id}/profiles`),
      get('/profiles'),
    ])
      .then(([mwu, sups, nts, allSups, profs, allProfs]) => {
        setMaterial(mwu.Material)
        setUses(mwu.Uses)
        setSuppliers(sups ?? [])
        setNotes(nts ?? [])
        setAllSuppliers(allSups ?? [])
        setLinkedProfiles(profs ?? [])
        setAllProfiles(allProfs ?? [])
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

  async function addSupplier() {
    if (!newSupplierID) { return }
    const ms = await post(`/materials/${id}/suppliers`, { SupplierID: Number(newSupplierID) })
    setSuppliers(s => [...s, ms])
    setNewSupplierID('')
  }

  async function saveSupplierLink(stmID, field, value) {
    const current = suppliers.find(s => s.ID === stmID)
    const updated = { ...current, [field]: value }
    await put(`/supplier-materials/${stmID}`, updated)
    setSuppliers(s => s.map(x => x.ID === stmID ? updated : x))
  }

  async function removeSupplierLink(stmID) {
    if (!confirm('Remove this supplier link?')) { return }
    await del(`/supplier-materials/${stmID}`)
    setSuppliers(s => s.filter(x => x.ID !== stmID))
  }

  async function addProfile() {
    if (!newProfileID) { return }
    const mp = await post(`/materials/${id}/profiles`, { ProfileID: Number(newProfileID) })
    setLinkedProfiles(p => [...p, mp])
    setNewProfileID('')
  }

  async function saveProfileLink(mpID, value) {
    const current = linkedProfiles.find(p => p.ID === mpID)
    const updated = { ...current, Price: value }
    await put(`/material-profiles/${mpID}`, updated)
    setLinkedProfiles(p => p.map(x => x.ID === mpID ? updated : x))
  }

  async function removeProfileLink(mpID) {
    if (!confirm('Remove this profile link?')) { return }
    await del(`/material-profiles/${mpID}`)
    setLinkedProfiles(p => p.filter(x => x.ID !== mpID))
  }

  async function addNote() {
    if (!newNote.trim()) { return }
    const note = await post(`/materials/${id}/notes`, { Content: newNote })
    setNotes(n => [note, ...n])
    setNewNote('')
  }

  async function deleteNote(noteID) {
    if (!confirm('Delete this note?')) { return }
    await del(`/materials/${id}/notes/${noteID}`)
    setNotes(n => n.filter(x => x.ID !== noteID))
  }

  async function saveNote(noteID, content) {
    const note = await put(`/materials/${id}/notes/${noteID}`, { Content: content })
    setNotes(n => n.map(x => x.ID === noteID ? note : x))
  }

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  let title = material.Name
  if (material.Nickname) { title = `${material.Name} (${material.Nickname})` }

  const usesFlags = [
    'Exterior', 'Interior', 'EarlyAccess', 'Stringers',
    'Risers', 'Treads', 'Timber', 'Panel', 'Handrail',
  ]

  return (
    <div className="page">
      <header>
        <button className="btn-back" onClick={() => navigate('/materials')}>←Back</button>
        <h1>{title}</h1>
        <div className="header-actions">
          <span className={`badge ${uses.Archived ? 'archived' : ''}`}>
            {uses.Archived ? 'Archived' : uses.Published ? 'Published' : 'Draft'}
          </span>
        </div>
      </header>

      <div className="tabs">
        {['specs', 'uses', 'suppliers', 'profiles', 'notes'].map(t => (
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
          <div className="field-row">
            <label>Color</label>
            <div className="color-field">
              <span className="color-swatch" style={{ background: material.Color }} />
              <InlineField value={material.Color} onSave={v => saveField('Color', v)} />
            </div>
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
        <div>
          <div className="link-add">
            <select value={newSupplierID} onChange={e => setNewSupplierID(e.target.value)}>
              <option value="">— select supplier —</option>
              {allSuppliers.map(s => (
                <option key={s.ID} value={s.ID}>{s.Name}</option>
              ))}
            </select>
            <button className="btn-primary" onClick={addSupplier} disabled={!newSupplierID}>
              Link Supplier
            </button>
          </div>
          <table>
            <thead>
              <tr>
                <th>Supplier</th>
                <th>Priority</th>
                <th>Price</th>
                <th>Lead Time (working days)</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {suppliers.map(s => (
                <tr key={s.ID}>
                  <td>{s.SupplierName}</td>
                  <td><InlineField value={s.Priority} type="number" onSave={v => saveSupplierLink(s.ID, 'Priority', v)} /></td>
                  <td><InlineField value={s.Price} type="number" onSave={v => saveSupplierLink(s.ID, 'Price', v)} /></td>
                  <td><InlineField value={s.LeadTime} type="number" onSave={v => saveSupplierLink(s.ID, 'LeadTime', v)} /></td>
                  <td><button className="btn-danger-sm" onClick={() => removeSupplierLink(s.ID)}>×</button></td>
                </tr>
              ))}
              {suppliers.length === 0 && (
                <tr><td colSpan={5}>No suppliers linked</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {tab === 'profiles' && (
        <div>
          <div className="link-add">
            <select value={newProfileID} onChange={e => setNewProfileID(e.target.value)}>
              <option value="">— select profile —</option>
              {allProfiles.map(p => (
                <option key={p.ID} value={p.ID}>{p.Name}</option>
              ))}
            </select>
            <button className="btn-primary" onClick={addProfile} disabled={!newProfileID}>
              Link Profile
            </button>
          </div>
          <table>
            <thead>
              <tr>
                <th>Profile</th>
                <th>Price</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {linkedProfiles.map(p => (
                <tr key={p.ID}>
                  <td>{p.ProfileName}</td>
                  <td><InlineField value={p.Price} type="number" onSave={v => saveProfileLink(p.ID, v)} /></td>
                  <td><button className="btn-danger-sm" onClick={() => removeProfileLink(p.ID)}>×</button></td>
                </tr>
              ))}
              {linkedProfiles.length === 0 && (
                <tr><td colSpan={3}>No profiles linked</td></tr>
              )}
            </tbody>
          </table>
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
