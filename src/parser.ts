import * as R from 'ramda'
import {
  ActorSchema,
  IssueSchema,
  RepositorySchema,
  TopIssues,
  RateLimit,
} from './schema'

export const parse = (res: any): TopIssues => {
  const rateLimit: any = R.path(['data', 'rateLimit'], res) || {}
  const search: any = R.path(['data', 'search'], res) || {}
  const pageInfo = search.pageInfo || {}

  return {
    rateLimit: parseRateLimit(rateLimit),
    endCursor: pageInfo.endCursor,
    issues: parseIssues(search.edges),
  }
}

const parseRateLimit = (rateLimit: any): RateLimit => {
  return {
    remaining: rateLimit.remaining,
    resetAt: rateLimit.resetAt && new Date(rateLimit.resetAt),
  }
}

const parseIssues = (edges: any[]): IssueSchema[] => {
  const issues = R.map(parseIssue, edges)
  return R.reject((issue: IssueSchema) => R.isNil(issue.githubId), issues)
}

const parseIssue = (edge: any): IssueSchema => {
  const node: any = R.path(['node'], edge)
  const groups = node.reactionGroups || []

  return {
    githubId: node.id,
    url: node.url,
    number: node.number,
    title: node.title,
    bodyText: node.bodyText,
    state: node.state,
    createdAt: node.createdAt && new Date(node.createdAt),
    updatedAt: node.updatedAt && new Date(node.updatedAt),
    heart: getGroupCount('HEART', groups),
    hooray: getGroupCount('HOORAY', groups),
    thumbsUp: getGroupCount('THUMBS_UP', groups),
    laugh: getGroupCount('LAUGH', groups),
    confused: getGroupCount('CONFUSED', groups),
    thumbsDown: getGroupCount('THUMBS_DOWN', groups),
    author: parseActor(node.author),
    repository: parseRepository(node.repository),
  }
}

const getGroupCount = (key: string, groups: any[]): number => {
  const group = R.find(g => Object.is(key, g.content), groups)
  return R.path(['users', 'totalCount'], group) || 0
}

const parseActor = (actor: any): ActorSchema => {
  if (!actor) {
    return undefined
  }

  return {
    githubId: actor.id,
    url: actor.url,
    login: actor.login,
    avatarUrl: actor.avatarUrl,
  }
}

const parseRepository = (repository: any): RepositorySchema => {
  if (!repository) {
    return undefined
  }

  return {
    githubId: repository.id,
    url: repository.url,
    name: repository.name,
    primaryLanguage: R.path(['primaryLanguage', 'name'], repository),
    forks: R.path(['forks', 'totalCount'], repository) || 0,
    stargazers: R.path(['stargazers', 'totalCount'], repository) || 0,
    owner: parseActor(repository.owner),
  }
}
