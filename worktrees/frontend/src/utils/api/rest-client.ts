import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { logger } from '../logger';
import { handleError, retryWithBackoff } from '../errors/error-handler';
import { createNetworkError, getErrorCodeFromStatus } from '../errors/error-factory';

/**
 * REST API 클라이언트
 * 철칙: 타입 안전성, 에러 처리, 재시도 로직
 */

/**
 * API 응답 타입
 */
export interface ApiResponse<T> {
  readonly success: boolean;
  readonly data: T;
  readonly message?: string;
  readonly errors?: readonly string[];
}

/**
 * 요청 설정 타입
 */
export interface RequestConfig extends AxiosRequestConfig {
  readonly retries?: number;
  readonly skipAuth?: boolean;
  readonly skipErrorHandling?: boolean;
}

/**
 * REST 클라이언트 클래스
 */
class RestClient {
  private instance: AxiosInstance;
  
  constructor() {
    this.instance = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    });
    
    this.setupInterceptors();
  }
  
  /**
   * 인터셉터 설정
   */
  private setupInterceptors(): void {
    // 요청 인터셉터
    this.instance.interceptors.request.use(
      (config) => {
        // 인증 토큰 추가
        const token = localStorage.getItem('auth_token');
        if (token && !config.headers?.skipAuth) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        
        // 추적 ID 추가
        config.headers['x-trace-id'] = `rest-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
        
        // 요청 로깅
        logger.debug('REST request', {
          method: config.method?.toUpperCase(),
          url: config.url,
          headers: config.headers,
        });
        
        return config;
      },
      (error) => {
        logger.error('Request interceptor error', { error });
        return Promise.reject(error);
      }
    );
    
    // 응답 인터셉터
    this.instance.interceptors.response.use(
      (response) => {
        // 응답 로깅
        logger.debug('REST response', {
          status: response.status,
          url: response.config.url,
          data: response.data,
        });
        
        return response;
      },
      (error: AxiosError) => {
        // 에러 로깅
        logger.error('REST response error', {
          status: error.response?.status,
          url: error.config?.url,
          message: error.message,
          data: error.response?.data,
        });
        
        return Promise.reject(error);
      }
    );
  }
  
  /**
   * 에러 변환
   */
  private transformError(error: AxiosError): Error {
    const status = error.response?.status || 0;
    const message = error.response?.data?.message || error.message;
    
    return createNetworkError(message, {
      status,
      statusText: error.response?.statusText,
      url: error.config?.url,
      code: getErrorCodeFromStatus(status),
    });
  }
  
  /**
   * 요청 실행
   */
  private async executeRequest<T>(
    config: RequestConfig
  ): Promise<AxiosResponse<ApiResponse<T>>> {
    const { retries = 3, skipErrorHandling = false, ...axiosConfig } = config;
    
    try {
      if (retries > 0) {
        return await retryWithBackoff(
          () => this.instance.request(axiosConfig),
          {
            maxRetries: retries,
            shouldRetry: (error) => {
              if ('status' in error && typeof error.status === 'number') {
                return error.status >= 500 || error.status === 0;
              }
              return false;
            },
          }
        );
      } else {
        return await this.instance.request(axiosConfig);
      }
    } catch (error) {
      const transformedError = this.transformError(error as AxiosError);
      
      if (!skipErrorHandling) {
        handleError(transformedError);
      }
      
      throw transformedError;
    }
  }
  
  /**
   * GET 요청
   */
  async get<T>(url: string, config?: RequestConfig): Promise<T> {
    const response = await this.executeRequest<T>({
      ...config,
      method: 'GET',
      url,
    });
    
    return response.data.data;
  }
  
  /**
   * POST 요청
   */
  async post<T, D = unknown>(
    url: string,
    data?: D,
    config?: RequestConfig
  ): Promise<T> {
    const response = await this.executeRequest<T>({
      ...config,
      method: 'POST',
      url,
      data,
    });
    
    return response.data.data;
  }
  
  /**
   * PUT 요청
   */
  async put<T, D = unknown>(
    url: string,
    data?: D,
    config?: RequestConfig
  ): Promise<T> {
    const response = await this.executeRequest<T>({
      ...config,
      method: 'PUT',
      url,
      data,
    });
    
    return response.data.data;
  }
  
  /**
   * PATCH 요청
   */
  async patch<T, D = unknown>(
    url: string,
    data?: D,
    config?: RequestConfig
  ): Promise<T> {
    const response = await this.executeRequest<T>({
      ...config,
      method: 'PATCH',
      url,
      data,
    });
    
    return response.data.data;
  }
  
  /**
   * DELETE 요청
   */
  async delete<T>(url: string, config?: RequestConfig): Promise<T> {
    const response = await this.executeRequest<T>({
      ...config,
      method: 'DELETE',
      url,
    });
    
    return response.data.data;
  }
  
  /**
   * 인증 토큰 설정
   */
  setAuthToken(token: string | null): void {
    if (token) {
      this.instance.defaults.headers.common.Authorization = `Bearer ${token}`;
    } else {
      delete this.instance.defaults.headers.common.Authorization;
    }
  }
}

/**
 * 싱글톤 REST 클라이언트 인스턴스
 */
export const restClient = new RestClient();