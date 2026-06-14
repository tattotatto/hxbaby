const fs = require('fs');
const path = require('path');
const archiver = require('archiver');
const crypto = require('crypto');

async function packageToZip(outputDir) {
  return new Promise((resolve, reject) => {
    const zipPath = outputDir + '.zip';
    const output = fs.createWriteStream(zipPath);
    const archive = archiver('zip', { zlib: { level: 9 } });

    output.on('close', () => {
      const buf = fs.readFileSync(zipPath);
      const md5 = crypto.createHash('md5').update(buf).digest('hex');
      const size = fs.statSync(zipPath).size;
      resolve({ zipPath, md5, size });
    });

    archive.on('error', reject);
    archive.pipe(output);
    archive.directory(outputDir, path.basename(outputDir));
    archive.finalize();
  });
}

module.exports = { packageToZip };
