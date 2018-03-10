export interface RateLimit {
  remaining: number
  resetAt: Date
}

export interface Variables {
  query: string
  first: number
  after?: string
}

export interface TopIssues {
  rateLimit: RateLimit
  issueCount: number
  endCursor: string
  issues: IssueSchema[]
}

export interface IssueSchema {
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
  author: ActorSchema
  repository: RepositorySchema
}

export interface ActorSchema {
  githubId: string
  url: string
  login: string
  avatarUrl: string
}

export interface RepositorySchema {
  githubId: string
  url: string
  name: string
  primaryLanguage: string
  forks: number
  stargazers: number
  owner: ActorSchema
}
