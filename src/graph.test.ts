import test from 'ava'
import { gql } from './graph'
import { getOrg } from './graph.test.fixture'

test('lookup organization by login', async t => {
  const res = await getOrg({ login: 'curated' })
  t.deepEqual(res, { name: 'Curated' })
})

test('fail to lookup organization by login', async t => {
  try {
    await getOrg({ login: '#' })
  } catch (res) {
    t.deepEqual(res.errors[0].type, 'NOT_FOUND')
  }
})

test('provide gql for inline graphql templates', t => {
  t.is(gql`query {\n}`, 'query {\n}')
})
