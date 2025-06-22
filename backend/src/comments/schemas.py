import uuid
from datetime import datetime
from pydantic import BaseModel


class CommentBase(BaseModel):
    description: str


class CommentCreate(CommentBase):
    pass


class CommentUpdate(CommentBase):
    pass


class Comment(CommentBase):
    id: uuid.UUID
    record_id: uuid.UUID
    created_at: datetime
    updated_at: datetime

    class Config:
        orm_mode = True 