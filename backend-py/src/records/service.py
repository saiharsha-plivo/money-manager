import uuid
from sqlmodel import select
from sqlmodel.ext.asyncio.session import AsyncSession
from fastapi import HTTPException

from src.db.models import Record
from src.utils.currency import Currency
from src.utils.types import RecordType
from .schema import RecordCreate


class RecordService:
    def __init__(self):
      self.currency = Currency()
      self.record_type = RecordType()

    async def get_records_for_account(self, account_id: uuid.UUID, session: AsyncSession):
      statement = select(Record).where(Record.account_id == account_id)
      result = await session.exec(statement)
      return result.all()

    async def create_record_for_account(
      self, record_details: RecordCreate, account_id: uuid.UUID, session: AsyncSession
    ) -> Record:
      if str(record_details.currency_id) not in self.currency.get_currency_list():
          raise HTTPException(status_code=400, detail="Invalid currency id")
      if str(record_details.type_id) not in self.record_type.get_types():
          raise HTTPException(status_code=400, detail="Invalid type id")

      new_record = Record(**record_details.model_dump(), account_id=account_id)
      session.add(new_record)
      await session.commit()
      return new_record

    async def delete_record(self, record_id: uuid.UUID, session: AsyncSession):
      statement = (
        select(Record).where(Record.id == record_id)
      )
      record = await session.scalar(statement)

      if not record:
        raise HTTPException(status_code=404, detail="Record not found")

      await session.delete(record)
      await session.commit()

      return {"detail": f"Record {record_id} deleted successfully"}
