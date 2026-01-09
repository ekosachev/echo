from .auth_handlers import router as auth_router
from .main_handlers import router as main_router
from .schedule_handlers import router as schedule_router

__all__ = ['auth_router', 'main_router', 'schedule_router']