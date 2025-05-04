# Glimmer UI

The frontend interface for Glimmer, a self-hosted learning assistant.

## Tech Stack

- Svelte
- TypeScript
- Vite
- PocketBase Client

## Development

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Build for production:
   ```bash
   npm run build
   ```

4. Preview production build:
   ```bash
   npm run preview
   ```

## Project Structure

- `src/` - Source files
  - `lib/` - Utility functions and shared code
  - `routes/` - Page components
  - `components/` - Reusable UI components
  - `types/` - TypeScript type definitions

## Features

- User authentication
- Role-based views (Instructor/Learner)
- Practice topic management
- Practice session tracking
- Real-time updates

## Configuration

The application expects a PocketBase instance running at `http://localhost:8787`. You can modify this in `src/lib/pocketbase.ts`.

## Deployment

Build the application and serve the `dist` directory using your preferred static file server.
