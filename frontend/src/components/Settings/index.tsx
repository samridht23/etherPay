import { useState } from "react"
import clsx from "clsx"

const Settings = () => {

  const [formData, setFormData] = useState({
    username: "",
    about: ""
  })


  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/profile`,
        {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData)
        });

      if (response.ok) {
        const data = await response.json();
        console.log('Success:', data);
      } else {
        const errorData = await response.json();
        console.error('Error:', errorData);
      }
    } catch (err) {
      console.log(err)
    }
  };


  return (
    <div className="w-full flex gap-4">
      <div className="w-3/5 m-auto bg-[var(--bg-muted)] rounded-[var(--radius-md)] p-4">
        <form className="py-4 flex flex-col gap-5" onSubmit={handleSubmit} >
          <div className="flex flex-col gap-1.5">
            <label className="font-semibold text-xs text-[var(--text-subtle)]">Username</label>
            <input
              id="username"
              type="text"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              placeholder="Username"
              className={clsx(
                'outline-none p-2 w-full bg-[var(--bg-subtle)] placeholder:font-[Geist] placeholder:text-sm',
                'text-base border border-[var(--border-subtle)] rounded-[var(--radius-sm)]',
              )}
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="font-semibold text-xs text-[var(--text-subtle)]">About</label>
            <textarea
              id="about"
              value={formData.about}
              onChange={(e) => setFormData({ ...formData, about: e.target.value })}
              placeholder="Tell us something about yourself..."
              className={clsx(
                'outline-none p-2 w-full bg-[var(--bg-subtle)] placeholder:font-[Geist] placeholder:text-sm',
                'text-base border border-[var(--border-subtle)] rounded-[var(--radius-sm)]',
              )}
            />
          </div>
          <div className="">
            <button
              type="submit"
              className="p-3 font-semibold text-sm w-full border rounded-md bg-[var(--bg-inverted)] text-[var(--text-inverted)]">
              Update Profile
            </button>
          </div>
        </form>
      </div>
    </div >
  )
}
export default Settings;


