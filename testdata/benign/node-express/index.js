const express = require('express');
const app = express();
const port = 3000;

app.get('/', (req, res) => {
  res.json({ 
    message: 'Hello from Cocoon sandbox!',
    success: true 
  });
});

app.get('/api/test', (req, res) => {
  res.json({ 
    data: 'This is a test endpoint',
    timestamp: new Date().toISOString()
  });
});

app.listen(port, () => {
  console.log(`Express server running on port ${port}`);
  console.log(`Open http://localhost:${port} to test`);
});