import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

function initializeApp() {
  const appElement = document.getElementById('app')
  
  if (!appElement) {
    console.error('Failed to find app element')
    return
  }

  try {
    const app = mount(App, {
      target: appElement,
    })

    return app
  } catch (error) {
    console.error('Failed to mount app:', error)
    appElement.innerHTML = `
      <div class="min-h-screen flex items-center justify-center bg-gray-50">
        <div class="text-center">
          <h1 class="text-2xl font-bold text-red-600">Error</h1>
          <p class="mt-2 text-gray-600">Failed to initialize the application. Please try refreshing the page.</p>
        </div>
      </div>
    `
  }
}

// Initialize the app
const app = initializeApp()

export default app
