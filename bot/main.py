import asyncio, logging
from aiogram import Bot, Dispatcher, types, F
from aiogram.filters.command import Command
from config_reader import config
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.storage.memory import MemoryStorage
import aiohttp

authenticated_users = {}
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
                authenticated_users[message.from_user.id] = login
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
    user_login = authenticated_users.get(message.from_user.id, "Не привязан")

    await message.answer(
        f"👤 Ваш аккаунт:\n"
        f"Telegram ID: {message.from_user.id}\n"
        f"Привязан к: {user_login}\n"
        f"Статус: {'✅ Активен' if message.from_user.id in authenticated_users else '❌ Не привязан'}",
        reply_markup=get_main_keyboard()
    )


@dp.message(F.text == "Today")
async def today(message: types.Message):
    if message.from_user.id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    from datetime import datetime, date
    import aiohttp

    today_date = date.today().isoformat()
    telegram_id = message.from_user.id
    await message.answer("Получаю задачи на сегодня...")

    async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=API_TIMEOUT)) as session:
        try:
            endpoints = [
                f"{API_BASE_URL}/api/tasks/?due_date={today_date}",
                f"{API_BASE_URL}/api/tasks/?date={today_date}",
                f"{API_BASE_URL}/tasks/?due_date={today_date}",
                f"{API_BASE_URL}/api/schedule/?date={today_date}",
            ]

            tasks_found = False
            today_tasks = []

            for endpoint in endpoints:
                try:
                    async with session.get(endpoint) as resp:
                        if resp.status == 200:
                            data = await resp.json()

                            # Обрабатываем разные форматы ответа API
                            if isinstance(data, list) and data:
                                today_tasks = data
                                tasks_found = True
                                break
                            elif isinstance(data, dict) and 'results' in data and data['results']:
                                today_tasks = data['results']
                                tasks_found = True
                                break
                            elif isinstance(data, dict) and 'tasks' in data and data['tasks']:
                                today_tasks = data['tasks']
                                tasks_found = True
                                break
                except Exception as e:
                    logging.debug(f"Endpoint {endpoint} failed: {e}")
                    continue

            if tasks_found and today_tasks:
                tasks_text = format_tasks_for_today(today_tasks, today_date)
                await message.answer(tasks_text, reply_markup=get_main_keyboard())
            else:
                await message.answer(
                    f"🎉 Отлично! На сегодня ({datetime.now().strftime('%d.%m.%Y')}) задач нет.\n"
                    "Можете отдохнуть или запланировать новые задачи!",
                    reply_markup=get_main_keyboard()
                )
        except aiohttp.ClientError as e:
            logging.error(f"Connection error: {e}")
            await message.answer(
                "❌ Ошибка подключения к серверу. Попробуйте позже.",
                reply_markup=get_main_keyboard()
            )
        except Exception as e:
            logging.error(f"Unexpected error: {e}")
            await message.answer(
                "❌ Произошла ошибка при получении задач.",
                reply_markup=get_main_keyboard()
            )

def format_tasks_for_today(tasks, today_date):
    from datetime import datetime

    date_obj = datetime.strptime(today_date, "%Y-%m-%d")
    formatted_date = date_obj.strftime("%d.%m.%Y")
    weekday = get_russian_weekday(today_date)

    if not tasks:
        return f"На {weekday} ({formatted_date}) задач нет."
    task_list = [f"{weekday}, {formatted_date}\n"]

    for i, task in enumerate(tasks, 1):
        if isinstance(task, dict):
            title = task.get('title', 'Без названия')
            description = task.get('description', '')
            completed = task.get('completed', False)
            due_time = task.get('due_time', '') or task.get('time', '')
            priority = task.get('priority', 'medium')
            status_emoji = "✅" if completed else "⏳"

            priority_emoji = "🔴"
            if priority == 'low':
                priority_emoji = "🟢"
            elif priority == 'medium':
                priority_emoji = "🟡"
            elif priority == 'high':
                priority_emoji = "🔴"

            task_line = f"{i}. {status_emoji} {priority_emoji} {title}"
            if due_time:
                task_line += f" 🕒 {due_time}"
            if description:
                short_desc = description[:100] + "..." if len(description) > 100 else description
                task_line += f"\n   📝 {short_desc}"
            task_list.append(task_line)
        else:
            task_list.append(f"{i}. 📋 {task}")
    completed_count = sum(1 for task in tasks if isinstance(task, dict) and task.get('completed'))
    total_count = len(tasks)
    task_list.append(f"\n📊 Итого: {completed_count}/{total_count} выполнено")

    if completed_count == total_count and total_count > 0:
        task_list.append("🎉 Все задачи выполнены! Отличная работа!")
    return "\n".join(task_list)

def get_russian_weekday(date_string=None):
    from datetime import datetime
    if date_string:
        date_obj = datetime.strptime(date_string, "%Y-%m-%d")
    else:
        date_obj = datetime.now()
    days = ["Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"]
    return days[date_obj.weekday()]


@dp.message(F.text == "Week")
async def week(message: types.Message):
    if message.from_user.id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    from datetime import datetime, date, timedelta
    import aiohttp
    today = date.today()
    start_of_week = today - timedelta(days=today.weekday())
    end_of_week = start_of_week + timedelta(days=6)
    await message.answer("📅 Получаю задачи на неделю...")

    async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=API_TIMEOUT)) as session:
        try:
            endpoints = [
                f"{API_BASE_URL}/api/tasks/?due_date_after={start_of_week}&due_date_before={end_of_week}",
                f"{API_BASE_URL}/api/tasks/?date_from={start_of_week}&date_to={end_of_week}",
                f"{API_BASE_URL}/tasks/?start_date={start_of_week}&end_date={end_of_week}",
                f"{API_BASE_URL}/api/tasks/",
                f"{API_BASE_URL}/tasks/",
            ]
            tasks_found = False
            all_tasks = []
            week_tasks = []

            for endpoint in endpoints:
                try:
                    async with session.get(endpoint) as resp:
                        if resp.status == 200:
                            data = await resp.json()
                            if isinstance(data, list):
                                all_tasks = data
                                tasks_found = True
                                break
                            elif isinstance(data, dict) and 'results' in data:
                                all_tasks = data['results']
                                tasks_found = True
                                break
                            elif isinstance(data, dict) and 'tasks' in data:
                                all_tasks = data['tasks']
                                tasks_found = True
                                break
                except Exception as e:
                    logging.debug(f"Endpoint {endpoint} failed: {e}")
                    continue
            if tasks_found and all_tasks:
                for task in all_tasks:
                    if isinstance(task, dict):
                        task_date = task.get('due_date') or task.get('date')
                        if task_date:
                            try:
                                task_date_obj = datetime.strptime(task_date.split('T')[0], "%Y-%m-%d").date()
                                if start_of_week <= task_date_obj <= end_of_week:
                                    week_tasks.append(task)
                            except:
                                continue
                if week_tasks:
                    tasks_text = format_tasks_for_week(week_tasks, start_of_week, end_of_week)
                    await message.answer(tasks_text, reply_markup=get_main_keyboard())
                else:
                    await message.answer(
                        f"🎉 На неделю с {start_of_week.strftime('%d.%m')} по {end_of_week.strftime('%d.%m.%Y')} задач нет!\n"
                        "Можете отдохнуть или запланировать новые задачи!",
                        reply_markup=get_main_keyboard()
                    )
            else:
                await message.answer(
                    f"📅 Неделя: {start_of_week.strftime('%d.%m')} - {end_of_week.strftime('%d.%m.%Y')}\n\n"
                    "Не удалось получить задачи. Возможно, API не настроено.",
                    reply_markup=get_main_keyboard()
                )
        except aiohttp.ClientError as e:
            logging.error(f"Connection error: {e}")
            await message.answer(
                "❌ Ошибка подключения к серверу. Попробуйте позже.",
                reply_markup=get_main_keyboard()
            )
        except Exception as e:
            logging.error(f"Unexpected error: {e}")
            await message.answer(
                "❌ Произошла ошибка при получении задач на неделю.",
                reply_markup=get_main_keyboard()
            )

def format_tasks_for_week(tasks, start_date, end_date):
    from datetime import datetime
    if not tasks:
        return f"На неделю с {start_date.strftime('%d.%m')} по {end_date.strftime('%d.%m.%Y')} задач нет."
    sorted_tasks = sorted(tasks, key=lambda x: (
        x.get('due_date', ''),
        not x.get('completed', False)  # Невыполненные задачи сначала
    ))
    task_list = [
        f"📅 Задачи на неделю:",
        f"{start_date.strftime('%d.%m')} - {end_date.strftime('%d.%m.%Y')}\n"
    ]
    for i, task in enumerate(sorted_tasks, 1):
        if isinstance(task, dict):
            title = task.get('title', 'Без названия')
            description = task.get('description', '')
            completed = task.get('completed', False)
            due_date = task.get('due_date', '') or task.get('date', '')
            due_time = task.get('due_time', '') or task.get('time', '')
            priority = task.get('priority', 'medium')
            status_emoji = "✅" if completed else "⏳"

            priority_emoji = "🔴"
            if priority == 'low':
                priority_emoji = "🟢"
            elif priority == 'medium':
                priority_emoji = "🟡"

            date_display = ""
            if due_date:
                try:
                    date_obj = datetime.strptime(due_date.split('T')[0], "%Y-%m-%d")
                    date_display = f"📅 {date_obj.strftime('%d.%m')}"
                except:
                    date_display = f"📅 {due_date}"
            task_line = f"{i}. {status_emoji} {priority_emoji} {title}"
            if date_display:
                task_line += f" {date_display}"
            if due_time:
                task_line += f" 🕒 {due_time}"
            if description:
                short_desc = description[:80] + "..." if len(description) > 80 else description
                task_line += f"\n   📝 {short_desc}"
            task_list.append(task_line)
        else:
            task_list.append(f"{i}. 📋 {task}")
    completed_count = sum(1 for task in tasks if isinstance(task, dict) and task.get('completed'))
    total_count = len(tasks)
    task_list.append(f"\n📊 Итого на неделю: {completed_count}/{total_count} выполнено")

    if completed_count == total_count and total_count > 0:
        task_list.append("🎉 Все задачи на неделю выполнены! Супер!")
    elif completed_count == 0 and total_count > 0:
        task_list.append("💪 Начните выполнять задачи! У вас всё получится!")
    elif completed_count > 0:
        progress = int((completed_count / total_count) * 100)
        task_list.append(f"📈 Прогресс: {progress}% выполнено")

    return "\n".join(task_list)

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