import { gql, graph } from './graph'

export interface Variables {
  login: string
}

export interface Organization {
  name: string
}

const query = gql`
  query Organization($login: String!) {
    organization(login: $login) {
      name
    }
  }
`

export const getOrg = (variables: Variables): Promise<Organization> => {
  return graph.query({ query, variables }, parse)
}

const parse = (res: any): Organization => {
  return { name: res.data.organization.name }
}
