export const edges = {
  data: {
    search: {
      pageInfo: {
        hasNextPage: true,
        endCursor: 'Y3Vyc29yOjM=',
      },
      edges: [
        {
          node: {},
        },
        {
          node: {
            url: 'https://github.com/curated/octograph/issues/2',
          },
        },
        {
          node: {
            url: 'https://github.com/curated/octograph/issues/1',
            number: 1,
            title: 'foo',
            bodyText: 'bar',
            state: 'CLOSED',
            createdAt: '2017-07-26T20:05:32Z',
            updatedAt: '2017-11-03T10:31:08Z',
            author: {
              url: 'https://github.com/inspectocat',
              login: 'inspectocat',
              avatarUrl:
                'https://avatars0.githubusercontent.com/u/36428176?v=4',
            },
            repository: {
              url: 'https://github.com/curated/octograph',
              name: 'octograph',
              owner: {
                url: 'https://github.com/curated',
                login: 'curated',
                avatarUrl:
                  'https://avatars2.githubusercontent.com/u/36278390?v=4',
              },
              primaryLanguage: {
                name: 'TypeScript',
              },
              forks: {
                totalCount: 123,
              },
              stargazers: {
                totalCount: 456,
              },
            },
            reactionGroups: [
              {
                content: 'THUMBS_DOWN',
                users: {
                  totalCount: 17,
                },
              },
              {
                content: 'CONFUSED',
                users: {
                  totalCount: 26,
                },
              },
              {
                content: 'LAUGH',
                users: {
                  totalCount: 35,
                },
              },
              {
                content: 'THUMBS_UP',
                users: {
                  totalCount: 44,
                },
              },
              {
                content: 'HOORAY',
                users: {
                  totalCount: 53,
                },
              },
              {
                content: 'HEART',
                users: {
                  totalCount: 62,
                },
              },
            ],
          },
        },
      ],
    },
  },
}
