import uuid
from datetime import datetime
from pydantic import BaseModel, Field

class RecordBase(BaseModel):
    amount: int
    description: str | None = None
    currency_id: int = Field(gt=0)
    type_id: int = Field(gt=0)

class RecordCreate(RecordBase):
    pass

class Record(RecordBase):
    id: uuid.UUID
    account_id: uuid.UUID
    created_at: datetime

    class Config:
        orm_mode = True
