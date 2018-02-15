export interface Variables {
  query: string
  first: number
  after?: string
}

export interface TopIssues {
  hasNextPage: boolean
  endCursor: string
  issues: Issue[]
}

export interface Issue {
  githubId: string
  url: string
  number: number
  title: string
  bodyText: string
  state: string
  createdAt: Date
  updatedAt: Date
  heart: number
  hooray: number
  thumbsUp: number
  laugh: number
  confused: number
  thumbsDown: number
  author: Actor
  repository: Repository
}

export interface Actor {
  githubId: string
  url: string
  login: string
  avatarUrl: string
}

export interface Repository {
  githubId: string
  url: string
  name: string
  primaryLanguage: string
  forks: number
  stargazers: number
  owner: Actor
}
