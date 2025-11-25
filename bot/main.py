import asyncio, logging
from aiogram import Bot, Dispatcher, types, F
from aiogram.filters.command import Command
from config_reader import config
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.storage.memory import MemoryStorage
import aiohttp


logging.basicConfig(level=logging.INFO)
bot = Bot(token=config.bot_token.get_secret_value())
storage = MemoryStorage()
dp = Dispatcher(storage=storage)


class AuthStates(StatesGroup):
    waiting_login = State()
    waiting_password = State()

API_BASE_URL = config.api_base_url
API_TIMEOUT = config.api_timeout

def get_start_keyboard():
    kb = [[
        types.KeyboardButton(text='Start')
    ]]
    return types.ReplyKeyboardMarkup(
        keyboard=kb,
        resize_keyboard=True,
    )

def get_main_keyboard():
    kb = [
        [types.KeyboardButton(text='Today'), types.KeyboardButton(text='Week')],
        [types.KeyboardButton(text='My Account')]
    ]
    return types.ReplyKeyboardMarkup(
        keyboard=kb,
        resize_keyboard=True,
    )

@dp.message(Command('start'))
async def cmd_start(message: types.Message, state: FSMContext):
    await message.answer(
        "Для привязки аккаунта введите ваш логин от сайта:",
        reply_markup=types.ReplyKeyboardRemove()
    )
    await state.set_state(AuthStates.waiting_login)


@dp.message(AuthStates.waiting_login)
async def process_login(message: types.Message, state: FSMContext):
    await state.update_data(login=message.text)
    await message.answer("Введите ваш пароль:")
    await state.set_state(AuthStates.waiting_password)

@dp.message(AuthStates.waiting_password)
async def process_password(message: types.Message, state: FSMContext):
    user_data = await state.get_data()
    login = user_data.get('login')
    password = message.text
    async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=API_TIMEOUT)) as session:
        try:
            endpoints = [
                f"{API_BASE_URL}/api/auth/login/",
                f"{API_BASE_URL}/auth/login/",
                f"{API_BASE_URL}/api/token/auth/",
                f"{API_BASE_URL}/api/auth/",
            ]
            success = False
            for endpoint in endpoints:
                try:
                    async with session.post(
                            endpoint,
                            json={'username': login, 'password': password},
                            headers={'Content-Type': 'application/json'}
                    ) as resp:
                        if resp.status == 200:
                            data = await resp.json()
                            success = True
                            break
                except Exception as e:
                    logging.debug(f"Endpoint {endpoint} failed: {e}")
                    continue
            if success:
                await message.answer(
                    "Аккаунт успешно привязан! Теперь вы можете использовать функции бота.",
                    reply_markup=get_main_keyboard()
                )
            else:
                await message.answer(
                    "Не удалось привязать аккаунт. Проверьте логин и пароль.\nПопробуйте снова: /start",
                    reply_markup=get_start_keyboard()
                )

        except aiohttp.ClientError as e:
            logging.error(f"Connection error: {e}")
            await message.answer(
                "Ошибка подключения к серверу. Попробуйте позже.",
                reply_markup=get_start_keyboard()
            )
        except Exception as e:
            logging.error(f"Unexpected error: {e}")
            await message.answer(
                "Произошла непредвиденная ошибка.",
                reply_markup=get_start_keyboard()
            )

    await state.clear()

async def try_authentication(login: str, password: str, telegram_id: int) -> bool:
    token = await get_auth_token(login, password)
    if token:
        return await update_profile_with_telegram_id(token, telegram_id)
    return await direct_authentication(login, password, telegram_id)

async def get_auth_token(login: str, password: str) -> str:
    try:
        async with aiohttp.ClientSession() as session:
            endpoints = [
                f"{API_BASE_URL}/api-token-auth/",
                f"{API_BASE_URL}/auth/token/",
                f"{API_BASE_URL}/token/",
                f"{API_BASE_URL}/login/",
            ]

            for endpoint in endpoints:
                try:
                    async with session.post(
                            endpoint,
                            json={'username': login, 'password': password},
                            headers={'Content-Type': 'application/json'}
                    ) as resp:
                        if resp.status == 200:
                            data = await resp.json()
                            return data.get('token') or data.get('access_token')
                except:
                    continue
    except Exception as e:
        logging.error(f"Token auth error: {e}")
    return None

async def update_profile_with_telegram_id(token: str, telegram_id: int) -> bool:
    try:
        async with aiohttp.ClientSession() as session:
            headers = {
                'Authorization': f'Token {token}',
                'Content-Type': 'application/json'
            }
            endpoints = [
                f"{API_BASE_URL}/profile/",
                f"{API_BASE_URL}/user/profile/",
                f"{API_BASE_URL}/me/",
            ]
            for endpoint in endpoints:
                try:
                    async with session.get(endpoint, headers=headers) as resp:
                        if resp.status == 200:
                            profile_data = await resp.json()
                            update_data = {**profile_data, 'telegram_id': telegram_id}
                            async with session.patch(
                                    endpoint,
                                    json=update_data,
                                    headers=headers
                            ) as update_resp:
                                return update_resp.status in [200, 201]
                except:
                    continue
    except Exception as e:
        logging.error(f"Profile update error: {e}")

    return False


async def direct_authentication(login: str, password: str, telegram_id: int) -> bool:
    try:
        async with aiohttp.ClientSession() as session:
            endpoints = [
                f"{API_BASE_URL}/auth/",
                f"{API_BASE_URL}/login/",
                f"{API_BASE_URL}/api/auth/",
            ]
            for endpoint in endpoints:
                try:
                    payload = {
                        'username': login,
                        'password': password,
                        'telegram_id': telegram_id
                    }
                    async with session.post(
                            endpoint,
                            json=payload,
                            headers={'Content-Type': 'application/json'}
                    ) as resp:
                        if resp.status in [200, 201]:
                            return True
                except:
                    continue
    except Exception as e:
        logging.error(f"Direct auth error: {e}")

    return False

@dp.message(Command('help'))
async def cmd_help(message: types.Message):
    await message.answer(
        "Доступные команды:\n"
        "/start - привязка аккаунта\n"
        "/help - помощь\n\n"
        "Кнопки:\n"
        "Today - расписание на сегодня\n"
        "Week - расписание на неделю\n"
        "My Account - информация об аккаунте",
        reply_markup=get_main_keyboard()
    )

@dp.message(F.text == "My Account")
async def my_account(message: types.Message):
    await message.answer(
        f"Ваш Telegram ID: {message.from_user.id}\n"
        "Для управления аккаунтом используйте сайт.",
        reply_markup=get_main_keyboard()
    )

@dp.message(F.text == "Today")
async def today(message: types.Message):
    await message.reply("Расписание на сегодня:\n")  # Возвращаем кнопку Start

@dp.message(F.text == "Week")
async def week(message: types.Message):
    await message.reply("Расписание на неделю:\n")  # Возвращаем кнопку Start

@dp.message(F.text == "Start")
async def start_button(message: types.Message,state: FSMContext):
    await cmd_start(message, state)

@dp.message()
async def any_message(message: types.Message):
    # При ЛЮБОМ сообщении показываем кнопку Start
    await message.answer(
        "Нажмите кнопку Start для начала работы!",
        reply_markup=get_start_keyboard()
    )

async def main():
    logging.info(f"Starting bot with API base URL: {API_BASE_URL}")
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())