import uuid
from datetime import date, datetime
from typing import List

from pydantic import BaseModel

class AccountCreateModel(BaseModel):
  name: str
  description: str | None

class Account(AccountCreateModel):
  id: uuid.UUID
  userid: uuid.UUID