import { BrowserRouter } from "react-router-dom";
import Router from "./pages/router";
import { create } from "zustand";
import AuthProvider from './AuthWrapper';
import { Toaster } from 'sonner';

export interface ViewingProfileData {
  address: string;
  username: string;
  profile_image_url: string;
  banner_image_url: string;
  about: string;
  cteated_at: string;
}

interface ContextProps {
  address: string | null
  viewingProfileData: ViewingProfileData | null
  updateAddress: (newValue: string | null) => void
  updateViewingProfileData: (newValue: ViewingProfileData | null) => void
}

export const zustStore = create<ContextProps>()((set) => (
  {
    address: null,
    viewingProfileData: null,
    updateAddress: (newValue) => set(() => ({ address: newValue })),
    updateViewingProfileData: (newValue) => set(() => ({ viewingProfileData: newValue }))
  }
))

function App() {

  return (
    <BrowserRouter>
      <AuthProvider>
        <Toaster
          toastOptions={{
            className:
              "p-4 flex gap-3 bg-[var(--bg-muted)] text-[var(--text)] border-[1px] rounded border-[var(--border-muted)]",
            duration: 1500,
          }}
        />
        <Router />
      </AuthProvider >
    </BrowserRouter>
  )
}

export default App
