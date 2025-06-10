/**
 * 로거 유틸리티
 * 철칙: 타입 안전성, 설정 가능, 성능 고려
 */

/**
 * 로그 레벨
 */
export const LogLevels = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
} as const;

export type LogLevel = keyof typeof LogLevels;

/**
 * 로그 메타데이터
 */
export interface LogMetadata {
  readonly [key: string]: unknown;
}

/**
 * 로그 항목
 */
export interface LogEntry {
  readonly level: LogLevel;
  readonly message: string;
  readonly timestamp: string;
  readonly metadata?: LogMetadata;
  readonly traceId?: string;
}

/**
 * 로거 설정
 */
export interface LoggerConfig {
  readonly level: LogLevel;
  readonly enableConsole: boolean;
  readonly enableRemote: boolean;
  readonly maxLocalEntries: number;
  readonly remoteEndpoint?: string;
}

/**
 * 기본 로거 설정
 */
const defaultConfig: LoggerConfig = {
  level: import.meta.env.PROD ? 'INFO' : 'DEBUG',
  enableConsole: true,
  enableRemote: import.meta.env.PROD,
  maxLocalEntries: 1000,
  remoteEndpoint: import.meta.env.VITE_LOGGING_ENDPOINT,
};

/**
 * 로컬 로그 저장소
 */
class LogStorage {
  private entries: LogEntry[] = [];
  private readonly maxEntries: number;
  
  constructor(maxEntries: number) {
    this.maxEntries = maxEntries;
  }
  
  add(entry: LogEntry): void {
    this.entries.push(entry);
    
    // 최대 개수 초과 시 오래된 로그 제거
    if (this.entries.length > this.maxEntries) {
      this.entries.shift();
    }
  }
  
  getEntries(): readonly LogEntry[] {
    return [...this.entries];
  }
  
  clear(): void {
    this.entries = [];
  }
}

/**
 * 로거 클래스
 * 철칙: 불변성, 타입 안전성, 성능 최적화
 */
class Logger {
  private config: LoggerConfig;
  private storage: LogStorage;
  
  constructor(config: LoggerConfig = defaultConfig) {
    this.config = config;
    this.storage = new LogStorage(config.maxLocalEntries);
  }
  
  /**
   * 로그 레벨 확인
   */
  private shouldLog(level: LogLevel): boolean {
    return LogLevels[level] >= LogLevels[this.config.level];
  }
  
  /**
   * 타임스탬프 생성
   */
  private createTimestamp(): string {
    return new Date().toISOString();
  }
  
  /**
   * 로그 항목 생성
   */
  private createLogEntry(
    level: LogLevel,
    message: string,
    metadata?: LogMetadata
  ): LogEntry {
    return {
      level,
      message,
      timestamp: this.createTimestamp(),
      metadata,
      traceId: metadata?.traceId as string | undefined,
    };
  }
  
  /**
   * 콘솔 출력
   */
  private logToConsole(entry: LogEntry): void {
    if (!this.config.enableConsole) return;
    
    const logMethod = entry.level === 'ERROR' 
      ? console.error 
      : entry.level === 'WARN' 
        ? console.warn 
        : console.log;
    
    const prefix = `[${entry.timestamp}] ${entry.level}:`;
    
    if (entry.metadata) {
      logMethod(prefix, entry.message, entry.metadata);
    } else {
      logMethod(prefix, entry.message);
    }
  }
  
  /**
   * 원격 로깅
   */
  private async logToRemote(entry: LogEntry): Promise<void> {
    if (!this.config.enableRemote || !this.config.remoteEndpoint) return;
    
    try {
      await fetch(this.config.remoteEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(entry),
      });
    } catch (error) {
      // 원격 로깅 실패 시 콘솔에만 출력
      console.error('Failed to send log to remote endpoint:', error);
    }
  }
  
  /**
   * 로그 메시지 출력
   */
  private log(
    level: LogLevel,
    message: string,
    metadata?: LogMetadata
  ): void {
    if (!this.shouldLog(level)) return;
    
    const entry = this.createLogEntry(level, message, metadata);
    
    // 로컬 저장
    this.storage.add(entry);
    
    // 콘솔 출력
    this.logToConsole(entry);
    
    // 원격 로깅 (비동기)
    if (level === 'ERROR' || level === 'WARN') {
      void this.logToRemote(entry);
    }
  }
  
  /**
   * 퍼블릭 로깅 메서드
   */
  debug(message: string, metadata?: LogMetadata): void {
    this.log('DEBUG', message, metadata);
  }
  
  info(message: string, metadata?: LogMetadata): void {
    this.log('INFO', message, metadata);
  }
  
  warn(message: string, metadata?: LogMetadata): void {
    this.log('WARN', message, metadata);
  }
  
  error(message: string, metadata?: LogMetadata): void {
    this.log('ERROR', message, metadata);
  }
  
  /**
   * 설정 업데이트
   */
  updateConfig(newConfig: Partial<LoggerConfig>): void {
    this.config = { ...this.config, ...newConfig };
  }
  
  /**
   * 로그 내역 조회
   */
  getLogs(): readonly LogEntry[] {
    return this.storage.getEntries();
  }
  
  /**
   * 로그 초기화
   */
  clearLogs(): void {
    this.storage.clear();
  }
}

/**
 * 싱글톤 로거 인스턴스
 */
export const logger = new Logger();

/**
 * 함수 실행 시간 측정 데코레이터
 * 철칙: 고차 함수, 타입 안전성
 */
export const withTiming = <T extends readonly unknown[], R>(
  fn: (...args: T) => R,
  name?: string
): ((...args: T) => R) => {
  return (...args: T): R => {
    const start = performance.now();
    const result = fn(...args);
    const end = performance.now();
    
    logger.debug(`Function ${name ?? fn.name} executed in ${end - start}ms`);
    
    return result;
  };
};

/**
 * 비동기 함수 실행 시간 측정
 */
export const withAsyncTiming = <T extends readonly unknown[], R>(
  fn: (...args: T) => Promise<R>,
  name?: string
): ((...args: T) => Promise<R>) => {
  return async (...args: T): Promise<R> => {
    const start = performance.now();
    const result = await fn(...args);
    const end = performance.now();
    
    logger.debug(`Async function ${name ?? fn.name} executed in ${end - start}ms`);
    
    return result;
  };
};