from aiogram import Router, types, F
from aiogram.filters.command import Command
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
import logging
from config_reader import config
from keyboards.main_kb import get_start_keyboard, get_main_keyboard
from services.api_service import api_service
from global_state import authenticated_users  # Изменено: импорт из глобального состояния

router = Router()
class AuthStates(StatesGroup):
    waiting_login = State()
    waiting_password = State()

@router.message(Command('start'))
async def cmd_start(message: types.Message, state: FSMContext):
    await message.answer(
        "Для привязки аккаунта введите ваш логин от сайта (email-адрес):",
        reply_markup=types.ReplyKeyboardRemove()
    )
    await state.set_state(AuthStates.waiting_login)

@router.message(AuthStates.waiting_login)
async def process_login(message: types.Message, state: FSMContext):
    await state.update_data(login=message.text)
    await message.answer("Введите ваш пароль:")
    await state.set_state(AuthStates.waiting_password)

@router.message(AuthStates.waiting_password)
async def process_password(message: types.Message, state: FSMContext):
    user_data = await state.get_data()
    login = user_data.get('login')
    password = message.text

    try:
        success = await api_service.authenticate_user(login, password)
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
    except Exception as e:
        logging.error(f"Authentication error: {e}")
        await message.answer(
            "Произошла ошибка при подключении к серверу. Попробуйте позже.",
            reply_markup=get_start_keyboard()
        )

    await state.clear()