import test from 'ava'
import { getTopIssues } from './topIssues'

test('search top issues by minimum number of reactions', async t => {
  const res = await getTopIssues({ query: 'reactions:>1000', first: 10 })
  t.true(res.hasNextPage)
  t.true(res.endCursor.length > 0)
  t.true(res.issues.length <= 10)
})

test('fail to search top issues with invalid variables', async t => {
  try {
    await getTopIssues({ query: undefined, first: undefined })
  } catch (err) {
    t.is(err.errors.length, 2)
    t.is(
      err.errors[0].message,
      'Variable query of type String! was provided invalid value',
    )
    t.is(
      err.errors[1].message,
      'Variable first of type Int! was provided invalid value',
    )
  }
})
