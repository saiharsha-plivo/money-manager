from sqlmodel import SQLModel, Field, Relationship, Column
import uuid
import datetime
import sqlalchemy.dialects.postgresql as pg
from typing import List, Optional


class User(SQLModel, table=True):
    __tablename__ = "users"
    id: uuid.UUID = Field(
        sa_column=Column(pg.UUID, primary_key=True, nullable=False, default=uuid.uuid4)
    )
    username: str = Field(sa_column=Column(pg.VARCHAR, nullable=False,unique=True))
    email: str = Field(sa_column=Column(pg.VARCHAR, nullable=False))
    role: str = Field(sa_column=Column(pg.VARCHAR, nullable=False, default="user"))
    verified: bool = Field(default=False)
    password_hash: str = Field(sa_column=Column(pg.VARCHAR, nullable=False),exclude=True)
    created_at: datetime.datetime = Field(
        sa_column=Column(pg.TIMESTAMP, default=datetime.datetime.utcnow)
    )

    accounts: List["Account"] = Relationship(back_populates="user", sa_relationship_kwargs={"lazy": "selectin"})


class Account(SQLModel, table=True):
    __tablename__ = "accounts"
    id: uuid.UUID = Field(
        sa_column=Column(pg.UUID, primary_key=True, nullable=False, default=uuid.uuid4)
    )
    name: str = Field(sa_column=Column(pg.VARCHAR, nullable=False))
    description: Optional[str] = Field(default=None)
    userid: uuid.UUID = Field(foreign_key="users.id", exclude=True)
    user: Optional["User"] = Relationship(back_populates="accounts")
    records: List["Record"] = Relationship(back_populates="account", sa_relationship_kwargs={"lazy": "selectin"})


class Record(SQLModel, table=True):
  __tablename__ = "records"
  id: uuid.UUID = Field(
    sa_column=Column(pg.UUID, primary_key=True, nullable=False, default=uuid.uuid4)
  )
  account_id: uuid.UUID = Field(foreign_key="accounts.id")
  account: Optional["Account"] = Relationship(back_populates="records")
  amount: float = Field(sa_column=Column(pg.FLOAT, nullable=False))
  description: Optional[str] = Field(default=None)
  type_id: int = Field(sa_column=Column(pg.INTEGER, nullable=False))
  currency_id: int = Field(sa_column=Column(pg.INTEGER, nullable=False))
  created_at: datetime.datetime = Field(
    sa_column=Column(pg.TIMESTAMP, default=datetime.datetime.utcnow)
  )
  updated_at: datetime.datetime = Field(
    sa_column=Column(pg.TIMESTAMP, default=datetime.datetime.utcnow)
  )

  comments: List["Comment"] = Relationship(back_populates="record", sa_relationship_kwargs={"lazy": "selectin"})


class Comment(SQLModel, table=True):
    __tablename__ = "comments"
    id: uuid.UUID = Field(
        sa_column=Column(pg.UUID, primary_key=True, nullable=False, default=uuid.uuid4)
    )
    description: str = Field(sa_column=Column(pg.TEXT, nullable=False))
    record_id: uuid.UUID = Field(foreign_key="records.id")
    record: Optional["Record"] = Relationship(back_populates="comments")
    created_at: datetime.datetime = Field(
        sa_column=Column(pg.TIMESTAMP, default=datetime.datetime.utcnow)
    )
    updated_at: datetime.datetime = Field(
        sa_column=Column(pg.TIMESTAMP, default=datetime.datetime.utcnow)
    )