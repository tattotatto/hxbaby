const express = require('express');
const { compose } = require('./scripts/compose');
const { packageToZip } = require('./scripts/packager');
const { v4: uuidv4 } = require('uuid');
const fs = require('fs');
const path = require('path');

const app = express();
app.use(express.json({ limit: '5mb' }));

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'codegen-service' });
});

app.get('/api/modules', (req, res) => {
  const dir = path.join(__dirname, 'templates', 'modules');
  const modules = fs.readdirSync(dir).filter(d => {
    const p = path.join(dir, d);
    return fs.statSync(p).isDirectory() && fs.existsSync(path.join(p, 'module.json'));
  }).map(d => JSON.parse(fs.readFileSync(path.join(dir, d, 'module.json'), 'utf-8')));
  res.json({ modules });
});

app.post('/api/build', async (req, res) => {
  try {
    const { project } = req.body;
    if (!project?.modules) {
      return res.status(400).json({ error: '缺少 modules 参数' });
    }
    project.build_task_id = uuidv4();

    const { outputDir, warnings } = await compose(project);
    const { zipPath, md5, size } = await packageToZip(outputDir);

    res.json({ task_id: project.build_task_id, status: 'done', zip_path: zipPath, md5, size_bytes: size, warnings });
  } catch (err) {
    res.status(500).json({ status: 'failed', error: err.message });
  }
});

const PORT = process.env.PORT || 3002;
app.listen(PORT, () => console.log(`CodeGen service on port ${PORT}`));
