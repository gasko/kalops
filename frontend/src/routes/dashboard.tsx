// src/routes/dashboard.tsx
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { supabase } from '../lib/supabaseClient'
import { Button } from "@/components/ui/button"

export const Route = createFileRoute('/dashboard')({
  beforeLoad: async () => {
    const { data } = await supabase.auth.getSession()
    if (!data.session) {
      throw redirect({ to: '/login' })
    }
  },
  component: Dashboard,
})

function Dashboard() {
  const navigate = useNavigate()

  const handleLogout = async () => {
    await supabase.auth.signOut()
    navigate({ to: '/login' })
  }

  return (
    <div className="flex flex-col items-center gap-6 p-6">
      <h1 className="text-2xl font-bold">Välkommen till dashboard!</h1>
      <Button variant="destructive" onClick={handleLogout}>
        Logga ut
      </Button>
    </div>
  )
}
