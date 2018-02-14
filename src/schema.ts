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
  url: string
  number: number
  title: string
  bodyText: string
  state: string
  createdAt: Date
  updatedAt: Date
  author: Actor
  repository: Repository
  reactionGroups: ReactionGroups
}

export interface Actor {
  url: string
  login: string
  avatarUrl: string
}

export interface Repository {
  url: string
  name: string
  owner: Actor
  primaryLanguage: string
  forks: number
  stargazers: number
}

export interface ReactionGroups {
  heart: number
  hooray: number
  thumbsUp: number
  laugh: number
  confused: number
  thumbsDown: number
}
