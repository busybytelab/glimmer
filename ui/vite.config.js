import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';
import { createRequire } from 'module';
import fs from 'fs';
var require = createRequire(import.meta.url);
// Plugin to copy src/assets/glimmer.svg to root
var copyGlimmerSvgPlugin = function () {
    return {
        name: 'copy-glimmer-svg',
        writeBundle: function () {
            var srcPath = path.resolve('./src/assets/glimmer.svg');
            var distPath = path.resolve('./dist/glimmer.svg');
            fs.copyFile(srcPath, distPath, function (err) {
                if (err)
                    console.error('Failed to copy glimmer.svg to dist root:', err);
            });
        }
    };
};
// https://vite.dev/config/
export default defineConfig({
    plugins: [
        svelte(),
        copyGlimmerSvgPlugin()
    ],
    resolve: {
        alias: {
            $lib: path.resolve('./src/lib')
        }
    },
    css: {
        postcss: './postcss.config.cjs'
    },
    base: '/',
    publicDir: 'public',
    server: {
        port: 3000,
        proxy: {
            '/api': {
                target: 'http://localhost:8787',
                changeOrigin: true,
                secure: false
            }
        }
    },
    build: {
        outDir: 'dist',
        rollupOptions: {
            input: {
                main: 'index.html'
            }
        }
    }
});
