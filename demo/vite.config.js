export default {
    proxy: {
        '/deploy': {
            target: 'http://127.0.0.1:16219/',
            changeOrigin: true,
            ws: true,
        }
    }
}
