import asyncio, logging
from aiogram import Bot, Dispatcher, types, F
from aiogram.filters.command import Command
from config_reader import config

logging.basicConfig(level=logging.INFO)
bot = Bot(token=config.bot_token.get_secret_value())

dp = Dispatcher()

def get_start_keyboard():
    kb = [[
        types.KeyboardButton(text='Start')
    ]]
    return types.ReplyKeyboardMarkup(
        keyboard=kb,
        resize_keyboard=True,
    )

@dp.message(Command('start'))
async def cmd_start(message: types.Message):
    await message.answer(
        "Привет! Это Echo, личный тайм-менеджер!\nЧтобы увидеть полный функционал, напиши команду /help"
    )

@dp.message(Command('help'))
async def cmd_help(message: types.Message):
    kb = [
        [types.KeyboardButton(text='Today'),
         types.KeyboardButton(text='Week')],
    ]
    keyboard = types.ReplyKeyboardMarkup(
        keyboard=kb,
        resize_keyboard=True,
        input_field_placeholder="Выберите одну из доступных команд:"
    )
    await message.answer("Какую команду используем?", reply_markup=keyboard)

@dp.message(F.text == "Today")
async def today(message: types.Message):
    await message.reply("День")  # Возвращаем кнопку Start

@dp.message(F.text == "Week")
async def week(message: types.Message):
    await message.reply("Неделя")  # Возвращаем кнопку Start

@dp.message(F.text == "Start")
async def start_button(message: types.Message):
    await message.answer(
        "Привет! Это Echo, личный тайм-менеджер!\nЧтобы увидеть полный функционал, напиши команду /help",
        reply_markup=types.ReplyKeyboardRemove()
    )

@dp.message()
async def any_message(message: types.Message):
    # При ЛЮБОМ сообщении показываем кнопку Start
    await message.answer(
        "Нажмите кнопку Start для начала работы!",
        reply_markup=get_start_keyboard()
    )

async def main():
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())