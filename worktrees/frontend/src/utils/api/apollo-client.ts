import { ApolloClient, InMemoryCache, from, HttpLink } from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { setContext } from '@apollo/client/link/context';
import { RetryLink } from '@apollo/client/link/retry';
import { logger } from '../logger';
import { handleError } from '../errors/error-handler';
import { createNetworkError } from '../errors/error-factory';

/**
 * Apollo Client 설정
 * 철칙: 타입 안전성, 에러 처리, 관찰 가능성
 */

/**
 * HTTP 링크 설정
 */
const httpLink = new HttpLink({
  uri: import.meta.env.VITE_GRAPHQL_ENDPOINT || 'http://localhost:8080/graphql',
  credentials: 'include',
});

/**
 * 인증 링크
 */
const authLink = setContext((_, { headers }) => {
  // 토큰이 필요한 경우 여기서 설정
  const token = localStorage.getItem('auth_token');
  
  return {
    headers: {
      ...headers,
      ...(token && { authorization: `Bearer ${token}` }),
      'x-client-version': import.meta.env.VITE_APP_VERSION || '1.0.0',
      'x-trace-id': `frontend-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    },
  };
});

/**
 * 에러 링크
 */
const errorLink = onError(({ graphQLErrors, networkError, operation, forward }) => {
  // GraphQL 에러 처리
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path, extensions }) => {
      const errorInfo = {
        message,
        locations,
        path,
        extensions,
        operation: operation.operationName,
      };
      
      logger.error('GraphQL error', errorInfo);
      
      // 사용자에게 표시할 에러 생성
      handleError(
        createNetworkError(message, {
          code: extensions?.code as string || 'GRAPHQL_ERROR',
        }),
        { showToast: true }
      );
    });
  }

  // 네트워크 에러 처리
  if (networkError) {
    logger.error('Network error', {
      message: networkError.message,
      operation: operation.operationName,
      networkError,
    });
    
    const status = 'statusCode' in networkError ? networkError.statusCode : undefined;
    handleError(
      createNetworkError(networkError.message, {
        status: status as number,
        url: operation.getContext().uri,
      }),
      { showToast: true }
    );
  }
});

/**
 * 재시도 링크
 */
const retryLink = new RetryLink({
  delay: {
    initial: 300,
    max: Infinity,
    jitter: true,
  },
  attempts: {
    max: 3,
    retryIf: (error, _operation) => {
      // 네트워크 에러이거나 5xx 서버 에러인 경우에만 재시도
      if (error?.networkError) {
        const status = 'statusCode' in error.networkError 
          ? error.networkError.statusCode 
          : 0;
        return status >= 500 || status === 0;
      }
      return false;
    },
  },
});

/**
 * 캐시 설정
 */
const cache = new InMemoryCache({
  typePolicies: {
    ObjectType: {
      fields: {
        properties: {
          merge: false, // 속성 배열은 완전 교체
        },
      },
    },
    Query: {
      fields: {
        objectTypes: {
          keyArgs: ['filter', 'sort'],
          merge(existing, incoming, { args }) {
            // 페이지네이션 병합 로직
            if (!existing || args?.pagination?.page === 1) {
              return incoming;
            }
            
            return {
              ...incoming,
              data: [...(existing.data || []), ...(incoming.data || [])],
            };
          },
        },
      },
    },
  },
});

/**
 * Apollo Client 인스턴스
 */
export const apolloClient = new ApolloClient({
  link: from([
    errorLink,
    authLink,
    retryLink,
    httpLink,
  ]),
  cache,
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
      errorPolicy: 'all',
    },
    query: {
      fetchPolicy: 'cache-first',
      errorPolicy: 'all',
    },
    mutate: {
      errorPolicy: 'all',
    },
  },
  connectToDevTools: !import.meta.env.PROD,
});

/**
 * 캐시 초기화
 */
export const clearCache = (): void => {
  apolloClient.clearStore().catch((error) => {
    logger.error('Failed to clear Apollo cache', { error });
  });
};

/**
 * 인증 토큰 설정
 */
export const setAuthToken = (token: string | null): void => {
  if (token) {
    localStorage.setItem('auth_token', token);
  } else {
    localStorage.removeItem('auth_token');
  }
  
  // 캐시 초기화하여 새로운 토큰으로 요청
  clearCache();
};