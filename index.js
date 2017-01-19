'use strict'
const app = require('express')()
const bodyParser = require('body-parser')
const jwt = require('jsonwebtoken')
const GITHUB_PRIVATEKEY = process.env.GITHUB_PRIVATEKEY
const GitHubApi = require('github')
const PORT = process.env.PORT || 8080 // eslint-disable-line no-magic-numbers, no-process-env
const INTEGRATION_ID = process.env.INTEGRATION_ID

if (!GITHUB_PRIVATEKEY || !INTEGRATION_ID) {
  throw new Error('missing needed environment variables INTEGRATION_ID or GITHUB_PRIVATEKEY')
}

app.use(bodyParser.json())

app.post('/github', (req, res) => {
  if (req.headers['x-github-event'] !== 'push') {
    return res.send('not a push event so ignoring')
  }

  const token = jwt.sign({iss: INTEGRATION_ID}, GITHUB_PRIVATEKEY, {algorithm: 'RS256', expiresIn: '2 minutes'})

  const github = new GitHubApi({
    Promise: require('bluebird'),
    timeout: 5000
  })

  github.authenticate({
    type: 'integration',
    token: token
  })

  const installationId = req.body.installation.id
  const repo = req.body.repository.name
  const owner = req.body.repository.owner.name

  return github.integrations.createInstallationToken({installation_id: installationId})
    .then(token =>
      github.authenticate({
        type: 'token',
        token: token.token
      })
    )
    .then(() => req.body.commits)
    .map(commit =>
      github.repos.getCommit({
        owner: owner,
        repo: repo,
        sha: commit.id,
        headers: {
          'Accept': 'application/vnd.github.diff'
        }
      })
    )
    .map(commit => commit.data)
    // @todo: pass all the commits through something
    .then(diffs => console.log(diffs))
    .then(() => res.send('done'))
})

module.exports = app.listen(PORT)
