"""api_service module handles interactions with the backend api"""
import logging
from typing import Optional

import aiohttp
from config_reader import config

class TokenExpiredError(Exception):
    pass

class APIService:
    def __init__(self):
        self.base_url = config.api_base_url
        self.timeout = config.api_timeout

    async def authenticate_user(self, login: str, password: str) -> tuple[bool, Optional[str]]:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoint = f'{self.base_url}/api/v1/auth/login'
            try:
                async with session.post(
                        endpoint,
                        json={'email': login, 'password': password},
                        headers={'Content-Type': 'application/json'}
                ) as resp:
                    if resp.status == 200:
                        body = await resp.json()
                        return True, body["token"]
                    if resp.status == 401:
                        return False, None
            except aiohttp.ClientError as e:
                logging.error(f"Auth endpoint {endpoint} failed: {e}")

            return False, None
        
    async def get_calendars(self, token: str) -> Optional[list]:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            endpoint = f'{self.base_url}/api/v1/calendars'
            try:
                async with session.get(endpoint, headers={'Authorization': f'Bearer {token}'}) as resp:
                    if resp.ok:
                        data = await resp.json()
                        return data['calendars']
                    if resp.status == 401:
                        raise TokenExpiredError
            except aiohttp.ClientError as e:
                logging.error(f"Failed to fetch calendars: {e}")
        return None
    
    async def _get_all_events(self, token: str) -> Optional[dict]:
        calendars = await self.get_calendars(token)
        if not calendars:
            return {} if calendars is not None else None
        
        base_endpoint = self.base_url + '/api/v1/calendars/{}/events'
        result = {}

        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            for calendar in calendars:
                endpoint = base_endpoint.format(calendar['id'])
                try:
                    async with session.get(endpoint, headers={'Authorization': f'Bearer {token}'}) as resp:
                        if resp.ok:
                            data = await resp.json()
                            result[(calendar["id"], calendar['name'])] = data['events']
                        if resp.status == 401:
                            raise TokenExpiredError

                except aiohttp.ClientError as e:
                    logging.error(f'Failed to get events for calendar {calendar["id"]}: {e}')
                    return None
        return result
            

    async def get_today_tasks(self, token: str) -> Optional[dict]:
        from datetime import datetime, date
        from zoneinfo import ZoneInfo
        
        events_by_calendar = await self._get_all_events(token)
        if not events_by_calendar:
            return events_by_calendar

        process_timestamp = lambda dt, tz: datetime.fromisoformat(dt).astimezone(tz=ZoneInfo(tz))

        events_filtered = {
            calendar_info: list(filter(
                lambda e: process_timestamp(e["start_at"], e["timezone"]).date() == date.today(),
                events
            )) for calendar_info, events in events_by_calendar.items()
        }

        return events_filtered

        

    async def get_week_tasks(self, token) -> Optional[dict]:
        from datetime import datetime, date, timedelta
        from zoneinfo import ZoneInfo
        
        events_by_calendar = await self._get_all_events(token)
        if not events_by_calendar:
            return events_by_calendar

        process_timestamp = lambda dt, tz: datetime.fromisoformat(dt).astimezone(tz=ZoneInfo(tz))

        events_filtered = {
            calendar_info: list(filter(
                lambda e: date.today() <= process_timestamp(e["start_at"], e["timezone"]).date() <= date.today() + timedelta(days=7),
                events
            )) for calendar_info, events in events_by_calendar.items()
        }

        return events_filtered

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