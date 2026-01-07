from aiogram import types
def get_start_keyboard():
    kb = [[types.KeyboardButton(text='Start')]]
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