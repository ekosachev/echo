import asyncio
import logging
from aiogram import Bot, Dispatcher
from aiogram.fsm.storage.memory import MemoryStorage

from config_reader import config
from handlers.auth_handlers import router as auth_router
from handlers.main_handlers import router as main_router
from handlers.schedule_handlers import router as schedule_router

logging.basicConfig(level=logging.INFO)
bot = Bot(token=config.bot_token.get_secret_value())

dp.include_router(auth_router)
dp.include_router(main_router)
dp.include_router(schedule_router)

async def main():
    logging.info(f"Starting bot with API base URL: {config.api_base_url}")
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())