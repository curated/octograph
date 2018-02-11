import test from 'ava'
import { gql, Graph, Query, RequestBody } from './Graph'

interface Variables {
  login: string
}

interface Organization {
  name: string
}

class OrganizationQuery implements Query<Organization, Variables> {
  public build(variables: Variables): RequestBody<Variables> {
    return {
      query: gql`
        query Organization($login: String!) {
          organization(login: $login) {
            name
          }
        }
      `,
      variables,
    }
  }

  public parse(response: any): Organization {
    return { name: response.data.organization.name }
  }
}

const graph = new Graph()
const query = new OrganizationQuery()

test('lookup organization by login', async t => {
  const response = await graph.query(query, { login: 'curated' })
  t.deepEqual(response, { name: 'Curated' })
})

test('fail to lookup organization by login', async t => {
  try {
    await graph.query(query, { login: '#' })
  } catch (response) {
    t.deepEqual(response.errors[0].type, 'NOT_FOUND')
  }
})

test('use github graphql endpoint', t => {
  t.is(Graph.URL, 'https://api.github.com/graphql')
})

test('use github token from environment variable', t => {
  const token = process.env.GITHUB_TOKEN
  t.is(Graph.HEADERS.Authorization, `bearer ${token}`)
  t.truthy(token)
})

test('provide gql for inline graphql templates', t => {
  t.is('query {\n}', gql`query {\n}`)
})
