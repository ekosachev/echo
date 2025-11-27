from aiogram import Router, types, F
import logging
from datetime import datetime, date, timedelta
from keyboards.main_kb import get_main_keyboard, get_start_keyboard
from services.api_service import api_service
from main import authenticated_users
from utils.date_utils import get_russian_weekday

router = Router()
@router.message(F.text == "Today")
async def today(message: types.Message):
    if message.from_user.id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    await message.answer("Получаю задачи на сегодня...")

    try:
        today_tasks = await api_service.get_today_tasks()
        if today_tasks:
            today_date = date.today().isoformat()
            tasks_text = format_tasks_for_today(today_tasks, today_date)
            await message.answer(tasks_text, reply_markup=get_main_keyboard())
        else:
            await message.answer(
                f"🎉 Отлично! На сегодня ({datetime.now().strftime('%d.%m.%Y')}) задач нет.\n"
                "Можете отдохнуть или запланировать новые задачи!",
                reply_markup=get_main_keyboard()
            )
    except Exception as e:
        logging.error(f"Today tasks error: {e}")
        await message.answer(
            "❌ Ошибка при получении задач.",
            reply_markup=get_main_keyboard()
        )
@router.message(F.text == "Week")
async def week(message: types.Message):
    if message.from_user.id not in authenticated_users:
        await message.answer(
            "❌ Сначала привяжите аккаунт через /start",
            reply_markup=get_start_keyboard()
        )
        return
    await message.answer("📅 Получаю задачи на неделю...")

    try:
        week_tasks = await api_service.get_week_tasks()
        today = date.today()
        start_of_week = today - timedelta(days=today.weekday())
        end_of_week = start_of_week + timedelta(days=6)

        if week_tasks:
            tasks_text = format_tasks_for_week(week_tasks, start_of_week, end_of_week)
            await message.answer(tasks_text, reply_markup=get_main_keyboard())
        else:
            await message.answer(
                f"🎉 На неделю с {start_of_week.strftime('%d.%m')} по {end_of_week.strftime('%d.%m.%Y')} задач нет!\n"
                "Можете отдохнуть или запланировать новые задачи!",
                reply_markup=get_main_keyboard()
            )
    except Exception as e:
        logging.error(f"Week tasks error: {e}")
        await message.answer(
            "❌ Ошибка при получении задач на неделю.",
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
