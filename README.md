# OMS (Order Management System)

## 프로젝트 개요
이 프로젝트는 주문 관리 시스템(OMS)을 구현한 것입니다.

## 프로젝트 구조
```
OMS/
├── README.md           # 프로젝트 문서
├── requirements.txt    # Python 의존성 파일
├── src/               # 소스 코드
│   ├── api/          # API 관련 코드
│   ├── core/         # 핵심 비즈니스 로직
│   ├── models/       # 데이터 모델
│   └── utils/        # 유틸리티 함수
├── tests/            # 테스트 코드
├── docs/             # 문서
└── scripts/          # 유틸리티 스크립트
```

## 기술 스택
- Backend: Python
- Database: PostgreSQL
- API: FastAPI
- Frontend: React (예정)

## 시작하기
1. 저장소 클론
```bash
git clone [repository-url]
cd OMS
```

2. 가상환경 설정
```bash
python -m venv venv
source venv/bin/activate  # Linux/Mac
# 또는
.\venv\Scripts\activate  # Windows
```

3. 의존성 설치
```bash
pip install -r requirements.txt
```

## 개발 가이드라인
- 모든 코드는 PEP 8 스타일 가이드를 준수합니다.
- 모든 함수와 클래스에는 한국어로 된 상세한 문서화 주석이 포함되어야 합니다.
- 테스트 커버리지는 최소 80% 이상을 유지해야 합니다.

## 라이선스
이 프로젝트는 MIT 라이선스 하에 배포됩니다.