"""config_reader module handles reading the .env file"""
from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import SecretStr

class Settings(BaseSettings):
    bot_token: SecretStr
    model_config = SettingsConfigDict(env_file='.env', env_file_encoding='utf-8',extra='ignore')
    api_base_url: str = "http://localhost:8000"
    api_timeout: int = 30
config = Settings()