import { gql } from './graph'

export const query = gql`
  query TopIssues($query: String!, $first: Int!, $after: String) {
    search(type: ISSUE, query: $query, first: $first, after: $after) {
      pageInfo {
        hasNextPage
        endCursor
      }
      edges {
        node {
          ... on Issue {
            id
            url
            number
            title
            bodyText
            state
            createdAt
            updatedAt
            reactionGroups {
              content
              users {
                totalCount
              }
            }
            author {
              url
              login
              avatarUrl
            }
            repository {
              id
              url
              name
              primaryLanguage {
                name
              }
              forks {
                totalCount
              }
              stargazers {
                totalCount
              }
              owner {
                id
                url
                login
                avatarUrl
              }
            }
          }
        }
      }
    }
  }
`
