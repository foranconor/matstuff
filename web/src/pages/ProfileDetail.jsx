import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { get, put, post, del } from '../api'
import InlineField from '../components/InlineField'

export default function ProfileDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [profile, setProfile] = useState(null)
  const [materials, setMaterials] = useState([])
  const [allMaterials, setAllMaterials] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [newMaterialID, setNewMaterialID] = useState('')

  useEffect(() => {
    Promise.all([
      get(`/profiles/${id}`),
      get(`/profiles/${id}/materials`),
      get('/materials?all=true'),
    ])
      .then(([p, mats, allMats]) => {
        setProfile(p)
        setMaterials(mats ?? [])
        setAllMaterials(allMats ?? [])
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  async function saveName(value) {
    await put(`/profiles/${id}`, { ...profile, Name: value })
    setProfile(p => ({ ...p, Name: value }))
  }

  async function deleteProfile() {
    if (!confirm('Delete this profile? All material links will be removed.')) { return }
    await del(`/profiles/${id}`)
    navigate('/profiles')
  }

  async function addMaterial() {
    if (!newMaterialID) { return }
    const mp = await post(`/profiles/${id}/materials`, { MaterialID: Number(newMaterialID) })
    setMaterials(m => [...m, mp])
    setNewMaterialID('')
  }

  async function savePrice(mpID, value) {
    const current = materials.find(m => m.ID === mpID)
    const updated = { ...current, Price: value }
    await put(`/material-profiles/${mpID}`, updated)
    setMaterials(m => m.map(x => x.ID === mpID ? updated : x))
  }

  async function removeMaterial(mpID) {
    if (!confirm('Remove this material link?')) { return }
    await del(`/material-profiles/${mpID}`)
    setMaterials(m => m.filter(x => x.ID !== mpID))
  }

  if (loading) return <p className="status">Loading...</p>
  if (error) return <p className="status error">{error}</p>

  return (
    <div className="page">
      <header>
        <button className="btn-back" onClick={() => navigate('/profiles')}>← Back</button>
        <h1>{profile.Name}</h1>
        <div className="header-actions">
          <button className="btn-danger" onClick={deleteProfile}>Delete</button>
        </div>
      </header>

      <section className="detail-section">
        <h2>Profile</h2>
        <div className="detail-grid">
          <div className="field-row">
            <label>Name</label>
            <InlineField value={profile.Name} onSave={saveName} />
          </div>
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
              <th>Price</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {materials.map(m => (
              <tr key={m.ID}>
                <td className="cell-link" onClick={() => navigate(`/materials/${m.MaterialID}`)}>{m.MaterialName}</td>
                <td><InlineField value={m.Price} type="number" onSave={v => savePrice(m.ID, v)} /></td>
                <td><button className="btn-danger-sm" onClick={() => removeMaterial(m.ID)}>×</button></td>
              </tr>
            ))}
            {materials.length === 0 && (
              <tr><td colSpan={3}>No materials linked</td></tr>
            )}
          </tbody>
        </table>
      </section>
    </div>
  )
}
