#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();
  
  let name = 'certctl';
  
  if (platform === 'win32') {
    name += '-windows';
  } else if (platform === 'darwin') {
    name += '-darwin';
  } else {
    name += '-linux';
  }
  
  if (arch === 'arm64') {
    name += '-arm64';
  } else {
    name += '-amd64';
  }
  
  if (platform === 'win32') {
    name += '.exe';
  }
  
  return name;
}

const binaryPath = path.join(__dirname, getBinaryName());

if (!fs.existsSync(binaryPath)) {
  console.error(`Binary not found: ${binaryPath}`);
  console.error('Please run: npm run postinstall');
  process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  env: process.env
});

child.on('exit', (code) => {
  process.exit(code || 0);
});
