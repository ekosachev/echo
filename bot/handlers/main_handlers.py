from aiogram import Router, types, F
from aiogram.filters.command import Command
from keyboards.main_kb import get_start_keyboard, get_main_keyboard
from global_state import authenticated_users

router = Router()
@router.message(Command('help'))
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

@router.message(F.text == "My Account")
async def my_account(message: types.Message):
    user_login = authenticated_users.get(message.from_user.id, "Не привязан")

    await message.answer(
        f"👤 Ваш аккаунт:\n"
        f"Telegram ID: {message.from_user.id}\n"
        f"Привязан к: {user_login}\n"
        f"Статус: {'✅ Активен' if message.from_user.id in authenticated_users else '❌ Не привязан'}",
        reply_markup=get_main_keyboard()
    )

@router.message(F.text == "Start")
async def start_button(message: types.Message, state):
    from handlers.auth_handlers import cmd_start
    await cmd_start(message, state)

@router.message()
async def any_message(message: types.Message):
    await message.answer(
        "Нажмите кнопку Start для начала работы!",
        reply_markup=get_start_keyboard()
    )