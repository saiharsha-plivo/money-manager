from fastapi import APIRouter, Depends, HTTPException, status
from src.logger import get_logger
from typing import List
import uuid

from src.db.main import get_session
from sqlmodel.ext.asyncio.session import AsyncSession

from .schema import Record, RecordCreate
from .service import RecordService
from src.auth.routes import get_current_user
from src.db.models import User, Account
from src.accounts.service import AccountService


record_router = APIRouter()
record_service = RecordService()
account_service = AccountService()
logger = get_logger(__name__)

logger.info("Record Router Initialized")


async def get_account_and_verify_access(
  account_id: uuid.UUID,
  user: User = Depends(get_current_user),
  session: AsyncSession = Depends(get_session),
) -> Account:
  if user is None:
      raise HTTPException(
          status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid token"
      )
  account = await session.get(Account, account_id)
  if account is None:
      raise HTTPException(
          status_code=status.HTTP_404_NOT_FOUND, detail="Account not found"
      )
  if account.userid != user.id:
      raise HTTPException(
          status_code=status.HTTP_403_FORBIDDEN,
          detail="User does not have access to this account",
      )
  return account


@record_router.get("/{account_id}/records",response_model=List[Record])
async def get_records(account: Account = Depends(get_account_and_verify_access),session: AsyncSession = Depends(get_session)):
  logger.debug(f"got a request for records for account {account.id}")
  records = await record_service.get_records_for_account(account.id, session)
  logger.debug(f"returned records for request for account {account.id}")
  return records


@record_router.post("/{account_id}/records", response_model=Record)
async def create_record(
  record_details: RecordCreate,
  account: Account = Depends(get_account_and_verify_access),
  session: AsyncSession = Depends(get_session),
):
  logger.debug(f"got a request for creating record for account {account.id}")
  record = await record_service.create_record_for_account(
      record_details, account.id, session
  )
  return record


@record_router.delete("/{account_id}/records/{record_id}", status_code=status.HTTP_200_OK)
async def delete_record(record_id: uuid.UUID,
  account: Account = Depends(get_account_and_verify_access),
  session=Depends(get_session),
):
  logger.debug(f"got a request to delete record : {record_id}")
  result = await record_service.delete_record(record_id, session)
  return result
