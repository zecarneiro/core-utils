from pydantic import BaseModel


class Config(BaseModel):
    promptStyle: dict[str, int]
