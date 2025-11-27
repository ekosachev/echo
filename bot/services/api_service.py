import aiohttp
import logging
from config_reader import config

class APIService:
    def __init__(self):
        self.base_url = config.api_base_url
        self.timeout = config.api_timeout

    async def authenticate_user(self, login: str, password: str) -> bool:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoints = [
                f"{self.base_url}/api/auth/login/",
                f"{self.base_url}/auth/login/",
                f"{self.base_url}/api/token/auth/",
                f"{self.base_url}/api/auth/",
            ]
            for endpoint in endpoints:
                try:
                    async with session.post(
                            endpoint,
                            json={'username': login, 'password': password},
                            headers={'Content-Type': 'application/json'}
                    ) as resp:
                        if resp.status == 200:
                            return True
                except Exception as e:
                    logging.debug(f"Auth endpoint {endpoint} failed: {e}")
                    continue

            return False

    async def get_today_tasks(self) -> list:
        from datetime import date
        today_date = date.today().isoformat()

        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoints = [
                f"{self.base_url}/api/tasks/?due_date={today_date}",
                f"{self.base_url}/api/tasks/?date={today_date}",
                f"{self.base_url}/tasks/?due_date={today_date}",
                f"{self.base_url}/api/schedule/?date={today_date}",
            ]

            for endpoint in endpoints:
                try:
                    async with session.get(endpoint) as resp:
                        if resp.status == 200:
                            data = await resp.json()
                            return self._parse_tasks_response(data)
                except Exception as e:
                    logging.debug(f"Tasks endpoint {endpoint} failed: {e}")
                    continue

            return []

    async def get_week_tasks(self) -> list:
        from datetime import date, timedelta, datetime
        today = date.today()
        start_of_week = today - timedelta(days=today.weekday())
        end_of_week = start_of_week + timedelta(days=6)

        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoints = [
                f"{self.base_url}/api/tasks/?due_date_after={start_of_week}&due_date_before={end_of_week}",
                f"{self.base_url}/api/tasks/?date_from={start_of_week}&date_to={end_of_week}",
                f"{self.base_url}/tasks/?start_date={start_of_week}&end_date={end_of_week}",
                f"{self.base_url}/api/tasks/",
                f"{self.base_url}/tasks/",
            ]

            for endpoint in endpoints:
                try:
                    async with session.get(endpoint) as resp:
                        if resp.status == 200:
                            data = await resp.json()
                            all_tasks = self._parse_tasks_response(data)
                            if "due_date_after" not in endpoint and "date_from" not in endpoint:
                                return self._filter_tasks_by_week(all_tasks, start_of_week, end_of_week)
                            else:
                                return all_tasks
                except Exception as e:
                    logging.debug(f"Week tasks endpoint {endpoint} failed: {e}")
                    continue

            return []

    def _parse_tasks_response(self, data) -> list:
        if isinstance(data, list):
            return data
        elif isinstance(data, dict) and 'results' in data:
            return data['results']
        elif isinstance(data, dict) and 'tasks' in data:
            return data['tasks']
        else:
            return []

    def _filter_tasks_by_week(self, tasks: list, start_date, end_date) -> list:
        from datetime import datetime

        week_tasks = []
        for task in tasks:
            if isinstance(task, dict):
                task_date = task.get('due_date') or task.get('date')
                if task_date:
                    try:
                        task_date_obj = datetime.strptime(task_date.split('T')[0], "%Y-%m-%d").date()
                        if start_date <= task_date_obj <= end_date:
                            week_tasks.append(task)
                    except:
                        continue
        return week_tasks

    async def get_user_profile(self, telegram_id: int) -> dict:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoints = [
                f"{self.base_url}/api/profile/?telegram_id={telegram_id}",
                f"{self.base_url}/profile/?telegram_id={telegram_id}",
                f"{self.base_url}/api/user/profile/",
            ]
            for endpoint in endpoints:
                try:
                    async with session.get(endpoint) as resp:
                        if resp.status == 200:
                            return await resp.json()
                except Exception as e:
                    logging.debug(f"Profile endpoint {endpoint} failed: {e}")
                    continue

            return {}

    async def update_telegram_id(self, login: str, password: str, telegram_id: int) -> bool:
        token = await self._get_auth_token(login, password)
        if not token:
            return False

        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            headers = {
                'Authorization': f'Token {token}',
                'Content-Type': 'application/json'
            }

            endpoints = [
                f"{self.base_url}/api/profile/",
                f"{self.base_url}/profile/",
                f"{self.base_url}/api/user/profile/",
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
                except Exception as e:
                    logging.debug(f"Update profile endpoint {endpoint} failed: {e}")
                    continue

            return False

    async def _get_auth_token(self, login: str, password: str) -> str:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoints = [
                f"{self.base_url}/api-token-auth/",
                f"{self.base_url}/auth/token/",
                f"{self.base_url}/token/",
                f"{self.base_url}/api/auth/token/",
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
                except Exception as e:
                    logging.debug(f"Token endpoint {endpoint} failed: {e}")
                    continue

            return None
api_service = APIService()