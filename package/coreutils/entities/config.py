from dataclasses import dataclass

@dataclass
class Config:
    promptStyle: dict[str, int]
