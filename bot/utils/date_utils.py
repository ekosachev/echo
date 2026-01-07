from datetime import datetime
def get_russian_weekday(date_string=None):
    if date_string:
        date_obj = datetime.strptime(date_string, "%Y-%m-%d")
    else:
        date_obj = datetime.now()

    days = ["Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"]
    return days[date_obj.weekday()]