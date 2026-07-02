const { spawn } = require('child_process')
const fs = require('fs')
const path = require('path')

const rootDir = __dirname
const webDir = path.join(rootDir, 'apps')
const tmpDir = path.join(rootDir, '.tmp')

fs.mkdirSync(path.join(tmpDir, 'localappdata'), { recursive: true })
fs.mkdirSync(path.join(tmpDir, 'appdata'), { recursive: true })
fs.mkdirSync(path.join(tmpDir, 'gocache'), { recursive: true })
fs.mkdirSync(path.join(tmpDir, 'gomodcache'), { recursive: true })

function buildChildEnv(extraEnv = {}) {
  return {
    ...process.env,
    HOME: rootDir,
    USERPROFILE: rootDir,
    HOMEDRIVE: 'D:',
    HOMEPATH: '\\final-production\\devserver',
    LOCALAPPDATA: path.join(tmpDir, 'localappdata'),
    APPDATA: path.join(tmpDir, 'appdata'),
    TEMP: tmpDir,
    TMP: tmpDir,
    GOCACHE: path.join(tmpDir, 'gocache'),
    GOMODCACHE: path.join(tmpDir, 'gomodcache'),
    ...extraEnv,
  }
}

function startProcess(command, args, cwd, label, env = {}) {
  const child = spawn(command, args, {
    cwd,
    env: buildChildEnv(env),
    shell: false,
    stdio: 'inherit',
  })

  child.on('exit', (code, signal) => {
    if (signal) {
      console.log(`[${label}] exited with signal ${signal}`)
      return
    }

    console.log(`[${label}] exited with code ${code ?? 0}`)
  })

  return child
}

function ensureWebDependencies() {
  const nodeModulesPath = path.join(webDir, 'node_modules')

  if (fs.existsSync(nodeModulesPath)) {
    return Promise.resolve()
  }

  return new Promise((resolve, reject) => {
    console.log('[web] node_modules not found, installing dependencies...')

    const installer = spawn('npm.cmd', ['install'], {
      cwd: webDir,
      env: buildChildEnv(),
      shell: false,
      stdio: 'inherit',
    })

    installer.on('exit', (code) => {
      if (code === 0) {
        resolve()
        return
      }

      reject(new Error(`[web] npm install failed with code ${code ?? 1}`))
    })
  })
}

async function main() {
  await ensureWebDependencies()

  const backend = startProcess('go', ['run', './cmd/devserver', 'serve'], rootDir, 'backend')
  const frontend = startProcess(
    'node',
    [path.join(webDir, 'node_modules', 'next', 'dist', 'bin', 'next'), 'dev'],
    webDir,
    'web',
  )
  let shuttingDown = false

  const shutdown = () => {
    if (shuttingDown) {
      return
    }

    shuttingDown = true

    if (!backend.killed) {
      backend.kill()
    }

    if (!frontend.killed) {
      frontend.kill()
    }
  }

  process.on('SIGINT', shutdown)
  process.on('SIGTERM', shutdown)

  backend.on('close', (code) => {
    if ((code ?? 1) !== 0) {
      shutdown()
      process.exitCode = code ?? 1
      return
    }

    console.log('[backend] completed successfully')
  })

  frontend.on('close', (code) => {
    shutdown()
    process.exitCode = code ?? 1
  })
}

main().catch((error) => {
  console.error(error.message)
  process.exitCode = 1
})
