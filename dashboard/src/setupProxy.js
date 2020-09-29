const { createProxyMiddleware } = require('http-proxy-middleware');
module.exports = function (app) {
  app.use(
    '/api',
    createProxyMiddleware('/api', {
      target: 'http://localhost:8090',
      changeOrigin: true,
    })
  );
  app.use(
    '/resources',
    createProxyMiddleware('/resources', {
      target: 'http://localhost:8090',
      ws: true,
    })
  );
};
