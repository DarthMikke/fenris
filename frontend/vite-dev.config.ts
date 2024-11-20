import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

/**
 * @todo For the production build, output the files into Django's
 * static directory.
 */

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../public',
    emptyOutDir: true,
    minify: false,
    sourcemap: 'inline',
  },
})
