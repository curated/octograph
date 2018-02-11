import fetch from 'node-fetch'

if (!process.env.GITHUB_TOKEN) {
  throw new Error('Missing environment variable: GITHUB_TOKEN')
}

export interface RequestBody<V> {
  query: string
  variables: V
}

export interface Query<T, V> {
  build(variables: V): RequestBody<V>
  parse(response: any): T
}

export class Graph {
  public static URL = 'https://api.github.com/graphql'

  public static HEADERS = {
    Authorization: `bearer ${process.env.GITHUB_TOKEN}`,
  }

  public query<T, V>(query: Query<T, V>, variables: V): Promise<T> {
    return new Promise<T>((resolve, reject) => {
      const options = {
        body: JSON.stringify(query.build(variables)),
        headers: Graph.HEADERS,
        method: 'POST',
      }

      fetch(Graph.URL, options)
        .then(responseText => responseText.json())
        .then(res => (res.errors ? reject(res) : resolve(query.parse(res))))
        .catch(reject)
    })
  }
}

export const gql = (template: TemplateStringsArray): string => {
  return template.join()
}
