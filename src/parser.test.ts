import test from 'ava'
import { parse } from './parser'
import { edges } from './parser.test.fixture'

const topIssues = parse(edges)

test('rate limit', async t => {
  t.is(topIssues.rateLimit.remaining, 4999)
  t.deepEqual(topIssues.rateLimit.resetAt, new Date('2018-03-06T20:05:51Z'))
})

test('page info', async t => {
  t.is(topIssues.issueCount, 3)
  t.is(topIssues.endCursor, 'Y3Vyc29yOjM=')
})

test('reject issues without an id', async t => {
  t.is(edges.data.search.edges.length, 3)
  t.is(topIssues.issues.length, 2)
  t.truthy(topIssues.issues[0].githubId)
  t.truthy(topIssues.issues[1].githubId)
})

test('parse issue with only an id', async t => {
  t.deepEqual(topIssues.issues[0], {
    githubId: 'BFWHIBFdeWIhFBIIFknIJWF=',
    url: undefined,
    number: undefined,
    title: undefined,
    bodyText: undefined,
    state: undefined,
    createdAt: undefined,
    updatedAt: undefined,
    heart: 0,
    hooray: 0,
    thumbsUp: 0,
    laugh: 0,
    confused: 0,
    thumbsDown: 0,
    author: undefined,
    repository: undefined,
  })
})

test('parse issue with all properties', async t => {
  t.deepEqual(topIssues.issues[1], {
    githubId: 'MDU6SXNzdWUyNDU4Mzg1MDY=',
    url: 'https://github.com/curated/octograph/issues/1',
    number: 1,
    title: 'foo',
    bodyText: 'bar',
    state: 'CLOSED',
    createdAt: new Date('2017-07-26T20:05:32.000Z'),
    updatedAt: new Date('2017-11-03T10:31:08.000Z'),
    heart: 62,
    hooray: 53,
    thumbsUp: 44,
    laugh: 35,
    confused: 26,
    thumbsDown: 17,
    author: {
      githubId: 'MDQ6VXNlcjI0MDMxMQ==',
      url: 'https://github.com/inspectocat',
      login: 'inspectocat',
      avatarUrl: 'https://avatars0.githubusercontent.com/u/36428176?v=4',
    },
    repository: {
      githubId: 'MDEwOlJlcG9zaXRvcnkxMDI3MDI1MA==',
      url: 'https://github.com/curated/octograph',
      name: 'octograph',
      primaryLanguage: 'TypeScript',
      forks: 123,
      stargazers: 456,
      owner: {
        githubId: 'MDEyOk9yZ2FuaXphdGlvbjY5NjMx',
        url: 'https://github.com/curated',
        login: 'curated',
        avatarUrl: 'https://avatars2.githubusercontent.com/u/36278390?v=4',
      },
    },
  })
})
