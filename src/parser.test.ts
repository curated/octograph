import test from 'ava'
import { parse } from './parser'
import { edges } from './parser.test.fixture'

const topIssues = parse(edges)

test('keep relevant page info', async t => {
  t.true(topIssues.hasNextPage)
  t.is(topIssues.endCursor, 'Y3Vyc29yOjM=')
})

test('reject issues without a url', async t => {
  t.is(edges.data.search.edges.length, 3)
  t.is(topIssues.issues.length, 2)
  t.truthy(topIssues.issues[0].url)
  t.truthy(topIssues.issues[1].url)
})

test('parse issue with only a url', async t => {
  t.deepEqual(topIssues.issues[0], {
    url: 'https://github.com/curated/octograph/issues/2',
    number: undefined,
    title: undefined,
    bodyText: undefined,
    state: undefined,
    createdAt: undefined,
    updatedAt: undefined,
    author: undefined,
    repository: undefined,
    reactionGroups: undefined,
  })
})

test('parse issue with all properties', async t => {
  t.deepEqual(topIssues.issues[1], {
    url: 'https://github.com/curated/octograph/issues/1',
    number: 1,
    title: 'foo',
    bodyText: 'bar',
    state: 'CLOSED',
    createdAt: new Date('2017-07-26T20:05:32.000Z'),
    updatedAt: new Date('2017-11-03T10:31:08.000Z'),
    author: {
      url: 'https://github.com/inspectocat',
      login: 'inspectocat',
      avatarUrl: 'https://avatars0.githubusercontent.com/u/36428176?v=4',
    },
    repository: {
      url: 'https://github.com/curated/octograph',
      name: 'octograph',
      owner: {
        url: 'https://github.com/curated',
        login: 'curated',
        avatarUrl: 'https://avatars2.githubusercontent.com/u/36278390?v=4',
      },
      primaryLanguage: 'TypeScript',
      forks: 123,
      stargazers: 456,
    },
    reactionGroups: {
      heart: 62,
      hooray: 53,
      thumbsUp: 44,
      laugh: 35,
      confused: 26,
      thumbsDown: 17,
    },
  })
})
