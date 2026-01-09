from aiogram import Router, types, F
import logging
from datetime import datetime, date, timedelta
from keyboards.main_kb import get_main_keyboard, get_start_keyboard
from services.api_service import TokenExpiredError, api_service
from global_state import authenticated_users
from utils.date_utils import get_russian_weekday

router = Router()
@router.message(F.text == "My Calendars")
async def my_calendars(message: types.Message):
    user_id = message.from_user.id
    if user_id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    
    await message.answer("Получаю список доступных календарей...")

    try:

        calendars = await api_service.get_calendars(authenticated_users[user_id])

        if calendars is None:
            await message.answer("Ошибка при получении календарей")
            return
        
        await message.answer(
            "Доступные календари:\n"
            + ("Ничего не нашлось :(" if not calendars else '\n'.join(
                ' - ' + calendar["name"] for calendar in calendars
            ))
        )

    except TokenExpiredError as _:
        await message.answer("Сессия истекла. Пожалуйста, войдите в аккаунт через /start")

@router.message(F.text == "Today")
async def today(message: types.Message):
    from datetime import datetime
    from zoneinfo import ZoneInfo
    user_id = message.from_user.id
    if user_id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    await message.answer("Получаю события на сегодня...")

    try:
        today_tasks = await api_service.get_today_tasks(authenticated_users[user_id])
        if today_tasks is None:
            await message.answer("Ошибка при получении событий")
            return

        process_timestamp = lambda dt, tz: datetime.fromisoformat(dt).astimezone(tz=ZoneInfo(tz)).strftime("%H:%M:%S")

        await message.answer(
            "События на сегодня:\n"
            + ( 
                "Ничего не нашлось :(" if not today_tasks else '\n'.join(
                    sum(
                        [[f'{calendar_name}:'] +
                        [f' - {event["title"]}' +
                        f' {process_timestamp(event["start_at"], event["timezone"])}' +
                        f' - {process_timestamp(event["end_at"], event["timezone"])}'
                            for event in events] if events else [f'{calendar_name}: Нет событий на сегодня']
                        for (_, calendar_name), events in today_tasks.items()],
                        start=[]
                    ) 
                )
            )
        )

    except TokenExpiredError as _:
        await message.answer("Сессия истекла. Пожалуйста, войдите в аккаунт через /start")
        
@router.message(F.text == "Week")
async def week(message: types.Message):
    from datetime import datetime
    from zoneinfo import ZoneInfo
    user_id = message.from_user.id
    if user_id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    await message.answer("Получаю события на неделю...")

    try:
        today_tasks = await api_service.get_week_tasks(authenticated_users[user_id])
        if today_tasks is None:
            await message.answer("Ошибка при получении событий")
            return

        process_timestamp = lambda dt, tz: datetime.fromisoformat(dt).astimezone(tz=ZoneInfo(tz)).strftime("%m/%d %H:%M:%S")

        await message.answer(
            "События на неделю:\n"
            + ( 
                "Ничего не нашлось :(" if not today_tasks else '\n'.join(
                    sum(
                        [[f'{calendar_name}:'] +
                        [f' - {event["title"]}' +
                        f' {process_timestamp(event["start_at"], event["timezone"])}' +
                        f' - {process_timestamp(event["end_at"], event["timezone"])}'
                            for event in events] if events else [f'{calendar_name}: Нет событий на сегодня']
                        for (_, calendar_name), events in today_tasks.items()],
                        start=[]
                    ) 
                )
            )
        )

    except TokenExpiredError as _:
        await message.answer("Сессия истекла. Пожалуйста, войдите в аккаунт через /start")

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
