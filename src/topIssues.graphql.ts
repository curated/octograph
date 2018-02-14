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
            url
            number
            title
            bodyText
            state
            createdAt
            updatedAt
            author {
              url
              login
              avatarUrl
            }
            repository {
              url
              name
              owner {
                url
                login
                avatarUrl
              }
              primaryLanguage {
                name
              }
              forks {
                totalCount
              }
              stargazers {
                totalCount
              }
            }
            reactionGroups {
              content
              users {
                totalCount
              }
            }
          }
        }
      }
    }
  }
`
