import fetch from 'node-fetch'

if (!process.env.GITHUB_TOKEN) {
  throw new Error('Missing environment variable: GITHUB_TOKEN')
}

const url = 'https://api.github.com/graphql'

const headers = {
  Authorization: `bearer ${process.env.GITHUB_TOKEN}`,
}

export interface Body<V> {
  query: string
  variables: V
}

const query = <T, V>(body: Body<V>, parse: (res: any) => T): Promise<T> => {
  return new Promise<T>((resolve, reject) => {
    const options = {
      body: JSON.stringify(body),
      headers,
      method: 'POST',
    }

    fetch(url, options)
      .then(text => text.json())
      .then(res => (res.errors ? reject(res) : resolve(parse(res))))
      .catch(reject)
  })
}

export const graph = { query }

export const gql = (template: TemplateStringsArray): string => {
  return template.join()
}
