import { gql } from './graph'

export const query = gql`
  query TopIssues($query: String!, $first: Int!, $after: String) {
    rateLimit {
      ...rateLimit
    }
    search(type: ISSUE, query: $query, first: $first, after: $after) {
      pageInfo {
        ...pageInfo
      }
      edges {
        node {
          ... on Issue {
            ...issue
          }
        }
      }
    }
  }

  fragment rateLimit on RateLimit {
    remaining
    resetAt
  }

  fragment pageInfo on PageInfo {
    endCursor
  }

  fragment issue on Issue {
    id
    url
    number
    title
    bodyText
    state
    createdAt
    updatedAt
    reactionGroups {
      ...reactionGroup
    }
    author {
      ...actor
    }
    # repository {
    #  ...repository
    # }
  }

  fragment reactionGroup on ReactionGroup {
    content
    users {
      totalCount
    }
  }

  fragment actor on Actor {
    ... on User {
      ...user
    }
    ... on Organization {
      ...organization
    }
    ... on Bot {
      ...bot
    }
  }

  # fragment repository on Repository {
  #   id
  #   url
  #   name
  #   primaryLanguage {
  #     name
  #   }
  #   forks {
  #     totalCount
  #   }
  #   stargazers {
  #     totalCount
  #   }
  #   owner {
  #     ...user
  #   }
  # }

  fragment user on User {
    id
    ...actorProps
  }

  fragment organization on Organization {
    id
    ...actorProps
  }

  fragment bot on Bot {
    id
    ...actorProps
  }

  fragment actorProps on Actor {
    url
    login
    avatarUrl
  }
`
