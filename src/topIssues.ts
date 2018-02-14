import { graph } from './graph'
import { parse } from './parser'
import { TopIssues, Variables } from './schema'
import { query } from './topIssues.graphql'

export const getTopIssues = (variables: Variables): Promise<TopIssues> => {
  return graph.query({ query, variables }, parse)
}
