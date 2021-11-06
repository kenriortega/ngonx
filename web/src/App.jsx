import { useState } from 'react'
import Header from './components/Header'
import Footer from './components/Footer'
import Endpoints from './components/Endpoints'
import { QueryClient, QueryClientProvider, useQuery } from 'react-query'
const queryClient = new QueryClient()
function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Header />
      <main className={"main"}>
        <Endpoints />
      </main>
      <Footer />
    </QueryClientProvider>
  )
}

export default App
